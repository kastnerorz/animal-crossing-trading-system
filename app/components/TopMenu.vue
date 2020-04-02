<template>
  <div class="menu-container">
    <div class="menu-header">
      <ICON type="logo" />
      <span @click="signOut">
        <ICON type="github" />
      </span>
    </div>
    <div class="opts">
      <n-link to="/">
        <ICON :type="gennerateIcon('buy')" />
      </n-link>
      <n-link to="/sell">
        <ICON :type="gennerateIcon('sell')" />
      </n-link>
      <n-link v-if="!isLogin" to="/login">
        <ICON :type="gennerateIcon('login')" />
      </n-link>
      <n-link v-if="!isLogin" to="/register">
        <ICON :type="gennerateIcon('reg')" />
      </n-link>
    </div>
  </div>
</template>

<script>
import jsCookie from "js-cookie";
import ICON from "./ICON";

export default {
  components: { ICON },
  data() {
    return {};
  },
  props: { opt: { type: String, required: true } },

  computed: {
    isLogin() {
      return !!this.$store.state.user.username;
    }
  },
  methods: {
    gennerateIcon(type) {
      return this.opt === type ? `${type}On` : type;
    },
    signOut() {
      if (!this.isLogin) {
        return;
      }
      this.$store.commit('setLoading')
      jsCookie.remove("auth", { path: "" });
      this.$buefy.toast.open({
        duration: 2000,
        message: "退出成功",
        position: "is-top",
        type: "is-success"
      });
      setTimeout(() => {
        this.$store.commit("setUser", {});
        this.$store.commit("setQuotation", {});
        this.$router.push("/login");
      }, 2000);
    }
  }
};
</script>

<style lang="scss" scoped>
.menu-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 5% 5% 0;
}
.menu-header {
  justify-content: space-between;
  display: flex;
  width: 100%;
  margin-bottom: 36px;
}
a {
  padding: 3px;
  border: none;
  outline: none;
  background: none;
  cursor: pointer;
}
</style>
