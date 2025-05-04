import { useEffect, useState } from "react";
import axios from "axios";

const API_BASE = "http://localhost:8081/admin";

const DBAdminDashboard = () => {
   const token = localStorage.getItem("token");
   const authHeader = { headers: { Authorization: `Bearer ${token}` } };

   const [status, setStatus] = useState({
      DatabaseSizeMB: 0,
      ActiveConnections: 0,
      LastBackupTime: "",
   });

   const [backupPath, setBackupPath] = useState("");
   const [message, setMessage] = useState("");
   const [loading, setLoading] = useState(true);

   const fetchDBStatus = async () => {
      try {
         const res = await axios.get(`${API_BASE}/db-status`, authHeader);
         setStatus(res.data);
      } catch (error) {
         console.error("Error get DB status:", error);
      }
   };

   const performBackup = async () => {
      try {
         const res = await axios.post(
            `${API_BASE}/backup`,
            { backup_path: backupPath },
            authHeader
         );
         setMessage(res.data.message);
         fetchDBStatus();
      } catch (error) {
         console.error("Backup error:", error);
      }
   };

   const performRestore = async () => {
      try {
         const res = await axios.post(
            `${API_BASE}/restore`,
            { backup_path: backupPath },
            authHeader
         );
         setMessage(res.data.message);
         fetchDBStatus();
      } catch (error) {
         console.error("Restore error:", error);
      }
   };

   const optimizeDB = async () => {
      try {
         const res = await axios.post(`${API_BASE}/optimize`, {}, authHeader);
         setMessage(res.data.message);
      } catch (error) {
         console.error("Optimize error:", error);
      }
   };

   useEffect(() => {
      const fetchAll = async () => {
         await fetchDBStatus();
         setLoading(false);
      };
      fetchAll();
   }, []);

   if (loading) return <div className="p-4 text-center text-lg">Loading...</div>;

   return (
      <div className="p-4 flex flex-col gap-4">
         <h1 className="text-2xl font-bold">Database Admin Dashboard</h1>

         {message && (
            <div className="p-3 bg-green-100 border border-green-300 rounded">
               ✅ {message}
            </div>
         )}

         <div className="border p-4 rounded-xl shadow bg-white">
            <h2 className="text-xl font-semibold mb-2">📊 Database Info</h2>
            <p><strong>Database Size:</strong> {status.DatabaseSizeMB} MB</p>
            <p><strong>Active Connections:</strong> {status.ActiveConnections}</p>
            <p><strong>Last Backup:</strong> {status.LastBackupTime ? new Date(status.LastBackupTime).toLocaleString() : "N/A"}</p>
            <button
               onClick={fetchDBStatus}
               className="mt-2 px-4 py-1 border rounded hover:bg-gray-100"
            >
               Refresh Info
            </button>
         </div>

         <div className="border p-4 rounded-xl shadow bg-white">
            <h2 className="text-xl font-semibold mb-2">💾 Backup & Restore</h2>
            <label className="block mb-2">
               Backup Path:
               <input
                  type="text"
                  value={backupPath}
                  onChange={(e) => setBackupPath(e.target.value)}
                  placeholder="Enter folder path"
                  className="ml-2 border px-2 py-1 w-1/2"
               />
            </label>
            <div className="flex gap-4 mt-2">
               <button
                  onClick={performBackup}
                  className="px-4 py-1 border rounded hover:bg-gray-100"
               >
                  Create Backup
               </button>
               <button
                  onClick={performRestore}
                  className="px-4 py-1 border rounded hover:bg-gray-100"
               >
                  Restore from Backup
               </button>
            </div>
         </div>

         <div className="border p-4 rounded-xl shadow bg-white">
            <h2 className="text-xl font-semibold mb-2">🛠️ Optimize Database</h2>
            <button
               onClick={optimizeDB}
               className="px-4 py-1 border rounded hover:bg-gray-100"
            >
               Run Optimization
            </button>
         </div>

 
      </div>
   );
};

export default DBAdminDashboard;
