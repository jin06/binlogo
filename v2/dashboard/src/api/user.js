import request from '@/utils/request'
import settings from '@/settings'

export function login(data) {
  return request({
    url: settings.host.api + '/api/user/login',
    method: 'post',
    data
  })
}

export function getInfo(token) {
  return request({
    url: settings.host.api + '/api/user/info',
    method: 'get',
    params: { token }
  })
}

export function logout() {
  return request({
    url: settings.host.api + '/api/user/logout',
    method: 'post'
  })
}

export function authType(query) {
  return request({
    url: settings.host.api + '/api/user/auth_type',
    method: 'get',
    params: query
  })
}
