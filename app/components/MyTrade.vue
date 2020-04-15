<template>
  <section class="form-wrapper">
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
          <option v-for="option in openTypes" :value="option.type" :key="option.type">
            {{ option.text }}
          </option>
        </b-select>
      </b-field>
    </div>
    <div class="lr-block" v-if="openType === 'PASS_CODE'">
      <b-field label="手续费(选填)">
        <div class="control is-clearfix">
          <div class="input-icon">
            <ICON type="money" />
          </div>
          <input class="input" placeholder="手续费" v-model.number="handlingFee" />
        </div>
      </b-field>
      <b-field label="岛屿密码">
        <div class="control is-clearfix">
          <input class="input" placeholder="岛屿密码" v-model="passCode" />
        </div>
      </b-field>
    </div>
    <b-field v-if="openType === 'FRIENDS'" label="Switch 好友编号">
      <div class="friendCode-wrap">
        <b-input class="friendCode" @input="friendCodeInput" maxlength="14" v-model="switchFriendCode"
          placeholder="XXXX-XXXX-XXXX"></b-input>
        <span
          :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': switchFriendCode.length === 0}]">SW-</span>
      </div>
    </b-field>
    <b-field v-if="openType === 'FRIENDS'" label="手续费(选填)">
      <div class="control is-clearfix">
        <input class="input" placeholder="手续费" v-model="handlingFee" />
      </div>
    </b-field>
    <b-button v-if="isAuth" class="btn-reg" type="is-primary" @click="validateAllData">{{hasQuotation ? '修改' : '发布'}}
    </b-button>
    <b-button v-else class="btn-reg" type="is-primary" @click="loginAni">登录后发布</b-button>
  </section>
</template>
<script>
import ICON from "./ICON";
import jsCookie from "js-cookie";
import { mapMutations } from "vuex";
import asyncValidator from "async-validator";
const validateRules = {
  price: [
    { required: true,  message: "收购价不能为空" },
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
  handlingFee: [
    { required: false, type: "number", message: "手续费必须为数字" },
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
  components: { ICON },
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
      hasQuotation: false
    };
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
        handlingFee: validateRules.handlingFee,
        openType: validateRules.openType
      };
      const vaildData = {
        price: this.price,
        handlingFee: this.handlingFee,
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
    async qryMyQuotation(force) {
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
      }
    },
    /**
     * 发布
     */
    async postQuotations() {
      this.$store.commit("setLoading");
      const quoParam = {
        type: this.tradeType,
        handlingFee: this.handlingFee || 0,
        price: this.price,
        openType: this.openType,
        passCode: this.passCode
      };

      if (this.openType === "PASS_CODE") {
        quoParam["passCode"] = this.passCode;
      } else if (this.openType === "FRIENDS") {
        quoParam["switchFriendCode"] = this.switchFriendCode;
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
      await this.qryMyQuotation(true);
      this.$emit("editMyApplication");
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
.btn-reg {
  background: #7bd9c2;
  height: 58px;
  width: 100%;
  margin-top: 8px;
  border-radius: 11px;
}
</style>