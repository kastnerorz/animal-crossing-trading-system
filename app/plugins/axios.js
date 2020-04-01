import Vue from "vue";
import jsCookie from "js-cookie";
const vm = new Vue({});

export default function({ $axios, redirect }) {
  $axios.onRequest(config => {
    const auth = jsCookie.get("auth");
    config.headers = {
      'Content-Type': `application/json`,
      Authorization: `Bearer ${auth}`
    };
  });
  $axios.onError(error => {
    let msg = "";
    try {
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
  });
}
