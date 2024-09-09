import Head from "next/head";
import { useRouter } from "next/router";
import { useState, useEffect } from "react";
import { message } from "antd";

import styles from "./styles.module.css";
import { OnlineUsers } from "./online-users/index";
import { OnlineStatus } from "../../components/status/index";
import {
  CheckTokenExpireTime,
  DecodeToken,
  GetToken,
} from "../../components/jwt/index";
import { handleStartChatAPI } from "@/apis/handlers/chat";

export default function Chat() {
  const [token, setToken] = useState("");
  const [socket, setSocket] = useState(null);
  const [messageInput, setMessageInput] = useState("");
  const [messages, setMessages] = useState([]);
  const router = useRouter();

  useEffect(() => {
    const storedToken = GetToken();
    if (!storedToken) {
      router.push("/");
      return;
    }
    setToken(storedToken);
    const decodedToken = DecodeToken(storedToken);
    CheckTokenExpireTime(handleSignOut, decodedToken);
  }, []);

  const handleSendMessage = () => {
    if (socket && messageInput.trim() !== "") {
      socket.emit("message", messageInput);
      console.log("Sent message:", messageInput);
      setMessageInput("");
    } else {
      console.log("Socket not available or message input is empty");
    }
  };

  const handleKeyPress = (event) => {
    if (event.key === "Enter") {
      handleSendMessage();
    }
  };

  const handleSignOut = () => {
    localStorage.removeItem("token");
    setToken("");
    router.push("/");
    message.success("Sign out successfully!");
  };

  const handleStartChat = (userId) => {
    handleStartChatAPI(token, userId, (data) => {
      setMessages(data.messages);
    });
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
      {/* <div> */}
      {/* <button onClick={handleSignOut}>Sign out</button> */}
      {/* </div> */}
      <OnlineUsers handleStartChat={handleStartChat} />

      <div className={styles.chatContainer}>
        <div className={styles.messages}>
          {messages &&
            [...messages].reverse().map((msg, index) => (
              <div className={styles.messageContent} key={index}>
                <span className={styles.username}>{msg.username}:</span>
                <span className={styles.messageText}>{msg.content}</span>
              </div>
            ))}
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
