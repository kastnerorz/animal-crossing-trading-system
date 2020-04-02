<template>
  <div class="buy-page sell-page">
    <TopMenu :opt="operaTypeLow" />
    <div class="item-container">
      <section>
        <p class="section-title">我的报价</p>
        <MyTrade :tradeType="operaType"></MyTrade>
      </section>
      <section>
        <p class="section-title">{{operaTypeText}}大厅</p>
        <div v-for="(good, gIndex) in goodsList" :key="gIndex" class="form-wrapper buy-wrapper info-item">
          <div class="lr-block">
            <b-field label="昵称">
              <span class="input verified-show default-color control">{{good.nickName}}</span>
            </b-field>
            <b-field label="收购价">
              <div class="control is-clearfix">
                <div class="input-icon">
                  <ICON type="money" />
                </div>
                <input disabled class="input" v-model="good.price" />
              </div>
            </b-field>
          </div>
          <div class="lr-block">
            <b-field label="岛屿开放类型">
              <span class="input verified-show default-color control">{{good.openType | openTypeTranslate}}</span>
            </b-field>
            <b-field label="有效性">
              <span class="input verified-show"
                :class="{ 'valid-color' : good.validCount > good.invalidCount }">{{good | verifiedTranslate}}</span>
            </b-field>
          </div>
          <div v-if="!good.isMine">
            <b-button class="btn-req" v-if="!isLogin" type="is-primary" @click="loginAni">登录后申请</b-button>
            <template v-else>
              <b-button class="btn-req" v-if="good.status === 'NORMAL'" type="is-primary"
                @click="requestApplications(good.id, gIndex)">{{good.status | applyBtnTextTranslate}}</b-button>
              <b-button class="btn-req btn-applyed" v-if="good.status === 'PENDING'" type="is-primary">
                {{good.status | applyBtnTextTranslate}}</b-button>
              <b-button class="btn-req btn-refused" v-if="good.status === 'REJECT'" type="is-primary">
                {{good.status | applyBtnTextTranslate}}</b-button>
            </template>
          </div>
          <p v-else class="my-realse">*我发布的</p>
          <!-- <div class="two-btn-block">
            <b-button class="btn-devalue" type="is-primary" @click="requestApplications">假的</b-button>
            <b-button class="btn-value" type="is-primary" @click="requestApplications">有效</b-button>
          </div> -->
        </div>
        <p class="form-wrapper-null" v-if="isLoading">{{loadingText}}</p>
        <p class="form-wrapper-null" v-if="!isLoading && goodsList.length === 0">啥也没有</p>
      </section>
    </div>
  </div>
</template>
<script>
import jsCookie from "js-cookie";
import MyTrade from "./MyTrade";
import TopMenu from "./TopMenu";
import ICON from "./ICON";
export default {
  name: "Opera",
  props: {
    operaType: {
      type: String,
      default: 'SELL'
    }
  },
  components: { TopMenu, ICON, MyTrade },
  data() {
    return {
      isLoading: true,
      loadTextClock: 0,
      loadTextTimer: null,
      loadingText: "载入中",
      myQuote: {
        price: "",
        verified: false,
        playerCount: 2
      },
      goodsList: [],
      applicationList: []
    };
  },
  computed: {
    operaTypeLow() {
      return this.operaType.toLowerCase();
    },
    operaTypeText() {
      return this.operaType === 'SELL' ? '买入' : '卖出';
    },
    isLogin() {
      return !!this.$store.state.user.id;
    },
    hasQuote() {
      return !!this.$store.state.quotation.openType;
    }
  },
  filters: {
    openTypeTranslate(val) {
      return val === "PASS_CODE" ? "密码" : "仅好友";
    },
    verifiedTranslate(good) {
      if (good.validCount === good.invalidCount) {
        return "待验证";
      }
      return good.validCount > good.invalidCount ? "有效" : "无效";
    },
    applyBtnTextTranslate(status) {
      let btnText = "";
      switch (status) {
        case "NORMAL":
          btnText = "申请";
          break;
        case "PENDING":
          btnText = "已申请";
          break;
        case "REJECT":
          btnText = "被拒绝";
          break;
        default:
          btnText = "申请";
          break;
      }
      return btnText;
    }
  },
  beforeDestroy() {
    clearInterval(this.loadTextTimer);
    this.loadTextClock = 0;
  },
  watch: {
    isLoading: {
      immediate: true,
      handler(val) {
        if (val) {
          this.loadTextTimer = setInterval(() => {
            this.loadingText =
              "载入中" + new Array(this.loadTextClock).fill(".").join("");
            this.loadTextClock++;
            if (this.loadTextClock > 3) {
              this.loadTextClock = 0;
            }
          }, 300);
        } else {
          clearInterval(this.loadTextTimer);
        }
      }
    },
    hasQuote: {
      immediate: true,
      handler(val) {
        if (val) {
          const quotation = this.$store.state.quotation;
          this.myQuote = {
            price: quotation.price || "",
            openType: quotation.openType || "",
            passCode: quotation.passCode || "",
            switchFriendCode: quotation.switchFriendCode || "",
            validCount: quotation.validCount || 0,
            invalidCount: quotation.invalidCount || 0
          };
        }
      }
    }
  },
  mounted() {
    this.checkAuth();
  },
  methods: {
    async checkAuth() {
      // this.$store.commit("setLoading");
      const user = this.$store.state.user;
      const auth = jsCookie.get("auth");
      if (!!auth) {
        await this.getMyApplications();
      }
      if (!user.username && !!auth) {
        let res = await this.$axios.$get("/me");
        this.$store.commit("setUser", res);
      }
      await this.qryQuotations();
    },
    async qryQuotations() {
      let res = await this.$axios.$get(`/quotations?type=${this.operaType}`);
      this.isLoading = false;
      // this.$store.commit("closeLoading");
      let goodsList = res.map(quo => {
        const isMine = this.isLogin
          ? quo.author.id === this.$store.state.user.id
          : false;
        return {
          id: quo.id,
          price: quo.price,
          nickName: quo.author.nickname,
          sellerId: quo.author.id,
          validCount: quo.validCount,
          invalidCount: quo.invalidCount,
          openType: quo.openType,
          isMine: isMine,
          lastModified: quo.lastModified
        };
      });
      goodsList.forEach(goods => {
        const cGood = this.applicationList.find(
          apl => goods.id === apl.quotationId
        );
        if (cGood) {
          goods.status = cGood.status;
        } else {
          goods.status = "NORMAL";
        }
      });
      this.goodsList = goodsList;
    },
    async requestApplications(qId, gIndex) {
      const reqData = {
        QuotationId: qId
      };
      this.$store.commit("setLoading");
      let trade = await this.$axios.$post("/applications", reqData);
      this.$store.commit("closeLoading");
      this.$set(this.goodsList[gIndex], "status", "PENDING");
      this.$buefy.toast.open({
        duration: 2000,
        message: "申请成功",
        position: "is-top",
        type: "is-success"
      });
    },
    async getMyApplications() {
      let trade = await this.$axios.$get("/applications"); // ?type=APPLY
      this.applicationList = trade.map(tra => {
        return {
          id: tra.id,
          status: tra.status,
          quotationId: tra.quotationId
        };
      });
    },
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
    }
  }
};
</script>
