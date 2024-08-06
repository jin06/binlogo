<script src="../../../lang/zh.js"></script>
<script src="../../../api/pipeline.js"></script>
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="8" :xs="24">
        <info-card :bitem="bitem" />
      </el-col>
      <el-col :span="14" :xs="24">
        <el-card>
          <el-tabs v-model="activeTab">
            <el-tab-pane :label="$t('pipe_detail.filter')" name="filter">
              <pfilter :list="bitem.pipeline.filters" :pipe_name="bitem.pipeline.name" />
            </el-tab-pane>
            <el-tab-pane :label="$t('pipe_detail.govern')" name="govern">
              <govern :govern_form="govern_form" :bitem="bitem" :instance="bitem.info.instance" :pipe_name="bitem.pipeline.name" />
            </el-tab-pane>
            <el-tab-pane :label="$t('pipe_detail.instance')" name="instance">
              <instance :instance="bitem.info.instance" />
            </el-tab-pane>
            <el-tab-pane :label="$t('pipe_detail.event')" name="event">
              <event :list="events" />
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import InfoCard from './components/InfoCard'
import Instance from './components/Instance'
import Govern from './components/Govern'
import Pfilter from './components/Pfilter'
import Event from './components/Event'
import { fetchGet } from '@/api/pipeline'
import { fetchGet as fetchGetPos } from '@/api/position'
import { fetchList as fetchListEvent } from '@/api/event'

export default {
  name: 'Profile',
  components: { InfoCard, Instance, Govern, Pfilter, Event },
  data() {
    return {
      pipe_name: undefined,
      bitem: {
        info: {
          instance:{}
        },
        pipeline: {
          status:''
        }
      },
      user: {},
      // activeTab: 'filter'
      activeTab: 'filter',
      govern_form: {
        mode:'',
        position:{
          binlog_file:'',
          binlog_position: {
            type: Number
          },
          gtid_set:''
        }
      },
      events: []
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'avatar',
      'roles'
    ])
  },
  created() {
    this.pipe_name = this.$route.params && this.$route.params.name
    this.getItem()
    this.getUser()
    this.getPosition()
    this.getEvent()
  },
  methods: {
    getUser() {
      this.user = {
        name: this.name,
        role: this.roles.join(' | '),
        email: 'admin@test.com',
        avatar: this.avatar
      }
    },
    getItem() {
      const req = {
        name:this.pipe_name
      }
      fetchGet(req).then(response => {
        this.bitem = response.data
        this.govern_form.mode = this.bitem.pipeline.mysql.mode
      })
    },
    getPosition() {
      const req = {
        pipe_name: this.pipe_name
      }
      fetchGetPos(req).then(response => {
        this.govern_form.position = Object.assign({}, response.data)
      })
    },
    getEvent() {
      const req = {
        res_type: 'pipeline',
        res_name: this.pipe_name
      }
      fetchListEvent(req).then(response => {
        this.events = response.data.items
      })
    }
  }
}
</script>
