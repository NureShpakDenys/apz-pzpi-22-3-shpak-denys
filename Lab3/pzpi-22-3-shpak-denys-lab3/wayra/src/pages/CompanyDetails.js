import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";

const CompanyDetails = () => {
  const { company_id } = useParams();
  const [data, setData] = useState(null);
  const [activeTab, setActiveTab] = useState("routes");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyIiwiZXhwIjoxNzQ2Mzc5ODQ4LCJpYXQiOjE3NDYyOTM0NDgsImp0aSI6IjIiLCJ1c2VybmFtZSI6IiJ9.L9IusAwGCgQawW_Wr-UT0NDi9W9i7A61nsS9wCL7qP8';

  useEffect(() => {
    const fetchCompany = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/company/${company_id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: 'application/json',
          },
        });
        setData(res.data);
      } catch (err) {
        setError("Ошибка при загрузке данных");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCompany();
  }, [company_id]);

  if (loading) return <div className="p-6 text-center">Загрузка...</div>;
  if (error) return <div className="p-6 text-center text-red-600">{error}</div>;
  if (!data) return null;

  const renderTable = () => {
    switch (activeTab) {
      case "routes":
        return (
          <table className="w-full table-auto">
            <thead>
              <tr><th>Name</th><th>Status</th><th>Details</th></tr>
            </thead>
            <tbody>
              {data.routes.map((route) => (
                <tr key={route.id}>
                  <td>{route.name}</td>
                  <td>{route.status}</td>
                  <td>{route.details}</td>
                </tr>
              ))}
            </tbody>
          </table>
        );
      case "deliveries":
        return (
          <table className="w-full table-auto">
            <thead>
              <tr><th>Status</th><th>Date</th><th>Duration</th><th>Route ID</th></tr>
            </thead>
            <tbody>
              {data.deliveries.map((d) => (
                <tr key={d.id}>
                  <td>{d.status}</td>
                  <td>{new Date(d.date).toLocaleDateString()}</td>
                  <td>{d.duration}</td>
                  <td>{d.route_id}</td>
                </tr>
              ))}
            </tbody>
          </table>
        );
      case "users":
        return (
          <table className="w-full table-auto">
            <thead>
              <tr><th>ID</th><th>Name</th></tr>
            </thead>
            <tbody>
              {data.users.map((u) => (
                <tr key={u.id}>
                  <td>{u.id}</td>
                  <td>{u.name}</td>
                </tr>
              ))}
            </tbody>
          </table>
        );
      default:
        return null;
    }
  };

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="text-center mb-4">
        <h1 className="text-3xl font-bold">{data.name}</h1>
        <p className="text-gray-600">{data.address}</p>
      </div>

      <div className="flex space-x-4 justify-center mb-4">
        <button onClick={() => setActiveTab("routes")} className="px-4 py-2 bg-blue-200 rounded">Routes</button>
        <button onClick={() => setActiveTab("deliveries")} className="px-4 py-2 bg-blue-200 rounded">Deliveries</button>
        <button onClick={() => setActiveTab("users")} className="px-4 py-2 bg-blue-200 rounded">Users</button>
      </div>

      <div className="overflow-x-auto bg-white shadow-md rounded p-4">
        {renderTable()}
      </div>
    </div>
  );
};

export default CompanyDetails;
