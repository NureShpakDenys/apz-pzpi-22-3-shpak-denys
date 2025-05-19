package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class SystemHealth(
    @SerializedName("db_status")
    val dbStatus: String?,
    @SerializedName("server_time")
    val serverTime: String?,
    @SerializedName("uptime")
    val uptime: String?
)

data class SystemConfigs(
    @SerializedName("auth_token_ttl_hours")
    var authTokenTtlHours: Int?,
    @SerializedName("encryption_key_exists")
    val encryptionKeyExists: Boolean?,
    @SerializedName("http_timeout_seconds")
    var httpTimeoutSeconds: Int?
)

data class SystemLog(
    @SerializedName("ID") 
    val id: Int?,
    @SerializedName("CreatedAt") 
    val createdAt: String?, 
    @SerializedName("UserID") 
    val userId: Int?,
    @SerializedName("ActionType") 
    val actionType: String?,
    @SerializedName("Description") 
    val description: String?,
    @SerializedName("Success") 
    val success: Boolean?
)

data class UpdateSystemConfigsRequest(
    @SerializedName("timeout_sec")
    val timeoutSec: Int,
    @SerializedName("token_ttl")
    val tokenTtl: Int
)

data class ClearLogsRequest(
    val days: Int
)