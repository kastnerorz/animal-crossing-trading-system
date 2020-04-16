<template>
  <section class="form-wrapper">
    <Dialog ref="dialog" class="recall-dlg">
      <h1>确定撤回吗？</h1>
      <button class="btn-confirm" @click="withdrawTrade">撤回</button>
    </Dialog>
    <div class="lr-block">
      <b-field label="报价">
        <div class="control is-clearfix">
          <div class="input-icon">
            <ICON type="money" />
          </div>
          <input class="input" placeholder="卖出价" v-model.number="price" />
        </div>
      </b-field>
      <b-field label="岛屿开放类型">
        <b-select v-model="openType" class="choose-open-type" placeholder="选择类型">
          <option
            v-for="option in openTypes"
            :value="option.type"
            :key="option.type"
          >{{ option.text }}</option>
        </b-select>
      </b-field>
    </div>
    <b-field label="手续费(选填)">
      <div class="control is-clearfix">
        <input class="input" placeholder="输入手续费,铃钱或物品" v-model="handlingFee" />
      </div>
    </b-field>
    <b-field v-if="openType === 'PASS_CODE'" label="岛屿密码">
      <div class="control is-clearfix">
        <input class="input" placeholder="岛屿密码" v-model.trim="passCode" />
      </div>
    </b-field>
    <b-field v-if="openType === 'FRIENDS'" label="Switch 好友编号">
      <div class="friendCode-wrap">
        <b-input
          class="friendCode"
          @input="friendCodeInput"
          maxlength="14"
          v-model="switchFriendCode"
          placeholder="XXXX-XXXX-XXXX"
        ></b-input>
        <span
          :class="[
            'friendCode-wrap-title',
            {'friendCode-wrap-title-gray': switchFriendCode.length === 0},
          ]"
        >SW-</span>
      </div>
    </b-field>
    <div class="opera-btn-wrap" v-if="isAuth && hasQuotation">
      <b-button class="btn-reg btn-refused" type="is-primary" @click="this.$refs.dialog.show">撤回</b-button>
      <b-button class="btn-reg" type="is-primary" @click="validateAllData">修改</b-button>
    </div>
    <b-button
      v-if="isAuth && !hasQuotation"
      class="btn-reg"
      type="is-primary"
      @click="validateAllData"
    >发布</b-button>
    <b-button v-if="!isAuth" class="btn-reg" type="is-primary" @click="loginAni">登录后发布</b-button>
  </section>
</template>
<script>
import ICON from "./ICON";
import jsCookie from "js-cookie";
import asyncValidator from "async-validator";
import Dialog from "../components/Dialog";

const validateRules = {
  price: [
    { required: true, message: "收购价不能为空" },
    { type: "number", message: "收购价必须为数字" },
    {
      validator: function() {
        if (Number(arguments[1]) === 0) {
          return new Error("收购价不可为 0");
        }
        return true;
      }
    }
  ],
  openType: [
    {
      required: true,
      message: "请选择岛屿开放类型"
    }
  ],
  passCode: [
    {
      required: true,
      message: "请输入岛屿密码"
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
  components: { ICON, Dialog },
  props: {
    tradeType: {
      type: String,
      default: "sell"
    }
  },
  data() {
    return {
      quoId: "",
      openTypes: [
        {
          type: "PASS_CODE",
          text: "密码"
        },
        {
          type: "FRIENDS",
          text: "仅好友"
        }
      ],
      openType: "PASS_CODE",
      price: "",
      verified: false,
      passCode: "",
      handlingFee: "",
      switchFriendCode: "",
      hasQuotation: false,
      isAlerting: false
    };
  },
  watch: {
    openType(val) {
      if (val === "FRIENDS") {
        let switchFriendCode = "";
        const userInfo = this.$store.state.user || {};
        if (userInfo && userInfo.switchFriendCode) {
          switchFriendCode = userInfo.switchFriendCode.substring(3);
        }
        this.switchFriendCode = switchFriendCode;
      }
    }
  },
  computed: {
    isLogin() {
      return !!this.$store.state.user.username;
    },
    isAuth() {
      return !!jsCookie.get("auth");
    }
  },
  mounted() {
    this.checkAuth();
  },
  methods: {
    async checkAuth() {
      if (this.isAuth) {
        this.$store.commit("setLoading");
        if (!this.isLogin) {
          let meRes = await this.$axios.$get("/me");
          this.$store.commit("setUser", meRes);
          this.$store.commit("closeLoading");
        }
        await this.qryMyQuotation();
      }
    },
    /**
     * 控制 switch 好友编号输入
     */
    friendCodeInput(val) {
      this.switchFriendCode = val
        .replace(/\s/g, "")
        .replace(/(\d{4})(?=\d)/g, "$1-");
    },
    /**
     * 校验所有输入数据
     */
    validateAllData(type) {
      const rules = {
        price: validateRules.price,
        openType: validateRules.openType
      };
      const vaildData = {
        price: this.price,
        openType: this.openType
      };
      if (this.openType === "PASS_CODE") {
        rules["passCode"] = validateRules.passCode;
        vaildData["passCode"] = this.passCode;
      }
      if (this.openType === "FRIENDS") {
        rules["switchFriendCode"] = validateRules.switchFriendCode;
        vaildData["switchFriendCode"] = this.switchFriendCode;
      }
      let validator = new asyncValidator(rules);
      validator
        .validate(vaildData)
        .then(() => {
          this.postQuotations();
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
     * 查询我的发布信息
     */
    async qryMyQuotation() {
      const switchFriendCode = this.$store.state.user.switchFriendCode || "";
      this.$store.commit("setLoading");
      let myQuo = await this.$axios.$get(
        `/my-quotations?type=${this.tradeType}`
      );
      this.$store.commit("closeLoading");
      if (myQuo && myQuo.length) {
        this.hasQuotation = true;
        this.quoId = myQuo[0].id || "";
        this.price = myQuo[0].price || "";
        this.openType = myQuo[0].openType || "";
        this.passCode = myQuo[0].passCode || "";
        this.handlingFee = myQuo[0].handlingFee || "";
        this.switchFriendCode = switchFriendCode
          ? switchFriendCode.substring(3)
          : "";
        this.$store.commit("setQuotation", {
          price: this.price,
          openType: this.openType,
          passCode: this.passCode,
          switchFriendCode: this.switchFriendCode,
          id: myQuo[0].id,
          validCount: myQuo[0].validCount,
          invalidCount: myQuo[0].invalidCount
        });
      } else {
        this.hasQuotation = false;
        this.quoId = "";
        this.price = "";
        this.openType = "PASS_CODE";
        this.passCode = "";
        this.handlingFee = "";
        this.switchFriendCode = "";
        this.$store.commit("setQuotation", {});
      }
    },
    /**
     * 发布
     */
    async postQuotations() {
      this.$store.commit("setLoading");
      const quoParam = {
        type: this.tradeType,
        handlingFee: this.handlingFee || "",
        price: this.price,
        openType: this.openType,
        passCode: this.passCode
      };
      if (this.openType === "PASS_CODE") {
        quoParam["passCode"] = this.passCode;
      } else if (this.openType === "FRIENDS") {
        quoParam["switchFriendCode"] = "SW-" + this.switchFriendCode;
      }
      if (this.hasQuotation) {
        await this.$axios.$put(`/quotations/${this.quoId}`, quoParam);
      } else {
        await this.$axios.$post("/quotations", quoParam);
      }
      this.$buefy.toast.open({
        duration: 3000,
        message: this.hasQuotation ? "修改成功" : "发布成功",
        position: "is-top",
        type: "is-success"
      });
      await this.qryMyQuotation();
      this.$emit("editMyApplication");
    },
    /**
     * 撤回报价
     */
    async withdrawTrade() {
      this.$refs.dialog.hide();
      this.$store.commit("setLoading");
      await this.$axios.$delete(`/quotations/${this.quoId}`);
      await this.qryMyQuotation();
      this.$emit("editMyApplication");
      this.$buefy.toast.open({
        duration: 3000,
        message: "撤回成功",
        position: "is-top",
        type: "is-success"
      });
    },
    /**
     * 转去登录
     */
    loginAni() {
      this.$router.push("/login");
    }
  }
};
</script>
<style lang="scss" scoped>
.recall-dlg {
  font-weight: bold;
  h1 {
    font-size: 1.8rem;
  }
  button {
    height: 50px;
    width: 100%;
    background-color: #d97b92;
    border: none;
    margin-top: 8px;
    font-size: 1.3rem;
    color: white;
  }
}
input[disabled] {
  background-color: #fffcf5;
  color: #67654a;
  box-shadow: none;
  border: none;
}
.control {
  display: flex;
  align-items: center;
}
.input-icon {
  position: absolute;
  left: 5px;
  z-index: 10;
  top: 7px;
  & ~ input {
    padding-left: 32px;
  }
}
.choose-open-type {
  /deep/ .select {
    width: 100%;
  }
  /deep/ select {
    width: 100%;
    color: #4a4a4a;
    &:not(.is-multiple):not(.is-loading)::after {
      top: 54%;
    }
  }
  .select.is-empty select {
    color: #d9d6cb;
  }
}
.opera-btn-wrap {
  display: flex;
  justify-content: space-between;
  button {
    width: 48%;
    height: 50px;
    border-radius: 11px;
  }
  .btn-refused {
    background: #d97b92;
  }
}
</style>
