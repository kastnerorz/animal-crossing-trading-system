<template>
  <div class="sell-page">
    <TopMenu opt="sell" />
    <!-- <section class="form-wrapper">
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
      <b-field  v-if="openType === '2'" label="Switch 好友编号">
        <div class="friendCode-wrap">
          <b-input class="friendCode" @input="friendCodeInput" maxlength="19" v-model="switchFriendCode"
            placeholder="XXXX-XXXX-XXXX-XXXX"></b-input>
          <span
            :class="['friendCode-wrap-title', {'friendCode-wrap-title-gray': switchFriendCode.length === 0}]">SW-</span>
        </div>
      </b-field>
      <b-button class="btn-reg" type="is-primary" @click="validateAllData">发布</b-button>
    </section> -->
    <Sell></Sell>
  </div>
</template>
<script>
import TopMenu from "../components/TopMenu";
import Sell from "../components/Sell";
export default {
  components: { TopMenu, Sell }
};
</script>
