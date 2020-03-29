<template>
  <div class="menu-container">
    <div class="menu-header">
      <ICON type="logo" />
      <ICON type="github" />
    </div>
    <div class="opts">
      <a href="/buy">
        <ICON v-bind:type="opt === 'buy' ? 'buyOn' : 'buy'" />
      </a>
      <a href="/sell">
        <ICON v-bind:type="opt === 'sell' ? 'sellOn' : 'sell'" />
      </a>
      <a href="/login">
        <ICON v-bind:type="opt === 'login' ? 'loginOn' : 'login'" />
      </a>
      <a href="/register">
        <ICON v-bind:type="opt === 'reg' ? 'regOn' : 'reg'" />
      </a>
    </div>
  </div>
</template>

<script>
import ICON from '../components/ICON'

export default {
  components: {ICON},
  props: {opt: {type: String, required: true}},
}
</script>

<style lang="scss" scoped>
.menu-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 5% 5% 0;
}
.menu-header {
  justify-content: space-between;
  display: flex;
  width: 100%;
  margin-bottom: 36px;
}
a {
  padding: 3px;
  border: none;
  outline: none;
  background: none;
  cursor: pointer;
}
</style>
