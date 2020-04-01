<template>
  <section class="form-wrapper">
    <div class="lr-block">
      <b-field label="卖出价">
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
    <b-field label="岛屿密码" v-if="openType === 'PASS_CODE'">
      <div class="control is-clearfix">
        <input class="input" placeholder="请输入岛屿密码" v-model="passCode" />
      </div>
    </b-field>
    <b-field v-if="openType === 'FRIENDS'" label="Switch 好友编号">
      <div class="friendCode-wrap">
        <b-input class="friendCode" @input="friendCodeInput" maxlength="19" v-model="switchFriendCode"
          placeholder="XXXX-XXXX-XXXX-XXXX"></b-input>
        <span
          :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': switchFriendCode.length === 0}]">SW-</span>
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
let loadingComponent = null;
const validateRules = {
  price: [
    { required: true, message: "收购价不能为空" },
    {
      validator: function() {
        if (Number(arguments[1]) === 0) {
          return new Error("收购价不可为 0");
        }
        if (!Number(arguments[1])) {
          return new Error("收购价必须为数字");
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
      pattern: /([0-9]{4})-([0-9]{4})-([0-9]{4})-([0-9]{4})$/,
      message: "Switch好友编号输入错误"
    }
  ]
};
export default {
  components: { ICON },
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
      openType: null,
      price: "",
      verified: false,
      passCode: "",
      switchFriendCode: ""
    };
  },
  computed: {
    isLogin() {
      return !!this.$store.state.user.username;
    },
    isAuth() {
      return !!jsCookie.get("username");
    },
    hasQuotation() {
      return !!this.$store.state.quotation.openType;
    }
  },
  mounted() {
    this.checkAuth();
  },
  methods: {
    async checkAuth() {
      if (this.isAuth) {
        loadingComponent = this.$buefy.loading.open();
        if (!this.isLogin) {
          let meRes = await this.$axios.$get("/me");
          this.$store.commit("setUser", meRes);
          this.qryMyQuotation();
        } else {
          this.qryMyQuotation();
        }
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
    async qryMyQuotation(force) {
      const switchFriendCode = this.$store.state.user.switchFriendCode || "";
      const selfQuotation = this.$store.state.quotation || {};
      if (this.hasQuotation && !force) {
        loadingComponent.close();
        this.quoId = selfQuotation.id || "";
        this.price = selfQuotation.price || "";
        this.openType = selfQuotation.openType || "";
        this.passCode = selfQuotation.passCode || "";
        this.switchFriendCode = switchFriendCode
          ? switchFriendCode.substring(3)
          : "";
        return;
      }
      let myQuo = await this.$axios.$get("/quotations/my");
      loadingComponent.close();
      if (myQuo && myQuo.length) {
        this.quoId = myQuo[0].id || "";
        this.price = myQuo[0].price || "";
        this.openType = myQuo[0].openType || "";
        this.passCode = myQuo[0].passCode || "";
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
      loadingComponent = this.$buefy.loading.open();
      const quoParam = {
        type: "SELL",
        handlingFee: 10000,
        price: this.price,
        openType: this.openType,
        passCode: this.passCode
      };

      if (this.openType === "PASS_CODE") {
        quoParam["passCode"] = this.passCode;
      }
      if (this.openType === "FRIENDS") {
        quoParam["switchFriendCode"] = this.switchFriendCode;
      }
      if (this.hasQuotation) {
        await this.$axios.$put(`/quotations/${this.quoId}`, quoParam);
      } else {
        await this.$axios.$post("/quotations", quoParam);
      }
      await this.qryMyQuotation(true);
      this.$buefy.toast.open({
        duration: 3000,
        message: this.hasQuotation ? "修改成功" : "发布成功",
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
  /deep/ select {
    width: 9rem;
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