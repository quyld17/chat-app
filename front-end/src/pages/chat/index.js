import Head from "next/head";
import { Inter } from "next/font/google";
import { useRouter } from "next/router";
import { useState, useEffect } from "react";
import { message } from "antd";

import styles from "./index.module.css";
// import { OnlineList } from "../../components/online-list/index";
import { OnlineStatus } from "../../components/status/index";
import {
  CheckTokenExpireTime,
  DecodeToken,
  GetToken,
} from "../../components/jwt/index";

const inter = Inter({ subsets: ["latin"] });

export default function Chat() {
  const [token, setToken] = useState("");
  const [username, setUsername] = useState("");
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
    CheckTokenExpireTime(handleSignOut, setUsername, decodedToken);
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

  return (
    <>
      <Head>
        <title>Chat App</title>
        <meta name="description" content="Real-time chatting website" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <OnlineStatus token={token} />
      <main className={`${styles.main} ${inter.className}`}>
        <div className={styles.chatContainer}>
          <h1>Chat App</h1>
          <p>{username}</p>
          <div className={styles.messages}>
            {messages.map((msg, index) => (
              <p key={index}>{msg}</p>
            ))}
          </div>
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
          <button onClick={handleSignOut} className={styles.button}>
            Sign out
          </button>
        </div>
      </main>
    </>
  );
}
