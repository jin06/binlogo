<script src="../../../api/node.js"></script>
<script src="../../../lang/zh.js"></script>
<template>
  <el-row>
    <aside>{{ $t('pipe_detail.eventTip') }} </aside>

    <el-table :data="list" border fit highlight-current-row style="width: 100%;margin-bottom: 10px">
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.type')">
        <template slot-scope="scope">
          <span v-if="scope.row.type==='info'"><i class="el-icon-info" style="color: #2ac06d" /> INFO</span>
          <span v-if="scope.row.type==='warn'"><i class="el-icon-warning" style="color: yellow" /> WARNING</span>
          <span v-if="scope.row.type==='error'"><i class="el-icon-error" style="color: red" /> ERROR</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.first_time')">
        <template slot-scope="scope">
          <span>{{ dateFormat(scope.row.first_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.last_time')">
        <template slot-scope="scope">
          <span>{{ dateFormat(scope.row.last_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.message')">
        <template slot-scope="scope">
          <span>{{ scope.row.message }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.count')">
        <template slot-scope="scope">
          <span>{{ scope.row.count }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.node_name')">
        <template slot-scope="scope">
          <span>{{ scope.row.node_name }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('pipe_detail.event_table.node_ip')">
        <template slot-scope="scope">
          <span>{{ scope.row.node_ip }}</span>
        </template>
      </el-table-column>
    </el-table>
  </el-row>
</template>

<script>
import { fetchList } from '@/api/event'
import { dateFormat } from '@/utils/my-time-format'

export default {
  props:{
    pipe_name: '',
    list:[]
  },
  data() {
    return {
      list: this.list,
      loading: false
    }
  },
  created() {
  },
  methods: {
    dateFormat(val) {
      return dateFormat(val);
    }
  }
}
</script>

