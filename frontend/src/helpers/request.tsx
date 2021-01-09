import {baseUrl} from "../const";

export function sendRequest(endpoint: string, body={}, headers={}, method: string="POST") {
  if (typeof(body) !== "string") {
    body = JSON.stringify(body)
  }
  return fetch(`${baseUrl}${endpoint}`, {
    method,
    body: (method === "GET") ? undefined : body as BodyInit,
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

export function get(endpoint, {headers={}}) {
  return sendRequest(endpoint, null, headers, "GET")
}

export function post(endpoint, {body={}, headers={}}) {
  return sendRequest(endpoint, body, headers, "POST")
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
