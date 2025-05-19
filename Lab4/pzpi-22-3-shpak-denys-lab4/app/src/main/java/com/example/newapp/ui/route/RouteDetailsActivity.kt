package com.example.newapp.ui.route

import android.annotation.SuppressLint
import android.app.Activity
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AlertDialog
import androidx.recyclerview.widget.LinearLayoutManager
import com.example.newapp.data.models.RouteResponse
import com.example.newapp.data.models.WeatherAlertResponse
import com.example.newapp.data.models.Waypoint
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityRouteDetailsBinding
import com.example.newapp.ui.adapter.WaypointAdapter
import com.example.newapp.ui.base.BaseActivity
import com.example.newapp.ui.waypoint.CreateWaypointActivity
import com.example.newapp.ui.waypoint.EditWaypointActivity
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import com.example.newapp.R

class RouteDetailsActivity : BaseActivity() {

    private lateinit var binding: ActivityRouteDetailsBinding
    private lateinit var adapter: WaypointAdapter
    private var route: RouteResponse? = null
    private var currentRouteId: Int = -1

    private val createWaypointLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        if (result.resultCode == RESULT_OK) {
            if (currentRouteId != -1) {
                fetchRoute(currentRouteId)
            }
        }
    }

    private val editWaypointLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
        if (result.resultCode == RESULT_OK) {
            if (currentRouteId != -1) {
                fetchRoute(currentRouteId)
                Toast.makeText(this, "Waypoint updated.", Toast.LENGTH_SHORT).show()
            }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_route_details)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)

        binding = ActivityRouteDetailsBinding.bind(contentFrame.getChildAt(0))

        adapter = WaypointAdapter(emptyList(), onDeleteClick = { waypoint ->
                showDeleteWaypointDialog(waypoint)
            }, onEditClick = { waypoint ->
                val intent = Intent(this, EditWaypointActivity::class.java)
                intent.putExtra("waypoint_id", waypoint.ID)
                editWaypointLauncher.launch(intent)
            }
        )

        binding.rvWaypoints.layoutManager = LinearLayoutManager(this)
        binding.rvWaypoints.adapter = adapter

        currentRouteId = intent.getIntExtra("route_id", -1)
        if (currentRouteId == -1) {
            Toast.makeText(this, "Error: Route ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        fetchRoute(currentRouteId)
        binding.btnDeleteRoute.setOnClickListener { confirmDelete(currentRouteId) }
        binding.btnWeatherAlert.setOnClickListener { fetchWeather(currentRouteId) }
        binding.btnAddWaypoint.setOnClickListener {
            val intent = Intent(this, CreateWaypointActivity::class.java)
            intent.putExtra("route_id", currentRouteId)
            createWaypointLauncher.launch(intent)
        }
    }

    private fun fetchRoute(id: Int) {
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@RouteDetailsActivity)
                val res = apiService.getRoute(id)
                Log.d("RouteDetails", "Fetched route: $res")

                withContext(Dispatchers.Main) {
                    route = res
                    binding.tvRouteName.text = res.Name
                    binding.tvRouteStatus.text = res.Status
                    binding.tvRouteDetails.text = res.Details
                    binding.tvCompanyName.text = res.company?.Name
                    adapter.updateData(res.waypoints)
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Log.e("RouteDetails", "Error fetching route: ", e)
                    Toast.makeText(this@RouteDetailsActivity, "Error loading route details: $e", Toast.LENGTH_LONG).show()
                }
            }
        }
    }

    private fun confirmDelete(id: Int) {
        AlertDialog.Builder(this)
            .setMessage("Confirm deletion?")
            .setPositiveButton("Yes") { _, _ -> deleteRoute(id) }
            .setNegativeButton("No", null)
            .show()
    }

    private fun deleteRoute(id: Int) {
        val apiService = RetrofitClient.getInstance(this@RouteDetailsActivity)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                apiService.deleteRoute(id)
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@RouteDetailsActivity, "Route deleted", Toast.LENGTH_SHORT).show()
                    setResult(RESULT_OK)
                    finish()
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@RouteDetailsActivity, "Error deleting route: ${e.message}", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun showDeleteWaypointDialog(waypoint: Waypoint) {
        AlertDialog.Builder(this)
            .setMessage("Confirm delete waypoint '${waypoint.Name}'?")
            .setPositiveButton("Yes") { _, _ -> deleteWaypoint(waypoint) }
            .setNegativeButton("No", null)
            .show()
    }

    private fun deleteWaypoint(wp: Waypoint) {
        val apiService = RetrofitClient.getInstance(this@RouteDetailsActivity)
        CoroutineScope(Dispatchers.IO).launch {
            try {
                apiService.deleteWaypoint(wp.ID)
                if (currentRouteId != -1) {
                    fetchRoute(currentRouteId)
                }
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@RouteDetailsActivity, "Waypoint '${wp.Name}' deleted", Toast.LENGTH_SHORT).show()
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@RouteDetailsActivity, "Error deleting waypoint: ${e.message}", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    @SuppressLint("SetTextI18n")
    private fun fetchWeather(id: Int) {
        binding.tvWeatherInfo.visibility = View.GONE
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val apiService = RetrofitClient.getInstance(this@RouteDetailsActivity)
                val res: WeatherAlertResponse = apiService.getWeatherAlert(id)
                withContext(Dispatchers.Main) {
                    if (res.alerts != null) {
                        binding.tvWeatherInfo.text = "Type: ${res.alerts.type}\n" +
                                "Message: ${res.alerts.message}\n" +
                                "Details: ${res.alerts.details}"
                        binding.tvWeatherInfo.visibility = View.VISIBLE
                    } else {
                        binding.tvWeatherInfo.text = "No weather alerts available."
                        binding.tvWeatherInfo.visibility = View.VISIBLE
                    }
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    binding.tvWeatherInfo.text = "Could not fetch weather information."
                    binding.tvWeatherInfo.visibility = View.VISIBLE
                    Log.e("RouteDetails", "Error fetching weather: ", e)
                }
            }
        }
    }
}
