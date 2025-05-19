package com.example.newapp.data.network

import android.content.Context
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.Response
import java.io.IOException
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class RedirectInterceptor : Interceptor {
    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        var response = chain.proceed(request)

        if (response.code() == 307) {
            val location = response.header("Location")
            if (location != null) {
                val newUrl = if (location.startsWith("http://") || location.startsWith("https://")) {
                    location
                } else {
                    request.url().newBuilder()
                        .encodedPath(location)
                        .build()
                        .toString()
                }

                val newRequest = request.newBuilder()
                    .url(newUrl)
                    .build()

                response.close()
                response = chain.proceed(newRequest)
            }
        }
        return response
    }
}

object RetrofitClient {
    private const val BASE_URL = "http://192.168.1.102:8081/"

    fun getInstance(context: Context): ApiService {
        val client = OkHttpClient.Builder()
            .addInterceptor(AuthInterceptor(context))
            .addInterceptor(RedirectInterceptor())
            .followRedirects(false)
            .build()

        return Retrofit.Builder()
            .baseUrl(BASE_URL)
            .addConverterFactory(GsonConverterFactory.create())
            .client(client)
            .build()
            .create(ApiService::class.java)
    }
}
