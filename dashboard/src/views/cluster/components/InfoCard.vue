<script src="../../../../router/index.js"></script>
<template>
  <el-card style="margin-bottom:20px;" shadow="always">
    <div slot="header" class="clearfix">
      <span>{{ $t('cluster.title') }}</span>
    </div>

    <el-card style="margin-bottom: 18px" shadow="never">
      <div class="user-profile">
        <div class="box-center">
          <div class="user-name text-center">{{ bcluster.name }}</div>
          <div class="user-role text-center text-muted">{{ $t('cluster.name' )}}</div>
        </div>
      </div>
    </el-card>
  </el-card>
</template>

<script>

import { dateFormat } from '@/utils/my-time-format'
import { fetchUpdateStatus} from "@/api/pipeline";

export default {
  components: {},
  props: {
    bcluster: {}
  },
  methods: {
    dateFormat(val) {
      return dateFormat(val);
    }
  }
}
</script>

<style lang="scss" scoped>
.box-center {
  margin: 0 auto;
  display: table;
}

.text-muted {
  color: #777;
}

.user-profile {
  .user-name {
    font-weight: bold;
  }

  .box-center {
    padding-top: 10px;
  }

  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }

  .box-social {
    padding-top: 30px;

    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }

  .user-follow {
    padding-top: 20px;
  }
}

.user-bio {
  margin-top: 20px;
  color: #606266;

  span {
    padding-left: 4px;
  }

  .user-bio-section {
    font-size: 14px;
    padding: 15px 0;

    .user-bio-section-header {
      border-bottom: 1px solid #dfe6ec;
      padding-bottom: 10px;
      margin-bottom: 10px;
      font-weight: bold;
    }
  }
}
</style>
