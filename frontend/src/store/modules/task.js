import request from '../../api/request'

const state = {
  // TaskList
  taskList: [],
  taskTotal: 0,
  taskForm: {},
  taskResultsData: [],
  taskResultsColumns: [],
  taskResultsTotal: 0,
  // log
  currentLogIndex: 0,
  logKeyword: '',
  errorLogData: [],
  isLogAutoScroll: false,
  isLogAutoFetch: false,
  isLogFetchLoading: false,
  taskLogList: [],
  taskLogTotal: 0,
  taskLogPage: 1,
  taskLogPageSize: 5000,
  activeErrorLogItem: {},
  // results
  resultsPageNum: 1,
  resultsPageSize: 10
}

const getters = {
  taskResultsColumns(state) {
    if (!state.taskResultsData || !state.taskResultsData.length) {
      return []
    }
    const keys = []
    const item = state.taskResultsData[0]
    for (const key of item) {
      keys.push(key)
    }
    return keys
  },
  logData(state) {
    const data = state.taskLogList
      .map((d, i) => {
        return {
          active: state.currentLogIndex === i + 1,
          ...d
        }
      })
    if (state.taskForm && state.taskForm.status === 'RUNNING') {
      data.push({
        line_num: data.length + 1,
        line_text: '###LOG_END###'
      })
    }
    return data
  },
  errorLogData(state, getters) {
    const idxList = getters.logData.map(d => d.id)
    return state.errorLogData.map(d => {
      const idx = idxList.indexOf(d.id)
      d.index = getters.logData[idx].index
      return d
    })
  }
}

const mutations = {
  SET_TASK_FORM(state, value) {
    state.taskForm = value
  },
  SET_TASK_LIST(state, value) {
    state.taskList = value
  },
  SET_TASK_LOG_LIST(state, value) {
    state.taskLogList = value
  },
  SET_TASK_LOG_TOTAL(state, value) {
    state.taskLogTotal = value
  },
  SET_CURRENT_LOG_INDEX(state, value) {
    state.currentLogIndex = value
  },
  SET_TASK_RESULTS_DATA(state, value) {
    state.taskResultsData = value
  },
  SET_TASK_RESULTS_COLUMNS(state, value) {
    state.taskResultsColumns = value
  },
  SET_TASK_TOTAL(state, value) {
    state.taskTotal = value
  },
  SET_TASK_RESULTS_TOTAL(state, value) {
    state.taskResultsTotal = value
  },
  SET_LOG_KEYWORD(state, value) {
    state.logKeyword = value
  },
  SET_ERROR_LOG_DATA(state, value) {
    state.errorLogData = value
  },
  SET_TASK_LOG_PAGE(state, value) {
    state.taskLogPage = value
  },
  SET_TASK_LOG_PAGE_SIZE(state, value) {
    state.taskLogPageSize = value
  },
  SET_IS_LOG_AUTO_SCROLL(state, value) {
    state.isLogAutoScroll = value
  },
  SET_IS_LOG_AUTO_FETCH(state, value) {
    state.isLogAutoFetch = value
  },
  SET_IS_LOG_FETCH_LOADING(state, value) {
    state.isLogFetchLoading = value
  },
  SET_ACTIVE_ERROR_LOG_ITEM(state, value) {
    state.activeErrorLogItem = value
  }
}

const actions = {
  getTaskData({ state, dispatch, commit }, id) {
    return request.get(`/tasks/${id}`)
      .then(response => {
        const data = response.data.data
        commit('SET_TASK_FORM', data)
        // dispatch('spider/getSpiderData', data.spider_id, { root: true })
      })
  },
  getTaskList({ state, commit }, params = {}) {
    return request.get('/tasks', params)
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_TASK_LIST', response.data.data.list || [])
        commit('SET_TASK_TOTAL', response.data.data.total || 0)
      })
  },
  deleteTask({ state, dispatch }, id) {
    return request.delete(`/tasks/${id}`)
  },
  deleteTaskMultiple({ state }, ids) {
    return request.delete(`/tasks`, {
      ids: ids
    })
  },
  restartTask({ state, dispatch }, id) {
    return request.post(`/tasks/${id}/restart`)
  },
  getTaskLogs({ state, commit }, { id, keyword }) {
    return request.get(`/tasks/${id}/logs`, {
      keyword,
      page_num: state.taskLogPage,
      page_size: state.taskLogPageSize
    })
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_TASK_LOG_LIST', response.data.data.list || [])
        commit('SET_TASK_LOG_TOTAL', response.data.data.total || 0)

        // auto switch to next page if not reaching the end
        if (state.isLogAutoScroll && state.taskLogTotal > (state.taskLogPage * state.taskLogPageSize)) {
          commit('SET_TASK_LOG_PAGE', Math.ceil(state.taskLogTotal / state.taskLogPageSize))
        }
      })
  },
  getTaskErrorLog({ state, commit }, id) {
    return request.get(`/tasks/${id}/error-log`, {})
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_ERROR_LOG_DATA', response.data.data.list || [])
      })
  },
  getTaskResults({ state, commit }, id) {
    return request.get(`/tasks/${id}/results`, {
      page_num: state.resultsPageNum,
      page_size: state.resultsPageSize
    })
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_TASK_RESULTS_DATA', response.data.data.list || [])
        commit('SET_TASK_RESULTS_TOTAL', response.data.data.total || 0)
      })
  },
  cancelTask({ state, dispatch }, id) {
    return new Promise(resolve => {
      request.post(`/tasks/${id}/cancel`)
        .then(res => {
          dispatch('getTaskData', id)
          resolve(res)
        })
    })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
