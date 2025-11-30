import axios from 'axios'

export const api = {
  getRoutes() {
    return axios.get('/admin/routes')
  },
  saveRoute(payload: any) {
    return axios.post('/admin/routes', payload)
  },
  updateRoute(payload: any) {
    return axios.post('/admin/routes/update', payload)
  },
  tailLogs(n: number) {
    return axios.get('/admin/logs', { params: { tail: n } })
  }
}

