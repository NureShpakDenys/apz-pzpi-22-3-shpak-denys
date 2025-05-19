package com.example.newapp.data

import android.content.Context
import android.content.SharedPreferences
import com.example.newapp.data.models.User

class SessionManager(context: Context) {

    private val prefs: SharedPreferences =
        context.getSharedPreferences("user_session", Context.MODE_PRIVATE)

    fun saveTokens(accessToken: String) {
        prefs.edit().apply {
            putString("ACCESS_TOKEN", accessToken)
            apply()
        }
    }

    fun saveUser(user: User) {
        prefs.edit().apply {
            putInt("USER_ID", user.id)
            putString("USER_NAME", user.name)
            putString("USER_ROLE", user.role)
            apply()
        }
    }

    fun getAccessToken(): String? {
        return prefs.getString("ACCESS_TOKEN", null)
    }

    fun getUser(): User? {
        val userId = prefs.getInt("USER_ID", -1)
        val userName = prefs.getString("USER_NAME", null)
        val userRole = prefs.getString("USER_ROLE", null)

        return if (userId != -1 && userName != null && userRole != null) {
            User(id = userId, name = userName, role = userRole)
        } else {
            null
        }
    }

    fun clearSession() {
        prefs.edit().clear().apply()
    }
}
