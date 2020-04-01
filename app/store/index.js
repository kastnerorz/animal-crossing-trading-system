import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

const store = () => new Vuex.Store({
  state: {
    user: {},
    quotation: {}
  },
  mutations: {
    setUser (state, user) {
      state.user = user
    },
    setQuotation (state, quotation) {
      state.quotation = quotation
    },
  }
})

export default store