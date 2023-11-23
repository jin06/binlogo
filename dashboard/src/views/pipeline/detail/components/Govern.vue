<script src="../../../../lang/zh.js"></script>
<template>
  <el-row :gutter="20">
    <el-col :span="24">
      <el-row>
        <aside>{{ $t('pipe_detail.positionTip') }} </aside>
      </el-row>
      <el-card>
        <el-form label-position="left" label-width="auto" :rules="rules">
          <el-form-item :label="$t('pipeline.mysqlFlavor')">
            <el-select v-model="govern_form.mode">
              <el-option v-for="item in $t('pipeline.mysqlModeOptions')" :key="item.key" :label="item.value" :value="item.key" />
            </el-select>
          </el-form-item>
          <el-form-item label="Binlog File">
            <el-input v-model="govern_form.position.binlog_file" placeholder="mysql-bin.012000" :disabled="govern_form.mode==='gtid'" />
          </el-form-item>
          <el-form-item label="Binlog Position">
            <el-input v-model.number="govern_form.position.binlog_position" type="number" placeholder="25110009" :disabled="govern_form.mode==='gtid'" />
          </el-form-item>
          <el-form-item label="Gtid Set">
            <el-input v-model="govern_form.position.gtid_set" placeholder="045c649a-408d-11ec-ae21-0242ac110006:1-38" :disabled="govern_form.mode==='position'" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updatePosition()">{{ $t('global.submit') }}</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </el-col>
  </el-row>
</template>
<script>

import { dateFormat } from '@/utils/my-time-format'
import { fetchUpdate } from "@/api/position";
import { fetchUpdateMode } from "@/api/pipeline";

export default {
  filters: {},
  props: {
    type: {
      type: String,
      default: 'CN'
    },
    instance:{},
    govern_form:{
      mode:'',
      position:{
        binlog_position: {
          type:Number
        }
      }
    },
    pipe_name:'',
    bitem: {}
  },
  data() {
    return {
      list: [],
      loading: false,
      rules: {
      }
    }
  },
  created() {
  },
  methods: {
    dateFormat(val) {
      return dateFormat(val);
    },
    updatePosition() {
      if (this.govern_form.mode !== this.bitem.mode) {
        fetchUpdateMode({name: this.bitem.pipeline.name, mode: this.govern_form.mode}).then(response => {
          this.$message({
            type: 'success',
            message: 'Update binlog mode success!'
          })
          this.bitem.pipeline.mysql.mode = this.govern_form.mode
        })
        let req = {
          mode: this.govern_form.mode,
          position: this.govern_form.position
        }
        req.position.pipeline_name = this.bitem.pipeline.name
        fetchUpdate(req).then(response => {
          this.$message({
            type: 'success',
            message: 'Update position success!'
          });
        })
      }
    }
  }
}
</script>

