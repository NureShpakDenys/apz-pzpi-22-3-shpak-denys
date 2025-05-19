package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class CreateWaypointRequest(
    @SerializedName("route_id")
    val routeId: Int,
    @SerializedName("name")
    val name: String,
    @SerializedName("device_serial")
    val deviceSerial: String,
    @SerializedName("send_data_frequency")
    val sendDataFrequency: Int,
    @SerializedName("get_weather_alerts")
    val getWeatherAlerts: Boolean
)

data class CreateWaypointResponse(
    val id: Int,
    val name: String,
    val device_serial: String,
    val latitude: Int,
    val longitude: Int,
    val send_data_frequency: Int,
    val get_weather_alerts: Boolean,
    val status: String,
    val details: String,
)

data class UpdateWaypointRequest(
    @SerializedName("name")
    val name: String,
    @SerializedName("device_serial")
    val deviceSerial: String,
    @SerializedName("send_data_frequency")
    val sendDataFrequency: Int,
    @SerializedName("get_weather_alerts")
    val getWeatherAlerts: Boolean
)