import Head from "next/head";
import { useRouter } from "next/router";
import { useState, useEffect, useRef } from "react";
import { message } from "antd";

import styles from "./styles.module.css";
import { OnlineUsers } from "./online-users/index";
import { OnlineStatus } from "../../components/status/index";
import {
  CheckTokenExpireTime,
  DecodeToken,
  GetToken,
} from "../../components/jwt/index";

export default function Chat() {
  const [token, setToken] = useState("");
  const [messageInput, setMessageInput] = useState("");
  const [messages, setMessages] = useState([]);
  const [receiverId, setReceiverId] = useState(0);
  const [receiverUsername, setReceiverUsername] = useState("");
  const [userId, setUserId] = useState(0);
  const router = useRouter();
  const socketRef = useRef(null);

  const messagesEndRef = useRef(null);

  useEffect(() => {
    const storedToken = GetToken();
    if (!storedToken) {
      router.push("/");
      return;
    }
    setToken(storedToken);
    const decodedToken = DecodeToken(storedToken);
    if (decodedToken && decodedToken.user_id) {
      setUserId(decodedToken.user_id);
    }

    CheckTokenExpireTime(handleSignOut, decodedToken);
  }, []);

  useEffect(() => {
    if (receiverId === 0) return;

    const ws = new WebSocket(
      `ws://localhost:8080/ws/chat?token=${token}&receiver_id=${receiverId}`
    );
    socketRef.current = ws;

    ws.onopen = () => {
      console.log("Connected to WebSocket server");
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (Array.isArray(data)) {
        setMessages(data);
      } else {
        setMessages((prevMessages) => [...prevMessages, data]);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
      message.error("Failed to connect to server");
    };

    ws.onclose = () => {
      console.log("Disconnected from WebSocket server");
    };

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [receiverId, token]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const handleKeyPress = (event) => {
    if (event.key === "Enter") {
      handleSendMessage();
    }
  };

  const handleSendMessage = () => {
    if (receiverId == 0) {
      message.error("Please select a user.");
    } else if (!messageInput.trim()) {
      return;
    } else {
      const trimmedMessage = messageInput.trim();
      if (
        socketRef.current &&
        socketRef.current.readyState === WebSocket.OPEN
      ) {
        const messageData = {
          message: trimmedMessage,
        };
        socketRef.current.send(JSON.stringify(messageData));
        setMessageInput("");
      } else {
        message.error("WebSocket connection is not open.");
      }
    }
  };

  const handleSignOut = () => {
    localStorage.removeItem("token");
    setToken("");
    router.push("/");
    message.success("Sign out successfully!");
  };

  return (
    <div className={styles.page}>
      <Head>
        <title>Chat App</title>
        <meta name="description" content="Real-time chatting website" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <OnlineStatus token={token} />
      {/* <div className={styles.onlineUsersAndSignOutContainer}> */}
      {/* <button onClick={handleSignOut} className={styles.signOutButton}>
          Sign out
        </button> */}
      <OnlineUsers
        className={styles.onlineUsers}
        setMessageInput={setMessageInput}
        setReceiverId={setReceiverId}
        setReceiverUsername={setReceiverUsername}
        userId={userId}
      />
      {/* </div> */}

      <div className={styles.chatContainer}>
        {receiverUsername && (
          <div className={styles.receiverUsername}>{receiverUsername}</div>
        )}

        <div className={styles.messages}>
          {messages &&
            [...messages].reverse().map((msg, index) => (
              <div className={styles.messageContent} key={index}>
                <span className={styles.username}>{msg.username}:</span>
                <span className={styles.messageText}>{msg.content}</span>
              </div>
            ))}
          <div ref={messagesEndRef} />
        </div>
        <div className={styles.inputContainer}>
          <input
            type="text"
            value={messageInput}
            onKeyUp={handleKeyPress}
            onChange={(e) => setMessageInput(e.target.value)}
            placeholder="Enter your message"
            className={styles.input}
          />
          <button onClick={handleSendMessage} className={styles.button}>
            Send
          </button>
        </div>
      </div>
    </div>
  );
}
