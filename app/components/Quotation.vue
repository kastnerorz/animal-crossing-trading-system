<template>
  <div class="form-wrapper buy-wrapper info-item">
    <div class="lr-block">
      <b-field label="昵称">
        <span class="input verified-show default-color control">{{quotation.nickName}}</span>
      </b-field>
      <b-field label="收购价">
        <div class="control is-clearfix">
          <div class="input-icon">
            <ICON type="money" />
          </div>
          <input disabled class="input" v-model="quotation.price" />
        </div>
      </b-field>
    </div>
    <div class="lr-block">
      <b-field label="岛屿开放类型">
        <span class="input verified-show default-color control">{{quotation.openType | openTypeTranslate}}</span>
      </b-field>
      <b-field label="发布时间">
        <div class="control is-clearfix">
          <span class="input verified-show default-color control">{{modifieTime}}</span>
        </div>
      </b-field>
    </div>
    <b-field label="手续费">
      <div class="control is-clearfix">
        <input disabled class="input" placeholder="随意" v-model="quotation.handlingFee" />
      </div>
    </b-field>
    <div v-if="!quotation.isMine">
      <b-button class="btn-req" v-if="!isLogin" type="is-primary" @click="loginAni">登录后申请</b-button>
      <template v-else>
        <b-button class="btn-req" v-if="quotation.status === 'NORMAL'" type="is-primary"
          @click="requestApplications(quotation.id, gIndex)">{{quotation.status | applyBtnTextTranslate}}</b-button>
        <b-button class="btn-req btn-applyed" disabled v-if="quotation.status === 'PENDING'" type="is-primary">
          {{quotation.status | applyBtnTextTranslate}}</b-button>
        <b-button class="btn-req btn-refused" disabled v-if="quotation.status === 'REJECT'" type="is-primary">
          {{quotation.status | applyBtnTextTranslate}}</b-button>
        <b-button class="btn-req" disabled v-if="quotation.status === 'ACCEPT'" type="is-primary">
          {{quotation.status | applyBtnTextTranslate}}</b-button>
      </template>
    </div>
    <p v-else class="my-realse">*我发布的</p>
  </div>
</template>
<script>
import ICON from "./ICON";
import jsCookie from "js-cookie";
export default {
  name: "Quotation",
  components: { ICON },
  props: {
    quotationsList: {
      type: Array,
      default() {
        return [];
      }
    },
    quotation: {
      type: Object,
      default() {
        return {};
      }
    },
    gIndex: {
      type: Number,
      default: 0
    },
  },
  computed: {
    isLogin() {
      return !!this.$store.state.user.username;
    },
    isAuth() {
      return !!jsCookie.get("auth");
    }
  },
  data() {
    return {
      timer: null,
      modifieTime: ""
    };
  },
  filters: {
    openTypeTranslate(val) {
      return val === "PASS_CODE" ? "密码" : "仅好友";
    },
    applyBtnTextTranslate(status) {
      let btnText = "";
      switch (status) {
        case "NORMAL":
          btnText = "申请";
          break;
        case "PENDING":
          btnText = "等待同意";
          break;
        case "REJECT":
          btnText = "被拒绝";
          break;
        case "ACCEPT":
          btnText = "已同意";
          break;
        default:
          btnText = "申请";
          break;
      }
      return btnText;
    }
  },
  mounted() {
    this.calcDateDiff();
    this.timer = setInterval(() => {
      this.calcDateDiff();
    }, 80000);
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },
  methods: {
    /**
     * 生成按钮文字
     * @param String qId 报价id
     */
    genereateBtnTxt(qId) {
      if (this.applicationList.length === 0) {
        return "申请";
      }
      const findQuo = this.applicationList.find(el => el.quotationId === qId);
      if (findQuo) {
        if (findQuo.status === "ACCEPT") {
          return "已同意";
        } else if (findQuo.status === "PENDING") {
          return "已申请";
        }
        return "已拒绝";
      } else {
        return "申请";
      }
    },
    /**
     * 转去登录
     */
    loginAni() {
      this.$router.push("/login");
    },
    /**
     * 提交申请
     * @param String qId 报价id
     * @param Number gIndex 报价序号
     */
    async requestApplications(qId, gIndex) {
      const reqData = {
        QuotationId: qId
      };
      this.$store.commit("setLoading");
      let trade = await this.$axios.$post("/applications", reqData);
      this.$set(this.quotationsList[gIndex], "status", "PENDING");
      this.$store.commit("closeLoading");
      this.$buefy.toast.open({
        duration: 2000,
        message: "申请成功!",
        position: "is-top",
        type: "is-success"
      });
    },
    /**
     * 计算发布时间
     */
    calcDateDiff() {
      let time = new Date(this.quotation.modifieTime);
      let minute = 1000 * 60;
      let hour = minute * 60;
      let day = hour * 24;
      let now = new Date().getTime();
      let diffValue = now - time;
      if (diffValue < 0) {
        this.modifieTime = "";
      }
      let dayC = diffValue / day;
      let hourC = diffValue / hour;
      let minC = diffValue / minute;
      if (dayC >= 1) {
        this.modifieTime = "" + parseInt(dayC) + "天前";
      } else if (hourC >= 1) {
        this.modifieTime = "" + parseInt(hourC) + "小时前";
      } else if (minC >= 1) {
        this.modifieTime = "" + parseInt(minC) + "分钟前";
      } else {
        this.modifieTime = "刚刚";
      }
    }
  }
};
</script>