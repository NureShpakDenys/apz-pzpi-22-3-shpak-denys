package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class CreateRouteRequest(
    @SerializedName("company_id")
    val companyId: Int,
    @SerializedName("name")
    val name: String
)

data class CreateRouteResponse(
    @SerializedName("id")
    val id: Int,
    @SerializedName("name")
    val name: String,
    @SerializedName("company_id")
    val companyId: Int
)

data class Route(
    val id: Int,
    val name: String,
    val status: String,
    val details: String
)

data class RouteResponse(
    val ID: Int,
    val Name: String,
    val Status: String,
    val Details: String,
    val company: RoutesCompany?,
    val waypoints: List<Waypoint>
)

data class RoutesCompany(
    val ID: Int,
    val Name: String,
    val Address: String,
    val CreatorID: String,
)

data class Waypoint(
    val ID: Int,
    val Name: String,
    val Latitude: Double,
    val Longitude: Double,
    val Status: String,
    val Details: String,
    val DeviceSerial: String,
    val SendDataFrequency: Int,
    val GetWeatherAlerts: Boolean,
    val RouteID: Int
)

data class WeatherAlert(
    val type: String,
    val message: String,
    val details: String
)

data class WeatherAlertResponse(
    val alerts: WeatherAlert
)