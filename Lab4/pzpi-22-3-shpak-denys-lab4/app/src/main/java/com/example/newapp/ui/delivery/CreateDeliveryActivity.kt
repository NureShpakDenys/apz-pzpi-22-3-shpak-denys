package com.example.newapp.ui.delivery

import android.app.DatePickerDialog
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.MenuItem
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import com.example.newapp.R
import com.example.newapp.data.models.CreateDeliveryRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityCreateDeliveryBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.launch
import java.text.SimpleDateFormat
import java.util.Calendar
import java.util.Locale

class CreateDeliveryActivity : BaseActivity() {

    private lateinit var binding: ActivityCreateDeliveryBinding
    private var companyId: Int = -1
    private val calendar: Calendar = Calendar.getInstance()

    companion object {
        const val EXTRA_COMPANY_ID = "company_id_for_delivery"
        const val RESULT_EXTRA_DELIVERY_ID = "new_delivery_id"
        private const val TAG = "CreateDeliveryActivity"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_create_delivery)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)

        binding = ActivityCreateDeliveryBinding.bind(contentFrame.getChildAt(0))
        companyId = intent.getIntExtra(EXTRA_COMPANY_ID, -1)
        if (companyId == -1) {
            Log.e(TAG, "Company ID not provided to CreateDeliveryActivity.")
            Toast.makeText(this, getString(R.string.error_company_id_missing_delivery), Toast.LENGTH_LONG).show()
            finish()
            return
        }

        setupDatePicker()

        binding.btnCreateDelivery.setOnClickListener {
            validateAndSubmit()
        }
    }

    private fun setupDatePicker() {
        val dateSetListener = DatePickerDialog.OnDateSetListener { _, year, monthOfYear, dayOfMonth ->
            calendar.set(Calendar.YEAR, year)
            calendar.set(Calendar.MONTH, monthOfYear)
            calendar.set(Calendar.DAY_OF_MONTH, dayOfMonth)
            updateDateInView()
        }

        binding.etDeliveryDate.setOnClickListener {
            DatePickerDialog(
                this,
                dateSetListener,
                calendar.get(Calendar.YEAR),
                calendar.get(Calendar.MONTH),
                calendar.get(Calendar.DAY_OF_MONTH)
            ).show()
        }

        updateDateInView()
    }

    private fun updateDateInView() {
        val myFormat = "yyyy-MM-dd"
        val sdf = SimpleDateFormat(myFormat, Locale.US)
        binding.etDeliveryDate.setText(sdf.format(calendar.time))
    }

    private fun validateAndSubmit() {
        val dateStr = binding.etDeliveryDate.text.toString()

        if (dateStr.isBlank()) {
            binding.etDeliveryDate.error = getString(R.string.error_date_required)
            Toast.makeText(this, getString(R.string.error_date_required), Toast.LENGTH_SHORT).show()
            return
        }
        binding.etDeliveryDate.error = null

        val deliveryRequest = CreateDeliveryRequest(
            companyId = companyId,
            date = dateStr
        )
        createDelivery(deliveryRequest)
    }

    private fun createDelivery(request: CreateDeliveryRequest) {
        showLoading(true)
        binding.tvErrorCreateDelivery.visibility = View.GONE

        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@CreateDeliveryActivity)
                val response = apiService.createDelivery(request)

                if (response.isSuccessful && response.body() != null) {
                    val newDelivery = response.body()!!
                    Log.i(TAG, "Delivery created successfully: ID = ${newDelivery.id}, Date = ${newDelivery.date}")
                    Toast.makeText(this@CreateDeliveryActivity, getString(R.string.delivery_created_successfully, newDelivery.id), Toast.LENGTH_LONG).show()

                    val resultIntent = Intent()
                    resultIntent.putExtra(RESULT_EXTRA_DELIVERY_ID, newDelivery.id)
                    setResult(RESULT_OK, resultIntent)
                    finish()
                } else {
                    val errorBody = response.errorBody()?.string() ?: "Unknown error from server"
                    Log.e(TAG, "Error creating delivery: ${response.code()} - $errorBody")
                    binding.tvErrorCreateDelivery.text = getString(R.string.error_creating_delivery_api, errorBody)
                    binding.tvErrorCreateDelivery.visibility = View.VISIBLE
                    Toast.makeText(this@CreateDeliveryActivity, getString(R.string.error_creating_delivery_api, errorBody), Toast.LENGTH_LONG).show()
                }
            } catch (e: Exception) {
                Log.e(TAG, "Exception during delivery creation: ${e.message}", e)
                binding.tvErrorCreateDelivery.text = getString(R.string.error_creating_delivery_api, e.localizedMessage ?: "Network error")
                binding.tvErrorCreateDelivery.visibility = View.VISIBLE
                Toast.makeText(this@CreateDeliveryActivity, getString(R.string.error_creating_delivery_api, e.localizedMessage ?: "Network error"), Toast.LENGTH_LONG).show()
            } finally {
                showLoading(false)
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBarCreateDelivery.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnCreateDelivery.text = if (isLoading) getString(R.string.button_creating_delivery) else getString(R.string.button_create_delivery)
        binding.btnCreateDelivery.isEnabled = !isLoading
        binding.etDeliveryDate.isEnabled = !isLoading
    }
}