import { Link } from "react-router-dom";
import { useEffect } from "react";

function Header({ user, setUser }) {
  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  };

  useEffect(() => {
    
  }, [user]);

  return (
    <header className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <Link to="/companies" className="text-xl font-bold">Companies</Link>
      <nav>
        <Link to="/companies" className="mr-4">Home</Link>

        {user ? (
          <>
            {user?.role == "admin" && (
              <Link to="/admin" className="mr-4">Admin</Link>
            )}
            {user?.role == "system_admin" && (
              <Link to="/system-admin" className="mr-4">System Admin</Link>
            )}
            {user?.role == "db_admin" && (
              <Link to="/db-admin" className="mr-4">DB Admin</Link>
            )}

            <button
              onClick={handleLogout}
              className="bg-red-500 px-3 py-1 rounded"
            >
              Logout
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className="mr-4">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </nav>
    </header>
  );
}

export default Header;
