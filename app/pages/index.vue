<template>
  <div class="buy-page">
    <TopMenu opt="buy" />
    <div class="item-container">
      <section>
        <p class="section-title">我发布的报价</p>
        <Sell v-show="!hasQuote"></Sell>
        <div v-show="hasQuote" class="form-wrapper sell-wrapper">
          <div class="lr-block">
            <b-field label="卖出价">
              <div class="control is-clearfix">
                <div class="input-icon">
                  <ICON type="money" />
                </div>
                <input disabled class="input" v-model="myQuote.price" />
              </div>
            </b-field>
            <b-field label="有效性">
              <span class="input verified-show"
                :class="{ 'valid-color' : myQuote.validCount > myQuote.invalidCount }">{{myQuote | verifiedTranslate}}</span>
            </b-field>
          </div>
          <b-field label="岛屿开放类型">
            <span class="input verified-show default-color">{{myQuote.openType | openTypeTranslate}}</span>
          </b-field>
        </div>
      </section>
      <section>
        <p class="section-title">买入大厅</p>
        <p class="form-wrapper-null" v-if="goodsList.length === 0">
          {{isLoading ? "载入中..." : "啥也没有" }}
        </p>
        <div v-for="(good, gIndex) in goodsList" :key="gIndex" class="form-wrapper buy-wrapper info-item">
          <div class="lr-block">
            <b-field label="昵称">
              <b-input disabled v-model="good.nickName"></b-input>
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
            <b-button class="btn-req btn-applyed" v-if="good.status === 'PENDING'" type="is-primary">{{good.status | applyBtnTextTranslate}}</b-button>
            <b-button class="btn-req btn-refused" v-if="good.status === 'REJECT'" type="is-primary">{{good.status | applyBtnTextTranslate}}</b-button>
            <b-button class="btn-req" v-if="isLogin && good.status === 'NORMAL'" type="is-primary" @click="requestApplications(good.id)">{{good.status | applyBtnTextTranslate}}</b-button>
            <b-button class="btn-req" v-if="!isLogin" type="is-primary" @click="loginAni">登录后申请</b-button>
          </div>
          <p v-else class="my-realse">*我发布的</p>
          <!-- <div class="two-btn-block">
            <b-button class="btn-devalue" type="is-primary" @click="requestApplications">假的</b-button>
            <b-button class="btn-value" type="is-primary" @click="requestApplications">有效</b-button>
          </div> -->
          <!-- <b-field label="好友编号">
            <b-input disabled class="friendCode" v-model="good.frNum"></b-input>
          </b-field> -->
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import jsCookie from "js-cookie";
import Sell from "../components/Sell";
import TopMenu from "../components/TopMenu";
import ICON from "../components/ICON";
let loadingComponent = null;
export default {
  components: { TopMenu, ICON, Sell },
  data() {
    return {
      isLoading: true,
      myQuote: {
        price: "",
        verified: false,
        playerCount: 2
      },
      goodsList: [],
      applicationList: [],
    };
  },
  computed: {
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
      let btnText = ''
      switch (status) {
        case 'NORMAL':
          btnText = '申请'
          break;
        case 'PENDING':
          btnText = '已申请'
          break;
        default:
          btnText = '申请'
          break;
      }
      return btnText
    }
  },
  watch: {
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
      loadingComponent = this.$buefy.loading.open();
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
      let res = await this.$axios.$get("/quotations?type=SELL");
      loadingComponent.close();
      let goodsList = res.map(quo => {
        const isMine = this.isLogin ? quo.author.id === this.$store.state.user.id : false
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
        const cGood = this.applicationList.find(apl => goods.id === apl.quotationId)
        if (cGood) {
          goods.status = cGood.status
        } else {
          goods.status = 'NORMAL'
        }
      })
      this.goodsList = goodsList
    },
    async requestApplications(qId) {
      loadingComponent = this.$buefy.loading.open();
      const reqData = {
        QuotationId: qId
      };
      let trade = await this.$axios.$post("/applications", reqData);
      await this.getMyApplications()
      loadingComponent.close();
    },
    async getMyApplications() {
      let trade = await this.$axios.$get("/applications"); // ?type=APPLY
      this.applicationList = trade.map(tra => {
        return {
          id: tra.id,
          status: tra.status,
          quotationId: tra.quotationId,
        }
      })
    },
    genereateBtnTxt(qId) {
      if (this.applicationList.length === 0) {
        return '申请'
      }
      const findQuo = this.applicationList.find(el => el.quotationId === qId)
      if (findQuo) {
        if (findQuo.status === 'ACCEPT'){
          return '已同意'
        } else if (findQuo.status === 'PENDING'){
          return '已申请'
        }
        return '已拒绝'
      } else {
        return '申请'
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