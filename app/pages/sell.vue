<template>
  <div class="sell-page">
    <TopMenu opt="sell" />
    <section class="form-wrapper">
      <div class="lr-block">
        <b-field label="收购价">
          <div class="control is-clearfix">
            <div class="input-icon"><ICON type="money" /></div>
            <input class="input" v-model="price" />
          </div>
        </b-field>
        <b-field label="已认证">
          <b-input v-model="verified"></b-input>
        </b-field>
      </div>
      <b-field label="岛上人数">
        <div class="counter-block">
          <button>-</button>
          <b-input v-model="playerCount"></b-input>
          <button>+</button>
        </div>
      </b-field>
    </section>
  </div>
</template>
i
<script>
import TopMenu from '../components/TopMenu'
import ICON from '../components/ICON'
export default {
  components: {TopMenu, ICON},
  data() {
    return {
      price: 99,
      verified: false,
      playerCount: 2,
    }
  },
}
</script>
