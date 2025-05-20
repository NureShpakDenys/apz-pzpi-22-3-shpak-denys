package com.example.newapp.ui.company

import android.app.Activity
import android.os.Bundle
import android.view.View
import android.widget.AdapterView
import android.widget.ArrayAdapter
import android.widget.FrameLayout
import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import androidx.lifecycle.map
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.models.AddUserToCompanyRequest
import com.example.newapp.data.models.User
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityAddUserToCompanyBinding
import com.example.newapp.databinding.ActivityRouteDetailsBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.async
import kotlinx.coroutines.launch

class AddUserToCompanyActivity : BaseActivity() {

    private lateinit var binding: ActivityAddUserToCompanyBinding
    private lateinit var sessionManager: SessionManager
    private var companyId: Int = -1

    private var companyUserIds: List<Int> = emptyList()
    private var availableUsersForCompany: List<User> = emptyList()

    private var selectedUserId: Int = -1
    private var selectedRole: String = "user"

    companion object {
        const val EXTRA_COMPANY_ID = "company_id"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_add_user_to_company)
        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)

        binding = ActivityAddUserToCompanyBinding.bind(contentFrame.getChildAt(0))

        sessionManager = SessionManager(this)

        companyId = intent.getIntExtra(EXTRA_COMPANY_ID, -1)
        if (companyId == -1) {
            Toast.makeText(this, getString(R.string.company_id_not_found), Toast.LENGTH_LONG).show()
            finish()
            return
        }

        setupRoleSpinner()
        fetchInitialData()

        binding.btnAddUser.setOnClickListener {
            handleSubmit()
        }
    }

    private fun setupRoleSpinner() {
        val rolesValues = resources.getStringArray(R.array.company_roles_array_values)
        val roleDisplayNames = resources.getStringArray(R.array.company_roles_array_display)

        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_item, roleDisplayNames)
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        binding.spinnerRole.adapter = adapter

        val defaultRoleIndex = rolesValues.indexOf("user")
        if (defaultRoleIndex != -1) {
            binding.spinnerRole.setSelection(defaultRoleIndex)
            selectedRole = rolesValues[defaultRoleIndex]
        }

        binding.spinnerRole.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
            override fun onItemSelected(parent: AdapterView<*>?, view: View?, position: Int, id: Long) {
                selectedRole = rolesValues[position]
            }
            override fun onNothingSelected(parent: AdapterView<*>?) {}
        }
    }

    private fun fetchInitialData() {
        showLoading(true)
        binding.tvErrorAddUser.visibility = View.GONE

        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@AddUserToCompanyActivity)


                val usersResponse = apiService.getAllUsers()
                val companyResponse = apiService.getCompany(companyId)

                companyUserIds = companyResponse.users.map { it.id }

                availableUsersForCompany = usersResponse.filterNot { companyUserIds.contains(it.id) }
                updateUserSpinner()
            } catch (e: Exception) {
                showError(getString(R.string.error_fetching_data_exception, e.localizedMessage ?: "Unknown error"))
                e.printStackTrace()
            } finally {
                showLoading(false)
            }
        }
    }

    private fun updateUserSpinner() {
        if (availableUsersForCompany.isEmpty()) {
            binding.spinnerUser.adapter = null
            binding.spinnerUser.isEnabled = false
            binding.spinnerUser.prompt = getString(R.string.no_users_available_to_add)
            selectedUserId = -1
            binding.btnAddUser.isEnabled = false
            return
        }

        binding.spinnerUser.isEnabled = true
        binding.btnAddUser.isEnabled = true

        val userNames = availableUsersForCompany.map { it.name }
        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_item, userNames)
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        binding.spinnerUser.adapter = adapter

        selectedUserId = if (availableUsersForCompany.isNotEmpty()) {
            availableUsersForCompany[0].id
        } else { -1 }


        binding.spinnerUser.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
            override fun onItemSelected(parent: AdapterView<*>?, view: View?, position: Int, id: Long) {
                if (position < availableUsersForCompany.size) {
                    selectedUserId = availableUsersForCompany[position].id
                }
            }
            override fun onNothingSelected(parent: AdapterView<*>?) {
                selectedUserId = -1
            }
        }
    }

    private fun handleSubmit() {
        if (selectedUserId == -1) {
            Toast.makeText(this, getString(R.string.please_select_a_user), Toast.LENGTH_SHORT).show()
            return
        }

        showLoading(true)
        binding.tvErrorAddUser.visibility = View.GONE

        lifecycleScope.launch {
            val addUserRequest = AddUserToCompanyRequest(userId = selectedUserId, role = selectedRole)

            try {
                val apiService = RetrofitClient.getInstance(this@AddUserToCompanyActivity)
                val response = apiService.addUserToCompany(companyId, addUserRequest)

                if (response.isSuccessful) {
                    Toast.makeText(this@AddUserToCompanyActivity, getString(R.string.user_added_successfully), Toast.LENGTH_LONG).show()
                    setResult(RESULT_OK)
                    finish()
                } else {
                    val errorBody = response.errorBody()?.string()
                    val errorMessage = if (!errorBody.isNullOrBlank()) {
                      errorBody
                    } else {
                        response.message()
                    }
                    showError(getString(R.string.error_adding_user_to_company_api, errorMessage))
                }
            } catch (e: Exception) {
                showError(getString(R.string.error_adding_user_to_company_exception, e.localizedMessage ?: "Unknown error"))
                e.printStackTrace()
            } finally {
                showLoading(false)
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBarAddUser.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnAddUser.isEnabled = !isLoading
        binding.spinnerUser.isEnabled = !isLoading
        binding.spinnerRole.isEnabled = !isLoading
    }

    private fun showError(message: String) {
        binding.tvErrorAddUser.text = message
        binding.tvErrorAddUser.visibility = View.VISIBLE
    }
}