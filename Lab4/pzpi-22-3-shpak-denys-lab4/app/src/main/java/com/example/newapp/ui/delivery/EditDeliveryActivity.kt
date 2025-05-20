package com.example.newapp.ui.delivery

import android.os.Bundle
import android.widget.*
import com.example.newapp.R
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.*
import retrofit2.HttpException

class EditDeliveryActivity : BaseActivity() {

    private lateinit var etDate: EditText
    private lateinit var etDuration: EditText
    private lateinit var spinnerStatus: Spinner
    private lateinit var btnSave: Button
    private var deliveryId: Int = 0
    private lateinit var token: String
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_edit_delivery)

        etDate = findViewById(R.id.etDate)
        etDuration = findViewById(R.id.etDuration)
        spinnerStatus = findViewById(R.id.spinnerStatus)
        btnSave = findViewById(R.id.btnSave)

        deliveryId = intent.getIntExtra("deliveryId", 0)
        token = intent.getStringExtra("token") ?: ""

        setupStatusSpinner()
        loadDelivery()

        btnSave.setOnClickListener {
            updateDelivery()
        }
    }

    private fun setupStatusSpinner() {
        val statuses = listOf("not_started", "in_progress", "completed")
        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_dropdown_item, statuses)
        spinnerStatus.adapter = adapter
    }

    private fun loadDelivery() {
        scope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditDeliveryActivity)
                val response = apiService.getDelivery(deliveryId)
                withContext(Dispatchers.Main) {
                    etDate.setText(response.date.toString().split("T")[0])
                    etDuration.setText(response.duration)
                    val statusPosition = (spinnerStatus.adapter as ArrayAdapter<String>).getPosition(response.status)
                    spinnerStatus.setSelection(statusPosition)
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditDeliveryActivity, "Error loading delivery", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun updateDelivery() {
        val date = etDate.text.toString()
        val duration = etDuration.text.toString()
        val status = spinnerStatus.selectedItem.toString()

        if (date.isEmpty() || duration.isEmpty()) {
            Toast.makeText(this, "All fields are required", Toast.LENGTH_SHORT).show()
            return
        }

        scope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditDeliveryActivity)
                apiService.updateDelivery(
                    deliveryId,
                    mapOf(
                        "date" to date,
                        "duration" to duration,
                        "status" to status
                    )
                )
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditDeliveryActivity, "Delivery updated", Toast.LENGTH_SHORT).show()
                    finish()
                }
            } catch (e: HttpException) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditDeliveryActivity, "Error updating delivery", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        scope.cancel()
    }
}
