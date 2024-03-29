<template>
  <div class="app-container">
    <el-card style="border-radius: 0">
      <!--filter-->
      <div class="filter">
        <div class="left">
          <el-form :model="filter" label-position="right" inline>
            <el-form-item prop="spider_id" :label="t('Spider')">
              <el-select
                v-model="filter.spider_id"
                size="small"
                :placeholder="t('All')"
                :disabled="isFilterSpiderDisabled"
                @focus="onSelectFilterSpider"
                @change="onFilterChange"
              >
                <el-option value="" :label="t('All')" />
                <el-option
                  v-for="spider in spiderList"
                  :key="spider.id"
                  :value="spider.id"
                  :label="spider.name"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="status" :label="t('Status')">
              <el-select
                v-model="filter.status"
                size="small"
                :placeholder="t('Status')"
                @change="onFilterChange"
              >
                <el-option value="" :label="t('All')" />
                <el-option value="PADDING" :label="t('Pending')" />
                <el-option value="RUNNING" :label="t('Running')" />
                <el-option value="FINISHED" :label="t('Finished')" />
                <el-option value="ERROR" :label="t('Error')" />
                <el-option value="CANCELLED" :label="t('Cancelled')" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>
        <div class="right">
          <el-button
            v-if="multipleSelection.length > 0"
            :icon="ElIconDelete"
            class="btn-delete"
            size="small"
            type="danger"
            @click="onRemoveMultipleTask"
          >
            {{ t('Remove') }}
          </el-button>
        </div>
      </div>
      <!--./filter-->

      <!--table list-->
      <el-table
        ref="table"
        :data="filteredTableData"
        class="table"
        :header-cell-style="{ background: 'rgb(48, 65, 86)', color: 'white' }"
        border
        row-key="id"
        @row-click="onRowClick"
        @selection-change="onSelectionChange"
      >
        <el-table-column
          type="selection"
          width="45"
          align="center"
          reserve-selection
        />
        <template v-for="col in columns">
          <el-table-column
            v-if="col.name === 'spider_name'"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
          >
            <template v-slot="scope">
              <!--<a href="javascript:" class="a-tag" @click="onClickSpider(scope.row)">{{ scope.row[col.name] }}</a>-->
              {{ scope.row[col.name] }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name.match(/_ts$/)"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ getTime(scope.row[col.name]) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'wait_duration'"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ getWaitDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'runtime_duration'"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ getRuntimeDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'total_duration'"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              {{ getTotalDuration(scope.row) }}
            </template>
          </el-table-column>
          <el-table-column
            v-else-if="col.name === 'status'"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          >
            <template v-slot="scope">
              <template v-if="scope.row[col.name] === 'ERROR'">
                <el-tooltip :content="scope.row.error" placement="top">
                  <status-tag :status="scope.row[col.name]" />
                </el-tooltip>
              </template>
              <status-tag v-else :status="scope.row[col.name]" />
            </template>
          </el-table-column>
          <el-table-column
            v-else
            :property="col.name"
            :label="t(col.label)"
            :sortable="col.sortable"
            :align="col.align"
            :width="col.width"
          />
        </template>
        <el-table-column
          :label="t('Action')"
          align="left"
          fixed="right"
          width="140px"
        >
          <template v-slot="scope">
            <el-tooltip :content="t('View')" placement="top">
              <el-button
                type="primary"
                :icon="ElIconSearch"
                size="small"
                @click="onView(scope.row)"
              />
            </el-tooltip>
            <el-tooltip :content="t('Restart')" placement="top">
              <el-button
                type="warning"
                :icon="ElIconRefresh"
                size="small"
                @click="onRestart(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip
              v-if="['PENDING', 'RUNNING'].includes(scope.row.status)"
              :content="t('Cancel')"
              placement="top"
            >
              <el-button
                type="danger"
                :icon="ElIconVideoPause"
                size="small"
                @click="onCancel(scope.row, $event)"
              />
            </el-tooltip>
            <el-tooltip v-else :content="t('Remove')" placement="top">
              <el-button
                type="danger"
                :icon="ElIconDelete"
                size="small"
                @click="onRemove(scope.row, $event)"
              />
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page_num"
          :page-sizes="[10, 20, 50, 100]"
          v-model:page-size="pagination.page_size"
          layout="sizes, prev, pager, next"
          :total="taskTotal"
          @current-change="getTaskList"
          @size-change="getTaskList"
        />
      </div>
      <!--./table list-->
    </el-card>
  </div>
</template>

<script>
import {
  Delete as ElIconDelete,
  Search as ElIconSearch,
  Refresh as ElIconRefresh,
  VideoPause as ElIconVideoPause,
} from '@element-plus/icons'
import * as Vue from 'vue'
import { mapState } from 'vuex'
import dayjs from 'dayjs'
import StatusTag from '../../components/Status/StatusTag'
import { useI18n } from 'vue-i18n'

export default {
  setup(props) {
    const { t } = useI18n()
    const inst = Vue.getCurrentInstance()
    inst.t = t
    return { t }
  },
  data() {
    return {
      filter: {},
      pagination: {
        page_num: 1,
        page_size: 10,
      },
      refreshHandle: undefined,
      isEditMode: false,
      dialogVisible: false,
      loadingSpiders: false,
      multipleSelection: [],
      // table columns
      columns: [
        { name: 'spider_name', label: 'Spider', width: '120' },
        { name: 'status', label: 'Status', width: '120' },
        { name: 'cmd', label: 'Cmd', width: '240' },
        { name: 'start_ts', label: 'Start Time', width: '140' },
        { name: 'finish_ts', label: 'Finish Time', width: '140' },
        { name: 'wait_duration', label: 'Wait Duration (sec)', align: 'right' },
        {
          name: 'runtime_duration',
          label: 'Runtime Duration (sec)',
          align: 'right',
        },
        {
          name: 'total_duration',
          label: 'Total Duration (sec)',
          align: 'right',
        },
        // { name: 'result_count', label: 'Results Count', align: 'right' }
      ],
      isFilterSpiderDisabled: false,
      ElIconDelete,
      ElIconSearch,
      ElIconRefresh,
      ElIconVideoPause,
    }
  },
  name: 'TaskList',
  components: {
    StatusTag,
  },
  computed: {
    ...mapState('task', ['taskList', 'taskTotal', 'taskForm']),
    ...mapState('spider', ['spiderList']),
    filteredTableData() {
      return this.taskList
        .map((d) => d)
        .sort((a, b) => (a.create_ts < b.create_ts ? 1 : -1))
        .filter((d) => {
          // keyword
          if (!this.filter.keyword) return true
          for (let i = 0; i < this.columns.length; i++) {
            const colName = this.columns[i].name
            if (
              d[colName] &&
              d[colName]
                .toLowerCase()
                .indexOf(this.filter.keyword.toLowerCase()) > -1
            ) {
              return true
            }
          }
          return false
        })
    },
  },
  created() {
    this.getTaskList()
  },
  mounted() {
    this.refreshHandle = setInterval(() => {
      this.getTaskList()
    }, 5000)
  },
  unmounted() {
    clearInterval(this.refreshHandle)
  },
  methods: {
    onCancel(row, ev) {
      ev.stopPropagation()
      this.$confirm(
        this.t('Are you sure to cancel this task?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.$store.dispatch('task/cancelTask', row.id).then(() => {
            this.$message({
              type: 'success',
              message: this.t('Canceled successfully'),
            })
            this.getTaskList()
          })
          this.$st.sendEv('任务列表', '取消任务')
        })
        .catch(() => {})
    },
    onRemove(row, ev) {
      ev.stopPropagation()
      this.$confirm(
        this.t('Are you sure to delete this task?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.$store.dispatch('task/deleteTask', row.id).then(() => {
            this.$message({
              type: 'success',
              message: this.t('Deleted successfully'),
            })
            this.getTaskList()
          })
          this.$st.sendEv('任务列表', '删除任务')
        })
        .catch(() => {})
    },
    onRemoveMultipleTask() {
      this.$confirm(
        this.t('Are you sure to delete these tasks?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          const ids = this.multipleSelection.map((item) => item.id)
          this.$store.dispatch('task/deleteTaskMultiple', ids).then((resp) => {
            this.$message({
              type: 'success',
              message: this.t('Deleted successfully'),
            })
            this.getTaskList()
            this.$refs['table'].clearSelection()
          })
          this.$st.sendEv('任务列表', '批量删除任务')
        })
        .catch(() => {})
    },
    onRestart(row, ev) {
      ev.stopPropagation()
      this.$confirm(
        this.t('Are you sure to restart this task?'),
        this.t('Notification'),
        {
          confirmButtonText: this.t('Confirm'),
          cancelButtonText: this.t('Cancel'),
          type: 'warning',
        },
      )
        .then(() => {
          this.$store.dispatch('task/restartTask', row.id).then(() => {
            this.$message({
              type: 'success',
              message: this.t('Restarted successfully'),
            })
            setTimeout(this.getTaskList, 500)
          })
          this.$st.sendEv('任务列表', '重新开始任务')
        })
        .catch(() => {})
    },
    onView(row) {
      this.$router.push(`/tasks/${row.id}`)
      this.$st.sendEv('任务列表', '查看任务')
    },
    onClickSpider(row) {
      this.$router.push(`/spiders/${row.spider_id}`)
      this.$st.sendEv('任务列表', '点击爬虫详情')
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
      if (column.label !== this.t('Action')) {
        this.onView(row)
      }
    },
    onSelectionChange(val) {
      this.multipleSelection = val
    },
    async getTaskList() {
      const params = Object.assign({}, this.pagination, this.filter)
      await this.$store.dispatch('task/getTaskList', params)
    },
    async onSelectFilterSpider() {
      this.loadingSpiders = true
      await this.$store.dispatch('spider/getSpiderList')
      this.loadingSpiders = false
    },
    async onFilterChange() {
      await this.getTaskList()
      this.$st.sendEv('任务列表', '筛选任务')
    },
  },
}
</script>

<style lang="scss" scoped>
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
