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
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.data.models.Route
import com.example.newapp.data.adapters.RoutesAdapter
import com.example.newapp.R
import com.example.newapp.ui.company.CompanyDetailsActivity
import com.example.newapp.ui.route.CreateRouteActivity
import com.example.newapp.ui.route.RouteDetailsActivity

class RoutesFragment : Fragment() {
    private lateinit var recyclerView: RecyclerView
    private lateinit var progressBar: ProgressBar
    private lateinit var noRoutesText: TextView
    private lateinit var addRouteButton: Button
    private lateinit var routesAdapter: RoutesAdapter

    private var routes: List<Route> = emptyList()
    private var companyId: Int = 0
    private var isCreator: Boolean = false

    private lateinit var createRouteLauncher: ActivityResultLauncher<Intent>
    private lateinit var routeDetailsLauncher: ActivityResultLauncher<Intent>

    companion object {
        private const val TAG = "RoutesFragment"

        fun newInstance(
            routes: List<Route>,
            isCreator: Boolean = false,
            companyId: Int = 0
        ): RoutesFragment {
            val fragment = RoutesFragment()
            fragment.routes = routes.toMutableList()
            fragment.isCreator = isCreator
            fragment.companyId = companyId
            return fragment
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Log.d(TAG, "onCreate called for Company ID: $companyId")

        createRouteLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            Log.d(TAG, "ActivityResult from CreateRouteActivity: resultCode=${result.resultCode}")
            if (result.resultCode == Activity.RESULT_OK) {
                val newRouteId = result.data?.getIntExtra(CreateRouteActivity.RESULT_EXTRA_ROUTE_ID, -1) ?: -1
                if (newRouteId != -1) {
                    Log.i(TAG, "New route created with ID: $newRouteId. Refreshing data.")
                    Toast.makeText(requireContext(), getString(R.string.route_created_successfully) + " Refreshing list...", Toast.LENGTH_SHORT).show()

                    (activity as? CompanyDetailsActivity)?.fetchCompanyData()
                } else {
                    Log.w(TAG, "New route ID not found in result. Refreshing data anyway.")
                    (activity as? CompanyDetailsActivity)?.fetchCompanyData()
                }
            } else {
                Log.d(TAG, "Create route was cancelled or failed.")
            }
        }

        routeDetailsLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            Log.d(TAG, "ActivityResult from RouteDetailsActivity: resultCode=${result.resultCode}")
            if (result.resultCode == Activity.RESULT_OK) {
                Log.i(TAG, "Route details changed or route deleted. Refreshing company data.")
                Toast.makeText(requireContext(), getString(R.string.route_list_updated_refreshing), Toast.LENGTH_SHORT).show()
                (activity as? CompanyDetailsActivity)?.fetchCompanyData()
            } else {
                Log.d(TAG, "RouteDetailsActivity returned with no changes requiring refresh or was cancelled.")
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_routes, container, false)

        recyclerView = view.findViewById(R.id.routesRecyclerView)
        progressBar = view.findViewById(R.id.routeProgressBar)
        noRoutesText = view.findViewById(R.id.tvNoRoutes)
        addRouteButton = view.findViewById(R.id.btnAddRoute)

        recyclerView.layoutManager = LinearLayoutManager(requireContext())

        if (isCreator) {
            addRouteButton.visibility = View.VISIBLE
            addRouteButton.setOnClickListener {
                if (companyId == 0) {
                    Toast.makeText(requireContext(), "Company ID is not set. Cannot create route.", Toast.LENGTH_LONG).show()
                    Log.e(TAG, "Cannot start CreateRouteActivity: companyId is not valid ($companyId)")
                    return@setOnClickListener
                }
                Log.d(TAG, "Add route button clicked for Company ID: $companyId")
                val intent = Intent(requireContext(), CreateRouteActivity::class.java)
                intent.putExtra(CreateRouteActivity.EXTRA_COMPANY_ID, companyId)
                createRouteLauncher.launch(intent)
            }
        } else {
            addRouteButton.visibility = View.GONE
        }

        routesAdapter = RoutesAdapter(
            routes = routes,
            onItemClick = { route ->
                val intent = Intent(requireContext(), RouteDetailsActivity::class.java)
                intent.putExtra("route_id", route.id)
                routeDetailsLauncher.launch(intent)
            },
        )
        recyclerView.adapter = routesAdapter

        updateUI()

        return view
    }

    private fun updateUI() {
        if (routes.isEmpty()) {
            recyclerView.visibility = View.GONE
            noRoutesText.visibility = View.VISIBLE
        } else {
            recyclerView.visibility = View.VISIBLE
            noRoutesText.visibility = View.GONE
        }
    }

}
