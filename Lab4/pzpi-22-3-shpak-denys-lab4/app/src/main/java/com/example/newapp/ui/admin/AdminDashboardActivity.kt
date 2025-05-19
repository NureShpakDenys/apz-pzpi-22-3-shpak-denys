package com.example.newapp.ui.admin

import android.content.Context
import android.content.Intent
import android.content.SharedPreferences
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.ProgressBar
import android.widget.TextView
import android.widget.Toast
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.adapters.UserAdminAdapter
import com.example.newapp.data.models.ChangeRoleRequest
import com.example.newapp.data.models.Role
import com.example.newapp.data.models.User
import com.example.newapp.data.network.ApiService
import com.example.newapp.data.SessionManager
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.login.LoginActivity
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.launch

class AdminDashboardActivity : BaseActivity() {

    private lateinit var rvUsers: RecyclerView
    private lateinit var progressBarAdmin: ProgressBar
    private lateinit var tvAdminError: TextView
    private lateinit var userAdminAdapter: UserAdminAdapter
    private lateinit var apiService: ApiService
    private lateinit var sharedPreferences: SharedPreferences

    private var currentAdminUser: User? = null
    private var allUsersList: MutableList<User> = mutableListOf()
    private val availableRoles: List<Role> = getHardcodedRoles()

    companion object {
        private const val TAG = "AdminDashboardActivity"
        private const val PREFS_NAME = "MyAppPrefs"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_admin_dashboard)

        supportActionBar?.setDisplayHomeAsUpEnabled(true)
        supportActionBar?.setDisplayShowHomeEnabled(true)

        apiService = RetrofitClient.getInstance(this)
        sharedPreferences = getSharedPreferences(PREFS_NAME, MODE_PRIVATE)

        rvUsers = findViewById(R.id.rvUsers)
        progressBarAdmin = findViewById(R.id.progressBarAdmin)
        tvAdminError = findViewById(R.id.tvAdminError)

        setupRecyclerView()

        val sessionManager = SessionManager(this)

        val currentUserId = sessionManager.getUser()?.id

        if (currentUserId == null) {
            navigateToLogin()
            return
        }
        fetchCurrentAdminDetailsAndThenUsers(currentUserId)
    }

    private fun getHardcodedRoles(): List<Role> {
        return listOf(
            Role(1, "admin"),
            Role(2, "user"),
            Role(3, "manager"),
            Role(4, "db_admin"),
            Role(5, "system_admin")
        )
    }

    private fun setupRecyclerView() {
        userAdminAdapter = UserAdminAdapter(
            this,
            allUsersList,
            availableRoles,
            currentAdminUser
        ) { userId, newRoleId, oldRoleName ->
            if (currentAdminUser?.role == "admin") {
                if (currentAdminUser?.id == userId) {
                    Toast.makeText(this, getString(R.string.action_not_allowed_for_self), Toast.LENGTH_SHORT).show()
                    val userIndex = allUsersList.indexOfFirst { it.id == userId }
                    if (userIndex != -1) {
                        userAdminAdapter.notifyItemChanged(userIndex)
                    }
                } else {
                    performChangeUserRole(userId, newRoleId, oldRoleName)
                }
            } else {
                Toast.makeText(this, getString(R.string.action_not_allowed_not_admin), Toast.LENGTH_SHORT).show()
                val userIndex = allUsersList.indexOfFirst { it.id == userId }
                if (userIndex != -1) {
                    userAdminAdapter.notifyItemChanged(userIndex)
                }
            }
        }
        rvUsers.adapter = userAdminAdapter
    }

    private fun fetchCurrentAdminDetailsAndThenUsers(adminUserId: Int) {
        showLoading(true)
        tvAdminError.visibility = View.GONE

        lifecycleScope.launch {
            try {
                val adminResponse = apiService.getUserDetails(adminUserId)
                if (adminResponse.isSuccessful && adminResponse.body() != null) {
                    currentAdminUser = adminResponse.body()
                    setupRecyclerView()
                    fetchAllUsers()
                } else {
                    Log.e(TAG, "Error fetching admin details: ${adminResponse.code()} - ${adminResponse.message()}")
                    showError(getString(R.string.error_loading_user_details) + " Code: ${adminResponse.code()}")
                }
            } catch (e: Exception) {
                Log.e(TAG, "Exception fetching admin details", e)
                showError(getString(R.string.error_loading_user_details) + ": ${e.localizedMessage}")
            } finally {
            }
        }
    }

    private fun fetchAllUsers() {
        lifecycleScope.launch {
            try {
                val usersResponse = apiService.getAllUsers()
                    allUsersList.clear()
                    allUsersList.addAll(usersResponse)
                    userAdminAdapter.updateUsers(allUsersList)
                    tvAdminError.visibility = View.GONE
            } catch (e: Exception) {
                Log.e(TAG, "Exception fetching users", e)
                showError(getString(R.string.error_loading_users) + ": ${e.localizedMessage}")
            } finally {
                showLoading(false)
            }
        }
    }

    private fun performChangeUserRole(userId: Int, newRoleId: Int, oldRoleName: String) {
        showLoading(true)
        lifecycleScope.launch {
            try {
                val request = ChangeRoleRequest(userId, newRoleId)
                val response = apiService.changeUserRole(request)

                if (response.isSuccessful) {
                    Toast.makeText(applicationContext, getString(R.string.role_changed_successfully), Toast.LENGTH_SHORT).show()
                    val newRoleName = availableRoles.find { it.id == newRoleId }?.name ?: oldRoleName
                    val userIndex = allUsersList.indexOfFirst { it.id == userId }
                    if (userIndex != -1) {
                        allUsersList[userIndex].role = newRoleName
                        userAdminAdapter.updateUserRoleLocally(userId, newRoleName)
                    }
                } else {
                    Log.e(TAG, "Error changing role: ${response.code()} - ${response.message()}")
                    Toast.makeText(applicationContext, getString(R.string.error_changing_role) + " Code: ${response.code()}", Toast.LENGTH_LONG).show()
                    val userIndex = allUsersList.indexOfFirst { it.id == userId }
                    if (userIndex != -1) {
                        allUsersList[userIndex].role = oldRoleName
                        userAdminAdapter.notifyItemChanged(userIndex)
                    }
                }
            } catch (e: Exception) {
                Log.e(TAG, "Exception changing role", e)
                Toast.makeText(applicationContext, getString(R.string.error_changing_role) + ": ${e.localizedMessage}", Toast.LENGTH_LONG).show()
                val userIndex = allUsersList.indexOfFirst { it.id == userId }
                if (userIndex != -1) {
                    allUsersList[userIndex].role = oldRoleName
                    userAdminAdapter.notifyItemChanged(userIndex)
                }
            } finally {
                showLoading(false)
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        progressBarAdmin.visibility = if (isLoading) View.VISIBLE else View.GONE
        rvUsers.visibility = if (isLoading) View.GONE else View.VISIBLE
    }

    private fun showError(message: String) {
        tvAdminError.text = message
        tvAdminError.visibility = View.VISIBLE
        progressBarAdmin.visibility = View.GONE
        rvUsers.visibility = View.GONE
    }

    private fun navigateToLogin() {
        val intent = Intent(this, LoginActivity::class.java)
        intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
        startActivity(intent)
        finish()
    }
}