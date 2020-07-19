<template>
  <div class="app-container schedule-list">
    <!--add schedule dialog-->
    <el-dialog
      :title="$t(dialogTitle)"
      :visible.sync="dialogVisible"
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
        <el-form-item :label="$t('Spider')" prop="spider_id" required>
          <el-select
            id="spider-id"
            v-model="scheduleForm.spider_id"
            :placeholder="$t('Spider')"
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
        <el-form-item v-if="scheduleForm.spider_id" :label="$t('Version')" inline-message prop="spiderVersionId">
          <el-select
            v-model="scheduleForm.spider_version_id"
            :loading="loadingVersions"
            :placeholder="$t('Latest Version')"
            @focus="onSelectSpiderVersion"
          >
            <el-option
              v-for="(version, index) in spiderVersionList"
              :key="version.id"
              :value="version.id"
              :label="getTime(version.create_ts).format('YYYY-MM-DD HH:mm:ss') + (index === 0 ? ` (${$t('Latest Version')})` : '')"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('Cron')" prop="cron" required>
          <el-input
            id="cron"
            ref="cron"
            v-model="scheduleForm.cron"
            class="cron"
            :placeholder="`${$t('[minute] [hour] [day] [month] [day of week]')}`"
            style="width: calc(100% - 100px)"
          />
          <el-button
            class="cron-edit"
            type="primary"
            icon="el-icon-edit"
            style="width: 100px"
            @click="onShowCronDialog"
          >
            {{ $t('Edit') }}
          </el-button>
        </el-form-item>
        <el-form-item :label="$t('Execute Command')" prop="cmd" required>
          <el-input
            id="cmd"
            v-model="scheduleForm.cmd"
            :placeholder="$t('Execute Command')"
          />
        </el-form-item>
        <el-form-item :label="$t('Schedule Description')" prop="description">
          <el-input
            id="schedule-description"
            v-model="scheduleForm.description"
            type="textarea"
            :placeholder="$t('Schedule Description')"
          />
        </el-form-item>
      </el-form>
      <!--取消、保存-->
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="onCancel">{{ $t('Cancel') }}</el-button>
        <el-button id="btn-submit" size="small" type="primary" :disabled="submitting" @click="onAddSubmit">{{ $t('Submit') }}</el-button>
      </span>
    </el-dialog>
    <!--./add schedule dialog-->

    <!--view tasks dialog-->
    <el-dialog
      :title="$t('Tasks')"
      :visible.sync="tasksDialogVisible"
      width="calc(100% - 240px)"
      :before-close="onCloseTasksDialog"
    >
      <schedule-task-list ref="schedule-task-list" />
    </el-dialog>
    <!--./view tasks dialog-->

    <!--cron generation dialog-->
    <el-dialog title="生成 Cron" :visible.sync="cronDialogVisible">
      <vue-cron-linux ref="vue-cron-linux" :data="scheduleForm.cron" :i18n="lang" @submit="onCronChange" />
      <span slot="footer" class="dialog-footer">
        <el-button size="small" @click="cronDialogVisible = false">{{ $t('Cancel') }}</el-button>
        <el-button size="small" type="primary" @click="onCronDialogSubmit">{{ $t('Confirm') }}</el-button>
      </span>
    </el-dialog>
    <!--./cron generation dialog-->

    <!--crawl confirm dialog-->
    <crawl-confirm-dialog
      :visible="crawlConfirmDialogVisible"
      :spider-id="scheduleForm.spider_id"
      @close="() => this.crawlConfirmDialogVisible = false"
      @confirm="() => this.crawlConfirmDialogVisible = false"
    />
    <!--./crawl confirm dialog-->

    <el-card style="border-radius: 0" class="schedule-list">
      <!--filter-->
      <div class="filter">
        <div class="right">
          <el-button
            size="small"
            type="primary"
            icon="el-icon-plus"
            class="btn-add"
            @click="onAdd"
          >
            {{ $t('Add Schedule') }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table
        :data="filteredTableData"
        class="table"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        border
      >
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'status'"
            :key="col.name"
            :property="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              <el-tooltip v-if="scope.row[col.name] === 'error'" :content="$t(scope.row['message'])" placement="top">
                <el-tag class="status-tag" type="danger">
                  {{ scope.row[col.name] ? $t(scope.row[col.name]) : $t('NA') }}
                </el-tag>
              </el-tooltip>
              <el-tag v-else class="status-tag">
                {{ scope.row[col.name] ? $t(scope.row[col.name]) : $t('NA') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column v-else-if="col.name === 'enable'" :key="col.name" :label="$t(col.label)" :width="col.width">
            <template slot-scope="scope">
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
            :key="col.name"
            :property="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ scope.row[col.name] }}
            </template>
          </el-table-column>
        </template>
        <el-table-column :label="$t('Action')" class="actions" align="left" width="170" fixed="right">
          <template slot-scope="scope">
            <!--edit-->
            <el-tooltip :content="$t('Edit')" placement="top">
              <el-button type="warning" icon="el-icon-edit" size="mini" @click="onEdit(scope.row)" />
            </el-tooltip>
            <!--./edit-->

            <!--delete-->
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button type="danger" icon="el-icon-delete" size="mini" @click="onRemove(scope.row)" />
            </el-tooltip>
            <!--./delete-->

            <!--view tasks-->
            <el-tooltip :content="$t('View Tasks')" placement="top">
              <el-button type="primary" icon="el-icon-search" size="mini" @click="onViewTasks(scope.row)" />
            </el-tooltip>
            <!--./view tasks-->

            <!--run-->
            <el-tooltip :content="$t('Run')" placement="top">
              <el-button type="success" icon="fa fa-bug" size="mini" @click="onRun(scope.row)" />
            </el-tooltip>
            <!--./run-->
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          :current-page.sync="pagination.page_num"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pagination.page_size"
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
  import VueCronLinux from '../../components/Cron'
  import { mapState } from 'vuex'
  import ScheduleTaskList from '../../components/Schedule/ScheduleTaskList'
  import CrawlConfirmDialog from '../../components/Common/CrawlConfirmDialog'
  import dayjs from 'dayjs'

  export default {
    name: 'ScheduleList',
    components: {
      CrawlConfirmDialog,
      ScheduleTaskList,
      VueCronLinux
    },
    data() {
      return {
        columns: [
          { name: 'spider_name', label: 'Spider', width: '150px' },
          { name: 'cron', label: 'Cron', width: '150px' },
          { name: 'cmd', label: 'Cmd', width: '200px' },
          { name: 'enable', label: 'Enable/Disable', width: '120px' },
          { name: 'description', label: 'Description' }
        // { name: 'status', label: 'Status', width: '100px' }
        ],
        pagination: {
          page_num: 1,
          page_size: 10
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
        crawlConfirmDialogVisible: false
      }
    },
    computed: {
      ...mapState('schedule', [
        'scheduleList',
        'scheduleTotal',
        'scheduleForm'
      ]),
      ...mapState('spider', [
        'spiderList',
        'spiderVersionList',
        'spiderForm'
      ]),
      lang() {
        const lang = this.$store.state.lang.lang || window.localStorage.getItem('lang')
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
      }
    },
    created() {
      this.getScheduleList()
    },
    mounted() {
    },
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
        this.$refs.scheduleForm.validate(res => {
          if (res) {
            const form = Object.assign({}, this.scheduleForm)
            if (form.cron.search(/^\S+?\s\S+?\s\S+?\s\S+?\s\S+?$/) !== -1) {
              this.$set(form, 'cron', '0 ' + form.cron)
            }
            this.$set(form, 'enabled', form.enabled ? 1 : 0)
            if (this.isEdit) {
              this.$store.dispatch('schedule/editSchedule', form).then(response => {
                if (response.data.code !== 200) {
                  this.$message.error(response.data.message)
                  return
                }
                this.dialogVisible = false
                this.getScheduleList()
                this.$message.success(this.$t('The schedule has been saved'))
              })
            } else {
              this.$store.dispatch('schedule/addSchedule', form).then(response => {
                if (response.data.code !== 200) {
                  this.$message.error(response.data.message)
                  return
                }
                this.dialogVisible = false
                this.getScheduleList()
                this.$message.success(this.$t('The schedule has been added'))
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
        await this.$store.dispatch('spider/getSpiderData', row.spider_id)
        this.submitting = false
      },
      onRemove(row) {
        this.$confirm(this.$t('Are you sure to delete the schedule task?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('schedule/removeSchedule', row.id)
            .then(() => {
              setTimeout(() => {
                this.getScheduleList()
                this.$message.success(this.$t('The schedule has been removed'))
              }, 100)
            })
        }).catch(() => {
        })
        this.$st.sendEv('定时任务', '删除定时任务')
      },
      async getScheduleList() {
        await this.$store.dispatch('schedule/getScheduleList', this.pagination)
      },
      async onSelectSpider() {
        this.loadingSpiders = true
        await this.$store.dispatch('spider/getSpiderList')
        this.loadingSpiders = false
      },
      async onSelectSpiderVersion() {
        this.loadingVersions = true
        await this.$store.dispatch('spider/getSpiderVersionList', { spider_id: this.scheduleForm.spider_id })
        this.loadingVersions = false
      },
      async onEnabledChange(row) {
        let res
        if (row.enabled) {
          res = await this.$store.dispatch('schedule/enableSchedule', row.id)
        } else {
          res = await this.$store.dispatch('schedule/disableSchedule', row.id)
        }
        if (!res || res.data.code !== 200) {
          this.$message.error(this.$t(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule unsuccessful`))
        } else {
          this.$message.success(this.$t(`${row.enabled ? 'Enabling' : 'Disabling'} the schedule successful`))
        }
        this.$st.sendEv('定时任务', '启用/禁用')
      },
      onCronChange(value) {
        this.$set(this.scheduleForm, 'cron', value)
        this.$st.sendEv('定时任务', '配置Cron')
      },
      onCronDialogSubmit() {
        const valid = this.$refs['vue-cron-linux'].submit()
        if (valid) {
          this.cronDialogVisible = false
        }
      },
      async onSpiderChange(spiderId) {
        await this.$store.dispatch('spider/getSpiderData', spiderId)
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
      }
    }
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
