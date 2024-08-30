import { useRouter } from "next/router";
import { useState } from "react";

import styles from "./sign-in.module.css";
import handleSignInAPI from "../../apis/handlers/sign-in";

export default function SignIn() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const handleSignIn = (e) => {
    e.preventDefault();
    handleSignInAPI(username, password)
      .then((data) => {
        localStorage.setItem("token", data.token);
        // router.push("/");
        alert("sign in successfull, token: ", data.token);
      })
      .catch((error) => {
        console.log("Error getting delivery address: ", error);
      });
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    setError("");

    try {
      const response = await fetch("/api/signin", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      const result = await response.json();

      if (response.ok) {
        alert("Sign in successful!");
      } else {
        // Handle error
        setError(result.message || "Something went wrong.");
      }
    } catch (error) {
      setError("Network error. Please try again.");
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.formContainer}>
        <h2>Sign In</h2>
        {error && <div className={styles.error}>{error}</div>}
        <form onSubmit={handleSignIn}>
          <div className={styles.inputGroup}>
            <label htmlFor="username">Username</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className={styles.inputGroup}>
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit" className={styles.button}>
            Sign In
          </button>
        </form>
      </div>
    </div>
  );
}
