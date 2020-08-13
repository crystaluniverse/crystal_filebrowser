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

export async function shareWithUsers(path, userList){
  const body = JSON.stringify(userList)
  const url = `/api/share${path}`
  return fetchJSON(url, {
    method: 'POST',
    body
  })
}
// export async function deleteUserAccess(path, user){
//   let body = { user}
//   body = JSON.stringify(body)
//   const url = `/api/share${path}`
//   return fetchJSON(url, {
//     method: 'POST',
//     body
//   })
// }
export async function deleteAllShares(path){
  return fetchURL(`/api/share${path}`, {
    method: 'DELETE'
  })
}
export async function listUserpermissions(path){
  return fetchJSON(`/api/share${path}`)
}



export async function getShareableLink(url, permissions){
  return fetchURL(`/api/share/link${url}?permission=${permissions}`)
}
export async function deleteSharableLink(uuid){
  return fetchURL(`/api/share/link${uuid}`, {
    method: 'DELETE'
  })
}
export async function listLinks(path){
  return fetchJSON(`/api/share/links${path}}`)
}