package com.example.newapp.data.network

import com.example.newapp.data.models.*
import retrofit2.Call
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    @POST("auth/login/")
    fun login(@Body request: LoginRequest): Call<LoginResponse>

    @POST("auth/register/")
    fun register(@Body request: RegisterRequest): Call<Void>

    @GET("user/{userId}")
    suspend fun getUserDetails(
        @Path("userId") userId: Int
    ): Response<User>

    @GET("company/")
    fun getCompanies(): Call<List<Company>>

    @GET("company/{id}")
    suspend fun getCompany(@Path("id") id: Int): CompanyResponse

    @GET("company/{id}/users")
    suspend fun getCompanyUsers(@Path("id") id: Int): List<CompanyUser>

    @PUT("company/{id}/update-user")
    suspend fun updateUserRole(@Path("id") id: Int, @Body body: UpdateUserRequest): Response<Unit>

    @HTTP(method = "DELETE", path = "company/{id}/remove-user", hasBody = true)
    suspend fun removeUser(@Path("id") id: Int, @Body body: RemoveUserRequest): Response<Unit>

    @GET("delivery/{id}")
    suspend fun getDelivery(@Path("id") id: Int, @Header("Authorization") token: String): CompanyDelivery

    @PUT("delivery/{id}")
    suspend fun updateDeliveryRoute(
        @Path("id") id: Int,
        @Body routeUpdate: DeliveryRouteUpdateRequest,
        @Header("Authorization") token: String
    ): Response<Unit>

    @POST("delivery/")
    suspend fun createDelivery(
        @Body deliveryRequest: CreateDeliveryRequest
    ): Response<DeliveryResponse>

    @GET("routes/{id}")
    suspend fun getRoute(
        @Path("id") id: Int,
    ): RouteResponse

    @GET("routes/{id}/weather-alert")
    suspend fun getWeatherAlert(
        @Path("id") id: Int,
    ): WeatherAlertResponse

    @DELETE("routes/{id}")
    suspend fun deleteRoute(
        @Path("id") id: Int,
    ): Response<Unit>

    @DELETE("waypoints/{id}")
    suspend fun deleteWaypoint(
        @Path("id") id: Int,
    ): Response<Unit>

    @POST("waypoints/")
    suspend fun createWaypoint(
        @Body waypointData: CreateWaypointRequest,
    ): CreateWaypointResponse

    @GET("waypoints/{id}")
    suspend fun getWaypoint(@Path("id") waypointId: Int): Waypoint

    @PUT("waypoints/{id}")
    suspend fun updateWaypoint(
        @Path("id") waypointId: Int,
        @Body waypointData: UpdateWaypointRequest
    ): Waypoint

    @GET("delivery/{delivery_id}")
    suspend fun getDeliveryDetails(
        @Path("delivery_id") deliveryId: Int,
    ): Delivery

    @DELETE("delivery/{delivery_id}")
    suspend fun deleteDelivery(
        @Path("delivery_id") deliveryId: Int,
    ): Response<Unit>

    @DELETE("products/{product_id}")
    suspend fun deleteProduct(
        @Path("product_id") productId: Int,
    ): Response<Unit>

    @GET("analytics/{delivery_id}/optimal-route")
    suspend fun getOptimalRoute(
        @Path("delivery_id") deliveryId: Int,
    ): OptimalRouteResponse

    @GET("analytics/{delivery_id}/optimal-back-route")
    suspend fun getOptimalBackRoute(
        @Path("delivery_id") deliveryId: Int,
    ): OptimalRouteResponse

    @POST("company/{company_id}/add-user")
    suspend fun addUserToCompany(
        @Path("company_id") companyId: Int,
        @Body request: AddUserToCompanyRequest,
    ): Response<Void>

    @GET("users")
    suspend fun getAllUsers(): List<User>

    @POST("routes/")
    suspend fun createRoute(
        @Body routeRequest: CreateRouteRequest,
    ): Response<CreateRouteResponse>

    @POST("admin/change-role")
    suspend fun changeUserRole(
        @Body request: ChangeRoleRequest
    ): Response<Unit>

    @GET("admin/health")
    suspend fun getSystemHealth(): Response<SystemHealth>

    @GET("admin/system-configs")
    suspend fun getSystemConfigs(): Response<SystemConfigs>

    @PUT("admin/system-configs")
    suspend fun updateSystemConfigs(
        @Body request: UpdateSystemConfigsRequest
    ): Response<Void>

    @GET("admin/logs")
    suspend fun getSystemLogs(): Response<List<SystemLog>>

    @POST("admin/clear-logs")
    suspend fun clearSystemLogs(
        @Body request: ClearLogsRequest
    ): Response<Void>

    @GET("/admin/db-status")
    suspend fun getDbStatus(): Response<DbStatusResponse>

    @POST("/admin/backup")
    suspend fun performBackup(
        @Body request: BackupRestoreRequest
    ): Response<AdminActionResponse>

    @POST("/admin/restore")
    suspend fun performRestore(
        @Body request: BackupRestoreRequest
    ): Response<AdminActionResponse>

    @POST("/admin/optimize")
    suspend fun optimizeDb(): Response<AdminActionResponse>
}