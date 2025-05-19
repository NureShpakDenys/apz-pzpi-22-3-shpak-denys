package com.example.newapp.data.fragments

import android.app.Activity
import android.app.AlertDialog
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.*
import android.widget.Button
import android.widget.ProgressBar
import android.widget.TextView
import android.widget.Toast
import androidx.activity.result.ActivityResultLauncher
import androidx.activity.result.contract.ActivityResultContracts
import androidx.fragment.app.Fragment
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.adapters.UsersAdapter
import com.example.newapp.data.models.CompanyUser
import com.example.newapp.data.models.RemoveUserRequest
import com.example.newapp.data.models.UpdateUserRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.company.AddUserToCompanyActivity
import com.example.newapp.ui.company.CompanyDetailsActivity
import kotlinx.coroutines.launch

class UsersFragment : Fragment() {

    private lateinit var recyclerView: RecyclerView
    private lateinit var progressBar: ProgressBar
    private lateinit var noUsersText: TextView
    private lateinit var addUserButton: Button
    private lateinit var usersAdapter: UsersAdapter

    private var users: MutableList<CompanyUser> = mutableListOf()
    private var companyId: Int = 0
    private var creatorId: Int = 0
    private var isCreator: Boolean = false

    private lateinit var addUserToCompanyLauncher: ActivityResultLauncher<Intent>

    companion object {
        private const val TAG = "UsersFragment"

        fun newInstance(
            users: List<CompanyUser>,
            creatorId: Int,
            companyId: Int,
            isCreator: Boolean = false,
        ): UsersFragment {
            val fragment = UsersFragment()
            fragment.users = users.toMutableList()
            fragment.creatorId = creatorId
            fragment.companyId = companyId
            fragment.isCreator = isCreator
            Log.d(TAG, "newInstance called. Users count: ${users.size}, Company ID: $companyId, isCreator: $isCreator")

            return fragment
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        Log.d(TAG, "onCreate called for Company ID: $companyId")
        addUserToCompanyLauncher = registerForActivityResult(ActivityResultContracts.StartActivityForResult()) { result ->
            Log.d(TAG, "ActivityResult received: resultCode=${result.resultCode} for Company ID: $companyId")
            if (result.resultCode == Activity.RESULT_OK) {
                Toast.makeText(requireContext(), "User added, refreshing data...", Toast.LENGTH_SHORT).show()
                Log.i(TAG, "User successfully added, attempting to refresh data in parent activity.")
                (activity as? CompanyDetailsActivity)?.fetchCompanyData()
            } else {
                Log.d(TAG, "Add user to company was cancelled or failed.")
                 Toast.makeText(requireContext(), "Add user cancelled.", Toast.LENGTH_SHORT).show()
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?,
    ): View {
        val view = inflater.inflate(R.layout.fragment_users, container, false)
        Log.d(TAG, "onCreateView called for Company ID: $companyId. Initial users count: ${users.size}")

        recyclerView = view.findViewById(R.id.usersRecyclerView)
        progressBar = view.findViewById(R.id.userProgressBar)
        noUsersText = view.findViewById(R.id.tvNoUsers)
        addUserButton = view.findViewById(R.id.btnAddUser)

        recyclerView.layoutManager = LinearLayoutManager(requireContext())

        if (isCreator) {
            addUserButton.visibility = View.VISIBLE
            addUserButton.setOnClickListener {
                Log.d(TAG, "Add user button clicked for Company ID: $companyId")
                val intent = Intent(requireContext(), AddUserToCompanyActivity::class.java)
                intent.putExtra(AddUserToCompanyActivity.EXTRA_COMPANY_ID, companyId)

                addUserToCompanyLauncher.launch(intent)
            }
        } else {
            addUserButton.visibility = View.GONE
        }

        usersAdapter = UsersAdapter(
            users = this.users,
            creatorId = creatorId,
            ::onRoleChangeClick,
            ::onRemoveClick
        )
        recyclerView.adapter = usersAdapter

        updateUI()

        return view
    }

    private fun updateUI() {
        if (!isAdded) {
            Log.w(TAG, "updateUI called but fragment not added. Company ID: $companyId")
            return
        }
        Log.d(TAG, "updateUI called. Users count: ${users.size} for Company ID: $companyId")
        if (users.isEmpty()) {
            recyclerView.visibility = View.GONE
            noUsersText.visibility = View.VISIBLE
        } else {
            recyclerView.visibility = View.VISIBLE
            noUsersText.visibility = View.GONE
        }
    }

    private fun onRoleChangeClick(companyUser: CompanyUser) {
        val roles = arrayOf("user", "admin", "manager")
        AlertDialog.Builder(requireContext())
            .setTitle("Change Role")
            .setItems(roles) { _, which ->
                val selectedRole = roles[which]
                lifecycleScope.launch {
                    try {
                        val body = UpdateUserRequest(userID = companyUser.UserID, role = selectedRole)
                        val response = RetrofitClient.getInstance(requireContext())
                            .updateUserRole(companyId, body)

                        if (response.isSuccessful) {
                            Toast.makeText(requireContext(), "Role updated", Toast.LENGTH_SHORT).show()

                            val index = users.indexOfFirst { it.UserID == companyUser.UserID }
                            if (index != -1) {
                                val updatedUser = companyUser.copy(Role = selectedRole)
                                users[index] = updatedUser
                                usersAdapter.notifyItemChanged(index)
                            }
                        }else {
                            Toast.makeText(requireContext(), "Failed to update role", Toast.LENGTH_SHORT).show()
                        }
                    } catch (e: Exception) {
                        Toast.makeText(requireContext(), "Error: ${e.message}", Toast.LENGTH_SHORT).show()
                    }
                }
            }.show()
    }

    private fun onRemoveClick(companyUser: CompanyUser) {
        AlertDialog.Builder(requireContext())
            .setTitle("Remove User")
            .setMessage("Are you sure you want to remove this user?")
            .setPositiveButton("Yes") { _, _ ->
                lifecycleScope.launch {
                    try {
                        val token = SessionManager(requireContext()).getAccessToken()
                        val body = RemoveUserRequest(userID = companyUser.UserID)
                        val response = RetrofitClient.getInstance(requireContext())
                            .removeUser(companyId, body)

                        if (response.isSuccessful) {
                            Toast.makeText(requireContext(), "User removed", Toast.LENGTH_SHORT).show()
                            val index = users.indexOfFirst { it.UserID == companyUser.UserID }
                            if (index != -1) {
                                users.removeAt(index)
                                usersAdapter.notifyItemRemoved(index)
                            }
                        } else {
                            Toast.makeText(requireContext(), "Failed to remove user", Toast.LENGTH_SHORT).show()
                        }
                    } catch (e: Exception) {
                        Toast.makeText(requireContext(), "Error: ${e.message}", Toast.LENGTH_SHORT).show()
                    }
                }
            }
            .setNegativeButton("Cancel", null)
            .show()
    }
}
