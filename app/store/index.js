import Vue from "vue";
import Vuex from "vuex";
const vm = new Vue({});
Vue.use(Vuex);

const store = () =>
  new Vuex.Store({
    state: {
      loadingComponent: "",
      hasApplicationNew: false,
      user: {},
      quotation: {},
    },
    mutations: {
      setUser(state, user) {
        state.user = user;
      },
      setQuotation(state, quotation) {
        state.quotation = quotation;
      },
      setHasApplicationNew(state, hasApplicationNew) {
        state.hasApplicationNew = hasApplicationNew;
      },
      setLoading(state) {
        if (typeof state.loadingComponent === "string") {
          state.loadingComponent = vm.$buefy.loading.open();
        }
      },
      closeLoading(state) {
        if (typeof state.loadingComponent === "object") {
          state.loadingComponent.close();
          state.loadingComponent = "";
        }
      },
    },
  });

export default store;
