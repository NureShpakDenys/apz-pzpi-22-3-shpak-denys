package com.example.newapp.ui.company

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.widget.Toast
import com.example.newapp.ui.base.BaseActivity
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.data.adapters.CompanyAdapter
import com.example.newapp.data.models.Company
import com.example.newapp.data.network.RetrofitClient
import retrofit2.Call
import retrofit2.Callback
import com.example.newapp.R
import retrofit2.Response

class CompaniesActivity : BaseActivity() {

    private lateinit var recyclerView: RecyclerView

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_companies)

        recyclerView = findViewById(R.id.companiesRecyclerView)

        val apiService = RetrofitClient.getInstance(this)

        apiService.getCompanies().enqueue(object : Callback<List<Company>> {
            override fun onResponse(call: Call<List<Company>>, response: Response<List<Company>>) {
                if (response.isSuccessful) {
                    val companies = response.body() ?: emptyList()
                    for (company in companies) {
                        Log.d("CompaniesActivity", "Company from API: id=${company.id}, name=${company.name}")
                    }
                    val adapter = CompanyAdapter(companies) { company ->
                        Log.d("CompaniesActivity", "Selected company ID: ${company.id}")
                        navigateToCompanyDetails(company)
                    }
                    recyclerView.adapter = adapter
                } else {
                    Toast.makeText(this@CompaniesActivity, "Ошибка загрузки данных", Toast.LENGTH_SHORT).show()
                }
            }

            override fun onFailure(call: Call<List<Company>>, t: Throwable) {
                Log.e("CompaniesActivity", "Ошибка при запросе компаний", t)
                Toast.makeText(this@CompaniesActivity, "Ошибка сети: ${t.message}", Toast.LENGTH_SHORT).show()
            }
        })
    }


    private fun navigateToCompanyDetails(company: Company) {
        val intent = Intent(this, CompanyDetailsActivity::class.java).apply {
            putExtra("company_id", company.id)          
        }
        startActivity(intent)
    }
}