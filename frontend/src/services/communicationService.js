
import api from './api'

export const communicationService = {
  getAllLogs() {
    return api.get('/communication_logs')
  },
  
  getLogById(logId) {
    return api.get(`/communication_logs/${logId}`)
  },
  
  getLogsByClientId(clientId) {
    return api.get(`/communication_logs/client/${clientId}`)
  }
}

export default communicationService
