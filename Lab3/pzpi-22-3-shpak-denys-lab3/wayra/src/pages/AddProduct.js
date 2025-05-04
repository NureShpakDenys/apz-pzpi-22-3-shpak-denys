import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";

const AddProduct = () => {
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

      try {
         const response = await axios.post(
            "http://localhost:8081/products/",
            {
               deliveryID: Number(delivery_id),
               name,
               product_type: productType,
               weight: parseFloat(weight),
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
         <h2 className="text-2xl font-bold text-center mb-4">Add Product</h2>

         {error && <p className="text-red-600 text-center">{error}</p>}

         <form onSubmit={handleSubmit} className="space-y-4">
            <div>
               <label className="block text-gray-700 font-medium">Product name</label>
               <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
               />
            </div>

            <div>
               <label className="block text-gray-700 font-medium">Product type</label>
               <select
                  value={productType}
                  onChange={(e) => setProductType(e.target.value)}
                  className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
               >
                  <option>Fruits</option>
                  <option>Vegetables</option>
                  <option>Frozen Foods</option>
                  <option>Dairy Products</option>
                  <option>Meat</option>
               </select>
            </div>

            <div>
               <label className="block text-gray-700 font-medium">Weight (kg)</label>
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
               {loading ? "Adding..." : "Add Product"}
            </button>
         </form>
      </div>
   );
};

export default AddProduct;
