import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const DeliveryDetails = ({ user }) => {
  const { delivery_id } = useParams();
  const [delivery, setDelivery] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
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
        setDelivery(res.data);
      } catch (err) {
        console.error(err);
        setError("Ошибка при загрузке данных доставки.");
      } finally {
        setLoading(false);
      }
    };

    fetchDelivery();
  }, [delivery_id]);

  if (loading) return <div className="p-6 text-center">Загрузка...</div>;
  if (error) return <div className="p-6 text-red-600">{error}</div>;
  if (!delivery) return null;

  const handleDeleteDelivery = async () => {
    if (!window.confirm("Confirm deletion?")) return;

    try {
      await axios.delete(`http://localhost:8081/delivery/${delivery_id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      navigate(`/company/${delivery.CompanyID}`);
    } catch (err) {
      console.error("Error while deleting delivery:", err);
      alert("Error while deleting company");
    }
  };


  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="bg-white p-6 shadow rounded mb-6 text-center">
        <h1 className="text-3xl font-bold">Delivery data</h1>
        <p className="mt-2 text-gray-600">
          <strong>Status:</strong> {delivery.Status}
        </p>
        <p className="text-gray-600">
          <strong>Date:</strong> {new Date(delivery.Date).toLocaleDateString()}
        </p>
        <p className="text-gray-600">
          <strong>Duration:</strong> {delivery.Duration}
        </p>
        {delivery.company.CreatorID == user.id && (
          <div className="flex justify-center space-x-4 mt-4">
            <button
              onClick={handleDeleteDelivery}
              className="px-4 py-2 bg-red-500 text-white rounded"
            >
              Delete
            </button>
            <button
              onClick={() => navigate(`/delivery/${delivery_id}/edit`)}
              className="px-4 py-2 bg-yellow-500 text-white rounded"
            >
              Edit
            </button>
          </div>
        )}
      </div>

      <div className="bg-white p-4 shadow rounded">
        <h2 className="text-xl font-semibold mb-4">Products</h2>
        <table className="w-full border table-auto">
          <thead className="bg-gray-100">
            <tr>
              <th className="border px-4 py-2">Name</th>
              <th className="border px-4 py-2">Weight</th>
              <th className="border px-4 py-2">Category</th>
              <th className="border px-4 py-2">Temperature (°C)</th>
              <th className="border px-4 py-2">Humidity (%)</th>
              <th className="border px-4 py-2">Perishable</th>
            </tr>
          </thead>
          <tbody>
            {delivery.products.map((product) => (
              <tr key={product.ID} className="text-center">
                <td className="border px-4 py-2">{product.Name}</td>
                <td className="border px-4 py-2">{product.Weight} kg</td>
                <td className="border px-4 py-2">{product.product_category.Name}</td>
                <td className="border px-4 py-2">
                  {product.product_category.MinTemperature}–{product.product_category.MaxTemperature}
                </td>
                <td className="border px-4 py-2">
                  {product.product_category.MinHumidity}–{product.product_category.MaxHumidity}
                </td>
                <td className="border px-4 py-2">
                  {product.product_category.IsPerishable ? "Yes" : "No"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DeliveryDetails;
