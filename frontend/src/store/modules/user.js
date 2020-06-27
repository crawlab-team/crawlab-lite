import request from '../../api/request'

const user = {
  namespaced: true,

  state: {
    name: '',
    avatar: '',
    userList: [],
    globalVariableList: [],
    globalVariableForm: {},
    userForm: {},
    userInfo: undefined,
    adminPaths: [
      '/users'
    ],
    pageNum: 1,
    pageSize: 10,
    totalCount: 0
  },

  getters: {
    token() {
      return window.localStorage.getItem('token')
    }
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_USER_INFO: (state, value) => {
      state.userInfo = value
    },
    SET_GLOBAL_VARIABLE_LIST: (state, value) => {
      state.globalVariableList = value
    }
  },

  actions: {
    // 登录
    async login({ commit }, userInfo) {
      const username = userInfo.username.trim()
      const res = await request.post('/login', { username, password: userInfo.password })
      if (res.status === 200) {
        const token = res.data.data
        commit('SET_TOKEN', token)
        window.localStorage.setItem('token', token)
      }
      return res
    },

    // 获取用户信息
    getInfo({ commit, state }) {
      return request.get('/me')
        .then(response => {
          // ensure compatibility
          if (response.data.data.setting && !response.data.data.setting.max_error_log) {
            response.data.data.setting.max_error_log = 1000
          }
        })
    },

    // 登出
    logout({ commit, state }) {
      return new Promise((resolve, reject) => {
        window.localStorage.removeItem('token')
        commit('SET_TOKEN', '')
        resolve()
      })
    },

    // 新增全局变量
    addGlobalVariable({ commit, state }) {
      return request.put(`/variable`, state.globalVariableForm)
        .then(() => {
          state.globalVariableForm = {}
        })
    },
    // 获取全局变量列表
    getGlobalVariable({ commit, state }) {
      request.get('/variables').then((response) => {
        commit('SET_GLOBAL_VARIABLE_LIST', response.data.data)
      })
    },
    // 删除全局变量
    deleteGlobalVariable({ commit, state }, id) {
      return request.delete(`/variable/${id}`)
    }
  }
}

export default user
