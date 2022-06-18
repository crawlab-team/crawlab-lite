<template>
  <div class="crawl-confirm-dialog-wrapper">
    <disclaimer-dialog v-model:value="disclaimerVisible" />
    <el-dialog
      :title="$t('Notification')"
      :visible="visible"
      class="crawl-confirm-dialog"
      width="480px"
      :before-close="beforeClose"
    >
      <div style="margin-bottom: 20px">
        {{ $t('Are you sure to run this spider?') }}
      </div>
      <el-form
        ref="form"
        :model="form"
        :label-width="lang === 'zh' ? '100px' : '150px'"
      >
        <el-form-item
          :label="$t('Execute Command')"
          inline-message
          prop="cmd"
          required
        >
          <el-input
            v-model:value="form.cmd"
            :placeholder="$t('Execute Command')"
          />
        </el-form-item>
        <el-form-item
          :label="$t('Version')"
          inline-message
          prop="spider_version_id"
          required
        >
          <el-select
            v-model:value="form.spider_version_id"
            :loading="loadingVersions"
            :placeholder="$t('Latest Version')"
            @focus="onSelectSpiderVersion"
          >
            <el-option
              value="00000000-0000-0000-0000-000000000000"
              :label="$t('Latest Version')"
            />
            <el-option
              v-for="version in spiderVersionList"
              :key="version.id"
              :value="version.id"
              :label="getTime(version.create_ts).format('YYYY-MM-DD HH:mm:ss')"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div>
        <el-checkbox v-model:value="isAllowDisclaimer" />
        <span v-if="lang === 'zh'" style="margin-left: 5px">
          我已阅读并同意
          <a href="javascript:" @click="onClickDisclaimer"> 《免责声明》 </a>
          所有内容
        </span>
        <span v-else style="margin-left: 5px">
          I have read and agree all content in
          <a href="javascript:" @click="onClickDisclaimer"> Disclaimer </a>
        </span>
      </div>
      <!--          <div>-->
      <!--            <el-checkbox v-model="isRedirect"/>-->
      <!--            <span style="margin-left: 5px">{{$t('Redirect to task detail')}}</span>-->
      <!--          </div>-->
      <template v-slot:footer>
        <el-button type="plain" size="small" @click="$emit('close')">{{
          $t('Cancel')
        }}</el-button>
        <el-button
          type="primary"
          size="small"
          :disabled="isConfirmDisabled"
          @click="onConfirm"
        >
          {{ $t('Confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { $on, $off, $once, $emit } from '../../utils/gogocodeTransfer'
import * as Vue from 'vue'
import { mapState } from 'vuex'
import dayjs from 'dayjs'
import DisclaimerDialog from '@/components/Disclaimer'

export default {
  name: 'CrawlConfirmDialog',
  components: {
    DisclaimerDialog,
  },
  props: {
    visible: {
      type: Boolean,
      default: false,
    },
    spiderId: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      form: this.defaultForm(),
      isAllowDisclaimer: true,
      isRetry: false,
      isRedirect: false,
      loadingVersions: false,
      versionList: [],
      disclaimerVisible: false,
    }
  },
  computed: {
    ...mapState('spider', ['spiderForm', 'spiderVersionList']),
    ...mapState('lang', ['lang']),
    isConfirmDisabled() {
      return !this.isAllowDisclaimer
    },
  },
  watch: {
    visible: function (current) {
      $emit(this, 'update:value', current)
      if (!this.visible) {
        this.$nextTick(() => {
          this.form = this.defaultForm()
        })
      }
    },
  },
  methods: {
    defaultForm() {
      return {
        spider: undefined,
        spider_version_id: '00000000-0000-0000-0000-000000000000',
        cmd: undefined,
      }
    },
    async onSelectSpiderVersion() {
      this.loadingVersions = true
      await this.$store.dispatch('spider/getSpiderVersionList', {
        spider_id: this.spiderId,
      })
      this.loadingVersions = false
    },
    beforeClose() {
      $emit(this, 'close')
    },
    async onConfirm() {
      this.$refs['form'].validate(async (valid) => {
        if (!valid) return

        const res = await this.$store.dispatch('spider/crawlSpider', {
          spiderId: this.spiderId,
          spiderVersionId: this.form.spiderVersionId,
          cmd: this.form.cmd,
        })

        // 消息提示
        this.$message.success(this.$t('A task has been scheduled successfully'))

        $emit(this, 'close')

        // 是否重定向
        if (this.isRedirect) {
          // 返回任务id
          const id = res.data.data[0]
          await this.$router.push(`/tasks/${id}`)
          this.$st.sendEv('爬虫确认', '跳转到任务详情')
        }

        $emit(this, 'confirm')
      })
    },
    onClickDisclaimer() {
      this.disclaimerVisible = true
    },
    getTime(str) {
      return dayjs(str)
    },
  },
  emits: ['close', 'update:value', 'confirm'],
}
</script>

<style scoped>
.crawl-confirm-dialog >>> .el-dialog__body {
  padding-left: 3rem;
  padding-right: 3rem;
}
.crawl-confirm-dialog >>> .el-form .el-form-item {
  margin-bottom: 20px;
}
.crawl-confirm-dialog >>> .checkbox-wrapper a {
  color: #409eff;
}
.crawl-confirm-dialog >>> .param-input {
  width: calc(100% - 56px);
}
.crawl-confirm-dialog >>> .param-input .el-input__inner {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
}
.crawl-confirm-dialog >>> .param-btn {
  width: 56px;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}
</style>
