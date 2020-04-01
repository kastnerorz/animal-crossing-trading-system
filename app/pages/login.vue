<template>
  <div class="login-page">
    <TopMenu opt="login" />
    <section class="form-wrapper">
      <b-field label="用户名">
        <b-input v-model="username" placeholder="请输入用户名"></b-input>
      </b-field>
      <b-field label="密码">
        <b-input v-model="password" placeholder="请输入密码" type="password" maxlength="30"></b-input>
      </b-field>
      <b-button class="btn-login" type="is-primary" @click="login">登录</b-button>
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
      username: "测试 1",
      password: "1234567"
    };
  },
  methods: {
    async login() {
      var hash = sha256.create();
      hash.update(this.password);
      const userInfo = {
        username: this.username,
        password: hash.hex()
      };
      const login = await this.$axios.$post("/login", userInfo);
      if (login) {
        jsCookie.set("auth", login.token, { expires: 1 });
        jsCookie.set("username", this.username, { expires: 1 });
        this.$router.push("/");
      }
    }
  }
};
</script>
