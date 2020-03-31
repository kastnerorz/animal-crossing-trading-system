<template>
  <div class="buy-page">
    <TopMenu opt="buy" />
    <div class="item-container">
      <section>
        <p class="section-title">我发布的报价</p>
        <div class="form-wrapper sell-wrapper">
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
              <span class="input verified-show">待验证</span>
            </b-field>
          </div>
          <b-field label="岛屿开放类型">
            <span class="input verified-show">密码</span>
          </b-field>
        </div>
      </section>
      <section>
        <p class="section-title">买入大厅</p>
        <p class="form-wrapper-null" v-if="goodsList.length === 0">
          啥也没有
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
          <b-button class="btn-req" type="is-primary" @click="requestTrading">申请</b-button>
          <div class="two-btn-block">
            <b-button class="btn-devalue" type="is-primary" @click="requestTrading">假的</b-button>
            <b-button class="btn-value" type="is-primary" @click="requestTrading">有效</b-button>
          </div>
          <!-- <b-field label="好友编号">
            <b-input disabled class="friendCode" v-model="good.frNum"></b-input>
          </b-field> -->
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import TopMenu from "../components/TopMenu";
import ICON from "../components/ICON";
export default {
  components: { TopMenu, ICON },
  data() {
    return {
      myQuote: {
        price: 99,
        verified: false,
        playerCount: 2
      },
      goodsList: []
    };
  },
  filters: {
    openTypeTranslate(val) {
      return val === 'PASS_CODE' ? '密码' : '仅好友'
    },
    verifiedTranslate(good) {
      if (good.validCount === good.invalidCount) {
        return "待验证"
      }
      return good.validCount > good.invalidCount ? "有效" : "无效";
    }
  },
  mounted() {
    this.qryQuotations();
  },
  methods: {
    async qryQuotations() {
      let res = await this.$axios.$get("/quotations?type=SELL");
      console.log("success", res);
      this.goodsList = res.map(quo => {
        return {
          id: quo.id,
          price: quo.price,
          nickName: quo.author.nickname,
          sellerId: quo.author.id,
          validCount: quo.validCount,
          invalidCount: quo.invalidCount,
          openType: quo.openType,
          lastModified: quo.lastModified
        };
      });
    },
    requestTrading() {
      console.log("requestTrading");
    }
  }
};
</script>