import axios from "axios";

const client = axios.create({
  baseURL: "http://localhost:5000"
});

function Request(options) {
  const onSuccess = function(response) {
    console.debug("Request Successful!", response);

    return response.data;
  };

  const onError = function(error) {
    console.error("Request Failed:", error.config);

    if (error.response) {
      // Request was made but server responded with something
      // other than 2xx
      console.error("Status:", error.response.status);
      console.error("Data:", error.response.data);
      console.error("Headers:", error.response.headers);

      // Clear the JWT if we have an old one and is being rejected.
      if (error.response.status === 401) {
        localStorage.removeItem("authorization");
      }
    } else {
      // Something else happened while setting up the request
      // triggered the error
      console.error("Error Message:", error.message);
    }

    return Promise.reject(error.response || error.message);
  };

  return client(options)
    .then(onSuccess)
    .catch(onError);
}

function GetCompetition(id) {
  return Request({
    method: "get",
    url: `/competition/${id}`
  });
}

const HttpService = {
  Request,
  GetCompetition
};

export default HttpService;

// vim: set ts=2 sw=2 et:
