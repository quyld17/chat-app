import { message } from "antd";
import postMethodAPI from "../methods/post-method-api";

export default function handleSignUpAPI(username, password) {
  return new Promise((resolve, reject) => {
    const credentials = {
      username,
      password,
    };

    const endpoint = "/sign-up";

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
