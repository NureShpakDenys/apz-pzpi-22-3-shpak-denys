package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class User(
    var id: Int,
    var name: String,
    var role: String
)

data class UserResponse(
    val id: String,
    val name: String,
    val bio: String,
)

data class UpdateUserRequest(
    val userID: Int,
    val role: String
)

data class RemoveUserRequest(
    val userID: Int
)

data class AddUserToCompanyRequest(
    @SerializedName("userID") val userId: Int,
    @SerializedName("role") val role: String
)

data class ChangeRoleRequest(
    @SerializedName("userId")
    val userId: Int,
    @SerializedName("roleId")
    val roleId: Int
)