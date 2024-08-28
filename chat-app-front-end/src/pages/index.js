import Head from "next/head";
import { Inter } from "next/font/google";
import { useState, useEffect } from "react";
import styles from "./index.module.css";

const inter = Inter({ subsets: ["latin"] });

export default function Home() {
  const [socket, setSocket] = useState(null);
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onmessage = function (event) {
      setMessages((prevMessages) => [...prevMessages, event.data]);
    };

    setSocket(ws);

    return () => {
      ws.close();
    };
  }, []);

  const sendMessage = () => {
    if (socket && message.trim() !== "") {
      socket.send(message);
      setMessage("");
    }
  };

  const handleKeyPress = (event) => {
    if (event.key === "Enter") {
      sendMessage();
    }
  };

  return (
    <>
      <Head>
        <title>Chat App</title>
        <meta name="description" content="Real time chatting website" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className={`${styles.main} ${inter.className}`}>
        <div className={styles.chatContainer}>
          <h1>WebSocket Chat</h1>
          <div className={styles.messages}>
            {messages.map((msg, index) => (
              <p key={index}>{msg}</p>
            ))}
          </div>
          <input
            type="text"
            value={message}
            onKeyUp={handleKeyPress}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Enter your message"
            className={styles.input}
          />
          <button onClick={sendMessage} className={styles.button}>
            Send
          </button>
        </div>
      </main>
    </>
  );
}
