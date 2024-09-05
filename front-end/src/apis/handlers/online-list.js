import { message } from "antd";
import getMethodAPI from "../methods/get-method-api";

export default function handleGetOnlineListAPI() {
  return new Promise((resolve, reject) => {
    const endpoint = "/online-list";
    getMethodAPI(
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
