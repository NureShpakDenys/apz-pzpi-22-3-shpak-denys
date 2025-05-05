import { Link } from "react-router-dom";
import { useEffect } from "react";

function Header({ user, setUser, i18n, t }) {
  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  };

  const handleLanguageChange = (e) => {
    i18n.changeLanguage(e.target.value);
  };

  useEffect(() => {

  }, [user]);

  return (
    <header className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <Link to="/companies" className="text-xl font-bold">{t("companies")}</Link>
      <nav>
        <Link to="/companies" className="mr-4">{t("home")}</Link>

        {user ? (
          <>
            {user?.role == "admin" && (
              <Link to="/admin" className="mr-4">{t("admin")}</Link>
            )}
            {user?.role == "system_admin" && (
              <Link to="/system-admin" className="mr-4">{t("systemAdmin")}</Link>
            )}
            {user?.role == "db_admin" && (
              <Link to="/db-admin" className="mr-4">{t("dbAdmin")}</Link>
            )}

            <button
              onClick={handleLogout}
              className="bg-red-500 px-3 py-1 rounded"
            >
              {t("logout")}
            </button>
          </>
        ) : (
          <>
            <Link to="/login" className="mr-4">{t("login")}</Link>
            <Link to="/register">{t("register")}</Link>
          </>
        )}

        <select
          onChange={handleLanguageChange}
          value={i18n.language}
          className="bg-gray-700 text-white ml-3 px-2 py-1 rounded"
        >
          <option value="en">EN</option>
          <option value="uk">UA</option>
        </select>
      </nav>
    </header>
  );
}

export default Header;
