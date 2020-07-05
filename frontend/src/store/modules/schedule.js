import request from '../../api/request'

const state = {
  scheduleList: [],
  scheduleTotal: 0,
  scheduleForm: {}
}

const getters = {}

const mutations = {
  SET_SCHEDULE_LIST(state, value) {
    state.scheduleList = value
  },
  SET_SCHEDULE_TOTAL(state, value) {
    state.scheduleTotal = value
  },
  SET_SCHEDULE_FORM(state, value) {
    state.scheduleForm = value
  }
}

const actions = {
  getScheduleList({ state, commit }, params = {}) {
    request.get('/schedules', params)
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_SCHEDULE_LIST', response.data.data.list || [])
        commit('SET_SCHEDULE_TOTAL', response.data.data.total || 0)
      })
  },
  addSchedule({ state }, payload) {
    return request.post('/schedules', payload || state.scheduleForm)
  },
  editSchedule({ state }, payload) {
    return request.put(`/schedules/${payload.id}`, payload || state.scheduleForm)
  },
  removeSchedule({ state }, id) {
    return request.delete(`/schedules/${id}`)
  },
  enableSchedule({ state, dispatch }, id) {
    return request.put(`/schedules/${id}`, { enabled: 1 })
  },
  disableSchedule({ state, dispatch }, id) {
    return request.put(`/schedules/${id}`, { enabled: 2 })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
