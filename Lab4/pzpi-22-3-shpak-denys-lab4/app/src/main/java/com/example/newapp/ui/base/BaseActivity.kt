package com.example.newapp.ui.base

import android.content.Intent
import android.os.Bundle
import android.widget.FrameLayout
import android.widget.TextView
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.ui.admin.AdminDashboardActivity
import com.example.newapp.ui.admin.DatabaseAdminActivity
import com.example.newapp.ui.admin.SystemAdminActivity
import com.example.newapp.ui.company.CompaniesActivity
import com.example.newapp.ui.login.LoginActivity
import com.google.android.material.bottomnavigation.BottomNavigationView

abstract class BaseActivity : AppCompatActivity() {
    private lateinit var contentFrame: FrameLayout
    private lateinit var sessionManager: SessionManager

    private var isBaseContentSet = false

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        sessionManager = SessionManager(this)
    }

    private fun setupBottomNavigation() {
        val bottomNavigation = findViewById<BottomNavigationView>(R.id.bottom_navigation)

        bottomNavigation.setOnItemSelectedListener { item ->
            if (item.itemId == bottomNavigation.selectedItemId && shouldPreventReload(item.itemId)) {
                return@setOnItemSelectedListener true
            }

            when (item.itemId) {
                R.id.nav_home -> {
                     if (this !is CompaniesActivity) {
                         val intent = Intent(this@BaseActivity, CompaniesActivity::class.java)
                         intent.flags = Intent.FLAG_ACTIVITY_REORDER_TO_FRONT
                         startActivity(intent)
                     }
                    true
                }
                R.id.nav_admin -> {
                    val userRole = sessionManager.getUser()?.role

                    if ("admin" == userRole) {
                        if (this !is AdminDashboardActivity) {
                            val adminIntent = Intent(this@BaseActivity, AdminDashboardActivity::class.java)
                            adminIntent.flags = Intent.FLAG_ACTIVITY_REORDER_TO_FRONT
                            startActivity(adminIntent)
                        }
                        true
                    } else if ("system_admin" == userRole) {
                        if (this !is SystemAdminActivity) {
                            val adminIntent = Intent(this@BaseActivity, SystemAdminActivity::class.java)
                            adminIntent.flags = Intent.FLAG_ACTIVITY_REORDER_TO_FRONT
                            startActivity(adminIntent)
                        }
                        true
                    } else if ("db_admin" == userRole) {
                        if (this !is DatabaseAdminActivity) {
                            val adminIntent = Intent(this@BaseActivity, DatabaseAdminActivity::class.java)
                            adminIntent.flags = Intent.FLAG_ACTIVITY_REORDER_TO_FRONT
                            startActivity(adminIntent)
                        }
                        true
                    } else {
                        Toast.makeText(this@BaseActivity, "Access for role '$userRole' will be implemented later.", Toast.LENGTH_LONG).show()
                        false
                    }
                }
                R.id.nav_logout -> {
                    sessionManager.clearSession()
                    val loginIntent = Intent(this@BaseActivity, LoginActivity::class.java)
                    loginIntent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
                    startActivity(loginIntent)
                    finish()
                    true
                }
                else -> false
            }
        }
    }

    private fun shouldPreventReload(itemId: Int): Boolean {
        return when (itemId) {
            R.id.nav_admin -> this is AdminDashboardActivity
            else -> false
        }
    }

    override fun setContentView(layoutResID: Int) {
        if (!isBaseContentSet) {
            super.setContentView(R.layout.activity_base)
            isBaseContentSet = true
            contentFrame = findViewById(R.id.content_frame)
            val header = findViewById<TextView>(R.id.header)
            header.text = "Wayra"
            setupBottomNavigation()
        }
        val content = layoutInflater.inflate(layoutResID, contentFrame, false)
        contentFrame.removeAllViews()
        contentFrame.addView(content)
    }

    override fun onResume() {
        super.onResume()
        updateBottomNavigationSelection()
    }

    private fun updateBottomNavigationSelection() {
        if (!::sessionManager.isInitialized) return

        val bottomNavigation = findViewById<BottomNavigationView>(R.id.bottom_navigation)
        when {
            this is AdminDashboardActivity && sessionManager.getUser()?.role == "admin" -> bottomNavigation.selectedItemId = R.id.nav_admin
            else -> { }
        }
    }
}