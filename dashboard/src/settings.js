switch (process.env.NODE_ENV) {
  case 'production':
    module.exports = require('./settings.pro')
    break
  case 'staging':
    module.exports = require('./settings.staging')
    break
  case 'development':
    module.exports = require('./settings.dev')
    break
  default:
    module.exports = require('./settings.dev')
}
