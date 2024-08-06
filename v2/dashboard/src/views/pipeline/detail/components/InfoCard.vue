<script src="../../../../router/index.js"></script>
<template>
  <el-card style="margin-bottom:20px;" shadow="always">
    <div slot="header" class="clearfix">
      <span>{{ $t('pipe_detail.runStatus') }}</span>
      <el-switch
        v-model="bitem.pipeline.status"
        style="float: right; padding: 3px 0;display: block"
        active-color="#13ce66"
        inactive-color="#ff4949"
        active-value="run"
        inactive-value="stop"
        @change="handleModifyStatus"
      />
    </div>

    <el-card style="margin-bottom: 18px" shadow="never">
      <div class="user-profile">
        <div class="box-center">
          <div class="user-name text-center">{{ bitem.pipeline.aliasName }}</div>
          <div class="user-role text-center text-muted">{{ bitem.pipeline.name }}</div>
        </div>
      </div>
    </el-card>

    <el-card style="margin-bottom: 18px">
      <div class="user-bio">
        <div class="user-education user-bio-section">
          <div class="user-bio-section-header"><svg-icon icon-class="education" /><span>{{ $t('pipe_detail.remark') }}</span></div>
          <div class="user-bio-section-body"><div class="text-muted">{{ bitem.pipeline.remark }}</div></div>
        </div>
      </div>
    </el-card>
    <el-card>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.createTime') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <span>{{ dateFormat(bitem.pipeline.create_time) }}</span>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlMode') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <div v-for="item in $t('pipeline.mysqlModeOptions')" class="text-muted">
            <span v-if="item.key === bitem.pipeline.mysql.mode">{{ item.value }}</span>
          </div>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.output.sender.type') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <div v-for="item in $t('pipeline.output.sender.typeOptions')" class="text-muted">
            <span v-if="item.key === bitem.pipeline.output.sender.type">{{ item.value }}</span>
          </div>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlAddress') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          {{ bitem.pipeline.mysql.address }}
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlPort') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <span>{{ bitem.pipeline.mysql.port }}</span>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlUser') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <span>{{ bitem.pipeline.mysql.user }}</span>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlServerId') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <span>{{ bitem.pipeline.mysql.server_id }}</span>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-bottom: 18px">
        <el-col :span="30">
          <span>{{ $t('pipeline.mysqlFlavor') }}</span>
        </el-col>
        <el-col :span="70" style="float: right">
          <span>{{ bitem.pipeline.mysql.flavor }}</span>
        </el-col>
      </el-row>
    </el-card>
  </el-card>
</template>

<script>

import { dateFormat } from '@/utils/my-time-format'
import { fetchUpdateStatus} from "@/api/pipeline";

export default {
  components: {},
  props: {
    bitem: {
      type: Object,
      default: () => {
        return {
          pipeline: {
            type: Object,
            default: () => {
              return {
                status: '',
                mysql: {
                  type: Object,
                  default: () => {
                    return {
                      address: ''
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  methods: {
    dateFormat(val) {
      return dateFormat(val);
    },
    handleModifyStatus() {
      let req = {
        name: this.bitem.pipeline.name,
        status: this.bitem.pipeline.status
      }

      fetchUpdateStatus(req).then(response => {
        this.$notify({
          title: 'success',
          message: 'Success',
          type: 'success',
          duration: 2000
        })
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.box-center {
  margin: 0 auto;
  display: table;
}

.text-muted {
  color: #777;
}

.user-profile {
  .user-name {
    font-weight: bold;
  }

  .box-center {
    padding-top: 10px;
  }

  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }

  .box-social {
    padding-top: 30px;

    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }

  .user-follow {
    padding-top: 20px;
  }
}

.user-bio {
  margin-top: 20px;
  color: #606266;

  span {
    padding-left: 4px;
  }

  .user-bio-section {
    font-size: 14px;
    padding: 15px 0;

    .user-bio-section-header {
      border-bottom: 1px solid #dfe6ec;
      padding-bottom: 10px;
      margin-bottom: 10px;
      font-weight: bold;
    }
  }
}
</style>
