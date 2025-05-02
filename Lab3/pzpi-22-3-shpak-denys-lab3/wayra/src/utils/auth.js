import { jwtDecode } from "jwt-decode";

export const getUserFromToken = () => {
  const token = localStorage.getItem('token')
  if (!token) return null;
  try {
    const decoded = jwtDecode(token);
    return decoded.user || decoded;
  } catch (error) {
    console.error("Invalid token:", error);
    return null;
  }
};