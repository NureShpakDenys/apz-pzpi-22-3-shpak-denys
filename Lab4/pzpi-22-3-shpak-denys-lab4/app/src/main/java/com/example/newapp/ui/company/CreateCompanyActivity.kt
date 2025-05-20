package com.example.newapp.ui.company

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import com.example.newapp.R
import com.example.newapp.data.models.CreateCompanyRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityCreateCompanyBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.IOException

class CreateCompanyActivity : BaseActivity() {

    private lateinit var binding: ActivityCreateCompanyBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_create_company)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)
        if (contentFrame != null) {
            binding = ActivityCreateCompanyBinding.bind(contentFrame.getChildAt(0))
        } else {
            Log.e("CreateCompanyActivity", "Content frame is null or empty. Cannot bind.")
            Toast.makeText(this, "Error initializing layout.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        binding.btnCreateCompany.setOnClickListener {
            handleCreateCompany()
        }
    }

    private fun handleCreateCompany() {
        val name = binding.etCompanyName.text.toString().trim()
        val address = binding.etCompanyAddress.text.toString().trim()

        if (name.isEmpty()) {
            binding.etCompanyName.error = "Company name is required"
            return
        } else {
            binding.etCompanyName.error = null
        }

        if (address.isEmpty()) {
            binding.etCompanyAddress.error = "Company address is required"
            return
        } else {
            binding.etCompanyAddress.error = null
        }

        val createCompanyRequest = CreateCompanyRequest(
            name = name,
            address = address
        )

        createCompany(createCompanyRequest)
    }

    private fun createCompany(request: CreateCompanyRequest) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@CreateCompanyActivity)
                apiService.createCompany(request)

                withContext(Dispatchers.Main) {
                    showLoading(false)
                    Toast.makeText(
                        this@CreateCompanyActivity,
                        "Company created successfully!",
                        Toast.LENGTH_SHORT
                    ).show()

                    val resultIntent = Intent()
                    setResult(RESULT_OK, resultIntent)
                    finish()
                }
            } catch (e: IOException) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Network Error: ${e.message}")
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Error creating company: ${e.message}")
                    Log.e("CreateCompany", "Error: ", e)
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnCreateCompany.isEnabled = !isLoading
        binding.etCompanyName.isEnabled = !isLoading
        binding.etCompanyAddress.isEnabled = !isLoading
    }

    private fun showErrorDialog(message: String) {
        AlertDialog.Builder(this)
            .setTitle("Error")
            .setMessage(message)
            .setPositiveButton("OK", null)
            .show()
    }
}
