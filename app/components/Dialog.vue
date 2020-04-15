<style lang="scss" scoped>
.dialog-container {
  opacity: 1;
  transition: all 0.25s ease-out;
}
.fk-hidden {
  opacity: 0;
  pointer-events: none;
}
</style>

<template>
  <div class="dialog-container" v-bind:class="{'fk-hidden': !this.isVisible}">
    <div class="bg-mask"></div>
    <div class="dialog">
      <i @click="hide">
        <ICON type="x" />
      </i>
      <slot></slot>
    </div>
  </div>
</template>

<script>
import ICON from "./ICON";

export default {
  name: "Dialog",
  components: { ICON },
  data() {
    return {
      isVisible: false
    };
  },
  methods: {
    hide() {
      this.isVisible = false;
    },
    show() {
      this.isVisible = true;
    }
  }
};
</script>
