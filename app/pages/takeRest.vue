<template>
  <div class="rest-page">
    <TopMenu opt="rest" />
    <section class="form-wrapper rest-wrapper">
       <p>{{restRes}}</p>
       <p>休息一下啦</p>
       <img src="../assets/rest.jpg" alt="">
    </section>
  </div>
</template>

<script>
import TopMenu from "../components/TopMenu";
export default {
  name: "REST",
  components: { TopMenu },
  computed: {
    restRes() {
      const type = this.$route.query.type || ''
      let restRes = ''
      switch (type) {
        case 'buy':
          restRes = '今天只开放卖出'
          break;
        case 'sell':
          restRes = '今天只开放买入'
          break;
        default:
          restRes = '休市了,不要过度劳累'
          break;
      }
      return restRes
    }
  },
};
</script>
<style lang="scss" scoped>
  .rest-wrapper{
    border-radius: 27px;
    background: #fcf9f2;
    margin: 5%;
    box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.25);
    padding:  5%;
    color: #a0a0a0;
    font-weight: bold;
    font-size: 20px;
    text-align: center;
    line-height: 3;
  }
</style>