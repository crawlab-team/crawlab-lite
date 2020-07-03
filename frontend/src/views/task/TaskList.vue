<template>
  <div class="app-container">
    <!--tour-->
    <v-tour
      name="task-list"
      :steps="tourSteps"
      :callbacks="tourCallbacks"
      :options="$utils.tour.getOptions(true)"
    />
    <!--./tour-->

    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form class="filter-form" :model="filter" label-position="right" inline>
            <el-form-item prop="spider_id" :label="$t('Spider')">
              <el-select
                v-model="filter.spider_id"
                size="small"
                :placeholder="$t('Spider')"
                :disabled="isFilterSpiderDisabled"
                @change="onFilterChange"
              >
                <el-option value="" :label="$t('All')" />
                <el-option v-for="spider in spiderList" :key="spider.id" :value="spider.id" :label="spider.name" />
              </el-select>
            </el-form-item>
            <el-form-item prop="status" :label="$t('Status')">
              <el-select v-model="filter.status" size="small" :placeholder="$t('Status')" @change="onFilterChange">
                <el-option value="" :label="$t('All')" />
                <el-option value="PADDING" :label="$t('Pending')" />
                <el-option value="RUNNING" :label="$t('Running')" />
                <el-option value="FINISHED" :label="$t('Finished')" />
                <el-option value="ERROR" :label="$t('Error')" />
                <el-option value="CANCELLED" :label="$t('Cancelled')" />
                <el-option value="ABNORMAL" :label="$t('Abnormal')" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button v-if="this.multipleSelection.length" class="btn-delete" size="small" type="danger" @click="onRemoveMultipleTask">
            删除任务
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table
        ref="table"
        :data="filteredTableData"
        class="table"
        :header-cell-style="{background:'rgb(48, 65, 86)',color:'white'}"
        border
        row-key="id"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="45" align="center" reserve-selection />
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'spider_name'"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
          >
            <template slot-scope="scope">
              <!--<a href="javascript:" class="a-tag" @click="onClickSpider(scope.row)">{{ scope.row[col.name] }}</a>-->
              {{ scope.row[col.name] }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getTime(scope.row[col.name]) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'wait_duration'"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getWaitDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'runtime_duration'"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getRuntimeDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'total_duration'"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              {{ getTotalDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'status'"
            :key="col.name"
            :label="$t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template slot-scope="scope">
              <template
                v-if="scope.row[col.name] === 'ERROR'"
              >
                <el-tooltip :content="scope.row.error" placement="top">
                  <status-tag :status="scope.row[col.name]" />
                </el-tooltip>
              </template>
              <status-tag v-else :status="scope.row[col.name]" />
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
          />
        </template>
        <el-table-column :label="$t('Action')" align="left" fixed="right" width="90px">
          <template slot-scope="scope">
            <el-tooltip :content="$t('Restart')" placement="top">
              <el-button
                type="warning"
                icon="el-icon-refresh"
                size="mini"
                @click="onRestart(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip :content="$t('Remove')" placement="top">
              <el-button
                type="danger"
                icon="el-icon-delete"
                size="mini"
                @click="onRemove(scope.row, $event)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          :current-page.sync="pageNum"
          :page-sizes="[10, 20, 50, 100]"
          :page-size.sync="pageSize"
          layout="sizes, prev, pager, next"
          :total="taskListTotalCount"
          @current-change="onPageChange"
          @size-change="onPageChange"
        />
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
  import { mapState } from 'vuex'
  import dayjs from 'dayjs'
  import StatusTag from '../../components/Status/StatusTag'

  export default {
    name: 'TaskList',
    components: {
      StatusTag
    },
    data() {
      return {
        // setInterval handle
        handle: undefined,

        // determine if is edit mode
        isEditMode: false,

        // dialog visibility
        dialogVisible: false,

        // table columns
        columns: [
          { name: 'spider_name', label: 'Spider', width: '120' },
          { name: 'status', label: 'Status', width: '120' },
          { name: 'cmd', label: 'Cmd', width: '120' },
          { name: 'start_ts', label: 'Start Time', width: '100' },
          { name: 'finish_ts', label: 'Finish Time', width: '100' },
          { name: 'wait_duration', label: 'Wait Duration (sec)', align: 'right' },
          { name: 'runtime_duration', label: 'Runtime Duration (sec)', align: 'right' },
          { name: 'total_duration', label: 'Total Duration (sec)', width: '80', align: 'right' }
        // { name: 'result_count', label: 'Results Count', width: '80', align: 'right' }
        ],

        multipleSelection: [],

        // tutorial
        tourSteps: [
          {
            target: '.filter-form',
            content: this.$t('You can filter tasks from this area.')
          },
          {
            target: '.table',
            content: this.$t('This is a list of spider tasks executed sorted in a time descending order.')
          },
          // {
          //   target: '.table .el-table__body-wrapper tr:nth-child(1)',
          //   content: this.$t('Click the row to or the view button to view the task detail.')
          // },
          {
            target: '.table tr td:nth-child(1)',
            content: this.$t('Tick and select the tasks you would like to delete in batches.'),
            params: {
              placement: 'right'
            }
          }
        ],

        tourCallbacks: {
          onStop: () => {
            this.$utils.tour.finishTour('task-list')
          },
          onPreviousStep: (currentStep) => {
            this.$utils.tour.prevStep('task-list', currentStep)
          },
          onNextStep: (currentStep) => {
            this.$utils.tour.nextStep('task-list', currentStep)
          }
        },

        isFilterSpiderDisabled: false
      }
    },
    computed: {
      ...mapState('task', [
        'filter',
        'taskList',
        'taskListTotalCount',
        'taskForm'
      ]),
      ...mapState('spider', [
        'spiderList'
      ]),
      pageNum: {
        get() {
          return this.$store.state.task.pageNum
        },
        set(value) {
          this.$store.commit('task/SET_PAGE_NUM', value)
        }
      },
      pageSize: {
        get() {
          return this.$store.state.task.pageSize
        },
        set(value) {
          this.$store.commit('task/SET_PAGE_SIZE', value)
        }
      },
      filteredTableData() {
        return this.taskList
          .map(d => d)
          .sort((a, b) => a.create_ts < b.create_ts ? 1 : -1)
          .filter(d => {
            // keyword
            if (!this.filter.keyword) return true
            for (let i = 0; i < this.columns.length; i++) {
              const colName = this.columns[i].name
              if (d[colName] && d[colName].toLowerCase().indexOf(this.filter.keyword.toLowerCase()) > -1) {
                return true
              }
            }
            return false
          })
      }
    },
    created() {
      this.$store.dispatch('task/resetFilter')
      this.$store.dispatch('task/getTaskList')
      this.$store.dispatch('spider/getSpiderList')
    },
    mounted() {
      this.handle = setInterval(() => {
        this.$store.dispatch('task/getTaskList')
      }, 5000)

      if (!this.$utils.tour.isFinishedTour('task-list')) {
        this.$utils.tour.startTour(this, 'task-list')
      }
    },
    destroyed() {
      clearInterval(this.handle)
    },
    methods: {
      onRefresh() {
        this.$store.dispatch('task/getTaskList')
        this.$st.sendEv('任务列表', '搜索')
      },
      onRemoveMultipleTask() {
        if (this.multipleSelection.length === 0) {
          this.$message({
            type: 'error',
            message: '请选择要删除的任务'
          })
          return
        }
        this.$confirm('确定删除任务', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          const ids = this.multipleSelection.map(item => item.id)
          this.$store.dispatch('task/deleteTaskMultiple', ids).then((resp) => {
            if (resp.data.status === 'ok') {
              this.$message({
                type: 'success',
                message: '删除任务成功'
              })
              this.$store.dispatch('task/getTaskList')
              this.$refs['table'].clearSelection()
              return
            }
            this.$message({
              type: 'error',
              message: resp.data.message
            })
          })
        }).catch(() => {
        })
      },
      onRemove(row, ev) {
        ev.stopPropagation()
        this.$confirm(this.$t('Are you sure to delete this task?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('task/deleteTask', row.id)
            .then(() => {
              this.$message({
                type: 'success',
                message: this.$t('Deleted successfully')
              })
            })
          this.$st.sendEv('任务列表', '删除任务')
        })
      },
      onRestart(row, ev) {
        ev.stopPropagation()
        this.$confirm(this.$t('Are you sure to restart this task?'), this.$t('Notification'), {
          confirmButtonText: this.$t('Confirm'),
          cancelButtonText: this.$t('Cancel'),
          type: 'warning'
        }).then(() => {
          this.$store.dispatch('task/restartTask', row.id)
            .then(() => {
              this.$message({
                type: 'success',
                message: this.$t('Restarted successfully')
              })
            })
          this.$st.sendEv('任务列表', '重新开始任务')
        })
      },
      onView(row) {
        this.$router.push(`/tasks/${row.id}`)
        this.$st.sendEv('任务列表', '查看任务')
      },
      onClickSpider(row) {
        this.$router.push(`/spiders/${row.spider_id}`)
        this.$st.sendEv('任务列表', '点击爬虫详情')
      },
      onPageChange() {
        setTimeout(() => {
          this.$store.dispatch('task/getTaskList')
        }, 0)
      },
      getTime(str) {
        if (str.match('^0001')) return 'NA'
        return dayjs(str).format('YYYY-MM-DD HH:mm:ss')
      },
      getWaitDuration(row) {
        if (row.start_ts.match('^0001')) return 'NA'
        return dayjs(row.start_ts).diff(row.create_ts, 'second')
      },
      getRuntimeDuration(row) {
        if (row.finish_ts.match('^0001')) return 'NA'
        return dayjs(row.finish_ts).diff(row.start_ts, 'second')
      },
      getTotalDuration(row) {
        if (row.finish_ts.match('^0001')) return 'NA'
        return dayjs(row.finish_ts).diff(row.create_ts, 'second')
      },
      onRowClick(row, event, column) {
        if (column.label !== this.$t('Action')) {
          this.onView(row)
        }
      },
      onSelectionChange(val) {
        this.multipleSelection = val
      },
      onFilterChange() {
        this.$store.dispatch('task/getTaskList')
        this.$st.sendEv('任务列表', '筛选任务')
      }
    }
  }
</script>

<style scoped lang="scss">
  .el-dialog {
    .el-select {
      width: 100%;
    }
  }

  .filter {
    display: flex;
    justify-content: space-between;

    .left {
      .filter-select {
        width: 180px;
        margin-right: 10px;
      }
    }

    .filter-search {
      width: 240px;
    }

    .add {
    }
  }

  .table {
    margin-top: 8px;
    border-radius: 5px;

    .el-button {
      padding: 7px;
    }
  }

  .delete-confirm {
    background-color: red;
  }

  .el-table .a-tag {
    text-decoration: underline;
  }

  .pagination {
    margin-top: 10px;
    text-align: right;
  }
</style>

<style scoped>
  .el-table >>> tr {
    cursor: pointer;
  }

  .el-table >>> .el-badge .el-badge__content {
    font-size: 7px;
  }
</style>
