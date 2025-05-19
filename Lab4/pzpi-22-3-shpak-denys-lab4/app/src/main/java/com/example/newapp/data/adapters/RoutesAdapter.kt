package com.example.newapp.data.adapters

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.Route

class RoutesAdapter(
    private var routes: List<Route>,
    private val onItemClick: (Route) -> Unit
) : RecyclerView.Adapter<RoutesAdapter.RouteViewHolder>() {

    inner class RouteViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val tvName: TextView = itemView.findViewById(R.id.tvRouteName)
        val tvStatus: TextView = itemView.findViewById(R.id.tvRouteStatus)
        val tvDetails: TextView = itemView.findViewById(R.id.tvRouteDetails)
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): RouteViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_route, parent, false)
        return RouteViewHolder(view)
    }

    override fun onBindViewHolder(holder: RouteViewHolder, position: Int) {
        val route = routes[position]
        holder.tvName.text = route.name
        holder.tvStatus.text = route.status
        holder.tvDetails.text = route.details

        holder.itemView.setOnClickListener {
            onItemClick(route)
        }
    }

    override fun getItemCount(): Int = routes.size

    fun updateData(newRoutes: List<Route>) {
        routes = newRoutes
        notifyDataSetChanged()
    }
}
