import Vue from "vue";
import jsCookie from "js-cookie";
const vm = new Vue({});

export default function({store, redirect, app: { $axios }}) {
  $axios.onRequest(config => {
    const auth = jsCookie.get("auth");
    config.headers = {
      'Content-Type': `application/json`,
      Authorization: `Bearer ${auth}`
    };
  });
  $axios.onError(error => {
    store.commit('closeLoading')
    let msg = "";
    let status = ""
    try {
      status = error.response.status
      msg = error.response.data.msg;
    } catch (error) {
      msg = "网络错误！";
    }
    vm.$buefy.toast.open({
      duration: 3000,
      message: msg,
      position: "is-top",
      type: "is-danger"
    });
    if (status === 401) {
      jsCookie.remove('auth', { path: '' });
      setTimeout(()=>{
        redirect('/login')
      }, 1000)
    }
  });
}
