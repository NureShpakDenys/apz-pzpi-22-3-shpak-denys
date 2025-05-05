import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const CreateWaypoint = ({ t }) => {
  const { route_id } = useParams();
  const [name, setName] = useState("");
  const [deviceSerial, setDeviceSerial] = useState("");
  const [sendDataFrequency, setSendDataFrequency] = useState(50);
  const [getWeatherAlerts, setGetWeatherAlerts] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const token = localStorage.getItem("token")
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await axios.post(
        "http://localhost:8081/waypoints/",
        {
          route_id: Number(route_id),
          name,
          device_serial: deviceSerial,
          send_data_frequency: Number(sendDataFrequency),
          get_weather_alerts: Boolean(getWeatherAlerts),
        },
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
      setError("Error while creating waypoint");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
      <h2 className="text-2xl font-bold text-center mb-4">{t("create_waypoint")}</h2>

      {error && <p className="text-red-600 text-center">{error}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium">{t("waypoint_name")}</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("device_serial")}</label>
          <input
            type="text"
            value={deviceSerial}
            onChange={(e) => setDeviceSerial(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("data_sending_frequancy")}</label>
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
          <label className="text-gray-700 font-medium">{t("get_weather_alerts")}</label>
        </div>

        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition"
          disabled={loading}
        >
          {loading ? t("adding") : t("create_waypoint")}
        </button>
      </form>
    </div>
  );
};

export default CreateWaypoint;
