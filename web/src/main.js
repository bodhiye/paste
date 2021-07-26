import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import Clipboard from "clipboard";
import VueCookie from "vue-cookie";
import VueQrcode from "@chenfengyuan/vue-qrcode";

import App from "./App.vue";
import router from "./router/router";
import store from "./store/store";
import i18n from "./assets/i18n";
import api from "./api/api";
import markdownIt from "./utils/markdown-it";

import "@/style/global.css";

Vue.use(BootstrapVue);
Vue.use(VueCookie);
Vue.component("QRCode", VueQrcode);

Vue.prototype.api = api;
Vue.prototype.clipboard = Clipboard;
Vue.prototype.markdown = markdownIt;

api
  .get("/config.json")
  .then((response) => {
    store.state.config = response;
    return api.get(store.state.config.api.backend + "health", {
      json: "true",
    });
  })
  .then(() => {
    return new Vue({
      render: (h) => h(App),
      i18n,
      router,
      store,
    }).$mount("#app");
  });

console.log(`
    欢迎使用代码便利贴~
    公众号: IT夜谈
    博客: https://yeqiongzhou.com
    Email: yeqiongzhou@whu.edu.cn
`);
