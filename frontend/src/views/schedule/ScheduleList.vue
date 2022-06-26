<template>
  <div class="app-container schedule-list">
    <!--add schedule dialog-->
    <el-dialog
      :title="t(dialogTitle)"
      v-model="dialogVisible"
      width="640px"
      :before-close="onDialogClose"
    >
      <el-form
        ref="scheduleForm"
        label-width="180px"
        class="add-form"
        :model="scheduleForm"
        :inline-message="true"
        label-position="right"
      >
        <el-form-item :label="t('Spider')" prop="spider_id" required>
          <el-select
            id="spider-id"
            v-model="scheduleForm.spider_id"
            :placeholder="t('Spider')"
            :loading="loadingSpiders"
            @focus="onSelectSpider"
            @change="onSpiderChange"
          >
            <el-option
              v-for="op in spiderList"
              :key="op.id"
              :value="op.id"
              :label="`${op.name}`"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="scheduleForm.spider_id"
          :label="t('Version')"
          inline-message
          prop="spider_version_id"
          required=""
        >
          <el-select
            v-model="scheduleForm.spider_version_id"
            :loading="loadingVersions"
            :placeholder="t('Latest Version')"
            @focus="onSelectSpiderVersion"
          >
            <el-option
              value="00000000-0000-0000-0000-000000000000"
              :label="t('Latest Version')"
            />
            <el-option
              v-for="version in spiderVersionList"
              :key="version.id"
              :value="version.id"
              :label="getTime(version.create_ts).format('YYYY-MM-DD HH:mm:ss')"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('Cron')" prop="cron" required>
          <el-input
            id="cron"
            ref="cron"
            v-model="scheduleForm.cron"
            class="cron"
            :placeholder="`${t(
              '[second] [minute] [hour] [day] [month] [day of week]',
            )}`"
            style="width: calc(100% - 100px)"
          />
          <el-button
            class="cron-edit"
            type="primary"
            :icon="ElIconEdit"
            style="width: 100px"
            @click="onShowCronDialog"
          >
            {{ t('Edit') }}
          </el-button>
        </el-form-item>
        <el-form-item :label="t('Execute Command')" prop="cmd" required>
          <el-input
            id="cmd"
            v-model="scheduleForm.cmd"
            :placeholder="t('Execute Command')"
          />
        </el-form-item>
        <el-form-item :label="t('Schedule Description')" prop="description">
          <el-input
            id="schedule-description"
            v-model="scheduleForm.description"
            type="textarea"
            :placeholder="t('Schedule Description')"
          />
        </el-form-item>
      </el-form>
      <!--取消、保存-->
      <template v-slot:footer>
        <span class="dialog-footer">
          <el-button size="small" @click="onCancel">{{
            t('Cancel')
          }}</el-button>
          <el-button
            id="btn-submit"
            size="small"
            type="primary"
            :disabled="submitting"
            @click="onAddSubmit"
            >{{ t('Submit') }}</el-button
          >
        </span>
      </template>
    </el-dialog>
    <!--./add schedule dialog-->

    <!--view tasks dialog-->
    <el-dialog
      :title="t('Tasks')"
      v-model="tasksDialogVisible"
      width="calc(100% - 240px)"
      :before-close="onCloseTasksDialog"
    >
      <schedule-task-list ref="schedule-task-list" />
    </el-dialog>
    <!--./view tasks dialog-->

    <!--cron generation dialog-->
    <el-dialog title="生成 Cron" v-model="cronDialogVisible">
      <vue-cron-linux
        ref="vue-cron-linux"
        :data="scheduleForm.cron"
        :i18n="lang"
        @submit="onCronChange"
      />
      <template v-slot:footer>
        <span class="dialog-footer">
          <el-button size="small" @click="cronDialogVisible = false">{{
            t('Cancel')
          }}</el-button>
          <el-button size="small" type="primary" @click="onCronDialogSubmit">{{
            t('Confirm')
          }}</el-button>
        </span>
      </template>
    </el-dialog>
    <!--./cron generation dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="scheduleForm.spider_id"
      @close="() => (this.crawlConfirmDialogVisible = false)"
      @confirm="() => (this.crawlConfirmDialogVisible = false)"
    />
    <!--./crawl confirm dialog-->

    <el-card style="border-radius: 0" class="schedule-list">
      <!--filter-->
      <div class="filter">
        <div class="right">
          <el-button
            size="small"
            type="primary"
            :icon="ElIconPlus"
            class="btn-add"
            @click="onAdd"
          >
            {{ t('Add Schedule') }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table
        :data="filteredTableData"
        class="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        border
      >
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'status'"
            :property="col.name"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              <el-tooltip
                v-if="scope.row[col.name] === 'error'"
                :content="t(scope.row['message'])"
                placement="top"
              >
                <el-tag class="status-tag" type="danger">
                  {{ scope.row[col.name] ? t(scope.row[col.name]) : t('NA') }}
                </el-tag>
              </el-tooltip>
              <el-tag v-else class="status-tag">
                {{ scope.row[col.name] ? t(scope.row[col.name]) : t('NA') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'enable'"
            :label="t(col.label)"
            :width="col.width"
          >
            <template v-slot="scope">
              <el-switch
                v-model="scope.row.enabled"
                active-color="#13ce66"
                inactive-color="#ff4949"
                @change="onEnabledChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :property="col.name"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ scope.row[col.name] }}
            </template>
          </el-table-column>
        </template>
        <el-table-column
          :label="t('Action')"
          class="actions"
          align="left"
          width="170"
          fixed="right"
        >
          <template v-slot="scope">
            <!--edit-->
            <el-tooltip :content="t('Edit')" placement="top">
              <el-button
                type="warning"
                :icon="ElIconEdit"
                size="mini"
                @click="onEdit(scope.row)"
              />
            </el-tooltip>
            <!--./edit-->

            <!--delete-->
            <el-tooltip :content="t('Remove')" placement="top">
              <el-button
                type="danger"
                :icon="ElIconDelete"
                size="mini"
                @click="onRemove(scope.row)"
              />
            </el-tooltip>
            <!--./delete-->

            <!--view tasks-->
            <el-tooltip :content="t('View Tasks')" placement="top">
              <el-button
                type="primary"
                :icon="ElIconSearch"
                size="mini"
                @click="onViewTasks(scope.row)"
              />
            </el-tooltip>
            <!--./view tasks-->

            <!--run-->
            <el-tooltip :content="t('Run')" placement="top">
              <el-button
                type="success"
                :icon="faBug"
                size="mini"
                @click="onRun(scope.row)"
              />
            </el-tooltip>
            <!--./run-->
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page_num"
          :page-sizes="[10, 20, 50, 100]"
          v-model:page-size="pagination.page_size"
          layout="sizes, prev, pager, next"
          :total="scheduleTotal"
          @current-change="getScheduleList"
          @size-change="getScheduleList"
        />
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import * as Vue from 'vue'
import VueCronLinux from '../../components/Cron'
import { mapState } from 'vuex'
import ScheduleTaskList from '../../components/Schedule/ScheduleTaskList'
import CrawlConfirmDialog from '../../components/Dialog/CrawlConfirmDialog'
import dayjs from 'dayjs'
import { useI18n } from 'vue-i18n'

export default {
  name: 'ScheduleList',
  components: {
    CrawlConfirmDialog,
    ScheduleTaskList,
    VueCronLinux,
  },
  setup(props) {
    const { t } = useI18n()
    return { t }
  },
  data() {
    return {
      columns: [
        { name: 'spider_name', label: 'Spider', width: '150px' },
        { name: 'cron', label: 'Cron', width: '150px' },
        { name: 'cmd', label: 'Cmd', width: '200px' },
        { name: 'enable', label: 'Enable/Disable', width: '120px' },
        { name: 'description', label: 'Description' },
        // { name: 'status', label: 'Status', width: '100px' }
      ],
      pagination: {
        page_num: 1,
        page_size: 10,
      },
      isEdit: false,
      dialogTitle: '',
      dialogVisible: false,
      cronDialogVisible: false,
      expression: '',
      isShowCron: false,
      submitting: false,
      loadingSpiders: false,
      loadingVersions: false,
      tasksDialogVisible: false,
      crawlConfirmDialogVisible: false,
    }
  },
  computed: {
    ...mapState('schedule', ['scheduleList', 'scheduleTotal', 'scheduleForm']),
    ...mapState('spider', ['spiderList', 'spiderVersionList', 'spiderForm']),
    lang() {
      const lang =
        this.$store.state.lang.lang || window.localStorage.getItem('lang')
      if (!lang) return 'cn'
      if (lang === 'zh') return 'cn'
      return 'en'
    },
    filteredTableData() {
      return this.scheduleList
    },
    spider() {
      for (let i = 0; i < this.spiderList.length; i++) {
        if (this.spiderList[i].id === this.scheduleForm.spider_id) {
          return this.spiderList[i]
        }
      }
      return {}
    },
  },
  created() {
    this.getScheduleList()
  },
  mounted() {},
  methods: {
    onDialogClose() {
      this.dialogVisible = false
    },
    onCancel() {
      this.dialogVisible = false
    },
    onAdd() {
      this.isEdit = false
      this.dialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', {})
      this.$st.sendEv('定时任务', '添加定时任务')
    },
    onAddSubmit() {
      this.$refs.scheduleForm.validate((res) => {
        if (res) {
          const form = Object.assign({}, this.scheduleForm)
          if (form.cron.search(/^\S+?\s\S+?\s\S+?\s\S+?\s\S+?$/) !== -1) {
            form['cron'] = '0 ' + form.cron
          }
          form['enabled'] = form.enabled ? 1 : 0
          if (this.isEdit) {
            this.$store
              .dispatch('schedule/editSchedule', form)
              .then((response) => {
                if (response.data.code !== 200) {
                  this.$message.error(response.data.message)
                  return
                }
                this.dialogVisible = false
                this.getScheduleList()
                this.$message.success(this.t('The schedule has been saved'))
              })
          } else {
            this.$store
              .dispatch('schedule/addSchedule', form)
              .then((response) => {
                if (response.data.code !== 200) {
                  this.$message.error(response.data.message)
                  return
                }
                this.dialogVisible = false
                this.getScheduleList()
                this.$message.success(this.t('The schedule has been added'))
              })
          }
        }
      })
      this.$st.sendEv('定时任务', '提交定时任务')
    },
    async onEdit(row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.dialogVisible = true
      this.isEdit = true
      this.$st.sendEv('定时任务', '修改定时任务')

      this.submitting = true
      await this.$store.dispatch('schedule/getSchedule', row.id)
      await this.getSpiderList()
      await this.getSpiderVersionList()
      this.submitting = false
    },
    onRemove(row) {
      this.$confirm(
        this.t('Are you sure to delete the schedule task?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.$store.dispatch('schedule/removeSchedule', row.id).then(() => {
            setTimeout(() => {
              this.getScheduleList()
              this.$message.success(this.t('The schedule has been removed'))
            }, 100)
          })
        })
        .catch(() => {})
      this.$st.sendEv('定时任务', '删除定时任务')
    },
    async getScheduleList() {
      await this.$store.dispatch('schedule/getScheduleList', this.pagination)
    },
    async getSpiderList() {
      this.loadingSpiders = true
      await this.$store.dispatch('spider/getSpiderList')
      this.loadingSpiders = false
    },
    async getSpiderVersionList() {
      this.loadingVersions = true
      await this.$store.dispatch('spider/getSpiderVersionList', {
        spider_id: this.scheduleForm.spider_id,
      })
      this.loadingVersions = false
    },
    async onSelectSpider() {
      if (this.spiderList.length === 0) {
        await this.getSpiderList()
      }
    },
    async onSelectSpiderVersion() {
      if (this.spiderVersionList.length === 0) {
        await this.getSpiderVersionList()
      }
    },
    async onEnabledChange(row) {
      let res
      if (row.enabled) {
        res = await this.$store.dispatch('schedule/enableSchedule', row.id)
      } else {
        res = await this.$store.dispatch('schedule/disableSchedule', row.id)
      }
      if (res.data && res.data.code === 200) {
        this.$message.success(
          this.t(
            `${row.enabled ? 'Enabling' : 'Disabling'} the schedule successful`,
          ),
        )
      } else {
        this.$message.error(
          this.t(
            `${
              row.enabled ? 'Enabling' : 'Disabling'
            } the schedule unsuccessful`,
          ),
        )
      }
      this.$st.sendEv('定时任务', '启用/禁用')
    },
    onCronChange(value) {
      this.scheduleForm['cron'] = value
      this.$st.sendEv('定时任务', '配置Cron')
    },
    onCronDialogSubmit() {
      const valid = this.$refs['vue-cron-linux'].submit()
      if (valid) {
        this.cronDialogVisible = false
      }
    },
    onSpiderChange() {
      this.scheduleForm['spider_version_id'] =
        '00000000-0000-0000-0000-000000000000'
    },
    onShowCronDialog() {
      this.cronDialogVisible = true
      this.$st.sendEv('定时任务', '点击编辑Cron')
    },
    onCloseTasksDialog() {
      this.$refs['schedule-task-list'].close()
      this.tasksDialogVisible = false
    },
    async onViewTasks(row) {
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.tasksDialogVisible = true
      setTimeout(() => {
        this.$refs['schedule-task-list'].open()
      }, 100)
      this.$st.sendEv('定时任务', '查看任务列表')
    },
    async onRun(row) {
      this.crawlConfirmDialogVisible = true
      this.$store.commit('schedule/SET_SCHEDULE_FORM', row)
      this.$st.sendEv('定时任务', '点击运行任务')
    },
    getTime(str) {
      return dayjs(str)
    },
  },
}
</script>

<style scoped>
.filter .right {
  text-align: right;
}
.table {
  margin-top: 8px;
  border-radius: 5px;
}
.table .el-button {
  width: 28px;
  height: 28px;
  padding: 0;
}
.status-tag {
  cursor: pointer;
}
.schedule-list >>> .param-input {
  width: calc(100% - 56px);
}
.schedule-list >>> .param-input .el-input__inner {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
}
.schedule-list >>> .param-btn {
  width: 56px;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}
.cron {
  width: calc(100% - 100px);
}
.cron >>> .el-input__inner {
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: none;
}
.cron-edit {
  width: 100px;
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
}
</style>
