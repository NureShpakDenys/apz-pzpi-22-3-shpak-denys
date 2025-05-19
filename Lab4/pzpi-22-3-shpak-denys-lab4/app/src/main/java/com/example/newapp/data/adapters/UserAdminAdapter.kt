package com.example.newapp.data.adapters

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.ArrayAdapter
import android.widget.TextView
import androidx.appcompat.widget.AppCompatSpinner
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.Role
import com.example.newapp.data.models.User

class UserAdminAdapter(
    private val context: Context,
    private var users: List<User>,
    private val availableRoles: List<Role>,
    private val currentAdminUser: User?,
    private val onRoleChangeAttempt: (userId: Int, newRoleId: Int, oldRoleName: String) -> Unit
) : RecyclerView.Adapter<UserAdminAdapter.UserViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): UserViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_user_admin, parent, false)
        return UserViewHolder(view)
    }

    override fun onBindViewHolder(holder: UserViewHolder, position: Int) {
        val user = users[position]
        holder.bind(user)
    }

    override fun getItemCount(): Int = users.size

    fun updateUserRoleLocally(userId: Int, newRoleName: String) {
        val userIndex = users.indexOfFirst { it.id == userId }
        if (userIndex != -1) {
            users[userIndex].role = newRoleName
            notifyItemChanged(userIndex)
        }
    }

    fun updateUsers(newUsers: List<User>) {
        this.users = newUsers
        notifyDataSetChanged()
    }

    inner class UserViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        private val tvUserId: TextView = itemView.findViewById(R.id.tvUserId)
        private val tvUserName: TextView = itemView.findViewById(R.id.tvUserName)
        private val tvUserRoleStatic: TextView = itemView.findViewById(R.id.tvUserRoleStatic)
        private val spinnerUserRole: AppCompatSpinner = itemView.findViewById(R.id.spinnerUserRole)

        fun bind(user: User) {
            tvUserId.text = user.id.toString()
            tvUserName.text = user.name

            val isAdminSuperUser = currentAdminUser?.role == "admin"
            val isThisUsersOwnEntry = currentAdminUser?.id == user.id

            if (isAdminSuperUser && !isThisUsersOwnEntry) {
                tvUserRoleStatic.visibility = View.GONE
                spinnerUserRole.visibility = View.VISIBLE

                val roleNames = availableRoles.map { it.name }
                val adapter = ArrayAdapter(context, android.R.layout.simple_spinner_item, roleNames)
                adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
                spinnerUserRole.adapter = adapter

                val currentRoleIndex = availableRoles.indexOfFirst { it.name.equals(user.role, ignoreCase = true) }
                if (currentRoleIndex != -1) {
                    spinnerUserRole.setSelection(currentRoleIndex, false)
                }

                spinnerUserRole.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
                    override fun onItemSelected(parent: AdapterView<*>?, view: View?, position: Int, id: Long) {
                        val selectedRole = availableRoles[position]
                        if (!selectedRole.name.equals(user.role, ignoreCase = true)) {
                            onRoleChangeAttempt(user.id, selectedRole.id, user.role)
                        }
                    }

                    override fun onNothingSelected(parent: AdapterView<*>?) {
                    }
                }
            } else {
                tvUserRoleStatic.visibility = View.VISIBLE
                spinnerUserRole.visibility = View.GONE
                tvUserRoleStatic.text = user.role
            }
        }
    }
}