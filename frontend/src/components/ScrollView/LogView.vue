<template>
  <div class="log-view-container">
    <div class="filter-wrapper">
      <div class="left">
        <el-switch
          v-model="isLogAutoScroll"
          :inactive-text="t('Auto-Scroll')"
          style="margin-right: 10px"
        />
        <!--        <el-input-->
        <!--          v-model="logKeyword"-->
        <!--          size="small"-->
        <!--          suffix-icon="el-icon-search"-->
        <!--          :placeholder="t('Search Log')"-->
        <!--          style="width: 240px; margin-right: 10px"-->
        <!--          clearable-->
        <!--        />-->
        <!--        <el-button-->
        <!--          size="small"-->
        <!--          type="primary"-->
        <!--          icon="el-icon-search"-->
        <!--          @click="onSearchLog"-->
        <!--        >-->
        <!--          {{ t('Search Log') }}-->
        <!--        </el-button>-->
      </div>
      <div class="right">
        <el-pagination
          size="small"
          :total="taskLogTotal"
          v-model:current-page="taskLogPage"
          :page-sizes="[100, 500, 1000, 5000]"
          v-model:page-size="taskLogPageSize"
          :page-count="3"
          layout="sizes, prev, pager, next"
        />
        <el-badge
          v-if="errorLogData.length > 0"
          :model-value="errorLogData.length"
        >
          <el-button
            type="danger"
            size="small"
            :icon="ElIconWarningOutline"
            @click="toggleErrors"
          >
            {{ t('Error Count') }}
          </el-button>
        </el-badge>
      </div>
    </div>
    <div class="content">
      <div
        :loading="isLogFetchLoading"
        class="log-view-wrapper"
        :class="isErrorsCollapsed ? 'errors-collapsed' : ''"
      >
        <virtual-list
          ref="log-view"
          class="log-view"
          :data-key="'index'"
          :data-sources="items"
          :data-component="itemComponent"
          :keeps="taskLogPageSize"
          :start="currentLogIndex - 1"
        />
      </div>
      <!--      <div-->
      <!--        v-show="!isErrorsCollapsed && !isErrorCollapsing"-->
      <!--        class="errors-wrapper"-->
      <!--        :class="isErrorsCollapsed ? 'collapsed' : ''"-->
      <!--      >-->
      <!--        <ul class="error-list">-->
      <!--          <li-->
      <!--            v-for="item in errorLogData"-->
      <!--            :key="item.index"-->
      <!--            class="error-item"-->
      <!--            :class="currentLogIndex === item.index ? 'active' : ''"-->
      <!--            @click="onClickError(item)"-->
      <!--          >-->
      <!--            <span class="line-content">-->
      <!--              {{ item.data }}-->
      <!--            </span>-->
      <!--          </li>-->
      <!--        </ul>-->
      <!--      </div>-->
    </div>
  </div>
</template>

<script>
import { WarningOutline as ElIconWarningOutline } from '@element-plus/icons'
import { $on, $off, $once, $emit } from '../../utils/gogocodeTransfer'
import * as Vue from 'vue'
import { mapState, mapGetters } from 'vuex'
import VirtualList from 'vue-virtual-scroll-list'
import Convert from 'ansi-to-html'
import hasAnsi from 'has-ansi'

import LogItem from './LogItem'
import { useI18n } from 'vue-i18n'

const convert = new Convert()
export default {
  setup(props) {
    const { t } = useI18n()
    return { t }
  },
  data() {
    return {
      itemComponent: LogItem,
      searchString: '',
      isScrolling: false,
      isScrolling2nd: false,
      errorRegex: this.$utils.log.errorRegex,
      currentOffset: 0,
      isErrorsCollapsed: true,
      isErrorCollapsing: false,
      ElIconWarningOutline,
    }
  },
  name: 'LogView',
  components: {
    VirtualList,
  },
  props: {
    data: {
      type: String,
      default: '',
    },
  },
  computed: {
    ...mapState('task', [
      'taskForm',
      'taskLogTotal',
      'logKeyword',
      'isLogFetchLoading',
      'errorLogData',
    ]),
    ...mapGetters('task', ['logData']),
    currentLogIndex: {
      get() {
        return this.$store.state.task.currentLogIndex
      },
      set(value) {
        this.$store.commit('task/SET_CURRENT_LOG_INDEX', value)
      },
    },
    logKeyword: {
      get() {
        return this.$store.state.task.logKeyword
      },
      set(value) {
        this.$store.commit('task/SET_LOG_KEYWORD', value)
      },
    },
    taskLogPage: {
      get() {
        return this.$store.state.task.taskLogPage
      },
      set(value) {
        this.$store.commit('task/SET_TASK_LOG_PAGE', value)
      },
    },
    taskLogPageSize: {
      get() {
        return this.$store.state.task.taskLogPageSize
      },
      set(value) {
        this.$store.commit('task/SET_TASK_LOG_PAGE_SIZE', value)
      },
    },
    isLogAutoScroll: {
      get() {
        return this.$store.state.task.isLogAutoScroll
      },
      set(value) {
        this.$store.commit('task/SET_IS_LOG_AUTO_SCROLL', value)
      },
    },
    isLogAutoFetch: {
      get() {
        return this.$store.state.task.isLogAutoFetch
      },
      set(value) {
        this.$store.commit('task/SET_IS_LOG_AUTO_FETCH', value)
      },
    },
    isLogFetchLoading: {
      get() {
        return this.$store.state.task.isLogFetchLoading
      },
      set(value) {
        this.$store.commit('task/SET_IS_LOG_FETCH_LOADING', value)
      },
    },
    items() {
      if (!this.logData || this.logData.length === 0) {
        return []
      }
      const filteredLogData = this.logData.filter((d) => {
        if (!this.searchString) return true
        return !!d.line_text
          .toLowerCase()
          .match(this.searchString.toLowerCase())
      })
      return filteredLogData.map((logItem) => {
        const isAnsi = hasAnsi(logItem.line_text)
        return {
          index: logItem.line_num,
          data: isAnsi ? convert.toHtml(logItem.line_text) : logItem.line_text,
          searchString: this.logKeyword,
          active: logItem.active,
          isAnsi,
        }
      })
    },
  },
  watch: {
    taskLogPage() {
      $emit(this, 'search')
      this.$st.sendEv('任务详情', '日志', '改变页数')
    },
    taskLogPageSize() {
      $emit(this, 'search')
      this.$st.sendEv('任务详情', '日志', '改变日志每页条数')
    },
    isLogAutoScroll() {
      if (this.isLogAutoScroll) {
        this.$store
          .dispatch('task/getTaskLogs', {
            id: this.$route.params.id,
            keyword: this.logKeyword,
          })
          .then(() => {
            this.toBottom()
          })
        this.$st.sendEv('任务详情', '日志', '点击自动滚动')
      } else {
        this.$st.sendEv('任务详情', '日志', '取消自动滚动')
      }
    },
  },
  mounted() {
    this.currentLogIndex = 0
    this.handle = setInterval(() => {
      if (this.isLogAutoScroll) {
        this.toBottom()
      }
    }, 200)
  },
  unmounted() {
    clearInterval(this.handle)
  },
  methods: {
    toBottom() {
      this.$el.querySelector('.log-view').scrollTo({ top: 99999999 })
    },
    toggleErrors() {
      this.isErrorsCollapsed = !this.isErrorsCollapsed
      this.isErrorCollapsing = true
      setTimeout(() => {
        this.isErrorCollapsing = false
      }, 300)
    },
    async onClickError(item) {
      const page = Math.ceil(item.line_num / this.taskLogPageSize)
      this.$store.commit('task/SET_LOG_KEYWORD', '')
      this.$store.commit('task/SET_TASK_LOG_PAGE', page)
      this.$store.commit('task/SET_IS_LOG_AUTO_SCROLL', false)
      this.$store.commit('task/SET_ACTIVE_ERROR_LOG_ITEM', item)
      $emit(this, 'search')
      this.$st.sendEv('任务详情', '日志', '点击错误日志')
    },
    onSearchLog() {
      $emit(this, 'search')
      this.$st.sendEv('任务详情', '日志', '搜索日志')
    },
  },
  emits: ['search'],
}
</script>

<style scoped>
.filter-wrapper {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
.content {
  display: block;
}
.log-view-wrapper {
  float: left;
  flex-basis: calc(100% - 240px);
  width: calc(100% - 300px);
  transition: width 0.3s;
}
.log-view-wrapper.errors-collapsed {
  flex-basis: 100%;
  width: 100%;
}
.log-view {
  margin-top: 0 !important;
  overflow-y: scroll !important;
  height: 600px;
  list-style: none;
  color: #a9b7c6;
  background: #2b2b2b;
  border: none;
}
.errors-wrapper {
  float: left;
  display: inline-block;
  margin: 0;
  padding: 0;
  flex-basis: 240px;
  width: 300px;
  transition: opacity 0.3s;
  border-top: 1px solid #dcdfe6;
  border-right: 1px solid #dcdfe6;
  border-bottom: 1px solid #dcdfe6;
  height: calc(100vh - 240px);
  font-size: 16px;
  overflow: auto;
}
.errors-wrapper.collapsed {
  width: 0;
}
.errors-wrapper .error-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.errors-wrapper .error-list .error-item {
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  border-bottom: 1px solid white;
  padding: 5px 0;
  background: #f56c6c;
  color: white;
  cursor: pointer;
}
.errors-wrapper .error-list .error-item.active {
  background: #e6a23c;
  font-weight: bolder;
  text-decoration: underline;
}
.errors-wrapper .error-list .error-item:hover {
  font-weight: bolder;
  text-decoration: underline;
}
.errors-wrapper .error-list .error-item .line-no {
  display: inline-block;
  text-align: right;
  width: 70px;
}
.errors-wrapper .error-list .error-item .line-content {
  display: inline;
  width: calc(100% - 70px);
  padding-left: 10px;
}
.right {
  display: flex;
  align-items: center;
}
.right .el-pagination {
  margin-right: 10px;
}
</style>
