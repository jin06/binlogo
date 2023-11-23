<script src="../../../api/node.js"></script>
<script src="../../../lang/zh.js"></script>
<template>
  <el-row>
    <aside style="margin-bottom: 50px;margin-top: 20px">{{$t('cluster.instanceTip')}}</aside>
    <el-table :data="list" border fit highlight-current-row style="width: 100%;margin-bottom: 10px">
      <el-table-column width="auto" align="center" :label="$t('instance.create_time')">
        <template slot-scope="scope">
          <span>{{ dateFormat(scope.row.create_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('instance.pipeline_name')">
        <template slot-scope="scope">
          <span>{{ scope.row.pipeline_name }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('instance.node_name')">
        <template slot-scope="scope">
          <span>{{ scope.row.node_name }}</span>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </el-row>
</template>

<script>
import { fetchList } from '@/api/instance'
import { dateFormat } from '@/utils/my-time-format'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'info',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  props: {
    type: {
      type: String,
      default: 'CN'
    }
  },
  data() {
    return {
      list: null,
      total: 0,
      listQuery: {
        page: 1,
        limit: 10,
        type: this.type,
        sort: '+id'
      },
      loading: false
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.loading = true
      this.$emit('create') // for test
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.loading = false
        this.total = response.data.total
      })
    },
    dateFormat(val) {
      return dateFormat(val);
    }
  }
}
</script>

