<template>
  <div class="buy-page">
    <TopMenu opt="buy" />
    <div class="item-container">
      <!-- ///////////////////// -->
      <!-- 会写Vue的把下面这个item用v-for 动态生成一下 -->
      <!-- //////////////////// -->
      <section class="form-wrapper info-item">
        <div class="lr-block">
          <b-field label="昵称">
            <b-input v-model="nickName"></b-input>
          </b-field>
          <b-field label="收购价">
            <div class="control is-clearfix">
              <div class="input-icon"><ICON type="money" /></div>
              <input class="input" v-model="price" />
            </div>
          </b-field>
        </div>
        <div class="lr-block">
          <b-field label="岛上人数">
            <b-input v-model="playerCount"></b-input>
          </b-field>
          <b-field label="已认证">
            <b-input v-model="verified"></b-input>
          </b-field>
        </div>
        <b-field label="好友编号">
          <b-input v-model="frNum"></b-input>
        </b-field>
      </section>
      <section class="form-wrapper info-item">
        <div class="lr-block">
          <b-field label="昵称">
            <b-input v-model="nickName"></b-input>
          </b-field>
          <b-field label="收购价">
            <div class="control is-clearfix">
              <div class="input-icon"><ICON type="money" /></div>
              <input class="input" v-model="price" />
            </div>
          </b-field>
        </div>
        <div class="lr-block">
          <b-field label="岛上人数">
            <b-input v-model="playerCount"></b-input>
          </b-field>
          <b-field label="已认证">
            <b-input v-model="verified"></b-input>
          </b-field>
        </div>
        <b-field label="好友编号">
          <b-input v-model="frNum"></b-input>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
import TopMenu from '../components/TopMenu'
import ICON from '../components/ICON'
export default {
  components: {TopMenu, ICON},
}
</script>
