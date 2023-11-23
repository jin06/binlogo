<script src="../../router/index.js"></script>
<script src="../../api/pipeline.js"></script>
<script src="../../lang/zh.js"></script>
<script src="../../lang/en.js"></script>
<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.name" :placeholder="$t('pipeline_table.name')" style="width: 200px;margin-right: 10px" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.status" :placeholder="$t('pipeline_table.status')" clearable class="filter-item" style="width: 130px; margin-right: 10px">
        <el-option v-for="item in $t('pipeline_table.statusOptions')" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 140px;margin-right: 10px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in $t('pipeline_table.sortOptions')" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        {{ $t('table.add') }}
      </el-button>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%"
      @sort-change="sortChange"
    >
      <el-table-column :label="$t('pipeline.createTime')" width="auto" align="center">
        <template slot-scope="{row}">
          <span>{{ dateFormat(row.pipeline.create_time) }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('pipeline.name')" prop="pipeline.name" width="auto">
        <template slot-scope="{row}">
          <router-link :to="'/pipeline/pipeline-detail/'+row.pipeline.name" class="link-type">
            <span>{{ row.pipeline.name }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('pipeline.aliasName')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.pipeline.aliasName }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('pipeline.status')" width="auto">
        <template slot-scope="{row}">
          <el-tag v-if="row.pipeline.status==='run' && row.info.instance.node_name !== ''" key="run" type="success" effect="dark">
            {{ $t('pipeline.statusValues.started') }}
          </el-tag>
          <el-tag v-if="row.pipeline.status==='run' && row.info.instance.node_name === ''" key="run" type="success" effect="">
            {{ $t('pipeline.statusValues.starting') }}
          </el-tag>
          <el-tag v-if="row.pipeline.status === 'stop' && row.info.instance.node_name === ''" key="stop" type="info" effect="dark">
            {{ $t('pipeline.statusValues.stopped') }}
          </el-tag>
          <el-tag v-if="row.pipeline.status === 'stop' && row.info.instance.node_name !== ''" key="stop" type="info" effect="">
            {{ $t('pipeline.statusValues.stopping') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('pipeline.bindNode')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.info.bind_node.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('pipeline.runNode')" width="auto">
        <template slot-scope="{row}">
          <span>{{ row.info.instance.node_name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.actions')" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            {{ $t('pipeline_table.edit') }}
          </el-button>
          <el-button v-if="row.pipeline.status!=='run'" size="mini" type="success" @click="handleModifyStatus(row,'run')">
            {{ $t('pipeline_table.run') }}
          </el-button>
          <el-button v-if="row.pipeline.status!=='stop'" size="mini" @click="handleModifyStatus(row,'stop')">
            {{ $t('pipeline_table.stop') }}
          </el-button>
          <el-button size="mini" type="danger" @click="handleDelete(row,$index)">
            {{ $t('pipeline_table.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="right" label-width="auto" style="width: 80%; ">
        <el-divider content-position="center">Pipeline</el-divider>
        <el-form-item :label="$t('pipeline.name')" prop="pipeline.name">
          <el-input v-model="temp.pipeline.name" :placeholder="$t('pipeline_table.pipeline.name')" :disabled="dialogStatus!=='create'" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.aliasName')">
          <el-input v-model="temp.pipeline.aliasName" :placeholder="$t('pipeline_table.pipeline.aliasName')" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.remark')">
          <el-input v-model="temp.pipeline.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" :placeholder="$t('pipeline_table.input.pleaseInput')" />
        </el-form-item>
        <el-divider content-position="center">MySQL</el-divider>
        <el-form-item :label="$t('pipeline.mysqlFlavor')" prop="pipeline.mysql.flavor">
          <el-radio-group v-model="temp.pipeline.mysql.flavor" size="small">
            <el-radio-button v-for="item in $t('pipeline.mysqlFlavorOptions')" :key="item.key" border :label="item.key">
              {{ item.value }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlMode')" prop="pipeline.mysql.mode">
          <el-radio-group v-model="temp.pipeline.mysql.mode" size="small">
            <el-radio-button v-for="item in $t('pipeline.mysqlModeOptions')" :key="item.key" border :label="item.key">
              {{ item.value }}
            </el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlAddress')" prop="pipeline.mysql.address">
          <el-input v-model="temp.pipeline.mysql.address" placeholder="127.0.0.1" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlPort')" prop="pipeline.mysql.port">
          <el-input v-model.number="temp.pipeline.mysql.port" placeholder="3306" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlUser')" prop="pipeline.mysql.user">
          <el-input v-model="temp.pipeline.mysql.user" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlPassword')" prop="pipeline.mysql.password">
          <el-input v-model="temp.pipeline.mysql.password" />
        </el-form-item>
        <el-form-item :label="$t('pipeline.mysqlServerId')" prop="pipeline.mysql.server_id">
          <el-input v-model.number="temp.pipeline.mysql.server_id" :placeholder="$t('pipeline_table.server_id')" />
        </el-form-item>
        <el-divider content-position="center">Output</el-divider>
        <el-form-item :label="$t('pipeline.output.sender.type')">
          <el-select v-model="temp.pipeline.output.sender.type" class="filter-item" :placeholder="$t('pipeline_table.select.pleaseSelect')" prop="pipeline.output.sender.type" style="margin-right: 15px">
            <el-option v-for="item in $t('pipeline.output.sender.typeOptions')" :key="item.key" :label="item.value" :value="item.key" />
          </el-select>
          <el-link v-if="temp.pipeline.output.sender.type==='kafka'" target="_blank" type="success" href="https://kafka.apache.org/documentation/#producerconfigs">Kafka Configs Doc</el-link>
          <el-link v-if="temp.pipeline.output.sender.type==='rabbitMQ'" target="_blank" type="success" href="https://www.rabbitmq.com/tutorials/tutorial-five-go.html">Using topic pattern of RabbitMQ</el-link>
          <el-link v-if="temp.pipeline.output.sender.type==='rocketMQ'" target="_blank" type="success" href="https://help.aliyun.com/document_detail/255810.html?spm=5176.rocketmq.help.dexternal.248c7d10NtDnwh">Using pattern of RocketMQ</el-link>
          <div v-if="temp.pipeline.output.sender.type==='rabbitMQ'" class="el-upload__tip">{{ $t('pipeline_table.rabbit.tips') }}</div>
          <div v-if="temp.pipeline.output.sender.type==='redis'" class="el-upload__tip">{{ $t('pipeline_table.redis.tips') }}</div>
          <div v-if="temp.pipeline.output.sender.type==='rocketMQ'" class="el-upload__tip">{{ $t('pipeline_table.rocket.tips') }}</div>
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="brokers">
          <el-input v-model="temp.pipeline.output.sender.kafka.brokers" type="textarea" :placeholder="$t('pipeline_table.kafka.brokers')" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="topic">
          <el-input v-model="temp.pipeline.output.sender.kafka.topic" placeholder="If it is blank, the name of pipeline is used as the topic name" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="acks">
          <el-select v-model="temp.pipeline.output.sender.kafka.require_acks" class="filter-item" :placeholder="$t('pipeline_table.select.pleaseSelect')">
            <el-option
              v-for="item in kafkaAcksOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="enable.idempotence">
          <el-select v-model="temp.pipeline.output.sender.kafka.idepotent" class="filter-item" :placeholder="$t('pipeline_table.select.pleaseSelect')">
            <el-option
              v-for="item in kafkaIdepotent"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="compression.type">
          <el-select v-model="temp.pipeline.output.sender.kafka.compression" class="filter-item" :placeholder="$t('pipeline_table.select.pleaseSelect')">
            <el-option
              v-for="item in kafkaCompressionOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'kafka'" label="retries">
          <el-input v-model.number="temp.pipeline.output.sender.kafka.retries" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'http'" label="API">
          <el-input v-model="temp.pipeline.output.sender.http.api" type="textarea" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'http'" :label="$t('pipeline_table.http.retries')">
          <el-input v-model.number="temp.pipeline.output.sender.http.retries" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rabbitMQ'" label="Exchange Url">
          <el-input v-model="temp.pipeline.output.sender.rabbitMQ.url" type="textarea" placeholder="amqp://guest:guest@localhost:5672/" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rabbitMQ'" label="Exchange Name">
          <el-input v-model="temp.pipeline.output.sender.rabbitMQ.exchange_name" placeholder="If it is blank, the name of pipeline is used as the exchange name" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'redis'" label="Address">
          <el-input v-model="temp.pipeline.output.sender.redis.address" placeholder="localhost:6379" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'redis'" label="User Name">
          <el-input v-model="temp.pipeline.output.sender.redis.username" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'redis'" label="Password">
          <el-input v-model="temp.pipeline.output.sender.redis.password" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'redis'" label="List">
          <el-input v-model="temp.pipeline.output.sender.redis.list" placeholder="If it is blank, the name of pipeline is used as the redis list name " />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rocketMQ'" label="Endpoint">
          <el-input v-model="temp.pipeline.output.sender.rocketMQ.endpoint" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rocketMQ'" label="TopicName">
          <el-input v-model="temp.pipeline.output.sender.rocketMQ.topic_name" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rocketMQ'" label="InstanceId">
          <el-input v-model="temp.pipeline.output.sender.rocketMQ.instance_id" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rocketMQ'" label="AccessKey">
          <el-input v-model="temp.pipeline.output.sender.rocketMQ.access_key" />
        </el-form-item>
        <el-form-item v-if="temp.pipeline.output.sender.type === 'rocketMQ'" label="SecretKey">
          <el-input v-model="temp.pipeline.output.sender.rocketMQ.secret_key" />
        </el-form-item>
        <el-divider content-position="center">Filter: <el-button size="small" @click="addFilter">{{ $t('pipeline_table.filter.addFilter') }}</el-button></el-divider>
        <el-form-item
          v-for="(filter, index) in temp.pipeline.filters"
          :key="'' + index"
          :label="$t('pipeline_table.filter.name') + ' '+ index"
        >
          <el-row :gutter="20">
            <el-col :span="5">
              <el-select v-model="filter.type" :placeholder="$t('pipeline_table.select.pleaseSelect')">
                <el-option :label="$t('pipeline_table.filter.whiteList')" value="white" />
                <el-option :label="$t('pipeline_table.filter.blackList')" value="black" />
              </el-select>
            </el-col>
            <el-col :span="11">
              <el-input v-model="filter.rule" :placeholder="$t('pipeline_table.filter.place')" />
            </el-col>
            <el-col :span="5">
              <el-button @click="removeFilter(filter)">{{ $t('pipeline_table.filter.delete') }}</el-button>
            </el-col>
          </el-row>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          {{ $t('table.cancel') }}
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          {{ $t('table.confirm') }}
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="dialogPvVisible" title="Reading statistics">
      <el-table :data="pvData" border fit highlight-current-row style="width: 100%">
        <el-table-column prop="key" label="Channel" />
        <el-table-column prop="pv" label="Pv" />
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialogPvVisible = false">{{ $t('table.confirm') }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { fetchPv } from '@/api/article'
import { fetchList, fetchUpdateStatus, fetchCreate, fetchUpdate, fetchDelete } from '@/api/pipeline'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  components: { Pagination },
  directives: { waves },
  // filters: {
  //   statusFilter(status) {
  //     const statusMap = {
  //       run: 'success',
  //       stop: 'danger'
  //     }
  //     return statusMap[status]
  //   }
  // },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        // importance: undefined,
        title: undefined,
        type: undefined,
        status: undefined,
        name: undefined,
        sort: '+id'
      },
      // importanceOptions: [1, 2, 3],
      kafkaAcksOptions:[
        { value:0, label:'0 (not wait any response, retries will not take effect)' },
        { value:1, label:'1 (write to kafka local log without awaiting full ack from fllowers.)'},
        { value: -1, label: '-1 (will wait for full set of in-sync replicas to acknowledge the record)'}
      ],
      kafkaCompressionOptions: [
        { value: 0, label:'no compression'},
        { value: 1, label:'gzip'},
        { value: 2, label:'snappy'},
        { value: 3, label:'lz4'},
        { value: 4, label:'zstd'}
      ],
      kafkaIdepotent:[
        { value: true, label:true},
        { value: false, label: false}
      ],
    sortOptions: [{ label: 'Time Ascending', key: '+id' }, { label: 'Time Descending', key: '-id' }],
      // statusOptions: ['published', 'draft', 'deleted'],
      showReviewer: false,
      temp: {
        pipeline: {
          mysql: {
            address: '',
            port: '',
            user: '',
            password: '',
            mode: 'gtid',
            flavor: ''
          },
          remark: '',
          alias_name: '',
          name: '',
          output: {
            sender:{
              type: 'kafka',
              kafka: {
                brokers: '',
                topic: '',
                require_acks: 1,
                compression: 0,
                retries: 3,
                idepotent:false
              },
              stdout: null,
              http: {
                api: '',
                retries: 3
              },
              rabbitMQ: {
                url: '',
                exchange_name: ''
              },
              redis: {
                address:'',
                username:'',
                password:'',
                list:''
              },
              rocketMQ: {
                endpoint:'',
                topic_name: '',
                instance_id: '',
                access_key: '',
                secret_key: ''
              }
            }
          },
          filters: [{
            "type": "black",
            "rule": "mysql"
          }]
        }
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: "Create"
      },
      dialogPvVisible: false,
      pvData: [],
      rules: {
        pipeline:{
          name: [{ required: true, message: "name is required", trigger: 'change'}],
          status:[{ required: true, message: 'field is required', trigger: 'blur'}],
          mysql: {
            flavor: [{required:true, message: 'field is required', trigger: 'blur'}],
            mode: [{required:true, message: 'field is required', trigger: 'blur'}],
            port: [{required:true, trigger: 'blur',
              validator : (rule, value, callback) => {
                if (!value) {
                  return callback(new Error("port is required"));
                }
                if (!Number.isInteger(value)){
                  callback(new Error('port must be number'));
                } else {
                  callback();
                }
            }}],
            address: [{required: true, message:'field is required', trigger:'blur'}],
            user: [{required: true, message:'field is required', trigger:'blur'}],
            password: [{required: true, message:'field is required', trigger:'blur'}],
            server_id: [{validator:(rule, value, callback) => {
                if (!Number.isInteger(value)){
                  callback(new Error('server_id must be number'));
                } else {
                  callback();
                }
              }, trigger:'blur'}]
          },
          output: {
            sender:{
              type:[{required: true, message:'field is required', trigger:'blur'}]
            }
          }
        }
      },
      downloadLoading: false
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
    handleModifyStatus(row, status) {
      let req = {
        name: row.pipeline.name,
        status: status
      }
      fetchUpdateStatus(req).then(response =>{
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
        // const index = this.list.findIndex(v => v.pipeline.name === req.name)
        // this.list.splice(index,1, response.data)
        row.pipeline.status = status
        this.$notify({
          title: 'success',
          message: 'Update success',
          type: 'danger',
          duration: 2000
        })
      })
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
    resetTemp() {
      this.temp = {
        pipeline: {
          mysql: {
            address: '',
            port: '',
            user: '',
            password: '',
            server_id: 0,
            flavor: 'MySQL',
            mode: 'gtid'
          },
          remark: '',
          alias_name: '',
          name: '',
          output: {
            sender:{
              type: 'kafka',
              kafka: {
                brokers: '',
                topic: '',
                require_acks: 1,
                compression:0,
                retries:3,
                idepotent:false
              },
              stdout: null,
              http: {
                api:'',
                retries: 3
              },
              rabbitMQ: {
                url: '',
                exchange_name: ''
              },
              redis: {
                address:'',
                username:'',
                password:'',
                list:''
              },
              rocketMQ: {
                endpoint:'',
                topic_name: '',
                instance_id: '',
                access_key: '',
                secret_key: ''
              }
            }
          },
          filters: [{
            "type": "black",
            "rule": "mysql"
          }]
        }
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          fetchCreate(this.temp.pipeline).then(response => {
            setTimeout(() => {
              this.listLoading = false
            }, 1.5 * 1000)
            })
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
          })
        }
      })
    },
    handleUpdate(row) {
      //this.temp.pipeline =  Object.assign({}, row.pipeline) // copy obj
      //this.temp = {}
      this.temp.pipeline = Object.assign({}, row.pipeline)
      this.temp.info = Object.assign({},row.info)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({},this.temp.pipeline)
          fetchUpdate(tempData).then(response => {
            setTimeout(() => {
              this.listLoading = false
            }, 3 * 1000)
            const index = this.list.findIndex(v => v.pipeline.name === this.temp.pipeline.name)
            this.list.splice(index,1, response.data)
            this.dialogFormVisible = false
            this.$notify({
              title: 'success',
              message: 'Update Success!',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      this.$confirm('Delete pipeline is dangerous. Please confirm to continue.','Delete', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'danger'
      }).then(() => {
        const req = {
          name: row.pipeline.name
        }
        fetchDelete(req).then(response => {
          this.$message({
            type: 'success',
            message: 'Delete success!'
          });
          this.list.splice(index, 1)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Delete cancel!'
        });
      });
    },
    // handleFetchPv(pv) {
    //   fetchPv(pv).then(response => {
    //     this.pvData = response.data.pvData
    //     this.dialogPvVisible = true
    //   })
    // },
    // handleDownload() {
    //   this.downloadLoading = true
    //   import('@/vendor/Export2Excel').then(excel => {
    //     const tHeader = ['timestamp', 'title', 'type', 'importance', 'status']
    //     const filterVal = ['timestamp', 'title', 'type', 'importance', 'status']
    //     const data = this.formatJson(filterVal)
    //     excel.export_json_to_excel({
    //       header: tHeader,
    //       data,
    //       filename: 'table-list'
    //     })
    //     this.downloadLoading = false
    //   })
    // },
    // formatJson(filterVal) {
    //   return this.list.map(v => filterVal.map(j => {
    //     if (j === 'timestamp') {
    //       return parseTime(v[j])
    //     } else {
    //       return v[j]
    //     }
    //   }))
    // },
    addFilter() {
      this.temp.pipeline.filters.push({
        type: '' ,
        rule: ''
      });
    },
    removeFilter(item) {
      var index = this.temp.pipeline.filters.indexOf(item)
      if (index !== -1) {
        this.temp.pipeline.filters.splice(index, 1)
      }
    },
    dateFormat(value){
      const t = new Date(value);
      return t.getFullYear()+"-"+(t.getMonth()+1)+"-"+t.getDate()+" "+t.getHours()+":"+t.getMinutes()+":"+t.getSeconds();
    }
  }
}
</script>
