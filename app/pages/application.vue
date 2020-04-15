<template>
  <div class="application">
    <Dialog ref="dialog" class="recall-dlg">
      <h1>确定撤回吗？</h1>
      <button class="btn-confirm" @click="withdrawApplication">撤回</button>
    </Dialog>
    <TopMenu opt="application" />
    <section class="application-send">
      <p class="section-title">我发出的申请</p>
      <div
        v-for="(applyInfo, gIndex) in applyList"
        :key="gIndex"
        class="form-wrapper buy-wrapper info-item"
      >
        <div class="lr-block">
          <b-field label="昵称">
            <span class="input verified-show default-color control">{{applyInfo.nickName}}</span>
          </b-field>
          <b-field label="报价">
            <div class="control is-clearfix">
              <div class="input-icon">
                <ICON type="money" />
              </div>
              <input disabled class="input" v-model="applyInfo.price" />
            </div>
          </b-field>
        </div>
        <div
          class="lr-block"
          v-if="applyInfo.status === 'ACCEPT' && applyInfo.quotationType === 'PASS_CODE'"
        >
          <b-field label="岛屿开放类型">
            <span
              class="input verified-show default-color control"
            >{{applyInfo.quotationType | openTypeTranslate}}</span>
          </b-field>
          <b-field label="岛屿密码">
            <div class="control is-clearfix">
              <span class="input verified-show default-color control">{{applyInfo.passCode}}</span>
            </div>
          </b-field>
        </div>
        <b-field v-else label="岛屿开放类型">
          <span
            class="input verified-show default-color control"
          >{{applyInfo.quotationType | openTypeTranslate}}</span>
        </b-field>
        <b-field
          v-if="applyInfo.status === 'ACCEPT' && applyInfo.quotationType === 'FRIENDS'"
          label="Switch 好友编号"
        >
          <div class="friendCode-wrap">
            <span class="input verified-show default-color control">{{applyInfo.switchFriendCode}}</span>
          </div>
        </b-field>
        <div :class="['opera-btn-wrap', {'opera-btn': applyInfo.status !== 'PENDING'}]">
          <b-button
            class="btn-req"
            disabled
            v-if="applyInfo.status === 'ACCEPT'"
            type="is-primary"
          >{{applyInfo.status | applyBtnTextTranslate}}</b-button>
          <b-button
            class="btn-req btn-refused"
            disabled
            v-if="applyInfo.status === 'REJECT'"
            type="is-primary"
          >{{applyInfo.status | applyBtnTextTranslate}}</b-button>
          <b-button
            class="btn-req btn-refused"
            @click="showDlg"
            v-if="applyInfo.status === 'PENDING'"
            type="is-primary"
          >撤回</b-button>
          <b-button
            class="btn-req btn-applyed"
            disabled
            v-if="applyInfo.status === 'PENDING'"
            type="is-primary"
          >{{applyInfo.status | applyBtnTextTranslate}}</b-button>
        </div>
      </div>
      <div
        class="application-form"
        v-if="applyList.length === 0"
      >{{isLoading ? loadingText : '还没有发出申请'}}</div>
    </section>
    <section>
      <p class="section-title">我收到的申请</p>
      <div
        v-for="(reviewInfo, gIndex) in reviewList"
        :key="gIndex"
        class="form-wrapper buy-wrapper info-item"
      >
        <b-field label="昵称">
          <span class="input verified-show default-color control">{{reviewInfo.nickName}}</span>
        </b-field>
        <b-field label="Switch 好友昵称" v-show="reviewInfo.switchNickname.length > 0">
          <span class="input verified-show default-color control">{{reviewInfo.switchNickname}}</span>
        </b-field>
        <b-field
          v-if="gIndex === showPassIndex && reviewInfo.quotationType === 'FRIENDS'"
          label="Switch 好友编号"
        >
          <div class="friendCode-wrap">
            <b-input
              class="friendCode"
              @input="friendCodeInput"
              maxlength="19"
              v-model="reviewInfo.switchFriendCode"
              placeholder="XXXX-XXXX-XXXX"
            ></b-input>
            <span
              :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': reviewInfo.switchFriendCode.length === 0}]"
            >SW-</span>
          </div>
        </b-field>
        <b-field label="即刻ID" v-show="reviewInfo.jikeId.length > 0">
          <span class="input verified-show default-color control">{{reviewInfo.jikeId}}</span>
        </b-field>
        <b-field
          label="岛屿密码"
          v-if="gIndex === showPassIndex && reviewInfo.quotationType === 'PASS_CODE'"
        >
          <div class="control is-clearfix">
            <input class="input" placeholder="请输入岛屿密码" v-model="reviewInfo.passCode" />
          </div>
        </b-field>
        <div class="opera-btn-wrap" v-if="reviewInfo.status === 'PENDING'">
          <b-button
            class="btn-req btn-refused"
            v-if="gIndex === showPassIndex"
            @click="showPassIndex = -1"
            type="is-primary"
          >取消</b-button>
          <b-button
            class="btn-req btn-refused"
            v-else
            @click="operaMyApplication(reviewInfo, 'REJECT')"
            type="is-primary"
          >拒绝</b-button>
          <b-button
            class="btn-req"
            v-if="gIndex === showPassIndex"
            @click="updateMyApplication(reviewInfo, 'ACCEPT', gIndex)"
            type="is-primary"
          >确认</b-button>
          <b-button
            class="btn-req"
            v-else
            @click="operaMyApplication(reviewInfo, 'ACCEPT', gIndex)"
            type="is-primary"
          >同意</b-button>
        </div>
        <div v-if="reviewInfo.status === 'REJECT'" class="opera-btn">
          <b-button class="btn-req btn-refused" disabled type="is-primary">已拒绝</b-button>
        </div>
        <div v-if="reviewInfo.status === 'ACCEPT'" class="opera-btn">
          <b-button class="btn-req" disabled type="is-primary">已同意</b-button>
        </div>
      </div>
      <div
        class="application-form"
        v-if="reviewList.length === 0"
      >{{isLoading ? loadingText : '没有收到申请'}}</div>
    </section>
  </div>
</template>
<script>
import TopMenu from "../components/TopMenu";
import jsCookie from "js-cookie";
import ICON from "../components/ICON";
import Dialog from "../components/Dialog";

export default {
  name: "APPLICATION",
  middleware: "curDay",
  components: { TopMenu, ICON, Dialog },
  data() {
    return {
      isLoading: true,
      loadTextClock: 0,
      loadTextTimer: null,
      loadingText: "载入中",
      reviewList: [],
      applyList: [],
      showPassIndex: -1
    };
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
    }
  },
  filters: {
    openTypeTranslate(val) {
      return val === "PASS_CODE" ? "密码" : "仅好友";
    },
    verifiedTranslate(apply) {
      if (apply.validCount === apply.invalidCount) {
        return "待验证";
      }
      return apply.validCount > apply.invalidCount ? "有效" : "无效";
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
  computed: {
    isUserInfo() {
      return !!this.$store.state.user.username;
    },
    isAuth() {
      return !!jsCookie.get("auth");
    },
    passCode() {
      return this.$store.state.quotation.passCode || "";
    }
  },
  mounted() {
    this.checkAuth();
  },
  beforeDestroy() {
    clearInterval(this.loadTextTimer);
    this.loadTextClock = 0;
  },
  methods: {
    showDlg() {
      this.$refs.dialog.show();
    },
    async checkAuth() {
      if (this.isAuth) {
        if (!this.isUserInfo) {
          this.$store.commit("setLoading");
          let meRes = await this.$axios.$get("/me");
          this.$store.commit("setUser", meRes);
        }
        await this.qryMyQuotation();
        await this.qryMyApplication();
      } else {
        this.$router.push("/login");
      }
    },
    qryMyApplication() {
      this.showPassIndex = -1;
      this.isLoading = true;
      this.$store.commit("setLoading");
      Promise.all([this.qryMyApply(), this.qryMyReview()]).then(
        applications => {
          this.isLoading = false;
          this.$store.commit("closeLoading");
          const cNum = applications[0].length + applications[1].length;
          this.applyList = applications[0].map(apply => {
            return {
              id: apply.id,
              price: apply.price || 0,
              status: apply.status,
              quotationType: apply.quotationType,
              applyId: apply.reviewerId,
              passCode: apply.passCode,
              nickName: apply.reviewerNickname,
              switchNickname: apply.applicant.switchNickname || "",
              switchFriendCode: apply.applicant.switchFriendCode
            };
          });
          const hasPending = applications[1].find(
            el => el.status === "PENDING"
          );
          this.$store.commit("setHasApplicationNew", hasPending);
          this.reviewList = applications[1].map(review => {
            return {
              id: review.id,
              price: review.price || 0,
              status: review.status,
              quotationType: review.quotationType,
              applyId: review.applicant.id,
              passCode: this.passCode,
              nickName: review.applicant.nickname,
              jikeId: review.applicant.jikeId || "",
              switchNickname: review.applicant.switchNickname || "",
              switchFriendCode: review.applicant.switchFriendCode
            };
          });
        }
      );
    },
    async qryMyApply() {
      const myApply = await this.$axios.$get("/applications?type=APPLY");
      return myApply;
    },
    async qryMyReview() {
      const myReview = await this.$axios.$get("/applications?type=REVIEW");
      return myReview;
    },
    async operaMyApplication(application, targetType, gIndex) {
      if (targetType === "REJECT") {
        await this.updateMyApplication(application, targetType);
      } else if (targetType === "ACCEPT") {
        this.showPassIndex = gIndex;
      }
    },
    async updateMyApplication(application, targetType) {
      this.$store.commit("setLoading");
      const putParams = {
        status: targetType
      };
      if (application.quotationType === "PASS_CODE") {
        putParams.passCode = application.passCode;
      } else {
        putParams.switchFriendCode = application.switchFriendCode;
      }
      const myAppli = await this.$axios.$put(
        `/applications/${application.id}`,
        putParams
      );
      this.$buefy.toast.open({
        duration: 2000,
        message: "处理成功",
        position: "is-top",
        type: "is-success"
      });
      await this.qryMyApplication();
      this.$store.commit("closeLoading");
    },
    /**
     * 查询我的发布信息
     */
    async qryMyQuotation(force) {
      const switchFriendCode = this.$store.state.user.switchFriendCode || "";
      let myQuo = await this.$axios.$get(`/my-quotations`);
      if (myQuo && myQuo.length) {
        this.$store.commit("setQuotation", {
          price: myQuo[0].price,
          openType: myQuo[0].openType,
          passCode: myQuo[0].passCode,
          switchFriendCode: switchFriendCode,
          id: myQuo[0].id,
          validCount: myQuo[0].validCount,
          invalidCount: myQuo[0].invalidCount
        });
      }
    },
    /**
     * 撤回申请
     */
    async withdrawApplication(applicationId) {
      this.$refs.dialog.hide();
      this.$store.commit("setLoading");
      await this.$axios.$delete(`/applications/${applicationId}`);
      await this.qryMyApplication();
      this.$buefy.toast.open({
        duration: 3000,
        message: "撤回成功",
        position: "is-top",
        type: "is-success"
      });
    },

    /**
     * 控制 switch 好友编号输入
     */
    friendCodeInput(val) {
      this.switchFriendCode = val
        .replace(/\s/g, "")
        .replace(/(\d{4})(?=\d)/g, "$1-");
    }
  }
};
</script>
<style lang="scss" scoped>
.dialog-container {
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
.application {
  section {
    position: relative;
    margin-top: 1rem;
    padding-top: 2rem;
    .input-icon + input {
      padding-left: 32px;
    }
    .input-icon + input[disabled] {
      background-color: #fffcf5;
    }
  }
  .opera-btn {
    button {
      width: 100%;
      height: 50px;
      border-radius: 11px;
      border: none;
      opacity: 1;
      color: #fff;
      -webkit-opacity: 1;
      -webkit-text-fill-color: #fff;
    }
    button[disabled] {
      opacity: 0.5;
      -webkit-opacity: 0.5;
    }
  }
}
.application-form {
  border-radius: 27px;
  background: #fcf9f2;
  margin: 5%;
  box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
  padding: 5%;
  color: #a0a0a0;
  font-weight: bold;
  font-size: 20px;
  text-align: center;
  line-height: 3;
}
.section-title {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  top: 0;
  display: inline-block;
  padding: 0.2rem 2rem;
  font-weight: bold;
  background-color: #f88d65;
  border-radius: 44px;
  color: #f7f0e2;
  font-size: 16px;
  text-align: center;
}
.info-item .input-icon {
  position: absolute;
  left: 5px;
  z-index: 10;
  margin-top: 7px;
}
.opera-btn-wrap {
  display: flex;
  justify-content: space-between;
  .btn-req,
  .btn-applyed,
  .btn-refused,
  button {
    width: 48%;
  }
  .btn-req {
    background: #7bd9c2;
    height: 50px;
    border-radius: 11px;
  }
  .btn-applyed,
  .btn-applyed:hover {
    background: #937bd9 !important;
  }
  .btn-refused,
  .btn-refused:hover {
    background: #d97b92 !important;
  }
}
</style>
