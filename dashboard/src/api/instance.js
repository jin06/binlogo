import request from '@/utils/request'
import settings from '@/settings'

export function fetchList(data) {
  return request({
    url: settings.host.api + '/api/instance/list',
    method: 'get',
    params: data
  })
}
