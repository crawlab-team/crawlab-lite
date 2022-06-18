import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  spiderTotal: 0,

  // list of spider versions
  spiderVersionList: [],

  spiderVersionTotal: 0,

  // active spider data
  spiderForm: {},

  // spider overview stats
  overviewStats: {},

  // spider daily stats
  dailyStats: [],

  // spider file tree
  fileTree: {},
}

const getters = {}

const mutations = {
  SET_SPIDER_FORM(state, value) {
    state.spiderForm = value
  },
  SET_SPIDER_LIST(state, value) {
    state.spiderList = value
  },
  SET_SPIDER_TOTAL(state, value) {
    state.spiderTotal = value
  },
  SET_SPIDER_VERSION_LIST(state, value) {
    state.spiderVersionList = value
  },
  SET_SPIDER_VERSION_TOTAL(state, value) {
    state.spiderVersionTotal = value
  },
  SET_OVERVIEW_STATS(state, value) {
    state.overviewStats = value
  },
  SET_DAILY_STATS(state, value) {
    state.dailyStats = value
  },
  SET_FILE_TREE(state, value) {
    state.fileTree = value
  },
}

const actions = {
  getSpiderList({ state, commit }, params = {}) {
    return request.get('/spiders', params).then((response) => {
      if (!response || !response.data || !response.data.data) {
        return
      }
      commit('SET_SPIDER_LIST', response.data.data.list || [])
      commit('SET_SPIDER_TOTAL', response.data.data.total || 0)
    })
  },
  editSpider({ state, dispatch }) {
    return request.post(`/spiders/${state.spiderForm.id}`, state.spiderForm)
  },
  deleteSpider({ state, dispatch }, id) {
    return request.delete(`/spiders/${id}`)
  },
  getSpider({ state, commit }, id) {
    return request.get(`/spiders/${id}`).then((response) => {
      const data = response.data.data
      commit('SET_SPIDER_FORM', data)
    })
  },
  crawlSpider({ state, dispatch }, payload) {
    const { spiderId, spiderVersionId, cmd } = payload
    return request.post(`/tasks`, {
      spider_id: spiderId,
      spider_version_id: spiderVersionId,
      cmd: cmd,
    })
  },
  getDir({ state, commit }, path) {
    const id = state.spiderForm.id
    return request.get(`/spiders/${id}/dir`).then((response) => {
      commit('')
    })
  },
  getSpiderStats({ state, commit }) {
    return request
      .get(`/spiders/${state.spiderForm.id}/stats`)
      .then((response) => {
        commit('SET_OVERVIEW_STATS', response.data.data.overview)
        // commit('SET_STATUS_STATS', response.data.task_count_by_status)
        commit('SET_DAILY_STATS', response.data.data.daily)
        // commit('SET_NODE_STATS', response.data.task_count_by_node)
      })
  },
  addSpider({ state }) {
    return request.post(`/spiders`, state.spiderForm)
  },
  getSpiderVersionList({ state, commit }, params = {}) {
    const { spider_id } = params
    return request.get(`/spiders/${spider_id}/versions`).then((response) => {
      if (!response || !response.data || !response.data.data) {
        return
      }
      commit('SET_SPIDER_VERSION_LIST', response.data.data.list || [])
      commit('SET_SPIDER_VERSION_TOTAL', response.data.data.total || 0)
    })
  },
  deleteSpiderVersion({ state, dispatch }, payload) {
    const { spider_id, version_id } = payload
    return request.delete(`/spiders/${spider_id}/versions/${version_id}`)
  },
  async getScheduleList({ state, commit }, payload) {
    const { id } = payload
    const res = await request.get(`/spiders/${id}/schedules`)
    let data = res.data.data
    if (data) {
      data = data.map((d) => {
        const arr = d.cron.split(' ')
        arr.splice(0, 1)
        d.cron = arr.join(' ')
        return d
      })
    }
    commit('schedule/SET_SCHEDULE_LIST', data, { root: true })
  },
  async getFileTree({ state, commit }, payload) {
    const id = payload ? payload.id : state.spiderForm.id
    const res = await request.get(`/spiders/${id}/file/tree`)
    commit('SET_FILE_TREE', res.data.data)
  },
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
}
