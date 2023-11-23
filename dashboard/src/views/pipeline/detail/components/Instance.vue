<script src="../../../../lang/zh.js"></script>

<template>
  <el-card>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-avatar v-if="instance.pipeline_name !== undefined" icon="el-icon-ship" style="background-color: #2ac06d" />
        <el-avatar v-if="instance.pipeline_name === undefined" icon="el-icon-ship" style="background-color: #7a807d" />
      </el-col>
      <el-col :span="6">
        <el-row><span style="font-size:12px">{{ $t('instance.create_time') }}</span></el-row>
        <el-row><span>{{ dateFormat(instance.create_time) }}</span></el-row>
      </el-col>
      <el-col :span="6">
        <el-row><span style="font-size: 12px">{{ $t('instance.pipeline_name') }}</span></el-row>
        <el-row><span>{{ instance.pipeline_name }}</span></el-row>
      </el-col>
      <el-col :span="6">
        <el-row><span style="font-size: 12px">{{ $t('instance.node_name') }}</span></el-row>
        <el-row><span> {{ instance.node_name }}</span></el-row>
      </el-col>
    </el-row>
  </el-card>
</template>
<script>

import { dateFormat } from '@/utils/my-time-format'

export default {
  filters: {},
  props: {
    type: {
      type: String,
      default: 'CN'
    },
    instance: undefined
  },
  data() {
    return {
      list: [],
      loading: false
    }
  },
  created() {
    // this.list.unshift(instance)
  },
  methods: {
    getList() {
      this.loading = true
      this.$emit('create') // for test
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.loading = false
      })
    },
    dateFormat(val) {
      return dateFormat(val);
    }
  }
}
</script>

