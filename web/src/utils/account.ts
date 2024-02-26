import axios from 'axios'

export async function login(username: string, password: string) {
  const res = await axios.request({
    method: 'POST',
    url: '/api/login',
    data: {
      username,
      password,
    },
  })
  if (res.status >= 300) {
    const { code, message } = res.data
    throw new Error(`login fail: ${code} - ${message}`)
  }
}

export async function logout() {
  const res = await axios.request({
    method: 'POST',
    url: '/api/logout',
    data: {},
  })
  if (res.status >= 300) {
    const { code, message } = res.data
    throw new Error(`logout fail: ${code} - ${message}`)
  }
}

export interface IMe {
  username: string
}

export async function me(): Promise<IMe> {
  const res = await axios.request({
    method: 'GET',
    url: '/api/me',
    data: {},
  })
  if (res.status >= 300) {
    const { code, message } = res.data
    throw new Error(`logout fail: ${code} - ${message}`)
  }

  return res.data
}
