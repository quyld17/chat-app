import { message } from "antd";
import postMethodAPI from "../methods/post-method-api";

export default function handleValidateGoogleToken(id_token) {
  return new Promise((resolve, reject) => {
    const endpoint = "/google-sign-in";
    const credentials = {
      id_token,
    };
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
