require("dotenv").config();

const URL = process.env.NEXT_PUBLIC_URL;
const PORT = process.env.NEXT_PUBLIC_PORT;

export function ServerURL() {
  const baseURL = `http://${URL}:${PORT}`;
  return baseURL;
}
