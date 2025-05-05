import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const EditCompany = ({ user, t }) => {
  const { company_id } = useParams();
  const [name, setName] = useState("");
  const [address, setAddress] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  
  const token = localStorage.getItem("token")
  useEffect(() => {
    const fetchCompany = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/company/${company_id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        });

        setName(res.data.name);
        setAddress(res.data.address);
      } catch (err) {
        setError("Error while loading company data");
        console.error(err);
      }
    };

    fetchCompany();
  }, [company_id, token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);


    try {
      await axios.put(
        `http://localhost:8081/company/${company_id}`,
        { name, address },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        }
      );

      navigate(`/company/${company_id}`);
    } catch (err) {
      setError("Error while updating company");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
      <h2 className="text-2xl font-bold text-center mb-4">{t("edit_company")}</h2>

      {error && <p className="text-red-600 text-center">{error}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium">{t("company_name")}</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("address")}</label>
          <input
            type="text"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
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

export default EditCompany;
