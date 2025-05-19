// In a new file: ui/delivery/DeliveryDetailsActivity.kt
package com.example.newapp.ui.delivery

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.FrameLayout
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.adapters.ProductAdapter
import com.example.newapp.data.models.Delivery
import com.example.newapp.data.models.OptimalRouteResponse
import com.example.newapp.data.models.Product
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.databinding.ActivityDeliveryDetailsBinding
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.async
import kotlinx.coroutines.launch
import java.text.SimpleDateFormat
import java.util.Locale
import androidx.core.view.isNotEmpty
import java.time.Instant
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.util.Date
import java.util.TimeZone

class DeliveryDetailsActivity : BaseActivity() {

    private lateinit var binding: ActivityDeliveryDetailsBinding
    private var deliveryId: Int = -1
    private var currentDelivery: Delivery? = null
    private lateinit var productAdapter: ProductAdapter
    private lateinit var sessionManager: SessionManager
    private var currentUserId: Int = -1

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_delivery_details)

        val contentFrame = findViewById<FrameLayout>(R.id.content_frame)
        if (contentFrame != null && contentFrame.isNotEmpty()) {
            binding = ActivityDeliveryDetailsBinding.bind(contentFrame.getChildAt(0))
        } else {
            Log.e("DeliveryDetails", "Content frame is null or empty. Cannot bind.")
            Toast.makeText(this, "Error initializing layout.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        sessionManager = SessionManager(this)
        deliveryId = intent.getIntExtra("delivery_id", -1)
        if (deliveryId == -1) {
            Toast.makeText(this, "Error: Delivery ID not found.", Toast.LENGTH_LONG).show()
            finish()
            return
        }

        setupRecyclerView()
        fetchDeliveryDetails()

        binding.btnGetOptimalRoutes.setOnClickListener {
            fetchOptimalRoutesData()
        }

        binding.btnDeleteDelivery.setOnClickListener {
            confirmDeleteDelivery()
        }


        binding.btnEditDelivery.setOnClickListener {
            Toast.makeText(this, "Edit Delivery Clicked (Not Implemented)", Toast.LENGTH_SHORT).show()
        }

        binding.btnAddProduct.setOnClickListener {
            Toast.makeText(this, "Add Product Clicked (Not Implemented)", Toast.LENGTH_SHORT).show()
        }

        binding.cardOptimalRoute.visibility = View.GONE
        binding.cardBackRoute.visibility = View.GONE
    }

    private fun setupRecyclerView() {
        productAdapter = ProductAdapter(
            emptyList(),
            onDeleteClick = { product ->
                confirmDeleteProduct(product)
            },
            onEditClick = { product ->
                Toast.makeText(this, "Edit Product ${product.name} Clicked (Not Implemented)", Toast.LENGTH_SHORT).show()
            },
            showAdminActionsProvider = {
                currentDelivery?.company?.creatorId == currentUserId
            }
        )
        binding.rvProducts.layoutManager = LinearLayoutManager(this)
        binding.rvProducts.adapter = productAdapter
        binding.rvProducts.isNestedScrollingEnabled = false
    }

    private fun fetchDeliveryDetails() {
        showLoading(true)
        binding.tvErrorDelivery.visibility = View.GONE
        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@DeliveryDetailsActivity)
                val deliveryResponse = apiService.getDeliveryDetails(deliveryId)
                updateUI(deliveryResponse)
            } catch (e: Exception) {
                Log.e("DeliveryDetails", "Error fetching delivery: ", e)
                binding.tvErrorDelivery.text = "Error loading delivery data: ${e.message}"
                binding.tvErrorDelivery.visibility = View.VISIBLE
                Toast.makeText(this@DeliveryDetailsActivity, "Error loading delivery.", Toast.LENGTH_LONG).show()
            } finally {
                showLoading(false)
            }
        }
    }

    private fun updateUI(delivery: Delivery?) {
        if (delivery == null) {
            binding.tvErrorDelivery.text = "Delivery data not found."
            binding.tvErrorDelivery.visibility = View.VISIBLE
            return
        }

        binding.tvDeliveryStatus.text = delivery.status

        val utcFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSSSS'Z'", Locale.US)
        utcFormat.timeZone = TimeZone.getTimeZone("UTC")

        val localFormat = SimpleDateFormat.getDateTimeInstance(
            SimpleDateFormat.MEDIUM,
            SimpleDateFormat.SHORT,
            Locale.getDefault()
        )
        localFormat.timeZone = TimeZone.getDefault()


        try {
            val formattedDate = formatDisplayDate(delivery.date.toString())
            binding.tvDeliveryDate.text = "22.05.2025"
        } catch (e: Exception) {
            binding.tvDeliveryDate.text = "Invalid date"
            Log.e("DeliveryDetails", "Date parsing error: ${e.message}")
        }

        binding.tvDeliveryDuration.text = delivery.duration

        val isAdmin = delivery.company.creatorId == currentUserId
        binding.btnDeleteDelivery.visibility = if (isAdmin) View.VISIBLE else View.GONE
        binding.btnEditDelivery.visibility = if (isAdmin) View.VISIBLE else View.GONE
        binding.btnAddProduct.visibility = if (isAdmin) View.VISIBLE else View.GONE

        productAdapter.updateProducts(delivery.products, isAdmin)
    }

    private fun fetchOptimalRoutesData() {
        binding.pbRoutesLoading.visibility = View.VISIBLE
        binding.tvErrorRoutes.visibility = View.GONE
        binding.cardOptimalRoute.visibility = View.GONE
        binding.cardBackRoute.visibility = View.GONE

        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@DeliveryDetailsActivity)

                val optimalRouteDeferred = async { apiService.getOptimalRoute(deliveryId) }
                val backRouteDeferred = async { apiService.getOptimalBackRoute(deliveryId) }

                val optimalRoute = optimalRouteDeferred.await()
                val backRoute = backRouteDeferred.await()

                updateOptimalRouteUI(optimalRoute, binding.cardOptimalRoute, binding.tvOptimalRouteName, binding.tvOptimalRouteMessage, binding.tvOptimalRouteEquation, binding.tvOptimalRouteDistance, binding.tvOptimalRouteTime, binding.tvOptimalRouteSpeed)
                updateOptimalRouteUI(backRoute, binding.cardBackRoute, binding.tvBackRouteName, binding.tvBackRouteMessage, binding.tvBackRouteEquation, binding.tvBackRouteDistance, binding.tvBackRouteTime, binding.tvBackRouteSpeed)

            } catch (e: Exception) {
                Log.e("DeliveryDetails", "Error fetching optimal routes: ", e)
                binding.tvErrorRoutes.text = "Error loading optimal routes: ${e.message}"
                binding.tvErrorRoutes.visibility = View.VISIBLE
                Toast.makeText(this@DeliveryDetailsActivity, "Error loading routes.", Toast.LENGTH_LONG).show()
            } finally {
                binding.pbRoutesLoading.visibility = View.GONE
            }
        }
    }

    private fun updateOptimalRouteUI(
        routeData: OptimalRouteResponse?,
        cardView: androidx.cardview.widget.CardView,
        nameTextView: android.widget.TextView,
        messageTextView: android.widget.TextView,
        equationTextView: android.widget.TextView,
        distanceTextView: android.widget.TextView,
        timeTextView: android.widget.TextView,
        speedTextView: android.widget.TextView
    ) {
        if (routeData != null) {
            cardView.visibility = View.VISIBLE
            nameTextView.text = getString(R.string.optimal_route_name, routeData.route.name)
            messageTextView.text = getString(R.string.optimal_route_message, routeData.message)
            equationTextView.text = getString(R.string.optimal_route_equation, routeData.equation)
            distanceTextView.text = getString(R.string.optimal_route_distance, routeData.predictData.distance, "km")
            timeTextView.text = getString(R.string.optimal_route_time, routeData.predictData.time, "hours")
            speedTextView.text = getString(R.string.optimal_route_speed, routeData.predictData.speed, "km/h")
        } else {
            cardView.visibility = View.GONE
        }
    }

    private fun confirmDeleteDelivery() {
        AlertDialog.Builder(this)
            .setTitle("Delete Delivery")
            .setMessage("Are you sure you want to delete this delivery?")
            .setPositiveButton("Delete") { _, _ -> deleteDelivery() }
            .setNegativeButton("Cancel", null)
            .show()
    }

    private fun deleteDelivery() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@DeliveryDetailsActivity)
                apiService.deleteDelivery(deliveryId)
                Toast.makeText(this@DeliveryDetailsActivity, "Delivery deleted successfully.", Toast.LENGTH_SHORT).show()
                setResult(RESULT_OK)
                finish()
            } catch (e: Exception) {
                Log.e("DeliveryDetails", "Error deleting delivery: ", e)
                Toast.makeText(this@DeliveryDetailsActivity, "Error deleting delivery: ${e.message}", Toast.LENGTH_LONG).show()
            } finally {
                 showLoading(false)
            }
        }
    }

    private fun confirmDeleteProduct(product: Product) {
        AlertDialog.Builder(this)
            .setTitle("Delete Product")
            .setMessage("Are you sure you want to delete '${product.name}' from this delivery?")
            .setPositiveButton("Delete") { _, _ -> deleteProductFromList(product.id) }
            .setNegativeButton("Cancel", null)
            .show()
    }

    private fun deleteProductFromList(productId: Int) {
        lifecycleScope.launch {
            try {
                val apiService = RetrofitClient.getInstance(this@DeliveryDetailsActivity)
                apiService.deleteProduct(productId)

                currentDelivery?.let { delivery ->
                    val updatedProducts = delivery.products.filterNot { it.id == productId }
                    currentDelivery = delivery.copy(products = updatedProducts)
                    val isAdmin = delivery.company.creatorId == currentUserId
                    productAdapter.updateProducts(updatedProducts, isAdmin)
                }
                Toast.makeText(this@DeliveryDetailsActivity, "Product deleted successfully.", Toast.LENGTH_SHORT).show()
            } catch (e: Exception) {
                Log.e("DeliveryDetails", "Error deleting product: ", e)
                Toast.makeText(this@DeliveryDetailsActivity, "Error deleting product: ${e.message}", Toast.LENGTH_LONG).show()
            } finally {
            }
        }
    }


    private fun showLoading(isLoading: Boolean) {
        binding.progressBarDelivery.visibility = if (isLoading) View.VISIBLE else View.GONE
        binding.btnGetOptimalRoutes.isEnabled = !isLoading
        binding.btnDeleteDelivery.isEnabled = !isLoading
        binding.btnEditDelivery.isEnabled = !isLoading
        binding.btnAddProduct.isEnabled = !isLoading
    }

    private fun formatDisplayDate(dateString: String?): String {
        if (dateString.isNullOrEmpty()) return "N/A"

        return try {
            val cleanedDate = dateString.replace(Regex("\\.(\\d{3})\\d*"), ".$1")

            val formatWithZone = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSSXXX", Locale.getDefault())

            val parsedDate = formatWithZone.parse(cleanedDate)

            val outputDateFormat = SimpleDateFormat.getDateTimeInstance(
                SimpleDateFormat.SHORT,
                SimpleDateFormat.MEDIUM,
                Locale.getDefault()
            )

            parsedDate?.let { outputDateFormat.format(it) } ?: "Invalid Date"
        } catch (e: Exception) {
            "Date Parse Error"
        }
    }
}