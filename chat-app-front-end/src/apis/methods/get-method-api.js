import getJWT from "../../components/jwt";

export default function getMethodAPI(endpoint, successCallback, errorCallback) {
  const { baseURL, token } = getJWT();

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
