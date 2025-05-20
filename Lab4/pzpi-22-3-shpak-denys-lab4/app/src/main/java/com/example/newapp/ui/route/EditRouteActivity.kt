package com.example.newapp.ui.route

import android.os.Bundle
import android.widget.*
import com.example.newapp.R
import com.example.newapp.data.network.RetrofitClient
import kotlinx.coroutines.*
import retrofit2.HttpException
import com.example.newapp.ui.base.BaseActivity

class EditRouteActivity : BaseActivity() {

    private lateinit var etRouteName: EditText
    private lateinit var btnSave: Button
    private var routeId: Int = 0
    private lateinit var token: String
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_edit_route)

        etRouteName = findViewById(R.id.etRouteName)
        btnSave = findViewById(R.id.btnSave)

        routeId = intent.getIntExtra("routeId", 0)
        token = intent.getStringExtra("token") ?: ""

        loadRoute()

        btnSave.setOnClickListener {
            updateRoute()
        }
    }

    private fun loadRoute() {
        scope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditRouteActivity)
                val response = apiService.getRoute(routeId)
                withContext(Dispatchers.Main) {
                    etRouteName.setText(response.Name)
                }
            } catch (e: Exception) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditRouteActivity, "Error loading route data", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun updateRoute() {
        val name = etRouteName.text.toString().trim()
        if (name.isEmpty()) {
            Toast.makeText(this, "Route name cannot be empty", Toast.LENGTH_SHORT).show()
            return
        }

        scope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@EditRouteActivity)
                apiService.updateRoute(
                    routeId,
                    mapOf("name" to name)
                )
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditRouteActivity, "Route updated", Toast.LENGTH_SHORT).show()
                    finish()
                }
            } catch (e: HttpException) {
                withContext(Dispatchers.Main) {
                    Toast.makeText(this@EditRouteActivity, "Error updating route", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        scope.cancel()
    }
}
