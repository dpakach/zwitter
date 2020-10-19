import {baseUrl} from "../const";

export function sendRequest(endpoint, body={}, headers={}, method="POST") {
  if (typeof(body) !== "string") {
    body = JSON.stringify(body)
  }
  return fetch(`${baseUrl}${endpoint}`, {
    method,
    body,
    headers
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

export function sendFileUploadRequest(endpoint, file, headers={}) {
  return fetch(`${baseUrl}${endpoint}`, {
    method: "POST",
    body: file,
    headers
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
