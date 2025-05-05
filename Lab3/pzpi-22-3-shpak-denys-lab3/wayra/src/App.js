import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Header from "./pages/Header";

import Register from "./pages/Register";
import Login from "./pages/Login";

import Companies from "./pages/Companies";
import CompanyDetails from "./pages/CompanyDetails";
import CreateCompany from "./pages/CreateCompany";
import EditCompany from "./pages/EditCompany";

import RouteDetails from "./pages/RouteDetails";
import CreateRoute from "./pages/CreateRoute";
import EditRoute from "./pages/EditRoute";

import CreateWaypoint from "./pages/CreateWaypoint";
import EditWaypoint from "./pages/EditWaypoint";

import DeliveryDetails from "./pages/DeliveryDetails";
import CreateDelivery from "./pages/CreateDelivery";
import EditDelivery from "./pages/EditDelivery";

import AddProduct from "./pages/AddProduct";
import EditProduct from "./pages/EditProduct";

import { getUserFromToken } from "./utils/auth";
import AddUserToCompany from "./pages/AddUserToCompany";

import AdminDashboard from "./pages/AdminDashboard";
import SystemAdminDashboard from "./pages/SystemAdminDashboard";
import DBAdminDashboard from "./pages/DBAdminDashboard";

import axios from "axios";
import { useTranslation } from "react-i18next";

function App() {
  const [user, setUser] = useState(null);
  const { i18n, t } = useTranslation();

  useEffect(() => {
    const fetchUser = async () => {
      const storedUser = getUserFromToken();
      const token = localStorage.getItem("token");
  
      if (!storedUser || !storedUser.id) return;
  
      try {
        const usersRes = await axios.get(`http://localhost:8081/user/${storedUser.id}`, {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json" },
        });
  
        setUser(usersRes.data);
      } catch (err) {
        console.error("Error fetching user data:", err);
      }
    };
  
    fetchUser();
  }, []);
  

  return (
    <Router>
      <Header user={user} setUser={setUser} i18n={i18n} t={t} />
      <div className="container mx-auto p-4">
        <Routes>
          <Route path="/login" element={user ? <Navigate to="/companies" /> : <Login setUser={setUser} t={t} />} />
          <Route path="/register" element={user ? <Navigate to="/companies" /> : <Register t={t} />} />

          <Route path="/companies" element={<Companies t={t} />} />
          <Route path="/company/:company_id" element={!user ? <Navigate to="/login" /> : <CompanyDetails user={user} t={t} />} />
          <Route path="/company/create" element={!user ? <Navigate to="/login" /> : <CreateCompany t={t} />} />
          <Route path="/company/:company_id/edit" element={!user ? <Navigate to="/login" /> : <EditCompany user={user} t={t} />} />
          <Route path="/company/:company_id/add-user" element={!user ? <Navigate to="/login" /> : <AddUserToCompany user={user} t={t} />} />

          <Route path="/route/:route_id" element={!user ? <Navigate to="/login" /> : <RouteDetails user={user} t={t} />} />
          <Route path="/company/:company_id/route/create" element={!user ? <Navigate to="/login" /> : <CreateRoute t={t} />} />
          <Route path="/route/:route_id/edit" element={!user ? <Navigate to="/login" /> : <EditRoute user={user} t={t} />} />

          <Route path="/route/:route_id/waypoint/create" element={!user ? <Navigate to="/login" /> : <CreateWaypoint t={t} />} />
          <Route path="/waypoint/:waypoint_id/edit" element={!user ? <Navigate to="/login" /> : <EditWaypoint user={user} t={t} />} />

          <Route path="/delivery/:delivery_id" element={!user ? <Navigate to="/login" /> : <DeliveryDetails user={user} t={t} />} />
          <Route path="/company/:company_id/delivery/create" element={!user ? <Navigate to="/login" /> : <CreateDelivery t={t} />} />
          <Route path="/delivery/:delivery_id/edit" element={!user ? <Navigate to="/login" /> : <EditDelivery user={user} t={t} />} />

          <Route path="delivery/:delivery_id/product/add" element={!user ? <Navigate to="/login" /> : <AddProduct t={t} />} />
          <Route path="/product/:product_id/edit" element={!user ? <Navigate to="/login" /> : <EditProduct t={t} />} />

          <Route path="/admin" element={user ? <AdminDashboard user_id={user.id} t={t} /> : <Navigate to="/login" />} />
          <Route path="/system-admin" element={user ? <SystemAdminDashboard t={t} /> : <Navigate to="/login" />} />
          <Route path="/db-admin" element={user ? <DBAdminDashboard t={t} /> : <Navigate to="/login" />} />

        </Routes>
      </div>
    </Router>
  );
}

export default App;
