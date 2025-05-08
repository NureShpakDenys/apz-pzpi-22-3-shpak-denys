import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Login = ({ setUser, t }) => {
   const [username, setUsername] = useState("");
   const [password, setPassword] = useState("");
   const [loading, setLoading] = useState(false);
   const [error, setError] = useState(null);
   const navigate = useNavigate();

   const handleLogin = async (e) => {
      e.preventDefault();
      setLoading(true);
      setError(null);
    
      try {
        const res = await axios.post(
          "http://localhost:8081/auth/login",
          { username, password },
          {
            headers: {
              Accept: "application/json",
              "Content-Type": "application/json",
            },
          }
        );
    
        if (res.data.token) {
          localStorage.setItem("token", res.data.token);
    
          setUser(res.data.user);
          navigate("/companies");
        }
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

   return (
      <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
         <h2 className="text-2xl font-bold text-center mb-4">{t("login")}</h2>

         {error && <p className="text-red-600 text-center">{error}</p>}

         <form onSubmit={handleLogin} className="space-y-4">
            <div>
               <label className="block text-gray-700 font-medium">{t("username")}</label>
               <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>
            <div>
               <label className="block text-gray-700 font-medium">{t("password")}</label>
               <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>

            <button
               type="submit"
               className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600 transition"
               disabled={loading}
            >
               {loading ? t("loading") : t("login")}
            </button>
         </form>
      </div>
   );
};

export default Login;
