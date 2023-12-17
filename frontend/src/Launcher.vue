<script lang="ts">
import StatusBar from "./components/StatusBar.vue";
import {OpenUrl} from "../wailsjs/go/app/Launcher";

export default {
  components: {
    StatusBar
  },
  setup() {
    const openLe7elWebsite = () => {
      OpenUrl("https://le7el.com");
    };

    return {openLe7elWebsite: openLe7elWebsite};
  }
}
</script>

<template>
  <main class="launcher">
    <router-view v-slot="{ Component, route }">
      <transition :name="route.meta.transition || 'fade'">
        <div :key="route.name" class="view">
          <component :is="Component" :key="route.fullPath"/>
        </div>
      </transition>
    </router-view>
    <StatusBar/>
    <div id="powered-by">
      Powered by <a @click="openLe7elWebsite()" style="text-decoration: underline">LE7EL</a>
    </div>
  </main>
</template>

<style lang="scss" scoped>
main.launcher {
  box-sizing: border-box;
  display: inline-flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
}

.view {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: auto;
  height: auto;
  padding: 0.25rem;
}

#powered-by {
  color: rgba(255, 255, 255, 0.75);
  position: absolute;
  bottom: 1rem;
  right: 1rem;
  font-size: 12px;
  z-index: 9999;
  font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  text-align: right;

  a {
    text-decoration: none;
  }
}
</style>
