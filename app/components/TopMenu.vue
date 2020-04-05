<template>
  <div class="menu-container">
    <div class="menu-header">
      <ICON type="logo" />
      <div class="logout-github">
        <span v-show="isLogin" @click="signOut">
          <ICON type="logout" />
        </span>
        <a target="_blank" href="https://github.com/kastnerorz/animal-crossing-trading-system">
          <ICON type="github" />
        </a>
      </div>
    </div>
    <div class="opts">
      <n-link v-show="showBuy" to="/">
        <ICON :type="gennerateIcon('buy')" />
      </n-link>
      <n-link v-show="showSell" to="/sell">
        <ICON :type="gennerateIcon('sell')" />
      </n-link>
      <n-link v-show="isLogin && showApplication" to="/application">
        <i v-show="hasNewApply" class="tip-light"></i>
        <ICON :type="gennerateIcon('application')" />
      </n-link>
      <n-link v-show="!isLogin" to="/login">
        <ICON :type="gennerateIcon('login')" />
      </n-link>
      <n-link v-show="!isLogin" to="/register">
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
    return {
      showBuy: false,
      showSell: false,
      showApplication: false
    };
  },
  props: {
    opt: { type: String, required: true }
  },
  computed: {
    isLogin() {
      return !!jsCookie.get("auth");
    },
    hasNewApply() {
      return this.$store.state.hasApplicationNew
    }
  },
  mounted() {
    this.calcDayLinkShow();
  },
  methods: {
    calcDayLinkShow() {
      const cDate = new Date();
      if (process.env.NODE_ENV === "development") {
        cDate.setDate(cDate.getDate() + 2);
      }
      const cDay = cDate.getDay();
      const cHour = cDate.getHours();
      if (cDay !== 0) {
        this.showBuy = false;
        this.showSell = true;
      } else {
        this.showSell = false;
        this.showBuy = cHour < 12;
      }
      if (this.showSell || this.showBuy) {
        this.showApplication = true;
      }
    },
    /**
     * 生成 ICON
     */
    gennerateIcon(type) {
      return this.opt === type ? `${type}On` : type;
    },
    /**
     * 登出
     */
    signOut() {
      if (!this.isLogin) {
        return;
      }
      this.$store.commit("setLoading");
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
  .logout-github {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    span {
      margin-right: 10px;
    }
  }
}
.opts {
  display: flex;
  justify-items: center;
  align-items: center;
}
a {
  padding: 3px;
  border: none;
  outline: none;
  background: none;
  cursor: pointer;
  position: relative;
}
.tip-light {
  display: block;
  width: 7px;
  height: 7px;
  background: #f00;
  position: absolute;
  right: 3px;
  animation: breathe-error 1000ms infinite alternate ease-in-out;
  border-radius: 50%;
}
@keyframes breathe-error {
  0% {
    background: rgba(255, 255, 255, 0);
    box-shadow: 0px 0px 0px red;
  }
  100% {
    background: #f00;
    box-shadow: 0px 0px 5px red;
  }
}
</style>
