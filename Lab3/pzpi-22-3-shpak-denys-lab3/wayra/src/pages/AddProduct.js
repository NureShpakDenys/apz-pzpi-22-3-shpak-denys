import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";
import convert from "../utils/convertors";

const AddProduct = ({ t, i18n}) => {
   const { delivery_id } = useParams();
   const [name, setName] = useState("");
   const [productType, setProductType] = useState("Fruits");
   const [weight, setWeight] = useState(1);
   const [loading, setLoading] = useState(false);
   const [error, setError] = useState(null);
   const navigate = useNavigate();

   const token = localStorage.getItem("token");

   const handleSubmit = async (e) => {
      e.preventDefault();
      setLoading(true);
      setError(null);
      const lang = i18n.language;

      const weightToSend = lang === "en" ? convert(parseFloat(weight), "weight", "metrical") : parseFloat(weight);
      
      try {
         const response = await axios.post(
            "http://localhost:8081/products/",
            {
               deliveryID: Number(delivery_id),
               name,
               product_type: productType,
               weight: weightToSend,
            },
            {
               headers: {
                  Authorization: `Bearer ${token}`,
                  Accept: "application/json",
                  "Content-Type": "application/json",
               },
            }
         );

         navigate("/delivery/" + delivery_id);
      } catch (err) {
         setError("Error while adding product");
         console.error(err);
      } finally {
         setLoading(false);
      }
   };

   return (
      <div className="p-6 max-w-lg mx-auto bg-white shadow-md rounded">
         <h2 className="text-2xl font-bold text-center mb-4">{t("add_product")}</h2>

         {error && <p className="text-red-600 text-center">{error}</p>}

         <form onSubmit={handleSubmit} className="space-y-4">
            <div>
               <label className="block text-gray-700 font-medium">{t("product_name")}</label>
               <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>

            <div>
               <label className="block text-gray-700 font-medium">{t("product_type")}</label>
               <select
                  value={productType}
                  onChange={(e) => setProductType(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
               >
                  <option value="Fruits">{t("product_type.Fruits")}</option>
                  <option value="Vegetables">{t("product_type.Vegetables")}</option>
                  <option value="Frozen Foods">{t("product_type.Frozen Foods")}</option>
                  <option value="Dairy Products">{t("product_type.Dairy Products")}</option>
                  <option value="Meat">{t("product_type.Meat")}</option>
               </select>

            </div>

            <div>
               <label className="block text-gray-700 font-medium">{t("weight")}</label>
               <input
                  type="number"
                  value={weight}
                  onChange={(e) => setWeight(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>

            <button
               type="submit"
               className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600 transition"
               disabled={loading}
            >
               {loading ? t("adding") : t("add_product")}
            </button>
         </form>
      </div>
   );
};

export default AddProduct;
