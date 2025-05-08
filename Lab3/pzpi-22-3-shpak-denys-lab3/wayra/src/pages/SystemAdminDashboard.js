import { useEffect, useState } from "react";
import axios from "axios";

const API_BASE = "http://localhost:8081/admin";

const SystemAdminDashboard = ({ t, i18n }) => {
   const [loading, setLoading] = useState(true);
   const token = localStorage.getItem("token");
   const authHeader = { headers: { Authorization: `Bearer ${token}` } };
   const [health, setHealth] = useState({});
   const [configs, setConfigs] = useState({
      auth_token_ttl_hours: 0,
      encryption_key_exists: false,
      http_timeout_seconds: 0,
   });
   const [logs, setLogs] = useState([]);
   const [clearDays, setClearDays] = useState(3);

   const [currentPage, setCurrentPage] = useState(1);
   const logsPerPage = 10;

   const indexOfLastLog = currentPage * logsPerPage;
   const indexOfFirstLog = indexOfLastLog - logsPerPage;
   const logsToDisplay = logs.slice(indexOfFirstLog, indexOfLastLog);
   const totalPages = Math.ceil(logs.length / logsPerPage);
   const lang = i18n.language;

   const fetchHealth = async () => {
      const res = await axios.get(`${API_BASE}/health`, authHeader);
      setHealth(res.data);
   };

   const fetchConfigs = async () => {
      const res = await axios.get(`${API_BASE}/system-configs`, authHeader);
      setConfigs(res.data);
   };

   const saveConfigs = async () => {
      await axios.put(
         `${API_BASE}/system-configs`,
         {
            timeout_sec: configs.http_timeout_seconds,
            token_ttl: configs.auth_token_ttl_hours,
         },
         authHeader
      );
   };

   const fetchLogs = async () => {
      const res = await axios.get(`${API_BASE}/logs`, authHeader);
      setLogs(res.data);
   };

   const clearLogs = async () => {
      await axios.post(`${API_BASE}/clear-logs`, { days: clearDays }, authHeader);
      fetchLogs();
   };

   useEffect(() => {
      const fetchAll = async () => {
         try {
            await Promise.all([
               fetchHealth(),
               fetchConfigs(),
               fetchLogs()
            ]);
         } catch (error) {
            console.error("loading error:", error);
         } finally {
            setLoading(false);
         }
      };

      fetchAll();
   }, []);

   if (loading) {
      return <div className="p-4 text-center text-lg">{t("loading")}</div>;
   }

   return (
      <div className="flex flex-row gap-4 p-4">
         <div className="flex flex-col gap-4 w-1/4">
            <div className="border p-4 rounded-xl shadow">
               <h2 className="text-xl font-bold mb-2">{t("health_check")}</h2>
               <p>{t("db_status")}: {health.db_status}</p>
               <p>{t("server_time")}: {
                  lang === "en"
                     ? new Date(health.server_time).toLocaleString("en-US")
                     : new Date(health.server_time).toLocaleString("uk-UA")}</p>
               <p>{t("uptime")}: {health.uptime}</p>
               <button
                  onClick={fetchHealth}
                  className="mt-2 px-4 py-1 border rounded hover:bg-gray-100"
               >
                  {t("check")}
               </button>
            </div>

            <div className="border p-4 rounded-xl shadow">
               <h2 className="text-xl font-bold mb-2">{t("system_configs")}</h2>
               <p>
                  {t("encryption_key_exists")}: {configs.encryption_key_exists ? "✔️" : "❌"}
               </p>
               <div className="my-2">
                  <label>{t("auth_token_ttl_hours")}:</label>
                  <input
                     type="number"
                     className="ml-2 border px-1 w-16"
                     value={configs.auth_token_ttl_hours}
                     onChange={(e) =>
                        setConfigs({
                           ...configs,
                           auth_token_ttl_hours: parseInt(e.target.value),
                        })
                     }
                  />
               </div>
               <div className="mb-2">
                  <label>{t("http_timeout_seconds")}:</label>
                  <input
                     type="number"
                     className="ml-2 border px-1 w-16"
                     value={configs.http_timeout_seconds}
                     onChange={(e) =>
                        setConfigs({
                           ...configs,
                           http_timeout_seconds: parseInt(e.target.value),
                        })
                     }
                  />
               </div>
               <button
                  onClick={saveConfigs}
                  className="px-4 py-1 border rounded hover:bg-gray-100"
               >
                  {t("save")}
               </button>
            </div>
         </div>

         <div className="border p-4 rounded-xl shadow w-3/4">
            <div className="flex items-center mb-4">
               <label className="mr-2">{t("clear_if_older_than")}</label>
               <input
                  type="number"
                  value={clearDays}
                  onChange={(e) => setClearDays(parseInt(e.target.value))}
                  className="border px-2 w-16 mr-2"
               />
               <button
                  onClick={clearLogs}
                  className="px-4 py-1 border rounded hover:bg-gray-100"
               >
                  {t("clear")}
               </button>
            </div>
            <table className="w-full border-t text-sm">
               <thead>
                  <tr className="bg-gray-100">
                     <th className="border px-2 py-1">{t("id")}</th>
                     <th className="border px-2 py-1">{t("created_at")}</th>
                     <th className="border px-2 py-1">{t("user_id")}</th>
                     <th className="border px-2 py-1">{t("actions")}</th>
                     <th className="border px-2 py-1">{t("description")}</th>
                     <th className="border px-2 py-1">{t("description")}</th>
                  </tr>
               </thead>
               <tbody>
                  {logsToDisplay.map((log) => (
                     <tr key={log.ID}>
                        <td className="border px-2 py-1">{log.ID}</td>
                        <td className="border px-2 py-1">
                           {lang === "en"
                              ? new Date(log.CreatedAt).toLocaleString("en-US")
                              : new Date(log.CreatedAt).toLocaleString("uk-UA")}
                        </td>
                        <td className="border px-2 py-1">{log.UserID}</td>
                        <td className="border px-2 py-1">{log.ActionType}</td>
                        <td className="border px-2 py-1">{log.Description}</td>
                        <td className="border px-2 py-1">{log.Success ? "✔️" : "❌"}</td>
                     </tr>
                  ))}
               </tbody>
            </table>
            <div className="flex justify-center mt-4 space-x-2">
               <button
                  onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                  className="px-3 py-1 border rounded hover:bg-gray-100"
                  disabled={currentPage === 1}
               >
                  {t("prev")}
               </button>

               {[...Array(totalPages)].map((_, idx) => {
                  const page = idx + 1;
                  if (
                     page === 1 ||
                     page === totalPages ||
                     (page >= currentPage - 1 && page <= currentPage + 1)
                  ) {
                     return (
                        <button
                           key={idx}
                           onClick={() => setCurrentPage(page)}
                           className={`px-3 py-1 border rounded ${currentPage === page ? "bg-blue-500 text-white" : "hover:bg-gray-100"}`}
                        >
                           {page}
                        </button>
                     );
                  } else if (
                     (page === currentPage - 2 || page === currentPage + 2)
                  ) {
                     return <span key={idx} className="px-2">...</span>;
                  } else {
                     return null;
                  }
               })}

               <button
                  onClick={() =>
                     setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                  }
                  className="px-3 py-1 border rounded hover:bg-gray-100"
                  disabled={currentPage === totalPages}
               >
                  {t("next")}
               </button>
            </div>
         </div>
      </div>
   );
};

export default SystemAdminDashboard;
