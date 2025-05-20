package com.example.newapp.ui.login

import android.content.Intent
import android.os.Bundle
import android.widget.Button
import android.widget.EditText
import android.widget.Toast
import com.example.newapp.ui.base.BaseActivity
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.models.LoginRequest
import com.example.newapp.data.models.LoginResponse
import com.example.newapp.data.models.User
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.company.CompaniesActivity
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import android.util.Log

class LoginActivity : BaseActivity() {

    private lateinit var sessionManager: SessionManager

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_login)

        sessionManager = SessionManager(this)

        val etName = findViewById<EditText>(R.id.etName)
        val etPassword = findViewById<EditText>(R.id.etPassword)
        val btnLogin = findViewById<Button>(R.id.btnLogin)
        val btnRegister = findViewById<Button>(R.id.btnRegister)

        btnLogin.setOnClickListener {
            val username = etName.text.toString().trim()
            val password = etPassword.text.toString().trim()

            if (username.isEmpty() || password.isEmpty()) {
                return@setOnClickListener
            }

            loginUser(username, password)
        }

        btnRegister.setOnClickListener {
            val intent = Intent(this, RegisterActivity::class.java)
            startActivity(intent)
            finish()
        }
    }

    private fun loginUser(username: String, password: String) {
        val apiService = RetrofitClient.getInstance(this)
        val request = LoginRequest(username, password)

        apiService.login(request).enqueue(object : Callback<LoginResponse> {
            override fun onResponse(call: Call<LoginResponse>, response: Response<LoginResponse>) {

                if (response.isSuccessful) {
                    val loginResponse = response.body()
                    if (loginResponse != null) {
                        sessionManager.saveTokens(loginResponse.token)
                        sessionManager.saveUser(User(
                            id = loginResponse.user.id,
                            name = loginResponse.user.name,
                            role = loginResponse.user.role
                        ))
                        Toast.makeText(this@LoginActivity, "Hi, ${loginResponse.user.name}!", Toast.LENGTH_LONG).show()

                        val intent = Intent(this@LoginActivity, CompaniesActivity::class.java)
                        startActivity(intent)
                        finish()
                    }
                    else {
                        Toast.makeText(this@LoginActivity, "Error, the body is empty", Toast.LENGTH_SHORT).show()
                    }
                } else {
                    Toast.makeText(this@LoginActivity, "Auth error", Toast.LENGTH_SHORT).show()
                }
            }

            override fun onFailure(call: Call<LoginResponse>, t: Throwable) {
                Toast.makeText(this@LoginActivity, "Network error: ${t.message}", Toast.LENGTH_SHORT).show()
            }
        })
    }
}
