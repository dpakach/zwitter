import {baseUrl} from "../const";

export function sendRequest(endpoint, body={}, headers={}) {
  if (typeof(body) !== "string") {
    body = JSON.stringify(body)
  }
  return fetch(`${baseUrl}${endpoint}`, {
    method: "POST",
    body,
    mode: 'cors',
    cache: 'no-cache',
    credentials: 'same-origin',
    headers: {
      Accept: "*/*",
      ...headers
    },
    redirect: 'follow',
    referrerPolicy: 'no-referrer',
    body
  })
  .then(res => {
    if (!res.ok) {
      return res.json().then(data => {
        throw Error("Request failed: " + data.message || "")
      })
    }
    return res
  })
}
