<template>
  <div class="sell-page">
    <TopMenu opt="sell" />
    <section class="form-wrapper">
      <div class="lr-block">
        <b-field label="卖出价">
          <div class="control is-clearfix">
            <div class="input-icon">
              <ICON type="money" />
            </div>
            <input class="input" placeholder="卖出价" v-model="price" />
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
      <b-field v-else label="Switch 好友编号">
        <div class="friendCode-wrap">
          <b-input class="friendCode" @input="friendCodeInput" maxlength="19" v-model="switchFriendCode"
            placeholder="XXXX-XXXX-XXXX-XXXX"></b-input>
          <span
            :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': switchFriendCode.length === 0}]">SW-</span>
        </div>
      </b-field>
      <!-- <b-field label="好友编号" v-if="openType === '2'">
        <div class="control is-clearfix">
          <input class="input friendCode" placeholder="请输入好友编号" disabled v-model="switchFriendCode" />
        </div>
      </b-field> -->
      <b-button class="btn-reg" type="is-primary" @click="validateAllData">发布</b-button>
    </section>
  </div>
</template>
i
<script>
import TopMenu from "../components/TopMenu";
import asyncValidator from "async-validator";
import ICON from "../components/ICON";
const validateRules = {
  price: [
    { required: true, message: "收购价不能为空" },
    { type: "number", message: "收购价必须为数字" }
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
  components: { TopMenu, ICON },
  data() {
    return {
      openTypes: [
        {
          type: "PASS_CODE",
          text: "密码"
        },
        {
          type: "2",
          text: "仅好友"
        }
      ],
      openType: null,
      price: "",
      verified: false,
      passCode: "",
      switchFriendCode: "1111-1111-1111"
    };
  },
  methods: {
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
    validateAllData() {
      const rules = {
        price: validateRules.price,
        openType: validateRules.openType,
        passCode: validateRules.passCode
      };
      let validator = new asyncValidator(rules);
      validator
        .validate({
          price: this.price,
          openType: this.openType,
          passCode: this.passCode
        })
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
     * 发布
     */
    postQuotations() {
      this.$axios
        .$post("/quotations", {
          type: "SELL",
          handlingFee: 10000,
          price: this.price,
          openType: this.openType,
          passCode: this.passCode
        })
        .then(res => {
          console.log("success", res);
          this.$buefy.toast.open({
            duration: 3000,
            message: `发布成功`,
            position: "is-top",
            type: "is-success"
          });
        });
    }
  }
};
</script>
