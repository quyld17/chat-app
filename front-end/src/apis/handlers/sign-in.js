import postMethodAPI from "../methods/post-method-api";
import { message } from "antd";

export default function handleSignInAPI(username, password) {
  return new Promise((resolve, reject) => {
    const credentials = {
      username,
      password,
    };

    const endpoint = "/sign-in";
    postMethodAPI(
      credentials,
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
