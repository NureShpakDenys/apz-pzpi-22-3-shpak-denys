package com.example.newapp.data.adapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.core.content.ContextCompat
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.SystemLog
import java.text.SimpleDateFormat
import java.util.Locale
import java.util.TimeZone

class SystemLogAdapter(
    private val context: Context,
    private var logs: List<SystemLog>
) : RecyclerView.Adapter<SystemLogAdapter.LogViewHolder>() {

    private val inputDateFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.US).apply {
        timeZone = TimeZone.getTimeZone("UTC")
    }
    private val outputDateFormat = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault())


    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): LogViewHolder {
        val view = LayoutInflater.from(context).inflate(R.layout.item_system_log, parent, false)
        return LogViewHolder(view)
    }

    override fun onBindViewHolder(holder: LogViewHolder, position: Int) {
        val log = logs[position]
        holder.bind(log)
    }

    override fun getItemCount(): Int = logs.size

    fun updateLogs(newLogs: List<SystemLog>) {
        this.logs = newLogs
        notifyDataSetChanged()
    }

    inner class LogViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        private val tvLogId: TextView = itemView.findViewById(R.id.tvLogId)
        private val tvLogCreatedAt: TextView = itemView.findViewById(R.id.tvLogCreatedAt)
        private val tvLogUserId: TextView = itemView.findViewById(R.id.tvLogUserId)
        private val tvLogActionType: TextView = itemView.findViewById(R.id.tvLogActionType)
        private val tvLogDescription: TextView = itemView.findViewById(R.id.tvLogDescription)
        private val tvLogSuccess: TextView = itemView.findViewById(R.id.tvLogSuccess)

        fun bind(log: SystemLog) {
            tvLogId.text = context.getString(R.string.log_id_prefix, log.id?.toString() ?: "N/A")
            tvLogCreatedAt.text = formatDisplayDate(log.createdAt)
            tvLogUserId.text = context.getString(R.string.log_user_id_prefix, log.userId?.toString() ?: "N/A")
            tvLogActionType.text = log.actionType ?: "N/A"
            tvLogDescription.text = log.description ?: "N/A"
            if (log.success == true) {
                tvLogSuccess.text = "✔️"
                tvLogSuccess.setTextColor(ContextCompat.getColor(context, android.R.color.holo_green_dark))
            } else {
                tvLogSuccess.text = "❌"
                tvLogSuccess.setTextColor(ContextCompat.getColor(context, android.R.color.holo_red_dark))
            }
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
}