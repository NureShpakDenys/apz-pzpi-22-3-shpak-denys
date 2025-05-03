import { Link } from "react-router-dom";

function Header({ user, setUser }) {
  const handleLogout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  return (
    <header className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <Link to="/" className="text-xl font-bold">Blog</Link>
      {/* <nav>
        <Link to="/" className="mr-4">Home</Link>
        {user ? (
          <>
            <Link to="/create" className="mr-4">Create Article</Link>
            <Link to={`/user/${user.id}`} className="mr-4">Profile</Link>
            <button onClick={handleLogout} className="bg-red-500 px-3 py-1 rounded">Logout</button>
          </>
        ) : (
          <>
            <Link to="/login" className="mr-4">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </nav> */}
    </header>
  );
}

export default Header;