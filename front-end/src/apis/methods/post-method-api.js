import { GetToken } from "../../components/jwt/index";
import { ServerURL } from "../urls";
export default function postMethodAPI(
  credentials,
  endpoint,
  successCallback,
  errorCallback
) {
  const baseURL = ServerURL();
  const token = GetToken();

  fetch(baseURL + endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: token ? `${token}` : undefined,
    },
    body: JSON.stringify(credentials),
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
