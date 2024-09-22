import { message } from "antd";
import getMethodAPI from "../methods/get-method-api";

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

export function handleGetChatHistory(receiver_id) {
  return new Promise((resolve, reject) => {
    const endpoint = `/chat-history?receiver_id=${receiver_id}`;
    getMethodAPI(
      endpoint,
      (data) => {
        resolve(data);
      },
      (error) => {
        reject(error);
        // message.error(error);
      }
    );
  });
}
