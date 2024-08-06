import request from '@/utils/request'
import settings from '@/settings'

export function fetchList(data) {
  return request({
    url: settings.host.api + '/api/pipeline/list',
    method: 'get',
    params: data
  })
}

export function fetchGet(data) {
  return request({
    url: settings.host.api + '/api/pipeline/get',
    method: 'get',
    params: data
  })
}

export function fetchCreate(data) {
  return request({
    url: settings.host.api + '/api/pipeline/create',
    method: 'post',
    data: data
  })
}

export function fetchUpdate(data) {
  return request({
    url: settings.host.api + '/api/pipeline/update',
    method: 'post',
    data: data
  })
}

export function fetchUpdateStatus(req) {
  return request({
    url: settings.host.api + '/api/pipeline/update/status',
    method: 'post',
    data: req
  })
}

export function fetchUpdateMode(req) {
  return request({
    url: settings.host.api + '/api/pipeline/update/mode',
    method: 'post',
    data: req
  })
}

export function fetchDelete(data) {
  return request({
    url: settings.host.api + '/api/pipeline/delete',
    method: 'post',
    data: data
  })
}

export function fetchIsFilter(data) {
  return request({
    url: settings.host.api + '/api/pipeline/is_filter',
    method: 'get',
    params: data
  })
}
