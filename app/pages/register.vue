<template>
  <section class="form-wrapper">
    <b-field label="用户名">
      <b-input v-model="username"></b-input>
    </b-field>
    <b-field label="密码">
      <b-input v-model="password" type="password" maxlength="30"></b-input>
    </b-field>
    <b-field label="再次输入密码">
      <b-input v-model="rePassword" type="password" maxlength="30"></b-input>
    </b-field>
    <b-field label="昵称">
      <b-input v-model="nickname"></b-input>
    </b-field>
    <b-field label="Switch 好友编号">
      <b-input v-model="switchFriendCode"></b-input>
    </b-field>
    <b-button class="btn-reg" type="is-primary" @click="register"
      >注册</b-button
    >
  </section>
</template>

<script>
import {sha256} from 'js-sha256'
export default {
  name: 'Register',
  mounted() {},
  data() {
    return {
      username: '',
      password: '',
      rePassword: '',
      nickname: '',
      switchFriendCode: '',
    }
  },
  methods: {
    async register() {
      var hash = sha256.create()
      hash.update(this.password)
      let res = await this.$axios.$post('/users', {
        username: this.username,
        password: hash.hex(),
        nickname: this.nickname,
        switchFriendCode: this.switchFriendCode,
      })
      if (res.code == 201) {
        console.log('success')
      }
    },
  },
}
</script>
