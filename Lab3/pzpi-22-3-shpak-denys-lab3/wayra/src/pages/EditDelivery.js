import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const EditDelivery = ({ user, t }) => {
  const { delivery_id } = useParams();
  const [date, setDate] = useState("");
  const [duration, setDuration] = useState("");
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const token = localStorage.getItem("token")
  
  useEffect(() => {
    const fetchDelivery = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/delivery/${delivery_id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        });

        setDate(res.data.Date.split("T")[0]);
        setDuration(res.data.Duration);
        setStatus(res.data.Status);
      } catch (err) {
        setError("Error while loading delivery data");
        console.error(err);
      }
    };

    fetchDelivery();
  }, [delivery_id, token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await axios.put(
        `http://localhost:8081/delivery/${delivery_id}`,
        { date, duration, status },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        }
      );

      navigate(`/delivery/${delivery_id}`);
    } catch (err) {
      setError("Error while updating delivery");
      console.error(err);
  } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
      <h2 className="text-2xl font-bold text-center mb-4">{t("edit_delivery")}</h2>

      {error && <p className="text-red-600 text-center">{error}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium">{t("edit_delivery")}</label>
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("duration")}</label>
          <input
            type="text"
            value={duration}
            onChange={(e) => setDuration(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Eg. 2 hours"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("status")}</label>
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            <option value="not_started">{t("not_started")}</option>
            <option value="in_progress">{t("in_progress")}</option>
            <option value="completed">{t("completed")}</option>
          </select>
        </div>

        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
          disabled={loading}
        >
          {loading ? t("adding") : t("save")}
        </button>
      </form>
    </div>
  );
};

export default EditDelivery;
