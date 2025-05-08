import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import convert from "../utils/convertors";

const DeliveryDetails = ({ user, t, i18n }) => {
  const { delivery_id } = useParams();
  const [delivery, setDelivery] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [optimalRoute, setOptimalRoute] = useState(null);
  const [backRoute, setBackRoute] = useState(null);
  const [routesLoading, setRoutesLoading] = useState(false);
  const [routesError, setRoutesError] = useState(null);
  const navigate = useNavigate();
  const system = i18n.language === "uk" ? "metric" : "imperial";

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
        setError("Error loading delivery data.");
      } finally {
        setLoading(false);
      }
    };

    fetchDelivery();
  }, [delivery_id]);

  const fetchOptimalRoutes = async () => {
    setRoutesLoading(true);
    setRoutesError(null);

    try {
      const [optimalRes, backRes] = await Promise.all([
        axios.get(`http://localhost:8081/analytics/${delivery_id}/optimal-route`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        }),
        axios.get(`http://localhost:8081/analytics/${delivery_id}/optimal-back-route`, {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        }),
      ]);
      setOptimalRoute(optimalRes.data);
      setBackRoute(backRes.data);
    } catch (err) {
      console.error("Error fetching routes", err);
      setRoutesError("Помилка при отриманні маршрутів");
    } finally {
      setRoutesLoading(false);
    }
  };

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


  const handleDeleteProduct = async (productId) => {
    if (!window.confirm("Confirm deletion?")) return;

    try {
      await axios.delete(`http://localhost:8081/products/${productId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setDelivery((prevDelivery) => ({
        ...prevDelivery,
        products: prevDelivery.products.filter((product) => product.ID !== productId),
      }));
    } catch (err) {
      console.error("Error while deleting product:", err);
      alert("Error while deleting product");
    }
  };
  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="bg-white p-6 shadow rounded mb-6 text-center">
        <h1 className="text-3xl font-bold">{t("delivery_data")}</h1>
        <p className="mt-2 text-gray-600">
          <strong>{t("status")}:</strong> {delivery.Status}
        </p>
        <p className="text-gray-600">
          <strong>{t("date")}:</strong> {new Date(delivery.Date).toLocaleDateString()}
        </p>
        <p className="text-gray-600">
          <strong>{t("duration")}:</strong> {delivery.Duration}
        </p>
        {delivery.company.CreatorID == user.id && (
          <div className="flex justify-center space-x-4 mt-4">
            <button
              onClick={handleDeleteDelivery}
              className="px-4 py-2 bg-red-500 text-white rounded"
            >
              {t("delete")}
            </button>
            <button
              onClick={() => navigate(`/delivery/${delivery_id}/edit`)}
              className="px-4 py-2 bg-yellow-500 text-white rounded"
            >
              {t("edit")}
            </button>
          </div>
        )}
      </div>

      <div className="bg-white p-4 shadow rounded mt-6 text-center">
        <button
          onClick={fetchOptimalRoutes}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          {t("get_optimal_routes")}
        </button>

        {routesLoading && <p className="mt-2 text-gray-500">{t("loading")}</p>}
        {routesError && <p className="mt-2 text-red-600">{routesError}</p>}
      </div>

      {optimalRoute && (
        <div className="bg-white p-4 shadow rounded mt-4">
          <h2 className="text-xl font-semibold mb-2 text-green-700">{t("optimal_route")}</h2>
          <p><strong>{t("route")}:</strong> {optimalRoute.route.name}</p>
          <p><strong>{t("message")}:</strong> {optimalRoute.message}</p>
          <p><strong>{t("equation")}:</strong> {optimalRoute.equation}</p>
          <p><strong>{t("distance")}:</strong> {
            system === "imperial"
              ? `${convert(optimalRoute.predict_data.Distance, "distance", "imperial").toFixed(2)}`
              : `${optimalRoute.predict_data.Distance}`
          } {t("km")}</p>
          <p><strong>{t("time")}:</strong> {optimalRoute.predict_data.Time.toFixed(2)} {t("hours")}</p>
          <p><strong>{t("speed")}:</strong> {
            system === "imperial"
              ? `${convert(optimalRoute.predict_data.Speed, "speed", "imperial").toFixed(2)}`
              : `${optimalRoute.predict_data.Speed.toFixed(2)}`
          } {t("km/hour")}</p>
        </div>
      )}

      {backRoute && (
        <div className="bg-white p-4 shadow rounded mt-4">
          <h2 className="text-xl font-semibold mb-2 text-blue-700">{t("back_route")}</h2>
          <p><strong>{t("route")}:</strong> {backRoute.route.name}</p>
          <p><strong>{t("message")}:</strong> {backRoute.message}</p>
          <p><strong>{t("equation")}:</strong> {backRoute.equation}</p>
          <p><strong>{t("distance")}:</strong> {
            system === "imperial"
              ? `${convert(backRoute.predict_data.Distance, "distance", "imperial").toFixed(2)}`
              : `${backRoute.predict_data.Distance}`
          } {t("km")}</p>
          <p><strong>{t("time")}:</strong> {backRoute.predict_data.Time.toFixed(2)} {t("hours")}</p>
          <p><strong>{t("speed")}:</strong> {
            system === "imperial"
              ? `${convert(backRoute.predict_data.Speed, "speed", "imperial").toFixed(2)}`
              : `${backRoute.predict_data.Speed.toFixed(2)}`
        } {t("km/hour")}</p>
        </div>
      )}

      <div className="bg-white p-4 shadow rounded">
        <h2 className="text-xl font-semibold mb-4">{t("products")}</h2>
        {delivery.company.CreatorID == user.id && (
          <button
            onClick={() => navigate(`/delivery/${delivery_id}/product/add`)}
            className="px-4 py-2 bg-green-500 text-white rounded mb-4"
          >
            {t("add_product")}
          </button>
        )}
        <table className="w-full border table-auto">
          <thead className="bg-gray-100">
            <tr>
              <th className="border px-4 py-2">{t("name")}</th>
              <th className="border px-4 py-2">{t("weight")}</th>
              <th className="border px-4 py-2">{t("category")}</th>
              <th className="border px-4 py-2">{t("temperature")} ({t("°C")})</th>
              <th className="border px-4 py-2">{t("humidity")} (%)</th>
              <th className="border px-4 py-2">{t("perishable")}</th>
              {delivery.company.CreatorID == user.id && <th className="px-4 py-2 border">{t("actions")}</th>}
            </tr>
          </thead>
          <tbody>
            {delivery.products.map((product) => (
              <tr key={product.ID} className="text-center">
                <td className="border px-4 py-2">{product.Name}</td>
                <td className="border px-4 py-2">
                  {system === "imperial"
                    ? `${convert(product.Weight, "weight", "imperial").toFixed(1)}`
                    : `${product.Weight.toFixed(1)}`}
                </td>
                <td className="border px-4 py-2">{t(`product_type.${product.product_category.Name}`)}</td>
                <td className="border px-4 py-2">
                  {system === "imperial"
                    ? `${convert(product.product_category.MinTemperature, "temperature", "imperial").toFixed(1)}    `
                    : `${product.product_category.MinTemperature}   `
                  }
                   - 
                  {system === "imperial"
                    ? `   ${convert(product.product_category.MaxTemperature, "temperature", "imperial").toFixed(1)}`
                    : `   ${product.product_category.MaxTemperature}`
                  }
                </td>
                <td className="border px-4 py-2">
                  {product.product_category.MinHumidity}–{product.product_category.MaxHumidity}
                </td>
                <td className="border px-4 py-2">
                  {product.product_category.IsPerishable ? t("yes") : t("no")}
                </td>
                {delivery.company.CreatorID == user.id && (
                  <td className="px-4 py-2 border">
                    <button
                      onClick={() => navigate(`/product/${product.ID}/edit`)}
                      className="px-3 py-1 bg-yellow-500 text-white rounded mr-2"
                    >
                      {t("edit")}
                    </button>
                    <button
                      onClick={() => handleDeleteProduct(product.ID)}
                      className="px-3 py-1 bg-red-500 text-white rounded"
                    >
                      {t("delete")}
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DeliveryDetails;
