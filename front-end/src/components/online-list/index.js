import { useEffect, useState } from "react";
import handleGetOnlineListAPI from "../../apis/handlers/online-list";

export const OnlineList = () => {
  const [list, setList] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchOnlineUsers = () => {
      handleGetOnlineListAPI()
        .then((data) => {
          setList(data);
        })
        .catch((err) => {
          console.error("Error getting online users: ", err);
          setError("Error fetching online users");
        });
    };

    fetchOnlineUsers();

    const intervalId = setInterval(() => {
      fetchOnlineUsers();
    }, 5000);

    return () => clearInterval(intervalId);
  }, []);

  return (
    <div>
      <h2>Online Users</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <ul>
        {list.map((user, index) => (
          <li key={index}>{user.username}</li>
        ))}
      </ul>
    </div>
  );
};
