import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const EditWaypoint = ({ user }) => {
  const { waypoint_id } = useParams();
  const [name, setName] = useState("");
  const [deviceSerial, setDeviceSerial] = useState("");
  const [sendDataFrequency, setSendDataFrequency] = useState(50);
  const [getWeatherAlerts, setGetWeatherAlerts] = useState(true);
  const [route_id, setRouteID] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
 
  const token = localStorage.getItem("token")
  
  useEffect(() => {
    const fetchWaypoint = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/waypoints/${waypoint_id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        });

        setName(res.data.Name);
        setDeviceSerial(res.data.DeviceSerial);
        setSendDataFrequency(res.data.SendDataFrequency);
        setGetWeatherAlerts(res.data.GetWeatherAlerts);
        setRouteID(res.data.RouteID);
      } catch (err) {
        setError("Error fetching waypoint data");
        console.error(err);
      }
    };

    fetchWaypoint();
  }, [waypoint_id, token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await axios.put(
        `http://localhost:8081/waypoints/${waypoint_id}`,
        { name, device_serial: deviceSerial, send_data_frequency: Number(sendDataFrequency), get_weather_alerts: Boolean(getWeatherAlerts) },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        }
      );

      navigate(`/route/${route_id}`);
    } catch (err) {
      setError("Error updating waypoint data");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
      <h2 className="text-2xl font-bold text-center mb-4">Edit Waypoint</h2>

      {error && <p className="text-red-600 text-center">{error}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium">Waypoint Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">Device Serial</label>
          <input
            type="text"
            value={deviceSerial}
            onChange={(e) => setDeviceSerial(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">Data Sending Frequency (seconds)</label>
          <input
            type="number"
            value={sendDataFrequency}
            onChange={(e) => setSendDataFrequency(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div className="flex items-center space-x-2">
          <input
            type="checkbox"
            checked={getWeatherAlerts}
            onChange={(e) => setGetWeatherAlerts(e.target.checked)}
            className="w-5 h-5"
          />
          <label className="text-gray-700 font-medium">Get Weather Alerts</label>
        </div>

        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
          disabled={loading}
        >
          {loading ? "Saving..." : "Save Changes"}
        </button>
      </form>
    </div>
  );
};

export default EditWaypoint;
