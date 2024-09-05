import { useEffect, useState } from "react";

export const OnlineList = ({ token }) => {
  const [list, setList] = useState([]);
  


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
