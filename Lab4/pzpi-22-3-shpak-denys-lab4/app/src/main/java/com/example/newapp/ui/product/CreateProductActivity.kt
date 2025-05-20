package com.example.newapp.ui.product

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.core.view.isNotEmpty
import com.example.newapp.R
import com.example.newapp.data.models.CreateProductRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityCreateProductBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.IOException

class CreateProductActivity : BaseActivity() {

    private lateinit var binding: ActivityCreateProductBinding
    private var deliveryId: Int = -1

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_create_product)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)
        if (contentFrame != null && contentFrame.isNotEmpty()) {
            binding = ActivityCreateProductBinding.bind(contentFrame.getChildAt(0))
        } else {
            Log.e("CreateProductActivity", "Content frame is null or empty.")
            Toast.makeText(this, "Error initializing layout.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        deliveryId = intent.getIntExtra("delivery_id", -1)
        if (deliveryId == -1) {
            Toast.makeText(this, "Delivery ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        binding.btnCreateProduct.setOnClickListener {
            handleCreateProduct()
        }
    }

    private fun handleCreateProduct() {
        val name = binding.etProductName.text.toString().trim()
        val productType = binding.spinnerProductType.selectedItem.toString()
        val weightStr = binding.etWeight.text.toString().trim()

        if (name.isEmpty()) {
            binding.etProductName.error = "Product name is required"
            return
        }

        if (weightStr.isEmpty()) {
            binding.etWeight.error = "Weight is required"
            return
        }

        val weight = weightStr.toFloatOrNull()
        if (weight == null || weight <= 0) {
            binding.etWeight.error = "Invalid weight"
            return
        }

        val request = CreateProductRequest(
            deliveryID = deliveryId,
            name = name,
            productType = productType,
            weight = weight
        )

        createProduct(request)
    }

    private fun createProduct(request: CreateProductRequest) {
        showLoading(true)

        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@CreateProductActivity)
                apiService.createProduct(request)

                withContext(Dispatchers.Main) {
                    showLoading(false)
                    Toast.makeText(this@CreateProductActivity, "Product created!", Toast.LENGTH_SHORT).show()
                    setResult(RESULT_OK, Intent())
                    finish()
                }
            } catch (e: IOException) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Network error: ${e.message}")
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Error: ${e.message}")
                    Log.e("CreateProduct", "Exception", e)
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnCreateProduct.isEnabled = !isLoading
        binding.etProductName.isEnabled = !isLoading
        binding.etWeight.isEnabled = !isLoading
        binding.spinnerProductType.isEnabled = !isLoading
    }

    private fun showErrorDialog(message: String) {
        AlertDialog.Builder(this)
            .setTitle("Error")
            .setMessage(message)
            .setPositiveButton("OK", null)
            .show()
    }
}
