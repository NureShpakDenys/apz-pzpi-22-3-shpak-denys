import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { ArrowRight, PlusCircle } from "lucide-react";

const CompanyList = () => {
   const [companies, setCompanies] = useState([]);
   const [loading, setLoading] = useState(true);
   const [error, setError] = useState(null);
   const navigate = useNavigate();

   const token = localStorage.getItem("token");
   useEffect(() => {
      const fetchCompanies = async () => {
         try {
            const res = await axios.get("http://localhost:8081/company/", {
               headers: {
                  Authorization: `Bearer ${token}`,
                  Accept: "application/json",
               },
            });
            setCompanies(res.data);
         } catch (err) {
            setError("Ошибка загрузки компаний");
         } finally {
            setLoading(false);
         }
      };

      fetchCompanies();
   }, []);

   if (loading) return <div className="p-6 text-center">loading...</div>;
   if (error) return <div className="p-6 text-red-600">{error}</div>;

   return (
      <div className="p-6 max-w-4xl mx-auto">
         <h1 className="text-3xl font-bold mb-4 text-center">Companies</h1>
         
         <div className="text-center mb-4">
            <button
               onClick={() => navigate("/company/create")}
               className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition flex items-center justify-center space-x-2"
            >
               <PlusCircle className="w-5 h-5" />
               <span>Create company</span>
            </button>
         </div>

         <div className="space-y-4">
            {companies.map((company) => (
               <div
                  key={company.id}
                  className="flex justify-between items-center border p-4 rounded shadow cursor-pointer hover:bg-gray-50"
                  onClick={() => navigate(`/company/${company.id}`)}
               >
                  <div>
                     <p className="text-lg font-semibold">{company.name}</p>
                     <p className="text-gray-600">{company.address}</p>
                  </div>
                  <ArrowRight className="w-6 h-6 text-gray-500" />
               </div>
            ))}
         </div>
      </div>
   );
};

export default CompanyList;
