<template>
  <div class="reg-page">
    <TopMenu opt="reg" />
    <section class="form-wrapper">
      <b-field label="昵称">
        <b-input v-model="nickname" @blur="validateMethod('nickname')" placeholder="请输入昵称"></b-input>
      </b-field>
      <b-field label="用户名">
        <b-input v-model="username" @blur="validateMethod('username')" placeholder="请输入用户名"></b-input>
      </b-field>
      <b-field label="密码">
        <b-input v-model="password" autocomplete="off"  @blur="validateMethod('password')" placeholder="请输入密码" type="password"
          maxlength="30"></b-input>
      </b-field>
      <b-field label="再次输入密码">
        <b-input v-model="rePassword" autocomplete="off"  @blur="validateMethod('rePassword')" placeholder="请再次输入密码" type="password"
          maxlength="30">
        </b-input>
      </b-field>
      <b-field label="Switch 好友昵称">
        <b-input v-model="switchNickname" @blur="validateMethod('switchNickname')" placeholder="请输入你的Switch 好友昵称"></b-input>
      </b-field>
      <b-field label="Switch 好友编号">
        <div class="friendCode-wrap">
          <b-input class="friendCode" @blur="validateMethod('switchFriendCode')" @input="friendCodeInput" maxlength="14"
            v-model="switchFriendCode" placeholder="XXXX-XXXX-XXXX"></b-input>
          <span
            :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': switchFriendCode.length === 0}]">SW-</span>
        </div>
      </b-field>
      <b-field label="即刻 ID">
        <b-input v-model="jikeId" placeholder="请输入你的即刻 ID（选填）"></b-input>
      </b-field>
      <b-button class="btn-reg" type="is-primary" @click="validateAllData">注册</b-button>
    </section>
  </div>
</template>

<script>
import TopMenu from "../components/TopMenu";
import asyncValidator from "async-validator";
import { sha256 } from "js-sha256";
let password = "";
const validateRules = {
  username: [
    {
      required: true,
      message: "请输入用户名"
    },
    {
      type: "string",
      min: 4,
      message: "用户名太短了"
    }
  ],
  nickname: [
    {
      required: true,
      message: "请输入昵称"
    },
    {
      type: "string",
      min: 1,
      message: "昵称太短了"
    }
  ],
  switchNickname: [
    {
      required: true,
      message: "请输入switch 好友昵称"
    }
  ],
  password: [
    {
      required: true,
      message: "请输入密码"
    },
    {
      type: "string",
      min: 6,
      message: "密码太短了"
    }
  ],
  rePassword: [
    {
      validator: function() {
        if (password === "") {
          return new Error("请先输入密码");
        }
        if (password !== arguments[1]) {
          return new Error("两次输入密码不一样");
        }
        return true;
      }
    }
  ],
  switchFriendCode: [
    {
      required: true,
      message: "请输入Switch好友编号"
    },
    {
      pattern: /([0-9]{4})-([0-9]{4})-([0-9]{4})$/,
      message: "Switch好友编号输入错误"
    }
  ]
};
export default {
  components: { TopMenu },
  name: "Register",
  mounted() {},
  data() {
    return {
      username: "",
      password: "",
      rePassword: "",
      nickname: "",
      switchFriendCode: "",
      switchNickname: "",
      jikeName: "",
      jikeId: ""
    };
  },
  watch: {
    password: {
      immediate: true,
      handler(val) {
        password = val;
      }
    }
  },
  methods: {
    /**
     * 校验用户输入
     * @param {String} modelName 数据校验名称
     */
    validateMethod(modelName) {
      const rules = {
        [modelName]: validateRules[modelName]
      };
      let validator = new asyncValidator(rules);
      validator
        .validate({
          [modelName]: this[modelName]
        })
        .then(res => {
          validator = null;
        })
        .catch(({ errors, fields }) => {
          const errMsg = errors
            .map(el => {
              return el.message;
            })
            .join("、");
          this.$buefy.toast.open({
            duration: 3000,
            message: errMsg,
            position: "is-top",
            type: "is-danger"
          });
          validator = null;
        });
    },
    /**
     * 校验所有输入数据
     */
    validateAllData() {
      const rules = {
        username: validateRules.username,
        nickname: validateRules.nickname,
        password: validateRules.password,
        rePassword: validateRules.rePassword,
        switchNickname: validateRules.switchNickname,
        switchFriendCode: validateRules.switchFriendCode
      };
      let validator = new asyncValidator(rules);
      validator
        .validate({
          username: this.username,
          nickname: this.nickname,
          rePassword: this.rePassword,
          password: this.password,
          switchNickname: this.switchNickname,
          switchFriendCode: this.switchFriendCode
        })
        .then(() => {
          this.register();
          validator = null;
        })
        .catch(({ errors, fields }) => {
          const errMsg = errors
            .map(el => {
              return el.message;
            })
            .join("、");
          this.$buefy.toast.open({
            duration: 3000,
            message: errMsg,
            position: "is-top",
            type: "is-danger"
          });
          validator = null;
        });
    },
    /**
     * 注册
     */
    async register() {
      this.$store.commit('setLoading');
      var hash = sha256.create();
      hash.update(this.password);
      let userInfo = {
        username: this.username,
        password: hash.hex(),
        nickname: this.nickname,
        jikeId: this.jikeId || "",
        switchNickname: this.switchNickname || "",
        switchFriendCode: "SW-" + this.switchFriendCode
      };
      const user = await this.$axios.$post("/users", userInfo);
      this.$buefy.toast.open({
        duration: 1000,
        message: "注册成功！",
        position: "is-top",
        type: "is-success"
      });
      if (user) {
        setTimeout(() => {
          this.$store.commit('closeLoading');
          this.$router.push("/login");
        }, 1000);
      }
    },
    /**
     * 控制 switch 好友编号输入
     */
    friendCodeInput(val) {
      this.switchFriendCode = val
        .replace(/\s/g, "")
        .replace(/(\d{4})(?=\d)/g, "$1-");
    }
  }
};
</script>
