package com.example.newapp.ui.waypoint

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import com.example.newapp.R
import com.example.newapp.data.models.CreateWaypointRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityCreateWaypointBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.IOException
import androidx.core.view.isNotEmpty

class CreateWaypointActivity : BaseActivity() {

    private lateinit var binding: ActivityCreateWaypointBinding
    private var routeId: Int = -1

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_create_waypoint)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)
        if (contentFrame != null && contentFrame.isNotEmpty()) {
            binding = ActivityCreateWaypointBinding.bind(contentFrame.getChildAt(0))
        } else {
            Log.e("CreateWaypointActivity", "Content frame is null or empty. Cannot bind.")
            Toast.makeText(this, "Error initializing layout.", Toast.LENGTH_LONG).show()
            finish()
            return
        }
        routeId = intent.getIntExtra("route_id", -1)
        if (routeId == -1) {
            Toast.makeText(this, "Error: Route ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        binding.btnCreateWaypoint.setOnClickListener {
            handleCreateWaypoint()
        }
    }

    private fun handleCreateWaypoint() {
        val name = binding.etWaypointName.text.toString().trim()
        val deviceSerial = binding.etDeviceSerial.text.toString().trim()
        val sendDataFrequencyStr = binding.etSendDataFrequency.text.toString().trim()
        val getWeatherAlerts = binding.cbGetWeatherAlerts.isChecked

        if (name.isEmpty()) {
            binding.etWaypointName.error = "Waypoint name is required"
            return
        } else {
            binding.etWaypointName.error = null
        }

        if (deviceSerial.isEmpty()) {
            binding.etDeviceSerial.error = "Device serial is required"
            return
        } else {
            binding.etDeviceSerial.error = null
        }

        if (sendDataFrequencyStr.isEmpty()) {
            binding.etSendDataFrequency.error = "Data sending frequency is required"
            return
        } else {
            binding.etSendDataFrequency.error = null
        }

        val sendDataFrequency = sendDataFrequencyStr.toIntOrNull()
        if (sendDataFrequency == null || sendDataFrequency <= 0) {
            binding.etSendDataFrequency.error = "Invalid frequency value"
            return
        } else {
            binding.etSendDataFrequency.error = null
        }

        val createWaypointRequest = CreateWaypointRequest(
            routeId = routeId,
            name = name,
            deviceSerial = deviceSerial,
            sendDataFrequency = sendDataFrequency,
            getWeatherAlerts = getWeatherAlerts
        )

        createWaypoint(createWaypointRequest)
    }

    private fun createWaypoint(request: CreateWaypointRequest) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@CreateWaypointActivity)
                val response = apiService.createWaypoint(request)

                withContext(Dispatchers.Main) {
                    showLoading(false)
                    Toast.makeText(
                        this@CreateWaypointActivity,
                        "Waypoint created successfully!",
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
                    showErrorDialog("Error creating waypoint: ${e.message}")
                    Log.e("CreateWaypoint", "Error: ", e)
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnCreateWaypoint.isEnabled = !isLoading
        binding.etWaypointName.isEnabled = !isLoading
        binding.etDeviceSerial.isEnabled = !isLoading
        binding.etSendDataFrequency.isEnabled = !isLoading
        binding.cbGetWeatherAlerts.isEnabled = !isLoading
    }

    private fun showErrorDialog(message: String) {
        AlertDialog.Builder(this)
            .setTitle("Error")
            .setMessage(message)
            .setPositiveButton("OK", null)
            .show()
    }
}