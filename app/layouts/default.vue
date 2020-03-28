<template>
  <div>
    <!-- <b-navbar>
        <template slot="brand">
            <b-navbar-item tag="router-link" :to="{ path: '/' }">
                <img
                    src="https://raw.githubusercontent.com/buefy/buefy/dev/static/img/buefy-logo.png"
                    alt="Lightweight UI components for Vue.js based on Bulma"
                >
            </b-navbar-item>
            <b-navbar-item ><b-icon pack="fab" icon="github"></b-icon></b-navbar-item>
        </template>
        <template slot="start">
            <b-navbar-item href="/">
                买盘
            </b-navbar-item>
            <b-navbar-item href="#">
                卖盘
            </b-navbar-item>
            <b-navbar-item href="#">
                GitHub
            </b-navbar-item>
        </template>

        <template slot="end">
            <b-navbar-item tag="div">
                <div class="buttons">
                    <a class="button is-primary" href="/register">
                        <strong>注册</strong>
                    </a>
                    <a class="button is-light" href="/login">
                        登录
                    </a>
                </div>
            </b-navbar-item>
        </template>
    </b-navbar> -->

    <section class="main-content columns">
      <div class="container column is-10">
        <nuxt />
      </div>
    </section>
  </div>
</template>

<script>
export default {
  data() {
    return {
      items: [
        {
          title: 'Home',
          icon: 'home',
          to: {name: 'index'},
        },
        {
          title: 'Inspire',
          icon: 'lightbulb',
          to: {name: 'inspire'},
        },
      ],
    }
  },
}
</script>
