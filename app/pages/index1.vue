<template>
  <section>
    <b-collapse class="card" animation="slide" aria-id="contentIdForA11y3">
      <div slot="trigger" slot-scope="props" class="card-header" role="button" aria-controls="contentIdForA11y3">
        <p class="card-header-title">
          我发布的报价
        </p>
        <a class="card-header-icon">
          <b-icon :icon="props.open ? 'menu-down' : 'menu-up'">
          </b-icon>
        </a>
      </div>
      <div class="card-content">
        <div class="content">
          收购价：{{ myQuotation.price }}铃钱
          已认证：{{ myQuotation.verified }}
          <b-field label="岛上人数">
            <b-numberinput v-model="myQuotation.participantCount" @input="updateMyQuotation" min="0" max="8"
              controls-position="compact"></b-numberinput>
          </b-field>
        </div>
      </div>
    </b-collapse>
    <b-table :data="quotations" :columns="columns"></b-table>
  </section>
</template>

<script>
export default {
  name: "HomePage",
  mounted() {
    this.fetchQuotations();
    this.fetchMyQuotation();
  },
  data() {
    return {
      myQuotation: {},
      quotations: [],
      columns: [
        {
          field: "author.nickname",
          label: "昵称"
        },
        {
          field: "author.switchFriendCode",
          label: "好友编号"
        },
        {
          field: "price",
          label: "报价",
          centered: true
        },
        {
          field: "participantCount",
          label: "岛上人数"
        },
        {
          field: "verified",
          label: "已认证"
        }
      ]
    };
  },
  methods: {
    async fetchQuotations() {
      this.quotations = await this.$axios.$get("/quotations?type=SELL");
    },
    async fetchMyQuotation() {
      this.myQuotation = await this.$axios.$get("/quotations/my");
    },
    async updateMyQuotation() {
      this.myQuotation = await this.$axios.$put(
        `/quotations/${this.myQuotation.id}`,
        {
          participantCount: this.myQuotation.participantCount
        }
      );
    }
  }
};
</script>
