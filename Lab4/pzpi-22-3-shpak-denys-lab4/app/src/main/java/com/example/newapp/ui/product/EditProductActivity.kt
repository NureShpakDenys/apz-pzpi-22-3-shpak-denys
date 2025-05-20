package com.example.newapp.ui.product

import android.os.Bundle
import android.widget.ArrayAdapter
import android.widget.Toast
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityEditProductBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import retrofit2.HttpException

class EditProductActivity : BaseActivity() {

    private lateinit var binding: ActivityEditProductBinding
    private var productId: Int = 0
    private var deliveryId: Int = 0
    private lateinit var token: String
    private var currentSystem: String = "metric" 

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityEditProductBinding.inflate(layoutInflater)
        setContentView(binding.root)

        
        productId = intent.getIntExtra("productId", 0)
        token = intent.getStringExtra("token") ?: ""
        currentSystem = intent.getStringExtra("system") ?: "metric"

        setupProductTypeDropdown()
        loadProduct()

        binding.btnSave.setOnClickListener {
            updateProduct()
        }
    }

    private fun setupProductTypeDropdown() {
        val types = listOf("Fruits", "Vegetables", "Frozen Foods", "Dairy Products", "Meat")
        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_dropdown_item, types)
        binding.spinnerProductType.adapter = adapter
    }

    private fun loadProduct() {
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditProductActivity)
                val res = apiService.getProductById(productId)
                withContext(Dispatchers.Main) {
                    binding.etProductName.setText(res.name)
                    deliveryId = res.deliveryId
                    val productType = res.productCategory.name
                    val spinnerPosition = (binding.spinnerProductType.adapter as ArrayAdapter<String>).getPosition(productType)
                    binding.spinnerProductType.setSelection(spinnerPosition)

                    val weight = if (currentSystem == "imperial") {
                        convertWeight(res.weight.toDouble(), "toImperial")
                    } else res.weight
                    binding.etWeight.setText(weight.toString())
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditProductActivity, "Failed to load product", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun updateProduct() {
        val name = binding.etProductName.text.toString().trim()
        val type = binding.spinnerProductType.selectedItem.toString()
        val weightInput = binding.etWeight.text.toString().toDoubleOrNull()

        if (name.isEmpty() || weightInput == null) {
            Toast.makeText(this, "Please fill in all fields", Toast.LENGTH_SHORT).show()
            return
        }

        val weightToSend = if (currentSystem == "imperial") {
            convertWeight(weightInput, "toMetric")
        } else weightInput

        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditProductActivity)
                apiService.updateProduct(
                    productId = productId,
                    request = mapOf(
                        "name" to name,
                        "product_type" to type,
                        "weight" to weightToSend
                    )
                )
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditProductActivity, "Product updated", Toast.LENGTH_SHORT).show()
                    finish() 
                }
            } catch (e: HttpException) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditProductActivity, "Error updating product", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun convertWeight(value: Double, direction: String): Double {
        return if (direction == "toImperial") value * 2.20462 else value / 2.20462
    }
}
