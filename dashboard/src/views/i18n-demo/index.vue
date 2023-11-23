<script src="../../lang/index.js"></script>
<script src="../../lang/zh.js"></script>
<script src="local.js"></script>
<template>
  <div>
    <el-card class="box-card" style="margin-top:40px;">
      <div slot="header" class="clearfix">
        <svg-icon icon-class="international" />
        <span style="margin-left:10px;">{{ $t('i18nView.title') }}</span>
      </div>
      <div>
        <el-radio-group v-model="lang" size="small">
          <el-radio label="zh" border>
            简体中文
          </el-radio>
          <el-radio label="en" border>
            English
          </el-radio>
          <el-radio label="es" border>
            Español
          </el-radio>
          <el-radio label="ja" border>
            日本語
          </el-radio>
        </el-radio-group>
        <el-tag style="margin-top:15px;display:block;" type="info">
          {{ $t('i18nView.note') }}
        </el-tag>
      </div>
    </el-card>
  </div>
</template>

<script>
import local from './local'
const viewName = 'i18nView'

export default {
  name: 'I18n',
  computed: {
    lang: {
      get() {
        return this.$store.state.app.language
      },
      set(lang) {
        this.$i18n.locale = lang
        this.$store.dispatch('app/setLanguage', lang)
      }
    }
  },
  watch: {
    lang() {
      this.setOptions()
    }
  },
  created() {
    if (!this.$i18n.getLocaleMessage('en')[viewName]) {
      this.$i18n.mergeLocaleMessage('en', local.en)
      this.$i18n.mergeLocaleMessage('zh', local.zh)
      this.$i18n.mergeLocaleMessage('es', local.es)
      this.$i18n.mergeLocaleMessage('ja', local.ja)
    }
    this.setOptions() // set default select options
  },
  methods: {
    setOptions() {
      this.options = [
        {
          value: '1',
          label: this.$t('i18nView.one')
        },
        {
          value: '2',
          label: this.$t('i18nView.two')
        },
        {
          value: '3',
          label: this.$t('i18nView.three')
        }
      ]
    }
  }
}
</script>

<style scoped>
.box-card {
  width: 600px;
  max-width: 100%;
  margin: 20px auto;
}
</style>
