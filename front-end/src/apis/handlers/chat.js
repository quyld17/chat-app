import { message } from "antd";

export function handleStartChatAPI(token, receiver_id, onMessage) {
  const ws = new WebSocket(
    `ws://localhost:8080/ws/chat?token=${token}&receiver_id=${receiver_id}`
  );

  ws.onopen = () => {
    console.log("Connected to WebSocket server");
  };

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (onMessage) {
      onMessage(data);
    }
  };

  ws.onerror = (err) => {
    console.error("WebSocket error:", err);
    message.error("Failed to connect to server");
  };

  ws.onclose = () => {
    console.log("Disconnected from WebSocket server");
  };

  return ws;
}
