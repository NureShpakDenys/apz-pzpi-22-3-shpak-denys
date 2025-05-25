import React, { useState } from "react";
import QRCode from "react-qr-code";

export default function InstallerPanel() {
  const [tab, setTab] = useState("exe");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleDownloadExe = async () => {
    setError("");

    if (!username || !password) {
      setError("Please enter both username and password.");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/get-creds", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ username, password })
      });

      if (response.ok) {
        const link = document.createElement("a");
        link.href = "/main.exe";
        link.download = "main.exe";
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } else {
        setError("Invalid credentials. Access denied.");
      }
    } catch (err) {
      setError("Server error. Please try again later.");
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-6xl mx-auto px-4 py-3 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-blue-700">Wayra company</h1>
          <nav className="space-x-4">
            <button
              className={`py-2 px-4 rounded-md font-medium ${
                tab === "exe"
                  ? "bg-blue-600 text-white"
                  : "text-blue-600 hover:bg-blue-100"
              }`}
              onClick={() => setTab("exe")}
            >
              EXE Installer
            </button>
            <button
              className={`py-2 px-4 rounded-md font-medium ${
                tab === "apk"
                  ? "bg-blue-600 text-white"
                  : "text-blue-600 hover:bg-blue-100"
              }`}
              onClick={() => setTab("apk")}
            >
              APK Installer
            </button>
          </nav>
        </div>
      </header>

      <div className="max-w-md mx-auto mt-10 bg-white shadow-md rounded-xl p-7">
        {tab === "exe" && (
          <>
            <h2 className="text-xl font-semibold mb-4 text-center">
              Windows Installer Access
            </h2>
            <div className="space-y-4">
              <input
                type="text"
                placeholder="Username"
                className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
              <input
                type="password"
                placeholder="Password"
                className="w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              {error && <p className="text-red-500 text-sm">{error}</p>}
              <button
                onClick={handleDownloadExe}
                className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
              >
                Download Installer (.exe)
              </button>
            </div>
          </>
        )}

        {tab === "apk" && (
          <div className="flex flex-col items-center">
            <h2 className="text-xl font-semibold mb-4 text-center">
              Mobile APK Installer
            </h2>
            <QRCode value="http://192.168.1.102:8080/get-apk" />
            <p className="mt-4 text-sm text-gray-600">
              Scan to download APK file for Android
            </p>
          </div>
        )}
      </div>
    </div>
  );
}