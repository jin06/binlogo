<template>
  <div class="block">
    <el-timeline style="margin-top: 55px">
      <el-timeline-item
        v-for="(active, index) in timeline"
        :key="index"
        :timestamp="active.create_revision"
        :color="active.color"
        :type="active.type"
        :icon="active.icon"
        :size="active.size"
      >
        {{ active.node }}
      </el-timeline-item>
    </el-timeline>
  </div>
</template>

<script>

import { fetchElectionList } from '@/api/cluster'

export default {
  data() {
    return {
      timeline: []
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      fetchElectionList().then(response => {
        this.timeline = response.data.items
        this.timeline.forEach((item, idx) => {
          // eslint-disable-next-line no-empty
          if (idx === 0) {
            item.type = 'primary'
            item.size = 'large'
            item.icon = 'el-icon-star-on"'
            item.remark = 'master'
          }
          if (idx === 1) {
            item.type = 'warning'
            item.size = 'normal'
            item.icon = 'el-icon-star-off"'
            item.remark = 'heir'
          }
        })
      })
    }
  }
}
</script>
