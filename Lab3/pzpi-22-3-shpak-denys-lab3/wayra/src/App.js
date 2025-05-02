import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Header from "./pages/Header";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import UserProfile from "./pages/UserProfile";
import Companies from "./pages/Companies";
import CompanyDetails from "./pages/CompanyDetails";
import CreateCompany from "./pages/CreateCompany";
import Deliveries from "./pages/Deliveries";
import DeliveryDetails from "./pages/DeliveryDetails";
import CreateDelivery from "./pages/CreateDelivery";
import Products from "./pages/Products";
import ProductDetails from "./pages/ProductDetails";
import CreateProduct from "./pages/CreateProduct";
import RoutesPage from "./pages/RoutesPage";
import RouteDetails from "./pages/RouteDetails";
import CreateRoute from "./pages/CreateRoute";
import SensorDataPage from "./pages/SensorDataPage";
import WaypointsPage from "./pages/WaypointsPage";
import WaypointDetails from "./pages/WaypointDetails";
import CreateWaypoint from "./pages/CreateWaypoint";
import WeatherAlerts from "./pages/WeatherAlerts";
import Analytics from "./pages/Analytics";
import AdminDashboard from "./pages/AdminDashboard";
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
          <Route path="/" element={<Home />} />
          <Route path="/login" element={user ? <Navigate to="/" /> : <Login setUser={setUser} />} />
          <Route path="/register" element={user ? <Navigate to="/" /> : <Register />} />
          <Route path="/user/:id" element={<UserProfile />} />

          <Route path="/companies" element={<Companies />} />
          <Route path="/company/:company_id" element={<CompanyDetails />} />
          <Route path="/company/create" element={<CreateCompany />} />

          <Route path="/deliveries" element={<Deliveries />} />
          <Route path="/delivery/:delivery_id" element={<DeliveryDetails />} />
          <Route path="/delivery/create" element={<CreateDelivery />} />

          <Route path="/products" element={<Products />} />
          <Route path="/product/:product_id" element={<ProductDetails />} />
          <Route path="/product/create" element={<CreateProduct />} />

          <Route path="/routes" element={<RoutesPage />} />
          <Route path="/route/:route_id" element={<RouteDetails />} />
          <Route path="/route/create" element={<CreateRoute />} />

          <Route path="/sensor-data/:sensor_data_id" element={<SensorDataPage />} />

          <Route path="/waypoints" element={<WaypointsPage />} />
          <Route path="/waypoint/:waypoint_id" element={<WaypointDetails />} />
          <Route path="/waypoint/create" element={<CreateWaypoint />} />

          <Route path="/weather-alerts/:route_id" element={<WeatherAlerts />} />

          <Route path="/analytics/:delivery_id/optimal-route" element={<Analytics />} />
          <Route path="/analytics/:delivery_id/optimal-back-route" element={<Analytics />} />

          <Route path="/admin" element={<AdminDashboard />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
