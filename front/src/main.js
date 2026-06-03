import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
//引入iconPark样式文件
import '@icon-park/vue/styles/index.css'
import '../src/assets/gloable.css'

import request from "../src/request";
Vue.prototype.request=request

Vue.config.productionTip = false

Vue.use(ElementUI);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
