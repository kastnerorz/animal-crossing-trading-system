<template>
  <Opera :operaType="operaType"> </Opera>
</template>
<script>
import Opera from '../components/Opera'
export default {
  name: 'SELL',
  middleware: 'curDay',
  components: {Opera},
  data() {
    return {
      operaType: 'SELL',
    }
  },
}
</script>
