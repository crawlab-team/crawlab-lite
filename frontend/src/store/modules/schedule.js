import request from '../../api/request'

const state = {
  scheduleList: [],
  scheduleForm: {}
}

const getters = {}

const mutations = {
  SET_SCHEDULE_LIST(state, value) {
    state.scheduleList = value
  },
  SET_SCHEDULE_FORM(state, value) {
    state.scheduleForm = value
  }
}

const actions = {
  getScheduleList({ state, commit }) {
    request.get('/schedules')
      .then(response => {
        if (!response || !response.data || !response.data.data) {
          return
        }
        commit('SET_SCHEDULE_LIST', response.data.data.list || [])
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
