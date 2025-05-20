package com.example.newapp.ui.company

import android.app.Activity
import android.os.Bundle
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import com.example.newapp.databinding.ActivityEditCompanyBinding
import com.example.newapp.ui.base.BaseActivity
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.data.models.Company
import com.example.newapp.data.models.UpdateCompanyRequest
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.IOException

class EditCompanyActivity : BaseActivity() {

    private lateinit var binding: ActivityEditCompanyBinding
    private var companyId: Int = -1
    private var currentCompany: Company? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityEditCompanyBinding.inflate(layoutInflater)
        setContentView(binding.root)

        companyId = intent.getIntExtra("company_id", -1)
        if (companyId == -1) {
            Toast.makeText(this, "Error: Company ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        fetchCompanyDetails(companyId)

        binding.btnSaveCompany.setOnClickListener {
            handleSaveChanges()
        }
    }

    private fun fetchCompanyDetails(id: Int) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditCompanyActivity)
                val company = apiService.getCompany(id)
                withContext(Dispatchers.Main) {
                    val simpleCompany = Company(
                        id = company.id,
                        name = company.name,
                        address = company.address
                    )
                    currentCompany = simpleCompany
                    populateUi(simpleCompany)
                    showLoading(false)
                }
            } catch (e: IOException) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Network Error: Could not load company details. ${e.message}")
                    finish()
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Error loading company: ${e.message}")
                    finish()
                }
            }
        }
    }

    private fun populateUi(company: Company) {
        binding.etCompanyName.setText(company.name)
        binding.etCompanyAddress.setText(company.address)
    }

    private fun handleSaveChanges() {
        val name = binding.etCompanyName.text.toString().trim()
        val address = binding.etCompanyAddress.text.toString().trim()

        if (name.isEmpty()) {
            binding.etCompanyName.error = "Company name is required"
            return
        }

        if (address.isEmpty()) {
            binding.etCompanyAddress.error = "Address is required"
            return
        }

        val updateRequest = UpdateCompanyRequest(name = name, address = address)

        updateCompany(updateRequest)
    }

    private fun updateCompany(request: UpdateCompanyRequest) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditCompanyActivity)
                apiService.updateCompany(companyId, request)
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    Toast.makeText(this@EditCompanyActivity, "Company updated successfully!", Toast.LENGTH_SHORT).show()
                    setResult(Activity.RESULT_OK)
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
                    showErrorDialog("Error updating company: ${e.message}")
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnSaveCompany.isEnabled = !isLoading
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
