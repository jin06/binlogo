import request from '@/utils/request'
import settings from '@/settings'

export function fetchGet(data) {
  return request({
    url: settings.host.api + '/api/position/get',
    method: 'get',
    params: data
  })
}

export function fetchUpdate(data) {
  return request({
    url: settings.host.api + '/api/position/update',
    method: 'post',
    data: data
  })
}

