import request from '@/utils/request'
import settings from '@/settings'

export function fetchGet(data) {
  return request({
    url: settings.host.api + '/api/cluster/get',
    method: 'get',
    params: data
  })
}

export function fetchRegisterList(data) {
  return request({
    url: settings.host.api + '/api/cluster/list/register',
    method: 'get',
    params: data
  })
}

export function fetchElectionList(data) {
  return request({
    url: settings.host.api + '/api/cluster/list/election',
    method: 'get',
    params: data
  })
}
