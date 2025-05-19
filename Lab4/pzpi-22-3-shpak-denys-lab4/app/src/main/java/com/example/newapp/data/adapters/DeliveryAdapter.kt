package com.example.newapp.data.adapters

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.ArrayAdapter
import android.widget.Spinner
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.CompanyDelivery
import com.example.newapp.data.models.Route
import java.text.SimpleDateFormat
import java.util.Locale

class DeliveryAdapter(
    private var deliveries: List<CompanyDelivery>,
    private val routes: List<Route> = emptyList(),
    private val isCreator: Boolean = false,
    private val onItemClick: (CompanyDelivery) -> Unit,
    private val onRouteChange: ((Int, Int) -> Unit)? = null
) : RecyclerView.Adapter<DeliveryAdapter.DeliveryViewHolder>() {

    inner class DeliveryViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val tvId: TextView = itemView.findViewById(R.id.tvDeliveryId)
        val tvStatus: TextView = itemView.findViewById(R.id.tvDeliveryStatus)
        val tvDate: TextView = itemView.findViewById(R.id.tvDeliveryDate)
        val tvDuration: TextView = itemView.findViewById(R.id.tvDeliveryDuration)
        val routeSpinner: Spinner = itemView.findViewById(R.id.spinnerRoute)
        val tvRouteId: TextView = itemView.findViewById(R.id.tvRouteId)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): DeliveryViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_delivery, parent, false)
        return DeliveryViewHolder(view)
    }

    override fun onBindViewHolder(holder: DeliveryViewHolder, position: Int) {
        val delivery = deliveries[position]
        holder.tvId.text = delivery.id.toString()
        holder.tvStatus.text = delivery.status

        // Format date based on locale
        val dateFormat = SimpleDateFormat("dd.MM.yyyy", Locale.getDefault())
        holder.tvDate.text = try {
            dateFormat.format(delivery.date)
        } catch (e: Exception) {
            delivery.date.toString()
        }

        holder.tvDuration.text = delivery.duration.toString()

        // Setup route selection
        if (isCreator && routes.isNotEmpty() && onRouteChange != null) {
            holder.routeSpinner.visibility = View.VISIBLE
            holder.tvRouteId.visibility = View.GONE

            val routeNames = routes.map { it.name }.toTypedArray()
            val adapter = ArrayAdapter(holder.itemView.context, android.R.layout.simple_spinner_item, routeNames)
            adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
            holder.routeSpinner.adapter = adapter

            // Set current selection
            val currentRouteIndex = routes.indexOfFirst { it.id == delivery.routeId }
            if (currentRouteIndex >= 0) {
                holder.routeSpinner.setSelection(currentRouteIndex)
            }

            holder.routeSpinner.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
                override fun onItemSelected(parent: AdapterView<*>?, view: View?, position: Int, id: Long) {
                    val selectedRouteId = routes[position].id
                    if (selectedRouteId != delivery.routeId) {
                        onRouteChange.invoke(delivery.id, selectedRouteId)
                    }
                }

                override fun onNothingSelected(parent: AdapterView<*>?) {}
            }
        } else {
            holder.routeSpinner.visibility = View.GONE
            holder.tvRouteId.visibility = View.VISIBLE
            holder.tvRouteId.text = delivery.routeId.toString()
        }

        // Set click listener for the delivery item
        holder.itemView.setOnClickListener {
            onItemClick(delivery)
        }
    }

    override fun getItemCount(): Int = deliveries.size

    fun updateData(newDeliveries: List<CompanyDelivery>) {
        deliveries = newDeliveries
        notifyDataSetChanged()
    }
}