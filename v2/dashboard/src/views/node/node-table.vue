<script src="../../api/node.js"></script>
<script src="../../lang/zh.js"></script>
<script src="../../lang/en.js"></script>
<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.name" :placeholder="$t('node_table.name')" style="width: 200px;margin-right: 10px" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.ready" :placeholder="$t('node_table.ready')" clearable class="filter-item" style="width: 130px;margin-right: 10px">
        <el-option v-for="item in $t('node_table.readyOptions')" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 140px;margin-right: 10px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in $t('node_table.sortOptions')" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange"
    >
      <el-table-column :label="$t('node.createTime')" width="auto" align="center">
        <template slot-scope="{row}">
          <span>{{ dateFormat(row.node.create_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('node.name')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.node.name }}</span>
        </template>
      </el-table-column>
      <el-table-column
        :label="$t('node.leader')"
        width="auto"
      >
        <template slot-scope="{row}">
          <el-tag v-if="row.info.role === 'leader'" type="primary" effect="plain">{{ $t('node.roleMap.leader') }}</el-tag>
          <el-tag v-else type="info" effect="plain">{{ $t('node.roleMap.follower') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('node.ip')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.node.ip }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('node.status')" width="auto">
        <template slot-scope="{row}" align="middle">
          <el-tag v-if="row.status.ready === true" type="primary" effect="dark">{{ $t('node_table.statusMap.ready.yes') }}</el-tag>
          <el-popover v-if="row.status.ready === false" trigger="hover" placement="top">
            <p v-if="row.status.network_unavailable === true">{{ $t('node_table.statusMap.network_unavailable.yes') }}</p>
            <div slot="reference" class="name-wrapper">
              <el-tag type="danger">
                {{ $t('node_table.statusMap.ready.no') }}
              </el-tag>
            </div>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column :label="$t('node.version')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.node.version }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('capacity.cpuCores')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.capacity.cpu_cores }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('capacity.memory')" width="auto">
        <template slot-scope="{row}">
          <span>{{ (row.capacity.memory/(1024*1024*1024)).toFixed(2) }}Gi</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('capacity.cpuUsage')" width="auto">
        <template slot-scope="{row}">
          <el-progress v-if="row.status.ready === true" :percentage="row.capacity.cpu_usage" :color="colors" />
        </template>
      </el-table-column>
      <el-table-column :label="$t('capacity.memoryUsage')" width="auto">
        <template slot-scope="{row}">
          <el-progress v-if="row.status.ready === true" :percentage="row.capacity.memory_usage" :color="colors" />
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { fetchList } from '@/api/node'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  filters: {
  },
  data() {
    return {
      colors: [
        {color: '#5cb87a', percentage: 20},
        {color: '#1989fa', percentage: 40},
        {color: '#6f7ad3', percentage: 60},
        {color: '#e6a23c', percentage: 80},
        {color: '#f56c6c', percentage: 100}
      ],
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        title: undefined,
        type: undefined,
        sort: '+id',
        ready: undefined,
        name: undefined
      },
      importanceOptions: [1, 2, 3],
      sortOptions: [{ label: 'ID Ascending', key: '+id' }, { label: 'ID Descending', key: '-id' }]
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    sortChange(data) {
      const { prop, order } = data
      if (prop === 'id') {
        this.sortByID(order)
      }
    },
    sortByID(order) {
      if (order === 'ascending') {
        this.listQuery.sort = '+id'
      } else {
        this.listQuery.sort = '-id'
      }
      this.handleFilter()
    },
    dateFormat(value){
      const t = new Date(value);
      return t.getFullYear()+"-"+(t.getMonth()+1)+"-"+t.getDate()+" "+t.getHours()+":"+t.getMinutes()+":"+t.getSeconds();
    },
    roleMap(mapping,value) {
      console.log(mapping);
      if (value === 'leader') {
        return mapping.leader;
      }else {
        return mapping.follower;
      }
    }
  }
}
</script>
