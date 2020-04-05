export default {
  server: {
    port: 3000, // default: 3000
    host: '0.0.0.0', // default: localhost
  },
  mode: 'spa',
  /*
   ** Headers of the page
   */
  head: {
    title: '大头菜交易市场 -- 集合啦，动物森友会！',
    meta: [
      {charset: 'utf-8'},
      {name: 'viewport', content: 'width=device-width, initial-scale=1'},
      {
        hid: '大头菜交易市场',
        name: '大头菜交易市场  动物之森',
        content: '集合啦，动物森友会！',
      },
    ],
    link: [{rel: 'icon', type: 'image/x-icon', href: '/favicon.ico'}],
  },
  /*
   ** Customize the progress-bar color
   */
  loading: {color: '#fff'},
  /*
   ** Global CSS
   */
  css: ['@/assets/css/main.scss'],
  /*
   ** Nuxt.js modules
   */
  modules: [
    // Doc: https://buefy.github.io/#/documentation
    'nuxt-buefy',
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    '@nuxtjs/pwa',
  ],
  /*
   ** Plugins to load before mounting the App
   */
  plugins: ['~plugins/axios'],
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: [],
  /*
   ** Axios module configuration
   ** See https://axios.nuxtjs.org/options
   */
  axios: {
    baseURL: 'http://ac.kastner.cn/api/v1',
  },
  /*
   ** Build configuration
   */
  build: {
    /*
     ** You can extend webpack config here
     */
    babel: {
      // plugins: ['transform-vue-jsx'],
    },
    extend(config, ctx) {},
  },
  router: {
    middleware: ['routerChange']
  }
}
