package com.example.newapp.ui.admin

import android.os.Bundle
import android.view.View
import android.widget.Button
import android.widget.EditText
import android.widget.ProgressBar
import android.widget.ScrollView
import android.widget.TextView
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.lifecycle.lifecycleScope
import com.example.newapp.R
import com.example.newapp.data.SessionManager
import com.example.newapp.data.models.BackupRestoreRequest
import com.example.newapp.data.models.DbStatusResponse
import com.example.newapp.data.network.ApiService
import com.example.newapp.data.network.RetrofitClient
import com.example.newapp.ui.base.BaseActivity
import kotlinx.coroutines.launch
import java.text.ParseException
import java.text.SimpleDateFormat
import java.util.Date
import java.util.Locale
import java.util.TimeZone

class DatabaseAdminActivity : BaseActivity() {

    private lateinit var apiService: ApiService
    private lateinit var sessionManager: SessionManager

    private lateinit var progressBar: ProgressBar
    private lateinit var contentScrollView: ScrollView
    private lateinit var tvDbAdminMessage: TextView

    private lateinit var tvDbSize: TextView
    private lateinit var tvActiveConnections: TextView
    private lateinit var tvLastBackup: TextView
    private lateinit var btnRefreshDbStatus: Button

    private lateinit var etBackupPath: EditText
    private lateinit var btnPerformBackup: Button
    private lateinit var btnPerformRestore: Button

    private lateinit var btnOptimizeDb: Button

    private val outputDateTimeFormat = SimpleDateFormat("yyyy-MM-dd HH:mm:ss", Locale.getDefault())
    private val inputDateTimeFormats = listOf(
        SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", Locale.US).apply { timeZone = TimeZone.getTimeZone("UTC") },
        SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.US).apply { timeZone = TimeZone.getTimeZone("UTC") }
    )

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_database_admin)

        sessionManager = SessionManager(this)
        apiService = RetrofitClient.getInstance(this)

        initializeUI()
        setupListeners()
        fetchDbStatusData()
    }

    private fun initializeUI() {
        progressBar = findViewById(R.id.progressBarDatabaseAdmin)
        contentScrollView = findViewById(R.id.contentScrollViewDatabaseAdmin)
        tvDbAdminMessage = findViewById(R.id.tvDbAdminMessage)

        tvDbSize = findViewById(R.id.tvDbSize)
        tvActiveConnections = findViewById(R.id.tvActiveConnections)
        tvLastBackup = findViewById(R.id.tvLastBackup)
        btnRefreshDbStatus = findViewById(R.id.btnRefreshDbStatus)

        etBackupPath = findViewById(R.id.etBackupPath)
        btnPerformBackup = findViewById(R.id.btnPerformBackup)
        btnPerformRestore = findViewById(R.id.btnPerformRestore)

        btnOptimizeDb = findViewById(R.id.btnOptimizeDb)

        tvDbAdminMessage.visibility = View.GONE
    }

    private fun setupListeners() {
        btnRefreshDbStatus.setOnClickListener { fetchDbStatusData() }
        btnPerformBackup.setOnClickListener { performBackup() }
        btnPerformRestore.setOnClickListener { performRestore() }
        btnOptimizeDb.setOnClickListener { optimizeDatabase() }
    }

    private fun showLoading(isLoading: Boolean) {
        progressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        contentScrollView.visibility = if (isLoading) View.GONE else View.VISIBLE
        if (isLoading) {
            tvDbAdminMessage.visibility = View.GONE
        }
    }

    private fun showMessage(message: String, isError: Boolean = false) {
        tvDbAdminMessage.text = message
        if (isError) {
            tvDbAdminMessage.setBackgroundResource(R.drawable.bg_message_error)
            tvDbAdminMessage.setTextColor(ContextCompat.getColor(this, android.R.color.white))
        } else {
            tvDbAdminMessage.setBackgroundResource(R.drawable.bg_message_success)
            tvDbAdminMessage.setTextColor(ContextCompat.getColor(this, android.R.color.black))
        }
        tvDbAdminMessage.visibility = View.VISIBLE
    }

    private fun formatLastBackupTime(timeString: String?): String {
        if (timeString.isNullOrEmpty()) {
            return getString(R.string.na)
        }
        for (format in inputDateTimeFormats) {
            try {
                val date: Date? = format.parse(timeString)
                if (date != null) {
                    return outputDateTimeFormat.format(date)
                }
            } catch (e: ParseException) {
            }
        }
        return timeString
    }


    private fun updateDbStatusUI(status: DbStatusResponse) {
        tvDbSize.text = getString(R.string.database_size_value, status.databaseSizeMB)
        tvActiveConnections.text = getString(R.string.active_connections_value, status.activeConnections)
        tvLastBackup.text = getString(R.string.last_backup_value, formatLastBackupTime(status.lastBackupTime))
    }

    private fun fetchDbStatusData() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                val response = apiService.getDbStatus()
                if (response.isSuccessful) {
                    response.body()?.let {
                        updateDbStatusUI(it)
                    } ?: showError(getString(R.string.error_fetching_db_status) + " (empty response)")
                } else {
                    showError(getString(R.string.error_fetching_db_status) + ": ${response.code()} ${response.message()}")
                }
            } catch (e: Exception) {
                showError(getString(R.string.error_fetching_db_status) + ": ${e.message}")
            } finally {
                showLoading(false)
            }
        }
    }

    private fun performBackup() {
        val path = etBackupPath.text.toString().trim()
        if (path.isEmpty()) {
            Toast.makeText(this, getString(R.string.backup_path_required), Toast.LENGTH_SHORT).show()
            return
        }
        showLoading(true)
        lifecycleScope.launch {
            try {
                val request = BackupRestoreRequest(backupPath = path)
                val response = apiService.performBackup(request)

                if (response.isSuccessful) {
                    response.body()?.let {
                        showMessage("✅ ${it.message}")
                        fetchDbStatusData()
                    } ?: showMessage(getString(R.string.operation_failed, "Backup: Empty response"), true)
                } else {
                    showMessage(getString(R.string.operation_failed, "Backup: ${response.code()} ${response.message()}"), true)
                }
            } catch (e: Exception) {
                showMessage(getString(R.string.operation_failed, "Backup: ${e.message}"), true)
            } finally {
                showLoading(false)
            }
        }
    }

    private fun performRestore() {
        val path = etBackupPath.text.toString().trim()
        if (path.isEmpty()) {
            Toast.makeText(this, getString(R.string.backup_path_required), Toast.LENGTH_SHORT).show()
            return
        }
        showLoading(true)
        lifecycleScope.launch {
            try {
                val request = BackupRestoreRequest(backupPath = path)
                val response = apiService.performRestore(request)

                if (response.isSuccessful) {
                    response.body()?.let {
                        showMessage("✅ ${it.message}")
                        fetchDbStatusData()
                    } ?: showMessage(getString(R.string.operation_failed, "Restore: Empty response"), true)
                } else {
                    showMessage(getString(R.string.operation_failed, "Restore: ${response.code()} ${response.message()}"), true)
                }
            } catch (e: Exception) {
                showMessage(getString(R.string.operation_failed, "Restore: ${e.message}"), true)
            } finally {
                showLoading(false)
            }
        }
    }

    private fun optimizeDatabase() {
        showLoading(true)
        lifecycleScope.launch {
            try {
                val response = apiService.optimizeDb()

                if (response.isSuccessful) {
                    response.body()?.let {
                        showMessage("✅ ${it.message}")
                        fetchDbStatusData()
                    } ?: showMessage(getString(R.string.operation_failed, "Optimize: Empty response"), true)
                } else {
                    showMessage(getString(R.string.operation_failed, "Optimize: ${response.code()} ${response.message()}"), true)
                }
            } catch (e: Exception) {
                showMessage(getString(R.string.operation_failed, "Optimize: ${e.message}"), true)
            } finally {
                showLoading(false)
            }
        }
    }

    private fun showError(message: String) {
        Toast.makeText(this, message, Toast.LENGTH_LONG).show()
    }
}