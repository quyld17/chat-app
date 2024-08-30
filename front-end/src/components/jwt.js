require("dotenv").config();

const URL = process.env.NEXT_PUBLIC_URL;
const PORT = process.env.NEXT_PUBLIC_PORT;

export default function getJWT() {
  const baseURL = `http://${URL}:${PORT}`;
  const token = localStorage.getItem("token");
  return { baseURL, token };
}
