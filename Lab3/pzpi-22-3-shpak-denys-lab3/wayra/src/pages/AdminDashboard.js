import { useEffect, useState } from "react";
import axios from "axios";

const AdminDashboard = ({ user_id, t }) => {
  const [user, setUser] = useState({});
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const token = localStorage.getItem("token");

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const usersRes = await axios.get("http://localhost:8081/users", {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json" },
        });

        const rolesRes = [
          { id: 1, name: "admin" },
          { id: 2, name: "user" },
          { id: 3, name: "manager" },
          { id: 4, name: "db_admin" },
          { id: 5, name: "system_admin" },
        ];

        setUsers(usersRes.data);
        setRoles(rolesRes);
      } catch (err) {
        setError("Error while loading data");
        console.error(err);
      }
    };

    fetchData();

    const fetchUser = async () => {
      setLoading(true);
      try {
        const res = await axios.get(`http://localhost:8081/user/${user_id}`, {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json" },
        });
        console.log(JSON.stringify(res.data));
        setUser(res.data);
      } catch (err) {
        setError("Error while loading user data");
        console.error(err);
      } finally {
        setLoading(false);
      }
    }

    fetchUser();
  }, [token]);

  const changeUserRole = async (userId, roleId) => {
    try {
      await axios.post(
        "http://localhost:8081/admin/change-role",
        { userId, roleId },
        {
          headers: { Authorization: `Bearer ${token}`, Accept: "application/json", "Content-Type": "application/json" },
        }
      );

      setUsers((prevUsers) =>
        prevUsers.map((user) =>
          user.id == userId ? { ...user, role: roles.find((r) => r.id == roleId).name } : user
        )
      );
    } catch (err) {
      console.error("Error updating user role:", err);
    }
  };

  if (loading) return <div className="p-6 text-center">{t("loading")}</div>;
  if (error) return <div className="p-6 text-center text-red-600">{error}</div>;

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <h1 className="text-3xl font-bold text-center mb-4">{t("admin_dashboard")}</h1>

      <div className="bg-white shadow-md rounded p-4">
        <table className="w-full table-auto border-collapse border border-gray-300">
          <thead>
            <tr className="bg-gray-100">
              <th className="border p-2">{t("id")}</th>
              <th className="border p-2">{t("name")}</th>
              <th className="border p-2">{t("password")}</th>
              {user.role == "admin" && <th className="border p-2">{t("actions")}</th>}
            </tr>
          </thead>
          <tbody>
            {users.map((u) => (
              <tr key={u.id} className="text-center">
                <td className="border p-2">{u.id}</td>
                <td className="border p-2">{u.name}</td>
                <td className="border p-2">
                  {user.role == "admin" ? (
                    <select
                      value={roles.find((r) => r.name == u.role)?.id || ""}
                      onChange={(e) => changeUserRole(u.id, Number(e.target.value))}
                      className="border rounded px-2 py-1"
                    >
                      {roles.map((role) => (
                        <option key={role.id} value={role.id}>
                          {role.name}
                        </option>
                      ))}
                    </select>
                  ) : (
                    <span>{u.role}</span>
                  )}
                </td>
                {user.role == "admin" && (
                  <td className="border p-2">
                    <button
                      onClick={() => changeUserRole(u.id, roles.find((r) => r.name == u.role)?.id)}
                      className="px-3 py-1 bg-blue-500 text-white rounded"
                    >
                      {t("save")}
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

export default AdminDashboard;
