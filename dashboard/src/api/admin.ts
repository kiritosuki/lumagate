import axios, { AxiosRequestConfig } from 'axios'

const gw = 'http://localhost:8080'

async function tryBoth<T>(relativeUrl: string, absUrl: string, cfg?: AxiosRequestConfig) {
  try {
    return await axios.request<T>({ url: relativeUrl, ...cfg })
  } catch {
    return await axios.request<T>({ url: absUrl, ...cfg })
  }
}

export const api = {
  getRoutes() {
    return tryBoth<any[]>('/admin/routes', `${gw}/admin/routes`)
  },
  saveRoute(payload: any) {
    return tryBoth('/admin/routes', `${gw}/admin/routes`, { method: 'POST', data: payload, headers: { 'Content-Type': 'application/json' } })
  },
  updateRoute(payload: any) {
    return tryBoth('/admin/routes/update', `${gw}/admin/routes/update`, { method: 'POST', data: payload, headers: { 'Content-Type': 'application/json' } })
  },
  tailLogs(n: number) {
    return tryBoth<string[]>('/admin/logs', `${gw}/admin/logs`, { params: { tail: n } })
  }
}
