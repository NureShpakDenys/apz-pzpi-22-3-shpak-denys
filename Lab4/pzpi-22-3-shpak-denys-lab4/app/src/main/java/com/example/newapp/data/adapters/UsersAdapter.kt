package com.example.newapp.data.adapters

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.*
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.CompanyUser

class UsersAdapter(
    private val users: List<CompanyUser>,
    private val creatorId: Int,
    private val onRoleChange: (CompanyUser) -> Unit,
    private val onRemove: (CompanyUser) -> Unit
) : RecyclerView.Adapter<UsersAdapter.UserViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): UserViewHolder {
        val view =
            LayoutInflater.from(parent.context).inflate(R.layout.item_user, parent, false)
        return UserViewHolder(view)
    }

    override fun onBindViewHolder(holder: UserViewHolder, position: Int) {
        holder.bind(users[position])
    }

    override fun getItemCount() = users.size

    inner class UserViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        private val username: TextView = itemView.findViewById(R.id.userName)
        private val role: TextView = itemView.findViewById(R.id.userRole)
        private val editRoleBtn: ImageButton = itemView.findViewById(R.id.editRoleButton)
        private val removeBtn: ImageButton = itemView.findViewById(R.id.removeUserButton)

        fun bind(companyUser: CompanyUser) {
            username.text = companyUser.user.name
            role.text = companyUser.Role
            editRoleBtn.visibility = if (companyUser.UserID != creatorId) View.VISIBLE else View.GONE
            removeBtn.visibility = if (companyUser.UserID != creatorId) View.VISIBLE else View.GONE

            editRoleBtn.setOnClickListener { onRoleChange(companyUser) }
            removeBtn.setOnClickListener { onRemove(companyUser) }
        }
    }
}
