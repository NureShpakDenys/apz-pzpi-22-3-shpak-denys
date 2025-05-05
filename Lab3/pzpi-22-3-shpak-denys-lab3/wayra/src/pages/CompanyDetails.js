import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const CompanyDetails = ({ user, t }) => {
  const { company_id } = useParams();
  const [data, setData] = useState(null);
  const [users, setUsers] = useState([]);
  const [activeTab, setActiveTab] = useState("routes");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const token = localStorage.getItem("token");
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
        setError("Error while loading company data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCompany();

    const fetchCompanyUsers = async () => {
      try {
        const res = await axios.get(`http://localhost:8081/company/${company_id}/users`, {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json" },
        });

        setUsers(res.data);
      } catch (err) {
        setError("Error while loading users data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchCompanyUsers();
  }, [company_id]);

  const handleDeleteCompany = async () => {
    if (!window.confirm("Confirm deletion?")) return;

    try {
      await axios.delete(`http://localhost:8081/company/${company_id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      navigate("/companies");
    } catch (err) {
      console.error("Error while deleting company:", err);
      alert("Error while deleting company");
    }
  };

  const updateUserRole = async (userId, newRole) => {
    try {
      await axios.put(
        `http://localhost:8081/company/${company_id}/update-user`,
        { userID: userId, role: newRole },
        {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json", "Content-Type": "application/json" },
        }
      );

      setUsers((prevUsers) =>
        prevUsers.map((user) => (user.UserID === userId ? { ...user, Role: newRole } : user))
      );
    } catch (err) {
      console.error("Error updating role:", err);
    }
  };

  const removeUserFromCompany = async (userId) => {
    if (!window.confirm("Are you sure?")) return;

    try {
      await axios.delete(`http://localhost:8081/company/${company_id}/remove-user`, {
        headers: { Authorization: `Bearer ${token}`, Accept: "application/json", "Content-Type": "application/json" },
        data: { userID: userId },
      });

      setUsers((prevUsers) => prevUsers.filter((user) => user.UserID !== userId));
    } catch (err) {
      console.error("Error removing user:", err);
    }
  };

  if (loading) return <div className="p-6 text-center">{t("loading")}</div>;
  if (error) return <div className="p-6 text-center text-red-600">{error}</div>;
  if (!data) return null;

  const renderTable = () => {
    switch (activeTab) {
      case "routes":
        return (
          <>
            <h2 className="text-xl font-bold mb-2">{t("routes")}</h2>
            {data.creator.id == user.id && (
              <button
                onClick={() => navigate(`/company/${company_id}/route/create`)}
                className="px-4 py-2 bg-green-500 text-white rounded mb-2"
              >
                {t("add_route")}
              </button>
            )}
            <table className="w-full table-auto">
              <thead>
                <tr><th>{t("name")}</th><th>{t("status")}</th><th>{t("details")}</th></tr>
              </thead>
              <tbody>
                {data.routes.map((route) => (
                  <tr key={route.id} className="text-center">
                    <td onClick={() => navigate(`/route/${route.id}`)} className="cursor-pointer hover:bg-gray-300 text-blue-500 underline">{route.name}</td>
                    <td>{route.status}</td>
                    <td>{route.details}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        );
      case "deliveries":
        return (
          <>
            <h2 className="text-xl font-bold mb-2">{t("deliveries")}</h2>
            {data.creator.id == user.id && (
              <button
                onClick={() => navigate(`/company/${company_id}/delivery/create`)}
                className="px-4 py-2 bg-green-500 text-white rounded mb-2"
              >
                {t("add_delivery")}
              </button>
            )}
            <table className="w-full table-auto">
              <thead>
                <tr>
                  <th>{t("id")}</th>
                  <th>{t("status")}</th>
                  <th>{t("date")}</th>
                  <th>{t("duration")}</th>
                  <th>{t("route_id")}</th>
                </tr>
              </thead>
              <tbody>
                {data.deliveries.map((d) => (
                  <tr key={d.id} className="text-center">
                    <td onClick={() => navigate(`/delivery/${d.id}`)} className="cursor-pointer hover:bg-gray-300 text-blue-500 underline">{d.id}</td>
                    <td>{d.status}</td>
                    <td>{new Date(d.date).toLocaleDateString()}</td>
                    <td>{d.duration}</td>
                    <td>{d.route_id}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        );
      case "users":
        return (
          <>
            <h2 className="text-xl font-bold mb-2">{t("users")}</h2>
            {data.creator.id == user.id && (
              <button
                onClick={() => navigate(`/company/${company_id}/add-user`)}
                className="px-4 py-2 bg-green-500 text-white rounded mb-2"
              >
                {t("add_user")}
              </button>

            )}
            <table className="w-full table-auto">
              <thead>
                <tr>
                  <th>{t("id")}</th>
                  <th>{t("name")}</th>
                  <th>{t("role")}</th>
                  {data.creator.id == user.id && <th className="border p-2">{t("actions")}</th>}
                </tr>
              </thead>
              <tbody>
                {users.map((u) => (
                  <tr key={u.id} className="text-center">
                    <td>{u.UserID}</td>
                    <td>{u.user.name}</td>
                    <td className="border p-2">
                      {data.creator.id == user.id ? (
                        <select value={u.Role} onChange={(e) => updateUserRole(u.UserID, e.target.value)} className="border rounded px-2 py-1">
                          <option value="user">User</option>
                          <option value="admin">Admin</option>
                          <option value="manager">Manager</option>
                        </select>
                      ) : (
                        <span>{u.Role}</span>
                      )}
                    </td>
                    {user.id == users[0]?.company.CreatorID && (
                      <td className="border p-2">
                        <button
                          onClick={() => updateUserRole(u.UserID, u.Role)}
                          className="px-3 py-1 bg-blue-500 text-white rounded mr-2"
                        >
                          {t("save")}
                        </button>
                        <button
                          onClick={() => removeUserFromCompany(u.UserID)}
                          className="px-3 py-1 bg-red-500 text-white rounded"
                        >
                          {t("remove")}
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </>
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
        {data.creator.id == user.id && (
          <div className="flex justify-center space-x-4 mt-4">
            <button
              onClick={handleDeleteCompany}
              className="px-4 py-2 bg-red-500 text-white rounded"
            >
              {t("delete")}
            </button>
            <button
              onClick={() => navigate(`/company/${company_id}/edit`)}
              className="px-4 py-2 bg-yellow-500 text-white rounded"
            >
              {t("edit")}
            </button>
          </div>
        )}
      </div>

      <div className="flex space-x-4 justify-center mb-4">
        <button onClick={() => setActiveTab("routes")} className="px-4 py-2 bg-blue-200 rounded">{t("routes")}</button>
        <button onClick={() => setActiveTab("deliveries")} className="px-4 py-2 bg-blue-200 rounded">{t("deliveries")}</button>
        <button onClick={() => setActiveTab("users")} className="px-4 py-2 bg-blue-200 rounded">{t("users")}</button>
      </div>

      <div className="overflow-x-auto bg-white shadow-md rounded p-4">
        {renderTable()}
      </div>
    </div>
  );
};

export default CompanyDetails;
