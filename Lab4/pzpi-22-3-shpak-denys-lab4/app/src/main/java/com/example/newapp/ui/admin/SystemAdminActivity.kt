package com.example.newapp.ui.admin

import android.content.Intent
import android.os.Bundle
import android.view.View
import android.widget.Button
import android.widget.EditText
import android.widget.LinearLayout
import android.widget.ProgressBar
import android.widget.ScrollView
import android.widget.TextView
import android.widget.Toast
import androidx.appcompat.widget.Toolbar
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.adapters.SystemLogAdapter
import com.example.newapp.data.models.ClearLogsRequest
import com.example.newapp.data.models.SystemConfigs
import com.example.newapp.data.models.SystemHealth
import com.example.newapp.data.models.SystemLog
import com.example.newapp.data.models.UpdateSystemConfigsRequest
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.data.network.ApiService
import com.example.newapp.ui.base.BaseActivity
import com.example.newapp.ui.login.LoginActivity
import com.google.android.material.textfield.TextInputEditText
import kotlinx.coroutines.launch
import java.text.SimpleDateFormat
import java.util.Locale
import java.util.TimeZone

class SystemAdminActivity : BaseActivity() {

    private lateinit var sessionManager: SessionManager
    private lateinit var apiService: ApiService

    private lateinit var progressBar: ProgressBar
    private lateinit var contentScrollView: ScrollView

    private lateinit var tvDbStatus: TextView
    private lateinit var tvServerTime: TextView
    private lateinit var tvUptime: TextView
    private lateinit var btnRefreshHealth: Button

    private lateinit var tvEncryptionKeyExists: TextView
    private lateinit var etAuthTokenTtl: EditText
    private lateinit var etHttpTimeout: EditText
    private lateinit var btnSaveConfigs: Button

    private lateinit var etClearLogsDays: EditText
    private lateinit var btnClearLogs: Button
    private lateinit var btnRefreshLogs: Button
    private lateinit var rvSystemLogs: RecyclerView
    private lateinit var logAdapter: SystemLogAdapter
    private var allLogs: List<SystemLog> = listOf()
    private var currentPage: Int = 1
    private val logsPerPage: Int = 10
    private lateinit var tvLogsPageInfo: TextView
    private lateinit var btnPrevLogsPage: Button
    private lateinit var btnNextLogsPage: Button
    private lateinit var layoutLogPagination: LinearLayout


    private val inputDateFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.US).apply {
        timeZone = TimeZone.getTimeZone("UTC")
    }
    private val outputDateFormat = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault())

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_system_admin)

        sessionManager = SessionManager(this)
        apiService = RetrofitClient.getInstance(this)

        initializeUI()
        setupListeners()
        fetchAllData()
    }

    private fun initializeUI() {
        progressBar = findViewById(R.id.progressBarSystemAdmin)
        contentScrollView = findViewById(R.id.contentScrollViewSystemAdmin)

        tvDbStatus = findViewById(R.id.tvDbStatus)
        tvServerTime = findViewById(R.id.tvServerTime)
        tvUptime = findViewById(R.id.tvUptime)
        btnRefreshHealth = findViewById(R.id.btnRefreshHealth)

        tvEncryptionKeyExists = findViewById(R.id.tvEncryptionKeyExists)
        etAuthTokenTtl = findViewById(R.id.etAuthTokenTtl)
        etHttpTimeout = findViewById(R.id.etHttpTimeout)
        btnSaveConfigs = findViewById(R.id.btnSaveConfigs)

        etClearLogsDays = findViewById(R.id.etClearLogsDays)
        btnClearLogs = findViewById(R.id.btnClearLogs)
        btnRefreshLogs = findViewById(R.id.btnRefreshLogs)
        rvSystemLogs = findViewById(R.id.rvSystemLogs)
        tvLogsPageInfo = findViewById(R.id.tvLogsPageInfo)
        btnPrevLogsPage = findViewById(R.id.btnPrevLogsPage)
        btnNextLogsPage = findViewById(R.id.btnNextLogsPage)
        layoutLogPagination = findViewById(R.id.layoutLogPagination)

        rvSystemLogs.layoutManager = LinearLayoutManager(this)
        logAdapter = SystemLogAdapter(this, listOf())
        rvSystemLogs.adapter = logAdapter
        rvSystemLogs.isNestedScrollingEnabled = false
    }

    private fun setupListeners() {
        btnRefreshHealth.setOnClickListener { fetchHealthData() }
        btnSaveConfigs.setOnClickListener { saveSystemConfigs() }
        btnClearLogs.setOnClickListener { clearSystemLogsData() }
        btnRefreshLogs.setOnClickListener { fetchLogsData() }

        btnPrevLogsPage.setOnClickListener {
            if (currentPage > 1) {
                currentPage--
                updateLogsDisplay()
            }
        }
        btnNextLogsPage.setOnClickListener {
            val totalPages = getTotalPages()
            if (currentPage < totalPages) {
                currentPage++
                updateLogsDisplay()
            }
        }
    }

    private fun showLoading(isLoading: Boolean) {
        progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        contentScrollView.visibility = if (isLoading) View.GONE else View.VISIBLE
    }

    private fun fetchAllData() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                fetchHealthDataInternal()
                fetchConfigsDataInternal()
                fetchLogsDataInternal()
            } catch (e: Exception) {
                Toast.makeText(this@SystemAdminActivity, getString(R.string.error_loading_data) + ": ${e.message}", Toast.LENGTH_LONG).show()
            } finally {
                showLoading(false)
            }
        }
    }

    private fun fetchHealthData() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                fetchHealthDataInternal()
            } catch (e: Exception) {
                 Toast.makeText(this@SystemAdminActivity, getString(R.string.error_loading_data) + ": ${e.message}", Toast.LENGTH_SHORT).show()
            } finally {
                showLoading(false)
            }
        }
    }

    private suspend fun fetchHealthDataInternal() {
        try {
            val response = apiService.getSystemHealth()
            if (response.isSuccessful) {
                response.body()?.let { updateHealthUI(it) }
            } else {
                showError("Health Check: ${response.code()} ${response.message()}")
            }
        } catch (e: Exception) {
            showError("Health Check Error: ${e.message}")
            throw e
        }
    }

    private fun fetchConfigsData() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                fetchConfigsDataInternal()
            } catch (e: Exception) {
                Toast.makeText(this@SystemAdminActivity, getString(R.string.error_loading_data) + ": ${e.message}", Toast.LENGTH_SHORT).show()
            } finally {
                 showLoading(false)
            }
        }
    }

    private suspend fun fetchConfigsDataInternal() {
        try {
            val response = apiService.getSystemConfigs()
            if (response.isSuccessful) {
                response.body()?.let { updateConfigsUI(it) }
            } else {
                 showError("System Configs: ${response.code()} ${response.message()}")
            }
        } catch (e: Exception) {
             showError("System Configs Error: ${e.message}")
            throw e
        }
    }

    private fun fetchLogsData() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                fetchLogsDataInternal()
            } catch (e: Exception) {
                Toast.makeText(this@SystemAdminActivity, getString(R.string.error_loading_data) + ": ${e.message}", Toast.LENGTH_SHORT).show()
            } finally {
                showLoading(false)
            }
        }
    }

    private suspend fun fetchLogsDataInternal() {
        try {
            val response = apiService.getSystemLogs()
            if (response.isSuccessful) {
                allLogs = response.body() ?: emptyList()
                currentPage = 1
                updateLogsDisplay()
            } else {
                showError("System Logs: ${response.code()} ${response.message()}")
            }
        } catch (e: Exception) {
            showError("System Logs Error: ${e.message}")
            throw e
        }
    }

    private fun updateHealthUI(health: SystemHealth) {
        tvDbStatus.text = getString(R.string.db_status_prefix, health.dbStatus ?: "N/A")
        tvServerTime.text = getString(R.string.server_time_prefix, formatDisplayDate(health.serverTime))
        tvUptime.text = getString(R.string.uptime_prefix, health.uptime ?: "N/A")
    }

    private fun updateConfigsUI(configs: SystemConfigs) {
        tvEncryptionKeyExists.text = getString(
            R.string.encryption_key_exists_prefix,
            if (configs.encryptionKeyExists == true) getString(R.string.yes_exists) else getString(R.string.no_not_exists)
        )
        etAuthTokenTtl.setText(configs.authTokenTtlHours?.toString() ?: "0")
        etHttpTimeout.setText(configs.httpTimeoutSeconds?.toString() ?: "0")
    }

    private fun saveSystemConfigs() {
        val ttlText = etAuthTokenTtl.text.toString()
        val timeoutText = etHttpTimeout.text.toString()

        if (ttlText.isBlank() || timeoutText.isBlank()) {
            Toast.makeText(this, getString(R.string.invalid_input_configs), Toast.LENGTH_SHORT).show()
            return
        }

        val ttl = ttlText.toIntOrNull()
        val timeout = timeoutText.toIntOrNull()

        if (ttl == null || timeout == null || ttl < 0 || timeout < 0) {
            Toast.makeText(this, getString(R.string.invalid_input_configs), Toast.LENGTH_SHORT).show()
            return
        }

        val request = UpdateSystemConfigsRequest(timeoutSec = timeout, tokenTtl = ttl)
        showLoading(true)
        lifecycleScope.launch {
            try {
                val response = apiService.updateSystemConfigs(request)
                if (response.isSuccessful) {
                    Toast.makeText(this@SystemAdminActivity, getString(R.string.configs_saved_successfully), Toast.LENGTH_SHORT).show()
                    fetchConfigsDataInternal()
                } else {
                    showError("Save Configs: ${response.code()} ${response.message()}")
                }
            } catch (e: Exception) {
                showError("Save Configs Error: ${e.message}")
            } finally {
                showLoading(false)
            }
        }
    }

    private fun clearSystemLogsData() {
        val daysText = etClearLogsDays.text.toString()

        if (daysText.isBlank()) {
            Toast.makeText(this, getString(R.string.invalid_input_clear_days), Toast.LENGTH_SHORT).show()
            return
        }
        val days = daysText.toIntOrNull()

        if (days == null || days < 0) {
             Toast.makeText(this, getString(R.string.invalid_input_clear_days), Toast.LENGTH_SHORT).show()
            return
        }

        val request = ClearLogsRequest(days = days)
        showLoading(true)
        lifecycleScope.launch {
            try {
                val response = apiService.clearSystemLogs(request)
                if (response.isSuccessful) {
                    Toast.makeText(this@SystemAdminActivity, getString(R.string.logs_cleared_successfully), Toast.LENGTH_SHORT).show()
                    fetchLogsDataInternal()
                } else {
                     showError("Clear Logs: ${response.code()} ${response.message()}")
                }
            } catch (e: Exception) {
                showError("Clear Logs Error: ${e.message}")
            } finally {
                showLoading(false)
            }
        }
    }


    private fun updateLogsDisplay() {
        val totalPages = getTotalPages()
        if (allLogs.isEmpty()) {
            layoutLogPagination.visibility = View.GONE
            logAdapter.updateLogs(emptyList())
            tvLogsPageInfo.text = getString(R.string.page_info, 0, 0)
            return
        }

        layoutLogPagination.visibility = View.VISIBLE
        val startIndex = (currentPage - 1) * logsPerPage
        val endIndex = minOf(startIndex + logsPerPage, allLogs.size)
        val logsToShow = if (startIndex < allLogs.size) allLogs.subList(startIndex, endIndex) else emptyList()

        logAdapter.updateLogs(logsToShow)
        tvLogsPageInfo.text = getString(R.string.page_info, currentPage, totalPages)
        btnPrevLogsPage.isEnabled = currentPage > 1
        btnNextLogsPage.isEnabled = currentPage < totalPages
    }

    private fun getTotalPages(): Int {
        return if (allLogs.isEmpty()) 0 else (allLogs.size + logsPerPage - 1) / logsPerPage
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

    private fun showError(message: String) {
        Toast.makeText(this, message, Toast.LENGTH_LONG).show()
    }
}