import { jwtDecode } from "jwt-decode";

export function GetToken() {
  const token = localStorage.getItem("token");
  return token;
}

export function DecodeToken(token) {
  const decodedToken = jwtDecode(token);
  return decodedToken;
}

export function CheckTokenExpireTime(handleSignOut, setUsername, token) {
  const expTime = token.exp;
  const currentTime = Date.now() / 1000;

  if (currentTime > expTime) {
    message.info("Session expired! Please sign in to continue");
    handleSignOut();
  } else {
    setUsername(token.username);
  }
}
