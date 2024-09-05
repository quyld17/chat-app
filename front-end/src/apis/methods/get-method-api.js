import { GetToken } from "../../components/jwt/index";
import { ServerURL } from "../urls";

export default function getMethodAPI(endpoint, successCallback, errorCallback) {
  const baseURL = ServerURL();
  const token = GetToken();

  fetch(baseURL + endpoint, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Authorization: token ? `${token}` : undefined,
    },
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.message) {
        errorCallback(data.message);
      } else {
        successCallback(data);
      }
    })
    .catch((error) => {
      errorCallback(error.message);
    });
}
