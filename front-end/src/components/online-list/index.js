import { useEffect, useState } from "react";

export const OnlineList = ({ token }) => {
  const [list, setList] = useState([]);

  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8080/ws/online-list?token=${token}`
    );

    ws.onmessage = function (event) {
      const onlineUsers = JSON.parse(event.data);
      setList(onlineUsers);
    };

    return () => {
      ws.close();
    };
  }, [token]);

  return (
    <div>
      <h2>Online Users</h2>
      <ul>
        {list.map((user, index) => (
          <li key={index}>{user.username}</li>
        ))}
      </ul>
    </div>
  );
};
