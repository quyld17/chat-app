import { message } from "antd";

export function handleChatAPI(token, receiverId, messageInput, onMessage) {
  const ws = new WebSocket(
    `ws://localhost:8080/ws/chat?token=${token}&receiver_id=${receiverId}`
  );

  ws.onopen = () => {
    console.log("Connected to WebSocket server");

    const trimmedMessage = messageInput.trim();
    if (!trimmedMessage) {
      return;
    }
    const messageData = {
      message: trimmedMessage,
    };
    ws.send(JSON.stringify(messageData));
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
