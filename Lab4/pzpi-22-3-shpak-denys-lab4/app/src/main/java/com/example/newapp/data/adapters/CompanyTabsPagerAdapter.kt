package com.example.newapp.data.adapters

import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentActivity
import androidx.viewpager2.adapter.FragmentStateAdapter
import com.example.newapp.data.fragments.DeliveriesFragment
import com.example.newapp.data.fragments.RoutesFragment
import com.example.newapp.data.fragments.UsersFragment
import com.example.newapp.data.models.CompanyResponse
import com.example.newapp.data.models.CompanyUser
import com.example.newapp.data.SessionManager

class CompanyTabsPagerAdapter(
    private val activity: FragmentActivity,
    private val company: CompanyResponse,
    private val companyUsers: List<CompanyUser>
) : FragmentStateAdapter(activity) {

    private val sessionManager = SessionManager(activity)
    private val isCreator = company.creator.id == sessionManager.getUser()?.id

    override fun getItemCount() = 3

    override fun createFragment(position: Int): Fragment {
        return when (position) {
            0 -> RoutesFragment.newInstance(
                routes = company.routes,
                isCreator = isCreator,
                companyId = company.id
            )
            1 -> DeliveriesFragment.newInstance(
                deliveries = company.deliveries,
                routes = company.routes,
                isCreator = isCreator,
                companyId = company.id
            )
            2 -> UsersFragment.newInstance(
                users = companyUsers,
                creatorId = company.creator.id,
                companyId = company.id,
                isCreator = isCreator
            )
            else -> Fragment()
        }
    }
}
