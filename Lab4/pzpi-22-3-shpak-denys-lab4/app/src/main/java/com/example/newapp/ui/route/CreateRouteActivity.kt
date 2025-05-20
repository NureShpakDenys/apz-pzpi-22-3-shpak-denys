package com.example.newapp.ui.route 

import android.app.Activity
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import com.example.newapp.R
import com.example.newapp.data.SessionManager 
import com.example.newapp.data.models.CreateRouteRequest
import com.example.newapp.data.models.CreateRouteResponse
import com.example.newapp.data.models.RouteResponse 
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityCreateRouteBinding 
import com.example.newapp.databinding.ActivityRouteDetailsBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.launch

class CreateRouteActivity : BaseActivity() {

    private lateinit var binding: ActivityCreateRouteBinding
    private var companyId: Int = -1

    companion object {
        const val EXTRA_COMPANY_ID = "company_id_for_route"
        const val RESULT_EXTRA_ROUTE_ID = "new_route_id"
        private const val TAG = "CreateRouteActivity"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_create_route)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)

        binding = ActivityCreateRouteBinding.bind(contentFrame.getChildAt(0))

        companyId = intent.getIntExtra(EXTRA_COMPANY_ID, -1)
        if (companyId == -1) {
            Log.e(TAG, "Company ID not passed to CreateRouteActivity.")
            Toast.makeText(this, getString(R.string.error_company_id_missing_route), Toast.LENGTH_LONG).show()
            finish()
            return
        }

        Log.d(TAG, "CreateRouteActivity started for company ID: $companyId")

        binding.btnCreateRoute.setOnClickListener {
            handleSubmit()
        }
    }

    private fun handleSubmit() {
        val routeName = binding.etRouteName.text.toString().trim()

        if (routeName.isEmpty()) {
            binding.etRouteName.error = getString(R.string.error_route_name_required)
            return
        } else {
            binding.etRouteName.error = null
        }

        showLoading(true)
        binding.tvErrorCreateRoute.visibility = View.GONE

        val createRouteRequest = CreateRouteRequest(companyId = companyId, name = routeName)

        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@CreateRouteActivity)
                val response = apiService.createRoute(createRouteRequest)

                if (response.isSuccessful && response.body() != null) {
                    val newRoute: CreateRouteResponse = response.body()!!
                    Log.i(TAG, "Route created successfully: ID=${newRoute.id}, Name=${newRoute.name}")
                    Toast.makeText(this@CreateRouteActivity, getString(R.string.route_created_successfully), Toast.LENGTH_LONG).show()

                    val resultIntent = Intent()
                    resultIntent.putExtra(RESULT_EXTRA_ROUTE_ID, newRoute.id)
                    setResult(RESULT_OK, resultIntent)
                    finish()
                } else {
                    val errorBody = response.errorBody()?.string() ?: response.message() ?: "Unknown error"
                    Log.e(TAG, "Failed to create route. Code: ${response.code()}, Error: $errorBody")
                    showError(getString(R.string.error_creating_route_api, errorBody))
                }
            } catch (e: Exception) {
                Log.e(TAG, "Exception while creating route: ${e.message}", e)
                showError(getString(R.string.error_creating_route_exception, e.localizedMessage ?: "Unknown error"))
            } finally {
                showLoading(false)
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBarCreateRoute.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnCreateRoute.isEnabled = !isLoading
        binding.etRouteName.isEnabled = !isLoading
        binding.btnCreateRoute.text = if (isLoading) getString(R.string.button_creating_route) else getString(R.string.button_create_route)
    }

    private fun showError(message: String) {
        binding.tvErrorCreateRoute.text = message
        binding.tvErrorCreateRoute.visibility = View.VISIBLE
    }
}