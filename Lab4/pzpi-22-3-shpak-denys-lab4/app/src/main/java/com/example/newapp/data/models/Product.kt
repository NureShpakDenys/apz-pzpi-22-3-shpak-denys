package com.example.newapp.data.models

import com.google.gson.annotations.SerializedName

data class CreateProductRequest(
    val deliveryID: Int,
    val name: String,
    val productType: String,
    val weight: Float
)

data class Product(
    @SerializedName("ID") val id: Int,
    @SerializedName("Name") val name: String,
    @SerializedName("Weight") val weight: Float,
    @SerializedName("product_category") val productCategory: ProductCategory,
    val deliveryId: Int
)
