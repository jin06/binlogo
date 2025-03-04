
export function successNotify(vm, msg) {
  vm.$notify({
    title: 'Success',
    message: msg || 'Operation Success',
    type: 'success',
    duration: 2000
  })
}

export function errorNotify(vm, msg) {
  vm.$notify({
    title: 'Request Error',
    message: msg || 'Operation failed',
    type: 'error',
    duration: 5000
  })
}

export function notify(vm, response) {
  if (response.code !== 20000) {
    vm.$notify({
      title: 'Request Error',
      message: response.msg,
      type: 'error',
      duration: 5000
    })
  } else {
    vm.$notify({
      title: 'Success',
      message: 'Operation Success',
      type: 'success',
      duration: 2000
    })
  }
}
