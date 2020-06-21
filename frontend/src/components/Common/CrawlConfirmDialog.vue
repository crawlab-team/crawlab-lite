<template>
  <div class="crawl-confirm-dialog-wrapper">
    <el-dialog
      :title="$t('Notification')"
      :visible="visible"
      class="crawl-confirm-dialog"
      width="480px"
      :before-close="beforeClose"
    >
      <div style="margin-bottom: 20px;">{{ $t('Are you sure to run this spider?') }}</div>
      <el-form ref="form" :model="form" :label-width="lang === 'zh' ? '100px' : '150px'">
        <el-form-item :label="$t('Execute Command')" inline-message prop="cmd" required>
          <template>
            <el-input v-model="form.cmd" :placeholder="$t('Execute Command')" />
          </template>
        </el-form-item>
      </el-form>
      <div>
        <el-checkbox v-model="isAllowDisclaimer" />
        <span v-if="lang === 'zh'" style="margin-left: 5px">
          我已阅读并同意
          <a href="javascript:" @click="onClickDisclaimer">
            《免责声明》
          </a>
          所有内容
        </span>
        <span v-else style="margin-left: 5px">
          I have read and agree all content in
          <a href="javascript:" @click="onClickDisclaimer">
            Disclaimer
          </a>
        </span>
      </div>
      <!--          <div>-->
      <!--            <el-checkbox v-model="isRedirect"/>-->
      <!--            <span style="margin-left: 5px">{{$t('Redirect to task detail')}}</span>-->
      <!--          </div>-->
      <template slot="footer">
        <el-button type="plain" size="small" @click="$emit('close')">{{ $t('Cancel') }}</el-button>
        <el-button type="primary" size="small" :disabled="isConfirmDisabled" @click="onConfirm">
          {{ $t('Confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
  import { mapState } from 'vuex'

  export default {
    name: 'CrawlConfirmDialog',
    props: {
      spiderId: {
        type: String,
        default: ''
      },
      spiders: {
        type: Array,
        default() {
          return []
        }
      },
      visible: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        form: {
          spider: undefined,
          cmd: ''
        },
        isAllowDisclaimer: true,
        isRetry: false,
        isRedirect: false,
        isLoading: false
      }
    },
    computed: {
      ...mapState('spider', [
        'spiderForm'
      ]),
      ...mapState('setting', [
        'setting'
      ]),
      ...mapState('lang', [
        'lang'
      ]),
      isConfirmDisabled() {
        if (this.isLoading) return true
        return !this.isAllowDisclaimer
      }
    },
    methods: {
      beforeClose() {
        this.$emit('close')
      },
      onConfirm() {
        this.$refs['form'].validate(async valid => {
          if (!valid) return

          // 请求响应
          let res

          if (!this.multiple) {
            // 运行单个爬虫

            // 参数
            const cmd = this.form.cmd

            // 发起请求
            res = await this.$store.dispatch('spider/crawlSpider', {
              spiderId: this.spiderId,
              cmd
            })
          } else {
            // 运行多个爬虫

            // 发起请求
            res = await this.$store.dispatch('spider/crawlSelectedSpiders', {
              taskParams: this.spiders.map(d => {
                // 参数
                const cmd = this.form.cmd

                return {
                  spider_id: d.id,
                  cmd
                }
              })
            })
          }

          // 消息提示
          this.$message.success(this.$t('A task has been scheduled successfully'))

          this.$emit('close')

          // 是否重定向
          if (
            this.isRedirect &&
            !this.spiderForm.is_long_task &&
            !this.multiple
          ) {
            // 返回任务id
            const id = res.data.data[0]
            this.$router.push('/tasks/' + id)
            this.$st.sendEv('爬虫确认', '跳转到任务详情')
          }

          this.$emit('confirm')
        })
      },
      onClickDisclaimer() {
        this.$router.push('/disclaimer')
      }
    }
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
