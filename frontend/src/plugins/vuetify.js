import Vue from "vue";
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import Vuetify from "vuetify";
import 'vuetify/src/styles/main.sass'

Vue.use(Vuetify);

export default new Vuetify({
  icons: {
    iconfont: "md"
  },
  theme: {
    themes: {
      light: {
        primary: "#16a085",
        secondary: '#2d4052'
      }
    }
  }
});
