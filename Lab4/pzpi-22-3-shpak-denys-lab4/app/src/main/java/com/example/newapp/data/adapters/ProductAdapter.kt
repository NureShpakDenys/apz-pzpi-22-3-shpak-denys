package com.example.newapp.data.adapters

import android.annotation.SuppressLint
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.example.newapp.R
import com.example.newapp.data.models.Product
import com.example.newapp.databinding.ItemProductBinding

class ProductAdapter(
    private var products: List<Product>,
    private val onDeleteClick: (Product) -> Unit,
    private val onEditClick: (Product) -> Unit,
    private var showAdminActionsProvider: () -> Boolean
) : RecyclerView.Adapter<ProductAdapter.ProductViewHolder>() {

    inner class ProductViewHolder(val binding: ItemProductBinding) : RecyclerView.ViewHolder(binding.root) {
        fun bind(product: Product, showAdminActions: Boolean) {

            val context = itemView.context
            val locale = context.resources.configuration.locales[0]
            val isUkrainian = locale.language == "uk"

            binding.tvProductName.text = product.name
            val weight = convertWeight(product.weight.toDouble(), isUkrainian)
            val weightUnit = context.getString(R.string.unit_weight)
            binding.tvProductWeight.text = context.getString(R.string.product_weight_format, weight, weightUnit)
            binding.tvProductCategory.text = product.productCategory.name

            val minTemp = convertTemperature(product.productCategory.minTemperature.toDouble(), isUkrainian)
            val maxTemp = convertTemperature(product.productCategory.maxTemperature.toDouble(), isUkrainian)
            val tempUnit = context.getString(R.string.unit_temperature)
            binding.tvProductTemperature.text = context.getString(
                R.string.temperature_range_format,
                minTemp, maxTemp, tempUnit
            )

            val humidityUnit = context.getString(R.string.unit_humidity)
            binding.tvProductHumidity.text = context.getString(
                R.string.humidity_range_format,
                product.productCategory.minHumidity,
                product.productCategory.maxHumidity,
                humidityUnit
            )

            binding.tvProductPerishable.text = itemView.context.getString(R.string.perishable) + " " + if (product.productCategory.isPerishable)
                itemView.context.getString(R.string.yes) else itemView.context.getString(R.string.no)

            if (showAdminActions) {
                binding.btnEditProduct.visibility = View.VISIBLE
                binding.btnDeleteProduct.visibility = View.VISIBLE
                binding.btnEditProduct.setOnClickListener { onEditClick(product) }
                binding.btnDeleteProduct.setOnClickListener { onDeleteClick(product) }
            } else {
                binding.btnEditProduct.visibility = View.GONE
                binding.btnDeleteProduct.visibility = View.GONE
            }
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ProductViewHolder {
        val binding = ItemProductBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return ProductViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ProductViewHolder, position: Int) {
        holder.bind(products[position], showAdminActionsProvider())
    }

    override fun getItemCount(): Int = products.size

    fun updateProducts(newProducts: List<Product>, showAdmin: Boolean) {
        this.products = newProducts
        notifyDataSetChanged()
    }



    private fun convertWeight(weightKg: Double, isUkrainian: Boolean): Double {
        return if (isUkrainian) weightKg else weightKg * 2.20462
    }

    private fun convertTemperature(tempC: Double, isUkrainian: Boolean): Double {
        return if (isUkrainian) tempC else tempC * 9 / 5 + 32
    }
}