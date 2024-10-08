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

import { handleGetChatHistoryAPI } from "../../apis/handlers/chat";

export default function Chat() {
  const [token, setToken] = useState("");
  const [messageInput, setMessageInput] = useState("");
  const [messages, setMessages] = useState([]);
  const [receiverId, setReceiverId] = useState(0);
  const [receiverUsername, setReceiverUsername] = useState("");
  const [username, setUsername] = useState("");
  const [offset, setOffset] = useState(20);
  const [hasMore, setHasMore] = useState(true);
  const [hasLoadedMessages, setHasLoadedMessages] = useState(false);
  const router = useRouter();
  const socketRef = useRef(null);
  const messagesContainerRef = useRef(null);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    const storedToken = GetToken();
    if (!storedToken) {
      router.push("/");
      return;
    }
    setToken(storedToken);
    const decodedToken = DecodeToken(storedToken);
    if (decodedToken) {
      setUsername(decodedToken.name);
    }
    CheckTokenExpireTime(handleSignOut, decodedToken);
  }, []);

  useEffect(() => {
    if (!token) return;

    const ws = new WebSocket(`ws://localhost:8080/ws/chat?token=${token}`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log("Connected to WebSocket server");
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (receiverId == data.sender_id) {
        setMessages((prevMessages) => [
          ...prevMessages,
          { content: data.message, username: receiverUsername },
        ]);
        scrollToBottomAfterNewMessage();
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
  }, [token, receiverId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    if (!hasLoadedMessages) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
      setHasLoadedMessages(true);
    }
  };

  const handleScroll = () => {
    const element = messagesContainerRef.current;
    if (element.scrollTop === 0 && hasMore) {
      loadMessages(offset);
    }
  };

  const loadMessages = async (newOffset) => {
    try {
      const olderMessages = await handleGetChatHistoryAPI(
        receiverId,
        newOffset
      );
      if (olderMessages.length === 0) {
        setHasMore(false);
      } else {
        const previousScrollHeight = messagesContainerRef.current.scrollHeight;

        setMessages((prevMessages) => [...olderMessages, ...prevMessages]);
        setOffset(newOffset + olderMessages.length);

        setTimeout(() => {
          const newScrollHeight = messagesContainerRef.current.scrollHeight;
          const scrollDiff = newScrollHeight - previousScrollHeight;

          messagesContainerRef.current.scrollTop = scrollDiff;
        }, 0);
      }
    } catch (error) {
      console.error("Error loading chat history:", error);
    }
  };

  const scrollToBottomAfterNewMessage = () => {
    setTimeout(() => {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, 0);
  };

  const handleSendMessage = () => {
    if (!messageInput.trim()) return;

    const trimmedMessage = messageInput.trim();
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      const messageData = {
        receiver_id: receiverId,
        message: trimmedMessage,
      };
      socketRef.current.send(JSON.stringify(messageData));
      setMessageInput("");
    } else {
      message.error("WebSocket connection is not open.");
    }
    setMessages((prevMessages) => [
      ...prevMessages,
      { content: trimmedMessage, username: username },
    ]);

    scrollToBottomAfterNewMessage();
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
      <div className={styles.onlineUsersAndSignOutContainer}>
        <OnlineUsers
          className={styles.onlineUsers}
          setMessageInput={setMessageInput}
          setReceiverId={setReceiverId}
          setReceiverUsername={setReceiverUsername}
          setMessages={setMessages}
          setHasLoadedMessages={setHasLoadedMessages}
        />
        <button onClick={handleSignOut} className={styles.signOutButton}>
          Sign out
        </button>
      </div>

      <div className={styles.chatContainer}>
        {receiverUsername && (
          <div className={styles.receiverUsername}>{receiverUsername}</div>
        )}

        <div
          className={styles.messages}
          onScroll={handleScroll}
          ref={messagesContainerRef}
        >
          {messages &&
            messages.map((msg, index) => (
              <div
                className={`${styles.messageContent} ${
                  msg.username === username
                    ? styles.myMessage
                    : styles.otherMessage
                }`}
                key={index}
              >
                <span className={styles.username}>{msg.username}:</span>
                <span className={styles.messageText}>{msg.content}</span>
              </div>
            ))}
          <div ref={messagesEndRef} />
        </div>
        {receiverId !== 0 && (
          <div className={styles.inputContainer}>
            <input
              type="text"
              value={messageInput}
              onKeyUp={(e) => e.key === "Enter" && handleSendMessage()}
              onChange={(e) => setMessageInput(e.target.value)}
              placeholder="Enter your message"
              className={styles.input}
            />
            <button onClick={handleSendMessage} className={styles.button}>
              Send
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
