import { BrowserRouter as Router, Route, Routes, Navigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Header from "./pages/Header";
import CompanyDetails from "./pages/CompanyDetails";
import { getUserFromToken } from "./utils/auth";
import RouteDetails from "./pages/RouteDetails";

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
          <Route path="/login" element={user ? <Navigate to="/" /> : <Login setUser={setUser} />} />
          <Route path="/register" element={user ? <Navigate to="/" /> : <Register />} />
          <Route path="/user/:id" element={<UserProfile />} />

          <Route path="/companies" element={<Companies />} /> */}
          <Route path="/company/:company_id" element={<CompanyDetails />} />
          <Route path="/route/:route_id" element={<RouteDetails />} />

          {/*<Route path="/company/create" element={<CreateCompany />} />

          <Route path="/deliveries" element={<Deliveries />} />
          <Route path="/delivery/:delivery_id" element={<DeliveryDetails />} />
          <Route path="/delivery/create" element={<CreateDelivery />} />

          <Route path="/products" element={<Products />} />
          <Route path="/product/create" element={<CreateProduct />} />

          <Route path="/route/:route_id" element={<RouteDetails />} />
          <Route path="/route/create" element={<CreateRoute />} />

          <Route path="/waypoint/create" element={<CreateWaypoint />} />

          <Route path="/admin" element={<AdminDashboard />} />*/}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
