import { fetchURL, fetchJSON, removePrefix } from './utils'

export async function getHash(hash) {
  return fetchJSON(`/api/public/share/${hash}`)
}

export async function get(url) {
  url = removePrefix(url)
  return fetchJSON(`/api/share${url}`)
}

export async function remove(hash) {
  const res = await fetchURL(`/api/share/${hash}`, {
    method: 'DELETE'
  })

  if (res.status !== 200) {
    throw new Error(res.status)
  }
}

export async function create(url, expires = '', unit = 'hours') {
  url = removePrefix(url)
  url = `/api/share${url}`
  if (expires !== '') {
    url += `?expires=${expires}&unit=${unit}`
  }

  return fetchJSON(url, {
    method: 'POST'
  })
}

export async function shareWithUsers(url, users, permissions){
  let body = { permissions, users}
  body = JSON.stringify(body)
  console.log("hereee")
  url = "/api/share/"
  return fetchJSON(url, {
    method: 'POST',
    body
  })
}

export async function getShareableLink(url, permissions){
  return fetchURL(`/api/share/link/${url}?permission=${permissions}`)
}
