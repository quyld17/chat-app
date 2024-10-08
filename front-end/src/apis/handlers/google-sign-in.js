import { message } from "antd";
import postMethodAPI from "../methods/post-method-api";

export default function handleGoogleSignIn(token) {
  return new Promise((resolve, reject) => {
    const endpoint = "/google-sign-in";
    postMethodAPI(
      token,
      endpoint,
      (data) => {
        resolve(data);
      },
      (error) => {
        reject(error);
        message.error(error);
      }
    );
  });
}
