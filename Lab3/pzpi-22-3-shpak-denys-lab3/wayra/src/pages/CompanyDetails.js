import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const CompanyDetails = ({ user, t, i18n }) => {
  const { company_id } = useParams();
  const [data, setData] = useState(null);
  const [users, setUsers] = useState([]);
  const [activeTab, setActiveTab] = useState("routes");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [routesSort, setRoutesSort] = useState({ column: 'name', direction: 'asc' });
  const [deliveriesSort, setDeliveriesSort] = useState({ column: 'id', direction: 'asc' });
  const [usersSort, setUsersSort] = useState({ column: 'id', direction: 'asc' });
  const [routes, setRoutes] = useState([]);
  const navigate = useNavigate();
  const lang = i18n.language;

  const sortData = (data, column, direction) => {
    const sorted = [...data].sort((a, b) => {
      if (direction === 'asc') {
        return a[column] > b[column] ? 1 : -1;
      } else {
        return a[column] < b[column] ? 1 : -1;
      }
    });
    return sorted;
  };

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
        setRoutes(res.data.routes);
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

  const updateDeliveryRoute = async (delivery_id, route_id) => {
    try {
      await axios.put(
        `http://localhost:8081/delivery/${delivery_id}`,
        { route_id: Number(route_id) },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        }
      );

      setData((prevData) => ({
        ...prevData,
        deliveries: prevData.deliveries.map((delivery) =>
          delivery.id === delivery_id ? { ...delivery, route_id: route_id } : delivery
        ),
      }));
    } catch (err) {
      console.error("Error updating delivery route:", err);
    }
  };


  if (loading) return <div className="p-6 text-center">{t("loading")}</div>;
  if (error) return <div className="p-6 text-center text-red-600">{error}</div>;
  if (!data) return null;

  const renderTable = () => {
    switch (activeTab) {
      case "routes":
        const sortedRoutes = sortData(data.routes, routesSort.column, routesSort.direction);
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
                <tr className="text-center h-10">
                  <th className="border" onClick={() => setRoutesSort({ column: 'name', direction: routesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('name')}
                    {routesSort.column === 'name' && (
                      <span>{routesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setRoutesSort({ column: 'status', direction: routesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('status')}
                    {routesSort.column === 'status' && (
                      <span>{routesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setRoutesSort({ column: 'details', direction: routesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('details')}
                    {routesSort.column === 'details' && (
                      <span>{routesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                </tr>
              </thead>
              <tbody>
                {sortedRoutes.map((route) => (
                  <tr key={route.id} className="border text-center h-10">
                    <td onClick={() => navigate(`/route/${route.id}`)} className="border cursor-pointer hover:bg-gray-300 text-blue-500 underline">{route.name}</td>
                    <td className="border">{route.status}</td>
                    <td className="border">{route.details}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        );
      case "deliveries":
        const sortedDeliveries = sortData(data.deliveries, deliveriesSort.column, deliveriesSort.direction);
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
            <table className="border w-full table-auto">
              <thead>
                <tr className=" border text-center h-10">
                  <th className="border" onClick={() => setDeliveriesSort({ column: 'id', direction: deliveriesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('id')}
                    {deliveriesSort.column === 'id' && (
                      <span>{deliveriesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setDeliveriesSort({ column: 'status', direction: deliveriesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('status')}
                    {deliveriesSort.column === 'status' && (
                      <span>{deliveriesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setDeliveriesSort({ column: 'date', direction: deliveriesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('date')}
                    {deliveriesSort.column === 'date' && (
                      <span>{deliveriesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setDeliveriesSort({ column: 'duration', direction: deliveriesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('duration')}
                    {deliveriesSort.column === 'duration' && (
                      <span>{deliveriesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setDeliveriesSort({ column: 'route_id', direction: deliveriesSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('route_id')}
                    {deliveriesSort.column === 'route_id' && (
                      <span>{deliveriesSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                </tr>
              </thead>
              <tbody>
                {sortedDeliveries.map((d) => (
                  <tr key={d.id} className="border text-center h-10">
                    <td onClick={() => navigate(`/delivery/${d.id}`)} className="border cursor-pointer hover:bg-gray-300 text-blue-500 underline">{d.id}</td>
                    <td className="border">{d.status}</td>
                    <td className="border">{
                      lang === "en" ? new Date(d.date).toLocaleDateString("en-US") : new Date(d.date).toLocaleDateString("uk-UA")
                    }
                    </td>
                    <td className="border">{d.duration}</td>
                    <td className="border">
                      {data.creator.id === user.id ? (
                        <select value={d.route_id} onChange={(e) => updateDeliveryRoute(d.id, e.target.value)} className="border rounded px-2 py-1">
                          {routes.map((route) => (
                            <option key={route.id} value={route.id}>
                              {route.name}
                            </option>
                          ))}
                        </select>
                      ) : (
                        <span>{d.route_id}</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </>
        );
      case "users":
        const sortedUsers = sortData(users, usersSort.column, usersSort.direction);
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
                <tr className="text-center h-10">
                  <th className="border" onClick={() => setUsersSort({ column: 'UserID', direction: usersSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('id')}
                    {usersSort.column === 'UserID' && (
                      <span>{usersSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setUsersSort({ column: 'user.name', direction: usersSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('name')}
                    {usersSort.column === 'user.name' && (
                      <span>{usersSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  <th className="border" onClick={() => setUsersSort({ column: 'Role', direction: usersSort.direction === 'asc' ? 'desc' : 'asc' })}>
                    {t('role')}
                    {usersSort.column === 'Role' && (
                      <span>{usersSort.direction === 'asc' ? ' ▲' : ' ▼'}</span>
                    )}
                  </th>
                  {data.creator.id == user.id && <th className="border">{t("actions")}</th>}
                </tr>
              </thead>
              <tbody>
                {sortedUsers.map((u) => (
                  <tr key={u.id} className="text-center h-10">
                    <td className="border">{u.UserID}</td>
                    <td className="border">{u.user.name}</td>
                    <td className="border">
                      {data.creator.id == user.id ? (
                        <select value={u.Role} onChange={(e) => updateUserRole(u.UserID, e.target.value)} className="border rounded px-2 py-1">
                          <option value="user">{t("company_role.user")}</option>
                          <option value="admin">{t("company_role.admin")}</option>
                          <option value="manager">{t("company_role.manager")}</option>
                        </select>
                      ) : (
                        <span>{u.Role}</span>
                      )}
                    </td>
                    {user.id == users[0]?.company.CreatorID && (
                      <td className="border">
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
