package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName
import java.util.Date

data class Delivery (
    @SerializedName("ID") val id: Int,
    @SerializedName("CompanyID") val companyId: Int,
    @SerializedName("Status") val status: String,
    @SerializedName("Date") val date: Date,
    @SerializedName("Duration") val duration: String,
    @SerializedName("company") val company: CompanyInDelivery,
    @SerializedName("products") val products: List<Product>
)

data class CompanyInDelivery(
    @SerializedName("ID") val id: Int,
    @SerializedName("Name") val name: String,
    @SerializedName("CreatorID") val creatorId: Int
)

data class Product(
    @SerializedName("ID") val id: Int,
    @SerializedName("Name") val name: String,
    @SerializedName("Weight") val weight: Float,
    @SerializedName("product_category") val productCategory: ProductCategory
)

data class ProductCategory(
    @SerializedName("ID") val id: Int,
    @SerializedName("Name") val name: String,
    @SerializedName("MinTemperature") val minTemperature: Float,
    @SerializedName("MaxTemperature") val maxTemperature: Float,
    @SerializedName("MinHumidity") val minHumidity: Int,
    @SerializedName("MaxHumidity") val maxHumidity: Int,
    @SerializedName("IsPerishable") val isPerishable: Boolean
)

data class OptimalRouteResponse(
    val route: RouteInfo,
    val message: String,
    val equation: String,
    @SerializedName("predict_data") val predictData: PredictionData
)

data class RouteInfo(
    @SerializedName("id") val id: Int,
    @SerializedName("name") val name: String
)

data class PredictionData(
    @SerializedName("Distance") val distance: Float,
    @SerializedName("Time") val time: Float,
    @SerializedName("Speed") val speed: Float
)

data class CreateDeliveryRequest(
    @SerializedName("company_id")
    val companyId: Int,
    @SerializedName("date")
    val date: String
)

data class DeliveryResponse(
    @SerializedName("id")
    val id: Int,
    @SerializedName("company_id")
    val companyId: Int? = null,
    @SerializedName("date")
    val date: String? = null,
    @SerializedName("status")
    val status: String? = null
)

data class DeliveryRouteUpdateRequest(
    val route_id: Int
)
