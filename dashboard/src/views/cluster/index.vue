<script src="../../api/instance.js"></script></script>
<script src="../../../api/pipeline.js"></script>
<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="6" :xs="24">
        <info-card :bcluster="bcluster" />
      </el-col>
      <el-col :span="18" :xs="24">
        <el-card>
          <el-tabs v-model="activeName">
            <el-tab-pane key="election" :label="$t('cluster.tabMap.election.label')" name="election">
              <keep-alive>
                <election type="election"></election>
              </keep-alive>
            </el-tab-pane>
            <el-tab-pane key="register" :label="$t('cluster.tabMap.register.label')" name="register">
              <keep-alive>
                <register type="register"></register>
              </keep-alive>
            </el-tab-pane>
            <el-tab-pane key="instance" :label="$t('cluster.tabMap.instance.label')" name="instance">
              <keep-alive>
                <instance type="instance"></instance>
              </keep-alive>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import InfoCard from './components/InfoCard'
import {fetchGet} from "@/api/cluster";
import Election from "./components/Election";
import Instance from "./components/Instance";
import Register from "./components/Register";

export default {
  name: 'Profile',
  components: { InfoCard, Election, Instance, Register},
  data() {
    return {
      bcluster: {},
      activeName: 'election'
    }
  },
  watch: {
    activeName(val) {
      this.$router.push(`${this.$route.path}?tab=${val}`)
    }
  },
  created() {
    const tab = this.$route.query.tab
    if (tab) {
      this.activeName = tab
    }
    this.getClusterInfo()
  },
  methods: {
    getClusterInfo() {
      fetchGet().then(response => {
        this.bcluster = response.data
      })
    }
  }
}
</script>
