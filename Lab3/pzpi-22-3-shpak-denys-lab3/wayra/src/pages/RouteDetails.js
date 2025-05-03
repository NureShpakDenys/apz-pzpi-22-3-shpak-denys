import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";

const RouteDetails = () => {
  const { route_id } = useParams();
  const [route, setRoute] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const token =
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyIiwiZXhwIjoxNzQ2Mzc5ODQ4LCJpYXQiOjE3NDYyOTM0NDgsImp0aSI6IjIiLCJ1c2VybmFtZSI6IiJ9.L9IusAwGCgQawW_Wr-UT0NDi9W9i7A61nsS9wCL7qP8";

  useEffect(() => {
    const fetchRoute = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/routes/${route_id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        });
        setRoute(res.data);
      } catch (err) {
        setError("Ошибка при загрузке маршрута");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchRoute();
  }, [route_id]);

  if (loading) return <div className="p-6 text-center">Загрузка...</div>;
  if (error) return <div className="p-6 text-center text-red-600">{error}</div>;
  if (!route) return null;

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="bg-white p-4 shadow rounded mb-4">
        <h2 className="text-2xl font-bold mb-2">Route Data</h2>
        <p><strong>Name:</strong> {route.Name}</p>
        <p><strong>Status:</strong> {route.Status}</p>
        <p><strong>Details:</strong> {route.Details}</p>
        <p><strong>Company:</strong> {route.company?.Name}</p>
      </div>

      <div className="bg-white p-4 shadow rounded">
        <h2 className="text-xl font-bold mb-4">Waypoints</h2>
        <table className="w-full table-auto border">
          <thead>
            <tr className="bg-gray-100">
              <th className="px-4 py-2 border">Name</th>
              <th className="px-4 py-2 border">Latitude</th>
              <th className="px-4 py-2 border">Longitude</th>
              <th className="px-4 py-2 border">Status</th>
              <th className="px-4 py-2 border">Details</th>
            </tr>
          </thead>
          <tbody>
            {route.waypoints.map((wp) => (
              <tr key={wp.ID}>
                <td className="px-4 py-2 border">{wp.Name}</td>
                <td className="px-4 py-2 border">{wp.Latitude}</td>
                <td className="px-4 py-2 border">{wp.Longitude}</td>
                <td className="px-4 py-2 border">{wp.Status}</td>
                <td className="px-4 py-2 border">{wp.Details}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default RouteDetails;
