package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName
import java.util.Date

data class Company(
    @SerializedName("id")
    val id: Int,

    @SerializedName("name")
    val name: String,

    val address: String
)

data class CompanyUser(
    val UserID: Int,
    val Role: String,
    val user: CompanyUsersUser
)

data class CreateCompanyRequest(
    val name: String,
    val address: String
)

data class CompanyUsersUser(
    val name: String
)

data class CompanyResponse(
    val id: Int,
    val name: String,
    val address: String,
    val creator: User,
    val routes: List<Route>,
    val deliveries: List<CompanyDelivery>,
    val users: List<User>
)

data class CompanyDelivery(
    val id: Int,
    val status: String,
    val date: Date,
    val duration: String,
    val routeId: Int
)

data class UpdateCompanyRequest(
    val name: String,
    val address: String
)