package com.example.newapp.ui.adapter

import android.annotation.SuppressLint
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.data.models.Waypoint
import com.example.newapp.databinding.ItemWaypointBinding
import java.util.Locale

class WaypointAdapter(
    private var waypoints: List<Waypoint>,
    private val onDeleteClick: (Waypoint) -> Unit,
    private val onEditClick: (Waypoint) -> Unit
) : RecyclerView.Adapter<WaypointAdapter.WaypointViewHolder>() {

    inner class ViewHolder(val binding: ItemWaypointBinding) : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): WaypointViewHolder {
        val binding = ItemWaypointBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return WaypointViewHolder(binding)
    }

    override fun onBindViewHolder(holder: WaypointViewHolder, position: Int) {
        val waypoint = waypoints[position]
        holder.bind(waypoint)
    }

    override fun getItemCount(): Int = waypoints.size

    fun updateData(newWaypoints: List<Waypoint>) {
        this.waypoints = newWaypoints
        notifyDataSetChanged()
    }

    inner class WaypointViewHolder(private val binding: ItemWaypointBinding) : RecyclerView.ViewHolder(binding.root) {
        @SuppressLint("SetTextI18n")
        fun bind(waypoint: Waypoint) {
            binding.tvWpName.text = waypoint.Name
            binding.tvWpStatus.text = "Status: ${waypoint.Status}"

            val latFormatted = String.format(Locale.US, "%.4f", waypoint.Latitude)
            val lonFormatted = String.format(Locale.US, "%.4f", waypoint.Longitude)
            binding.tvWpCoords.text = "Coords: $latFormatted, $lonFormatted"

            binding.tvWpDeviceSerial.text = "Device: ${waypoint.DeviceSerial}"
            binding.tvWpFrequency.text = "Frequency: ${waypoint.SendDataFrequency}s"
            binding.tvWpWeatherAlerts.text = "Weather Alerts: ${if (waypoint.GetWeatherAlerts) "Enabled" else "Disabled"}"

            if (waypoint.Details.isNotBlank()) {
                binding.tvWpDetails.text = "Details: ${waypoint.Details}"
                binding.tvWpDetails.visibility = View.VISIBLE
            } else {
                binding.tvWpDetails.visibility = View.GONE
            }

            binding.btnDeleteWp.setOnClickListener {
                onDeleteClick(waypoint)
            }

            binding.btnEditWp.setOnClickListener {
                onEditClick(waypoint)
            }
        }
    }
}
