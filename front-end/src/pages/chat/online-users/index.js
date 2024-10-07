import { useEffect, useState } from "react";
import handleGetOnlineListAPI from "../../../apis/handlers/online-list";
import styles from "./styles.module.css";
import { handleGetChatHistoryAPI } from "@/apis/handlers/chat";

export const OnlineUsers = ({
  setMessageInput,
  setReceiverId,
  setReceiverUsername,
  setMessages,
  setHasLoadedMessages,
}) => {
  const [list, setList] = useState([]);

  useEffect(() => {
    const fetchOnlineUsers = () => {
      handleGetOnlineListAPI().then((data) => {
        setList(data);
      });
    };

    fetchOnlineUsers();

    const intervalId = setInterval(() => {
      fetchOnlineUsers();
    }, 5000);

    return () => clearInterval(intervalId);
  }, []);

  const handleClick = (userId, username) => {
    handleGetChatHistoryAPI(userId, 0).then((data) => {
      setMessages(data);
      setHasLoadedMessages(false);
    });

    setMessageInput("");
    setReceiverId(userId);
    setReceiverUsername(username);
  };

  return (
    <div className={styles.container}>
      <h2 className={styles.headline}>Online Users</h2>
      <ul style={{ listStyleType: "none" }} className={styles.usersList}>
        {list.map((user, index) => (
          <li
            onClick={() => handleClick(user.id, user.username)}
            className={styles.user}
            key={index}
          >
            {user.username}
          </li>
        ))}
      </ul>
    </div>
  );
};
