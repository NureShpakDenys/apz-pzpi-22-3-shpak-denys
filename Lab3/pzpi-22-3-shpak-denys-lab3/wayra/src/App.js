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

import { getUserFromToken } from "./utils/auth";

function App() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedUser = getUserFromToken();
    setUser(storedUser);
  }, []);

  return (
    <Router>
      <Header user={user} setUser={setUser} />
      <div className="container mx-auto p-4">
        <Routes>
          { /*<Route path="/" element={<Home />} />

        
          <Route path="/user/:id" element={<UserProfile />} />

          */}
          <Route path="/login" element={user ? <Navigate to="/companies" /> : <Login setUser={setUser} />} />
          <Route path="/register" element={user ? <Navigate to="/companies" /> : <Register />} />

          <Route path="/companies" element={<Companies />} />
          <Route path="/company/:company_id" element={!user ? <Navigate to="/login" /> : <CompanyDetails user={user} />} />
          <Route path="/company/create" element={!user ? <Navigate to="/login" /> : <CreateCompany />} />
          <Route path="/company/:company_id/edit" element={!user ? <Navigate to="/login" /> : <EditCompany user={user} />} />

          <Route path="/route/:route_id" element={!user ? <Navigate to="/login" /> : <RouteDetails user={user} />} />
          <Route path="/company/:company_id/route/create" element={!user ? <Navigate to="/login" /> : <CreateRoute />} />
          <Route path="/route/:route_id/edit" element={!user ? <Navigate to="/login" /> : <EditRoute user={user} />} />

          <Route path="/route/:route_id/waypoint/create" element={!user ? <Navigate to="/login" /> : <CreateWaypoint />} />
          <Route path="/waypoint/:waypoint_id/edit" element={!user ? <Navigate to="/login" /> : <EditWaypoint user={user} />} />

          <Route path="/delivery/:delivery_id" element={!user ? <Navigate to="/login" /> : <DeliveryDetails user={user} />} />
          <Route path="/company/:company_id/delivery/create" element={!user ? <Navigate to="/login" /> : <CreateDelivery />} />
          <Route path="/delivery/:delivery_id/edit" element={!user ? <Navigate to="/login" /> : <EditDelivery user={user} />} />

          {/*
          <Route path="/products" element={<Products />} />
          <Route path="/product/create" element={<CreateProduct />} />

          <Route path="/admin" element={<AdminDashboard />} />*/}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
