<template>
  <div class="buy-page sell-page">
    <TopMenu :opt="operaTypeLow" />
    <div class="item-container">
      <section>
        <p class="section-title">我的报价</p>
        <MyTrade
          :tradeType="operaType"
          @editMyApplication="refreshQuotas"
        ></MyTrade>
      </section>
      <section>
        <p class="section-title">{{ operaTypeText }}大厅</p>
        <Quotation
          v-for="(quotation, gIndex) in quotationsList"
          :key="quotation.id"
          :quotationsList="quotationsList"
          :gIndex="gIndex"
          :quotation="quotation"
        ></Quotation>
        <p class="form-wrapper-null" v-if="isLoading">{{ loadingText }}</p>
        <p
          class="form-wrapper-null"
          v-if="!isLoading && quotationsList.length === 0"
        >
          啥也没有
        </p>
      </section>
    </div>
  </div>
</template>
<script>
import jsCookie from "js-cookie";
import MyTrade from "./MyTrade";
import Quotation from "./Quotation";
import TopMenu from "./TopMenu";
import ICON from "./ICON";

export default {
  name: "Opera",
  props: {
    operaType: {
      type: String,
      default: "SELL"
    }
  },
  components: { TopMenu, ICON, MyTrade, Quotation },
  data() {
    return {
      timer: null,
      isLoading: true,
      loadTextClock: 0,
      loadTextTimer: null,
      loadingText: "载入中",
      myQuote: {
        price: "",
        verified: false,
        playerCount: 2
      },
      quotationsList: [],
      applicationList: []
    };
  },
  computed: {
    operaTypeLow() {
      return this.operaType.toLowerCase();
    },
    operaTypeText() {
      return this.operaType === "SELL" ? "买入" : "卖出";
    },
    isLogin() {
      return !!this.$store.state.user.id;
    },
    hasQuote() {
      return !!this.$store.state.quotation.openType;
    }
  },
  filters: {
    verifiedTranslate(quotation) {
      if (quotation.validCount === quotation.invalidCount) {
        return "待验证";
      }
      return quotation.validCount > quotation.invalidCount ? "有效" : "无效";
    }
  },
  beforeDestroy() {
    clearInterval(this.timer);
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
    /**
     * 获取我的申请
     */
    async checkAuth() {
      // this.$store.commit("setLoading");
      const user = this.$store.state.user;
      const auth = jsCookie.get("auth");
      if (!!auth) {
        await this.getMyApplications();
        this.getApplicationCountQuiet();
        this.timer = setInterval(() => {
          this.getApplicationCountQuiet();
        }, 60000);
      }
      if (!user.username && !!auth) {
        let res = await this.$axios.$get("/me");
        this.$store.commit("setUser", res);
      }
      await this.qryQuotations();
    },
    /**
     * 刷新报价
     */
    async refreshQuotas() {
      this.quotationsList = [];
      this.isLoading = true;
      this.$store.commit("setLoading");
      await this.qryQuotations();
      this.$store.commit("closeLoading");
    },
    /**
     * 查询报价
     */
    async qryQuotations() {
      let res = await this.$axios.$get(`/quotations?type=${this.operaType}`);
      this.isLoading = false;
      let quotationsList = res.map(quo => {
        const isMine = this.isLogin
          ? quo.author.id === this.$store.state.user.id
          : false;
        return {
          id: quo.id,
          price: quo.price,
          nickName: quo.author.nickname,
          sellerId: quo.author.id,
          modifieTime: quo.lastModified,
          validCount: quo.validCount,
          invalidCount: quo.invalidCount,
          handlingFee: quo.handlingFee,
          openType: quo.openType,
          isMine: isMine,
          status: "NORMAL",
          lastModified: quo.lastModified
        };
      });
      quotationsList.forEach(goods => {
        const cGood = this.applicationList.find(
          apl => goods.id === apl.quotationId
        );
        if (cGood) {
          goods.status = cGood.status;
        }
      });
      this.quotationsList = quotationsList;
      // this.quotationsList = [...this.quotationsList, ...quotationsList];
    },
    /**
     * 查询我的申请，比对后修改报价状态
     */
    async getMyApplications() {
      let trade = await this.$axios.$get(`/applications?type=APPLY`); // ?type=APPLY
      this.applicationList = trade.map(tra => {
        return {
          id: tra.id,
          status: tra.status,
          quotationId: tra.quotationId
        };
      });
    },
    /**
     * 静默查询申请更新
     */
    async getApplicationCountQuiet() {
      const reveiwList = await this.$axios.$get("/applications?type=REVIEW");
      const hasPending = reveiwList.find(el => el.status === "PENDING");
      this.$store.commit("setHasApplicationNew", hasPending);
    }
  }
};
</script>
