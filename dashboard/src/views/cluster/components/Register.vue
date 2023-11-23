<script src="../../../lang/zh.js"></script>
<template>
  <el-row>
    <aside style="margin-bottom: 50px;margin-top: 20px">{{ $t('cluster.registerTip') }}</aside>
    <el-table :data="list" border fit highlight-current-row style="width: 100%;margin-bottom: 10px">
      <el-table-column width="auto" align="center" :label="$t('register.create_time')">
        <template slot-scope="scope">
          <span>{{ dateFormat(scope.row.create_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('register.name')">
        <template slot-scope="scope">
          <span>{{ scope.row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('register.ip')">
        <template slot-scope="scope">
          <span>{{ scope.row.ip }}</span>
        </template>
      </el-table-column>
      <el-table-column width="auto" align="center" :label="$t('register.version')">
        <template slot-scope="scope">
          <span>{{ scope.row.version }}</span>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </el-row>
</template>

<script>
import { fetchRegisterList } from '@/api/cluster'
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
      total: 0,
      list: null,
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
    this.listQuery.page = 1
    this.getList()
  },
  methods: {
    getList() {
      this.loading = true
      this.$emit('create') // for test
      fetchRegisterList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.loading = false
      })
    },
    dateFormat(val) {
      return dateFormat(val);
    }
  }
}
</script>

