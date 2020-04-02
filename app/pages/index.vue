<template>
  <Opera :operaType="operaType"></Opera>
</template>
<script>
import Opera from "../components/Opera";
export default {
  name: "BUY",
  middleware: "curDay",
  components: { Opera },
  data() {
    return {
      operaType: "BUY"
    };
  }
};
</script>
