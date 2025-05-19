package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class DbStatusResponse(
    @SerializedName("DatabaseSizeMB")
    val databaseSizeMB: Double = 0.0,
    @SerializedName("ActiveConnections")
    val activeConnections: Int = 0,
    @SerializedName("LastBackupTime")
    val lastBackupTime: String? = null
)

data class BackupRestoreRequest(
    @SerializedName("backup_path")
    val backupPath: String
)

data class AdminActionResponse(
    val message: String
)