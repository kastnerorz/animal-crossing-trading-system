<template>
  <div class="login-page">
    <TopMenu opt="login" />
    <section class="form-wrapper">
      <b-field label="用户名">
        <b-input v-model="username" placeholder="请输入用户名"></b-input>
      </b-field>
      <b-field label="密码">
        <b-input v-model="password" placeholder="请输入密码" @keyup.enter="login" type="password" maxlength="30"></b-input>
      </b-field>
      <b-button class="btn-login" type="is-primary" @click="login">登录</b-button>
      <b-button class="btn-login btn-register" type="is-primary" @click="goRegister">还没有账号</b-button>
    </section>
  </div>
</template>

<script>
import { sha256 } from "js-sha256";
import TopMenu from "../components/TopMenu";
import jsCookie from "js-cookie";
export default {
  components: { TopMenu },
  name: "Login",
  mounted() {},
  data() {
    return {
      username: "",
      password: ""
    };
  },
  methods: {
    async login() {
      this.$store.commit("setLoading");
      var hash = sha256.create();
      hash.update(this.password);
      const userInfo = {
        username: this.username,
        password: hash.hex()
      };
      const login = await this.$axios.$post("/login", userInfo);
      jsCookie.set("auth", login.token, { expires: 1 });
      const getMe = await this.$axios.$get("/me");
      this.$buefy.toast.open({
        duration: 2000,
        message: "登录成功，即将跳转",
        position: "is-top",
        type: "is-success"
      });
      setTimeout(() => {
        this.$store.commit("setUser", getMe);
        const cDate = new Date();
        // cDate.setDate(cDate.getDate() + 2);
        const cDay = cDate.getDay();
        const cHour = cDate.getHours();
        if (cDay !== 0) {
          this.$router.push("/sell");
        } else {
          this.$router.push("/");
        }
      }, 2000);
    },
    /**
     * 转去登录
     */
    goRegister() {
      this.$router.push("/register");
    }
  }
};
</script>
