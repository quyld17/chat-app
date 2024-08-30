import getJWT from "../../components/jwt";

export default function postMethodAPI(
  credentials,
  endpoint,
  successCallback,
  errorCallback
) {
  const { baseURL, token } = getJWT();

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
