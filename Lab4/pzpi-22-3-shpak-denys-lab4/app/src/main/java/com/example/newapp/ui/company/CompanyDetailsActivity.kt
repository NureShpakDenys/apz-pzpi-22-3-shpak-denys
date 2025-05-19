package com.example.newapp.ui.company

import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import androidx.viewpager2.widget.ViewPager2
import com.example.newapp.ui.base.BaseActivity
import com.example.newapp.R
import com.example.newapp.data.models.CompanyResponse
import com.example.newapp.data.models.CompanyUser
import com.example.newapp.data.adapters.CompanyTabsPagerAdapter
import com.google.android.material.tabs.TabLayout
import com.google.android.material.tabs.TabLayoutMediator
import kotlinx.coroutines.launch
import com.example.newapp.data.network.RetrofitClient
import android.os.Bundle
import android.util.Log

class CompanyDetailsActivity : BaseActivity() {
    private lateinit var viewPager: ViewPager2
    private lateinit var tabLayout: TabLayout
    private var companyId: Int = 0
    private lateinit var company: CompanyResponse
    private lateinit var companyUsers: List<CompanyUser>
    private var tabLayoutMediator: TabLayoutMediator? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_company_details)

        val extras = intent.extras
        if (extras != null && extras.containsKey("company_id")) {
            companyId = extras.getInt("company_id")
            Log.d("CompanyDetailsActivity", "Extracted company ID from extras: $companyId")
        } else {
            Log.e("CompanyDetailsActivity", "No company_id found in intent extras")
        }

        viewPager = findViewById(R.id.viewPager)
        tabLayout = findViewById(R.id.tabLayout)

        fetchCompanyData()
    }

    fun fetchCompanyData() {
        val apiService = RetrofitClient.getInstance(this)
        lifecycleScope.launch {
            try {
                val response = apiService.getCompany(companyId)
                val companyUsersResponse = apiService.getCompanyUsers(companyId)
                company = response
                companyUsers = companyUsersResponse
                setupViewPager()
            } catch (e: Exception) {
                Toast.makeText(this@CompanyDetailsActivity, "Error loading company: ${e.message}", Toast.LENGTH_SHORT).show()
            }
        }
    }

    private fun setupViewPager() {
        tabLayoutMediator?.detach()

        val adapter = CompanyTabsPagerAdapter(this, company, companyUsers)
        viewPager.adapter = adapter

        tabLayoutMediator = TabLayoutMediator(tabLayout, viewPager) { tab, position ->
            tab.text = when (position) {
                0 -> getString(R.string.routes_title)
                1 -> getString(R.string.deliveries_title)
                2 -> getString(R.string.users_title)
                else -> ""
            }
        }

        tabLayoutMediator?.attach()
    }

    override fun onDestroy() {
        tabLayoutMediator?.detach()
        super.onDestroy()
    }
}