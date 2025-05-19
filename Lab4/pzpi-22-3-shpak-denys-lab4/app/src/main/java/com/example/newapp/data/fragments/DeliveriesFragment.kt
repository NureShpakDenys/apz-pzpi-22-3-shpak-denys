package com.example.newapp.data.fragments

import android.app.Activity
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.ProgressBar
import android.widget.TextView
import android.widget.Toast
import androidx.activity.result.ActivityResultLauncher
import androidx.activity.result.contract.ActivityResultContracts
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.adapters.DeliveryAdapter
import com.example.newapp.data.models.CompanyResponse
import com.example.newapp.data.models.DeliveryRouteUpdateRequest
import com.example.newapp.data.models.CompanyDelivery
import com.example.newapp.data.models.Route
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.company.CompanyDetailsActivity
import com.example.newapp.ui.delivery.CreateDeliveryActivity
import com.example.newapp.ui.delivery.DeliveryDetailsActivity
import kotlinx.coroutines.launch

class DeliveriesFragment : Fragment() {
    private lateinit var recyclerView: RecyclerView
    private lateinit var progressBar: ProgressBar
    private lateinit var noDeliveriesText: TextView
    private lateinit var addDeliveryButton: Button
    private lateinit var deliveryAdapter: DeliveryAdapter

    private var deliveries: List<CompanyDelivery> = emptyList()
    private var routes: List<Route> = emptyList()
    private var companyId: Int = 0
    private var isCreator: Boolean = false

    private lateinit var createDeliveryLauncher: ActivityResultLauncher<Intent>

    companion object {
        private const val TAG = "DeliveriesFragment"

        fun newInstance(
            deliveries: List<CompanyDelivery>,
            routes: List<Route> = emptyList(),
            isCreator: Boolean = false,
            companyId: Int = 0
        ): DeliveriesFragment {
            val fragment = DeliveriesFragment()
            fragment.deliveries = deliveries
            fragment.routes = routes
            fragment.isCreator = isCreator
            fragment.companyId = companyId
            return fragment
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Log.d(TAG, "onCreate called. Company ID: $companyId, isCreator: $isCreator")

        createDeliveryLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            Log.d(TAG, "ActivityResult from CreateDeliveryActivity: resultCode=${result.resultCode}")
            if (result.resultCode == Activity.RESULT_OK) {
                val newDeliveryId = result.data?.getIntExtra(CreateDeliveryActivity.RESULT_EXTRA_DELIVERY_ID, -1) ?: -1
                if (newDeliveryId != -1) {
                    Log.i(TAG, "New delivery created with ID: $newDeliveryId. Refreshing company data.")
                    Toast.makeText(requireContext(), "Refreshing list ", Toast.LENGTH_SHORT).show()

                    (activity as? CompanyDetailsActivity)?.fetchCompanyData()
                } else {
                    Log.w(TAG, "New delivery ID not found in result, but result was OK. Refreshing company data anyway.")
                    (activity as? CompanyDetailsActivity)?.fetchCompanyData()
                }
            } else {
                Log.d(TAG, "Create delivery was cancelled or failed.")
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_delivery, container, false)

        recyclerView = view.findViewById(R.id.deliveryRecyclerView)
        progressBar = view.findViewById(R.id.deliveryProgressBar)
        noDeliveriesText = view.findViewById(R.id.tvNoDeliveries)
        addDeliveryButton = view.findViewById(R.id.btnAddDelivery)

        recyclerView.layoutManager = LinearLayoutManager(requireContext())

        if (isCreator) {
            addDeliveryButton.visibility = View.VISIBLE
            addDeliveryButton.setOnClickListener {
                if (companyId == 0) {
                    Toast.makeText(requireContext(), getString(R.string.error_company_id_missing_delivery), Toast.LENGTH_LONG).show()
                    Log.e(TAG, "Cannot start CreateDeliveryActivity: companyId is not valid ($companyId)")
                    return@setOnClickListener
                }
                Log.d(TAG, "Add delivery button clicked for Company ID: $companyId")
                val intent = Intent(requireContext(), CreateDeliveryActivity::class.java)
                intent.putExtra(CreateDeliveryActivity.EXTRA_COMPANY_ID, companyId)
                createDeliveryLauncher.launch(intent)
            }
        } else {
            addDeliveryButton.visibility = View.GONE
        }

        deliveryAdapter = DeliveryAdapter(
            deliveries = deliveries,
            routes = routes,
            isCreator = isCreator,
            onItemClick = { delivery ->
                val intent = Intent(requireContext(), DeliveryDetailsActivity::class.java)
                intent.putExtra("delivery_id", delivery.id)
                startActivity(intent)
            },
            onRouteChange = { deliveryId, routeId ->
                updateDeliveryRoute(deliveryId, routeId)
            }
        )

        recyclerView.adapter = deliveryAdapter

        updateUI()

        return view
    }

    private fun updateUI() {
        if (deliveries.isEmpty()) {
            recyclerView.visibility = View.GONE
            noDeliveriesText.visibility = View.VISIBLE
        } else {
            recyclerView.visibility = View.VISIBLE
            noDeliveriesText.visibility = View.GONE
        }
    }

    private fun updateDeliveryRoute(deliveryId: Int, routeId: Int) {
        progressBar.visibility = View.VISIBLE

        lifecycleScope.launch {
            try {
                val sessionManager = SessionManager(requireContext())
                val token = sessionManager.getAccessToken()

                val response = RetrofitClient.getInstance(requireContext())
                    .updateDeliveryRoute(deliveryId, DeliveryRouteUpdateRequest(routeId), "Bearer $token")

                if (response.isSuccessful) {
                    deliveries = deliveries.map {
                        if (it.id == deliveryId) it.copy(routeId = routeId) else it
                    }
                    deliveryAdapter.updateData(deliveries)

                    Toast.makeText(requireContext(), "Route updated successfully", Toast.LENGTH_SHORT).show()
                } else {
                    Toast.makeText(requireContext(), "Failed to update route", Toast.LENGTH_SHORT).show()
                }
            } catch (e: Exception) {
                Toast.makeText(requireContext(), "Error: ${e.message}", Toast.LENGTH_SHORT).show()
            } finally {
                progressBar.visibility = View.GONE
            }
        }
    }
}