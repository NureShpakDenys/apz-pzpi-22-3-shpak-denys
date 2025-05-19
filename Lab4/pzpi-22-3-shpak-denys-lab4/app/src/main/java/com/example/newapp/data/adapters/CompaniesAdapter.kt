package com.example.newapp.data.adapters

import android.annotation.SuppressLint
import android.view.LayoutInflater
import android.view.View

import android.view.ViewGroup
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.Company

class CompanyAdapter(
    private val companies: List<Company>,
    private val onItemClickListener: (Company) -> Unit
) : RecyclerView.Adapter<CompanyAdapter.CompanyViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): CompanyViewHolder {
        val view = LayoutInflater.from(parent.context).inflate(R.layout.item_company, parent, false)
        return CompanyViewHolder(view)
    }

    override fun onBindViewHolder(
        holder:
        CompanyViewHolder, position: Int
    ) {
        val company = companies[position]
        holder.bind(company)
    }

    override fun getItemCount(): Int {
        return companies.size
    }

    inner class CompanyViewHolder(itemView: View) :
        RecyclerView.ViewHolder(itemView) {
        private val nameTextView: TextView = itemView.findViewById(R.id.companyName)
        private val addressTextView: TextView = itemView.findViewById(R.id.companyAddress)

        init {
            itemView.setOnClickListener {
                val company =
                    companies[bindingAdapterPosition]
                onItemClickListener(company)
            }
        }

        @SuppressLint("SetTextI18n")
        fun bind(company: Company) {
            nameTextView.text = company.name
            val truncatedContent = if (company.address.length > 100) {
                company.address.substring(0, 100) + "..."
            } else {
                company.address
            }

            addressTextView.text = truncatedContent
        }
    }
}