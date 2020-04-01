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
    login() {
      var hash = sha256.create();
      hash.update(this.password);
      this.$axios
        .$post("/login", {
          username: this.username,
          password: hash.hex()
        })
        .then(res => {
          this.$router.push("/index");
        })
        .catch(err => {
          this.$buefy.toast.open({
            message: "用户名或密码错误",
            type: "is-danger"
          });
        });
    }
  }
};
</script>
