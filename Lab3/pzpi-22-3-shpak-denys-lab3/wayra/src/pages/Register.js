import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const Register = () => {
   const [username, setUsername] = useState("");
   const [password, setPassword] = useState("");
   const [loading, setLoading] = useState(false);
   const [error, setError] = useState(null);
   const navigate = useNavigate();

   const handleRegister = async (e) => {
      e.preventDefault();
      setLoading(true);
      setError(null);

      try {
         const res = await axios.post(
            "http://localhost:8081/auth/register",
            { username, password },
            {
               headers: {
                  Accept: "application/json",
                  "Content-Type": "application/json",
               },
            }
         );

         if (res.data.message === "User registered successfully") {
            navigate("/login");
         }
      } catch (err) {
         setError("Error registering user. Please try again.");
         console.error(err);
      } finally {
         setLoading(false);
      }
   };

   return (
      <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
         <h2 className="text-2xl font-bold text-center mb-4">Register</h2>

         {error && <p className="text-red-600 text-center">{error}</p>}

         <form onSubmit={handleRegister} className="space-y-4">
            <div>
               <label className="block text-gray-700 font-medium">User name</label>
               <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>
            <div>
               <label className="block text-gray-700 font-medium">Password </label>
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
               className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
               disabled={loading}
            >
               {loading ? "Register..." : "Register  "}
            </button>
         </form>
      </div>
   );
};

export default Register;
