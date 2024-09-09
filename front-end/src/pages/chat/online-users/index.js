import { useEffect, useState } from "react";
import handleGetOnlineListAPI from "../../../apis/handlers/online-list";
import styles from "./styles.module.css";

export const OnlineUsers = ({ handleStartChat }) => {
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

  const handleClick = (userId) => {
    handleStartChat(userId);
  };

  return (
    <div className={styles.container}>
      <h2 className={styles.headline}>Online Users</h2>
      <ul style={{ listStyleType: "none" }} className={styles.usersList}>
        {list.map((user, index) => (
          <li
            onClick={() => handleClick(user.id)}
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
