import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const RouteDetails = ({ user }) => {
  const { route_id } = useParams();
  const [route, setRoute] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const token = localStorage.getItem("token")
  
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
        setError("Error while loading route data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchRoute();
  }, [route_id]);


  const handleDeleteRoute = async () => {
    if (!window.confirm("Confirm deletion?")) return;

    try {
      await axios.delete(`http://localhost:8081/routes/${route_id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      navigate(`/company/${route.CompanyID}`);
    } catch (err) {
      console.error("Error while deleting route:", err);
      alert("Error while deleting route");
    }
  };

  const handleDeleteWaypoint = async (waypointId) => {
    if (!window.confirm("Confirm deletion?")) return;

    try {
      await axios.delete(`http://localhost:8081/waypoints/${waypointId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setRoute((prevRoute) => ({
        ...prevRoute,
        waypoints: prevRoute.waypoints.filter((wp) => wp.ID !== waypointId),
      }));
    } catch (err) {
      console.error("Error while deleting waypoint:", err);
      alert("Error while deleting waypoint");
    }
  };

  if (loading) return <div className="p-6 text-center">Loading...</div>;
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
        {route.company.CreatorID == user.id && (
          <div className="flex space-x-4 mt-4">
            <button
              onClick={handleDeleteRoute}
              className="px-4 py-2 bg-red-500 text-white rounded"
            >
              Delete
            </button>
            <button
              onClick={() => navigate(`/route/${route_id}/edit`)}
              className="px-4 py-2 bg-yellow-500 text-white rounded"
            >
              Edit
            </button>
          </div>
        )}
      </div>

      <div className="bg-white p-4 shadow rounded">
        <h2 className="text-xl font-bold mb-4">Waypoints</h2>
        {route.company.CreatorID == user?.id && (
          <button
            onClick={() => navigate(`/route/${route_id}/waypoint/create`)}
            className="px-4 py-2 bg-green-500 text-white rounded mt-4"
          >
            Create waypoint
          </button>
        )}

        <table className="w-full table-auto border">
          <thead>
            <tr className="bg-gray-100">
              <th className="px-4 py-2 border">Name</th>
              <th className="px-4 py-2 border">Latitude</th>
              <th className="px-4 py-2 border">Longitude</th>
              <th className="px-4 py-2 border">Status</th>
              <th className="px-4 py-2 border">Details</th>
              {route.company.CreatorID == user.id && <th className="px-4 py-2 border">Actions</th>}
            </tr>
          </thead>
          <tbody>
            {route.waypoints.map((wp) => (
              <tr key={wp.ID} className="text-center">
                <td className="px-4 py-2 border">{wp.Name}</td>
                <td className="px-4 py-2 border">{wp.Latitude}</td>
                <td className="px-4 py-2 border">{wp.Longitude}</td>
                <td className="px-4 py-2 border">{wp.Status}</td>
                <td className="px-4 py-2 border">{wp.Details}</td>
                {route.company.CreatorID == user.id && (
                  <td className="px-4 py-2 border">
                    <button
                      onClick={() => navigate(`/waypoint/${wp.ID}/edit`)}
                      className="px-3 py-1 bg-yellow-500 text-white rounded mr-2"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteWaypoint(wp.ID)}
                      className="px-3 py-1 bg-red-500 text-white rounded"
                    >
                      Delete
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
};

export default RouteDetails;

