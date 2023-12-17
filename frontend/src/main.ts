import {createApp} from 'vue'
import Launcher from './Launcher.vue'
import './style.scss'
import router from "./router";

const launcher = createApp(Launcher as any);
launcher.use(router)
launcher.mount('#launcher')
