import Vue from 'vue'
import request from '../../api/request'

const state = {
  // list of spiders
  spiderList: [],

  spiderTotal: 0,

  // active spider data
  spiderForm: {},

  // upload form for importing spiders
  importForm: {
    url: '',
    type: 'github'
  },

  // spider overview stats
  overviewStats: {},

  // spider status stats
  statusStats: [],

  // spider daily stats
  dailyStats: [],

  // filters
  filterSite: '',

  // preview crawl data
  previewCrawlData: [],

  // template list
  templateList: [],

  // spider file tree
  fileTree: {},

  // config list ts
  configListTs: undefined
}

const getters = {}

const mutations = {
  SET_SPIDER_TOTAL(state, value) {
    state.spiderTotal = value
  },
  SET_SPIDER_FORM(state, value) {
    state.spiderForm = value
  },
  SET_SPIDER_LIST(state, value) {
    state.spiderList = value
  },
  SET_IMPORT_FORM(state, value) {
    state.importForm = value
  },
  SET_OVERVIEW_STATS(state, value) {
    state.overviewStats = value
  },
  SET_STATUS_STATS(state, value) {
    state.statusStats = value
  },
  SET_DAILY_STATS(state, value) {
    state.dailyStats = value
  },
  SET_FILTER_SITE(state, value) {
    state.filterSite = value
  },
  SET_PREVIEW_CRAWL_DATA(state, value) {
    state.previewCrawlData = value
  },
  SET_SPIDER_FORM_CONFIG_SETTINGS(state, payload) {
    const settings = {}
    payload.forEach(row => {
      settings[row.name] = row.value
    })
    Vue.set(state.spiderForm.config, 'settings', settings)
  },
  SET_TEMPLATE_LIST(state, value) {
    state.templateList = value
  },
  SET_FILE_TREE(state, value) {
    state.fileTree = value
  },
  SET_SPIDER_SCRAPY_SETTINGS(state, value) {
    state.spiderScrapySettings = value
  },
  SET_SPIDER_SCRAPY_ITEMS(state, value) {
    state.spiderScrapyItems = value
  },
  SET_SPIDER_SCRAPY_PIPELINES(state, value) {
    state.spiderScrapyPipelines = value
  },
  SET_CONFIG_LIST_TS(state, value) {
    state.configListTs = value
  }
}

const actions = {
  getSpiderList({ state, commit }, params = {}) {
    return request.get('/spiders', params)
      .then(response => {
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
  getSpiderData({ state, commit }, id) {
    return request.get(`/spiders/${id}`)
      .then(response => {
        const data = response.data.data
        commit('SET_SPIDER_FORM', data)
      })
  },
  crawlSpider({ state, dispatch }, payload) {
    const { spiderId, cmd } = payload
    return request.post(`/tasks`, {
      spider_id: spiderId,
      cmd: cmd
    })
  },
  crawlSelectedSpiders({ state, dispatch }, payload) {
    const { taskParams, runType, nodeIds } = payload
    return request.post(`/spiders-run`, {
      task_params: taskParams,
      run_type: runType,
      node_ids: nodeIds
    })
  },
  getTaskList({ state, commit }, id) {
    return request.get(`/tasks`, { 'spider_id': id })
      .then(response => {
        commit('task/SET_TASK_LIST',
          response.data.data ? response.data.data.map(d => {
            return d
          }).sort((a, b) => a.create_ts < b.create_ts ? 1 : -1) : [],
          { root: true })
      })
  },
  getDir({ state, commit }, path) {
    const id = state.spiderForm.id
    return request.get(`/spiders/${id}/dir`)
      .then(response => {
        commit('')
      })
  },
  importGithub({ state }) {
    const url = state.importForm.url
    return request.post('/spiders/import/github', { url })
  },
  getSpiderStats({ state, commit }) {
    return request.get(`/spiders/${state.spiderForm.id}/stats`)
      .then(response => {
        commit('SET_OVERVIEW_STATS', response.data.data.overview)
        // commit('SET_STATUS_STATS', response.data.task_count_by_status)
        commit('SET_DAILY_STATS', response.data.data.daily)
        // commit('SET_NODE_STATS', response.data.task_count_by_node)
      })
  },
  getPreviewCrawlData({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm.id}/preview_crawl`)
      .then(response => {
        commit('SET_PREVIEW_CRAWL_DATA', response.data.items)
      })
  },
  extractFields({ state, commit }) {
    return request.post(`/spiders/${state.spiderForm.id}/extract_fields`)
  },
  postConfigSpiderConfig({ state }) {
    return request.post(`/config_spiders/${state.spiderForm.id}/config`, state.spiderForm.config)
  },
  saveConfigSpiderSpiderfile({ state, rootState }) {
    const content = rootState.file.fileContent
    return request.post(`/config_spiders/${state.spiderForm.id}/spiderfile`, { content })
  },
  addConfigSpider({ state }) {
    return request.post(`/config_spiders`, state.spiderForm)
  },
  addSpider({ state }) {
    return request.post(`/spiders`, state.spiderForm)
  },
  async getTemplateList({ state, commit }) {
    const res = await request.get(`/config_spiders_templates`)
    commit('SET_TEMPLATE_LIST', res.data.data)
  },
  async getScheduleList({ state, commit }, payload) {
    const { id } = payload
    const res = await request.get(`/spiders/${id}/schedules`)
    let data = res.data.data
    if (data) {
      data = data.map(d => {
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
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
