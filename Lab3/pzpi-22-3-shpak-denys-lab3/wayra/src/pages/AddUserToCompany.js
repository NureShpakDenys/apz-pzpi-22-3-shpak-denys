import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const AddUserToCompany = ({ user, t}) => {
  const { company_id } = useParams();
  const [users, setUsers] = useState([]);
  const [companyUsers, setCompanyUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState("");
  const [role, setRole] = useState("user");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const token = localStorage.getItem("token");

  useEffect(() => {
    const fetchUsers = async () => {
      try {
         const resUsers = await axios.get("http://localhost:8081/users", {
           headers: {
             Authorization: `Bearer ${token}`,
             Accept: "application/json",
           },
         });
 
         const resCompany = await axios.get(`http://localhost:8081/company/${company_id}`, {
           headers: {
             Authorization: `Bearer ${token}`,
             Accept: "application/json",
           },
         });
 
         setUsers(resUsers.data);
         setCompanyUsers(resCompany.data.users.map((user) => user.id));
       } catch (err) {
         setError("Error fetching users or company data");
         console.error(err);
       }
    };

    fetchUsers();
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await axios.post(
        `http://localhost:8081/company/${company_id}/add-user`,
        { userID: Number(selectedUser), role },
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
      setError("Error adding user to company");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
      <h2 className="text-2xl font-bold text-center mb-4">{t("add_user_to_company")}</h2>

      {error && <p className="text-red-600 text-center">{error}</p>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-gray-700 font-medium">{t("select_user")}</label>
          <select
            value={selectedUser}
            onChange={(e) => setSelectedUser(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            <option value="">{t("select_user")}</option>
            {users
              .filter((user) => !companyUsers.includes(user.id))
              .map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name}
                </option>
              ))}
          </select>
        </div>

        <div>
          <label className="block text-gray-700 font-medium">{t("select_role")}</label>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option>user</option>
            <option>admin</option>
            <option>manager</option>
          </select>
        </div>

        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
          disabled={loading}
        >
          {loading ? t("adding") : t("add_user_to_company")}
        </button>
      </form>
    </div>
  );
};

export default AddUserToCompany;
