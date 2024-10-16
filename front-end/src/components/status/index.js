import { useEffect } from "react";
import { message } from "antd";

export default function OnlineStatus({ token }) {
  useEffect(() => {
    if (!token) return;

    const ws = new WebSocket(`ws://localhost:8080/ws/status?token=${token}`);

    ws.onopen = () => {
      console.log("Connected to WebSocket server");
      message.success("You are now online");
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
      message.error("Failed to connect to server");
    };

    ws.onclose = () => {
      console.log("Disconnected from WebSocket server");
    };

    return () => {
      ws.close();
      console.log("WebSocket connection closed");
    };
  }, [token]);

  return null;
}
