package com.example.newapp.ui.waypoint

import android.app.Activity
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Toast
import com.example.newapp.R
import androidx.appcompat.app.AlertDialog
import com.example.newapp.data.models.UpdateWaypointRequest
import com.example.newapp.data.models.Waypoint // Your Waypoint data class
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityEditWaypointBinding
import com.example.newapp.ui.base.BaseActivity // Assuming you have a BaseActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.IOException
import android.widget.FrameLayout
import androidx.core.view.isNotEmpty

class EditWaypointActivity : BaseActivity() {

    private lateinit var binding: ActivityEditWaypointBinding
    private var waypointId: Int = -1
    private var currentWaypoint: Waypoint? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_edit_waypoint)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)
        if (contentFrame != null && contentFrame.isNotEmpty()) {
            binding = ActivityEditWaypointBinding.bind(contentFrame.getChildAt(0))
        } else {
            Log.e("EditWaypointActivity", "Content frame is null or empty. Cannot bind.")
            Toast.makeText(this, "Error initializing layout.", Toast.LENGTH_LONG).show()
            finish()
            return
        }
        waypointId = intent.getIntExtra("waypoint_id", -1)
        if (waypointId == -1) {
            Toast.makeText(this, "Error: Waypoint ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        fetchWaypointDetails(waypointId)

        binding.btnSaveWaypoint.setOnClickListener {
            handleSaveChanges()
        }
    }

    private fun fetchWaypointDetails(id: Int) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditWaypointActivity)
                val waypoint = apiService.getWaypoint(id)
                withContext(Dispatchers.Main) {
                    currentWaypoint = waypoint
                    populateUi(waypoint)
                    showLoading(false)
                }
            } catch (e: IOException) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Network Error: Could not load waypoint details. ${e.message}")
                    finish()
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    showLoading(false)
                    showErrorDialog("Error loading waypoint: ${e.message}")
                    android.util.Log.e("EditWaypoint", "Error fetching waypoint: ", e)
                    finish()
                }
            }
        }
    }

    private fun populateUi(waypoint: Waypoint) {
        binding.etWaypointName.setText(waypoint.Name)
        binding.etDeviceSerial.setText(waypoint.DeviceSerial)
        binding.etSendDataFrequency.setText(waypoint.SendDataFrequency.toString())
        binding.cbGetWeatherAlerts.isChecked = waypoint.GetWeatherAlerts
    }

    private fun handleSaveChanges() {
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

        val updateWaypointRequest = UpdateWaypointRequest(
            name = name,
            deviceSerial = deviceSerial,
            sendDataFrequency = sendDataFrequency,
            getWeatherAlerts = getWeatherAlerts
        )

        updateWaypoint(updateWaypointRequest)
    }

    private fun updateWaypoint(request: UpdateWaypointRequest) {
        showLoading(true)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditWaypointActivity)
                apiService.updateWaypoint(waypointId, request)

                withContext(Dispatchers.Main) {
                    showLoading(false)
                    Toast.makeText(this@EditWaypointActivity, "Waypoint updated successfully!", Toast.LENGTH_SHORT).show()
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
                    showErrorDialog("Error updating waypoint: ${e.message}")
                    android.util.Log.e("EditWaypoint", "Error: ", e)
                }
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        binding.progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnSaveWaypoint.isEnabled = !isLoading
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