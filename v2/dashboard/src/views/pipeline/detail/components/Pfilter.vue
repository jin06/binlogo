<script src="../../../../api/pipeline.js"></script>
<script src="../../../../lang/zh.js"></script>

<template>
  <el-row :gutter="20">
    <el-col :span="10">
      <el-row>
        <aside>{{ $t('pipe_detail.filterTip') }}</aside>
      </el-row>
      <el-row :model="validData">
        <el-form :label-position="top" @keyup.enter.native="validFilter(validData)">
          <el-form-item :label="$t('pipe_detail.validRule')">
            <el-input v-model="validData.rule" @input="clearValidResult(validData)" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="validFilter(validData)">{{ $t('pipe_detail.button_valid') }}</el-button>
          </el-form-item>
        </el-form>
      </el-row>
      <el-row>
        <span v-if="validResult === true" :model="validResult">{{ validData.rule }} {{ $t('pipe_detail.validTipTrue') }}</span>
        <span v-if="validResult === false" :model="validResult">{{ validData.rule }} {{ $t('pipe_detail.validTipFalse') }}</span>
      </el-row>
    </el-col>
    <el-col :span="10">
      <el-table :data="list" border fit highlight-current-row width="auto">
        <el-table-column width="auto" align="center" :label="$t('pipeline.filter.type')">
          <template slot-scope="{row}">
            <span v-if="row.type === 'white' ">{{ $t('pipeline.filterValues.white') }}</span>
            <span v-if="row.type === 'black' ">{{ $t('pipeline.filterValues.black') }}</span>
          </template>
        </el-table-column>
        <el-table-column width="auto" align="center" :label="$t('pipeline.filter.rule')">
          <template slot-scope="{row}">
            <span>{{ row.rule }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-col>
  </el-row>

</template>
<script>

import { dateFormat } from '@/utils/my-time-format'
import { fetchIsFilter } from "@/api/pipeline";

export default {
  filters: {},
  props: {
    type: {
      type: String,
      default: 'CN'
    },
    filter: undefined,
    list: [],
    pipe_name:''
  },
  data() {
    return {
      loading: false,
      validData:{},
      validResult:''
    }
  },
  created() {
    // this.list.unshift(instance)
  },
  methods: {
    dateFormat(val) {
      return dateFormat(val);
    },
    clearValidResult() {
      this.validResult = '';
    },
    validFilter(val) {
      const req = {
        name: this.pipe_name,
        rule: val.rule
      }
      this.validResult = '';
      fetchIsFilter(req).then(response => {
        this.validResult = response.data
      })
    }
  }
}
</script>

