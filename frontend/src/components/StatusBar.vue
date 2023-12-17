<script lang="ts">

import * as runtime from "../../wailsjs/runtime";
import {events} from "../common/events";
import {ref} from "vue";
import {model} from "../../wailsjs/go/models";
import AppV2 = model.AppV2;

export default {
  setup() {
    /**
     * @description The message to display.
     */
    const message = ref("");

    /**
     * @description The progress of the update.
     */
    const progress = ref(0);

    /**
     * @description Whether the update is in progress.
     */
    const loading = ref(false);

    // Listen for events.LauncherUpdateAvailable and update the message accordingly.
    runtime.EventsOn(events.LauncherUpdateProgress, (current: number, total: number) => {
      message.value = "Downloading update...";
      progress.value = Number.parseFloat((current / total).toFixed(2));
      loading.value = progress.value > 0 && progress.value < 1;
    });

    // Listen for events.LauncherUpdateAvailable and update the message accordingly.
    runtime.EventsOn(events.AppUpdateProgress, (app: AppV2, current: number, total: number) => {
      message.value = `Downloading ${app.name}...`;
      progress.value = Number.parseFloat((current / total).toFixed(2));
      loading.value = progress.value > 0 && progress.value < 1;
    });

    // Listen for events.LauncherUpdateAvailable and update the message accordingly.
    runtime.EventsOn(events.AppUpdateExtracting, (app: AppV2) => {
      message.value = `Extracting ${app.name}...`;
      progress.value = 1;
      loading.value = false;
    });

    // Listen for events.LauncherUpdateAvailable and update the message accordingly.
    runtime.EventsOn(events.AppUpdateCompleted, (_: AppV2) => {
      message.value = ``;
      progress.value = 0;
      loading.value = false;
    });

    return {message, progress, loading};
  }
}
</script>

<template>
  <main>
    <div class="wrapper">
      <div class="loading">
        <div class="dot" v-if="loading"></div>
      </div>
      <div class="content">
        {{ message }}
      </div>
    </div>
    <div class="bar">
      <div class="fill" v-bind:style="{'width': `${progress*100.0}%`}"></div>
    </div>
  </main>
</template>

<style lang="scss" scoped>
main {
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  font-size: 12px;
  font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  text-align: center;
  align-items: flex-start;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  height: 32px;

  .wrapper {
    display: flex;
    flex-direction: row;
    justify-items: flex-start;
    justify-content: flex-start;
    align-items: center;
    align-content: center;

    .loading {
      width: 20px;
      height: 20px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;

      .dot {
        height: 10px;
        width: 10px;
        position: relative;
        border-radius: 50%;
        margin: 0;
        background-color: #00aaaa;
        animation: pulse 1.2s infinite;
      }
    }

    .content {
      text-align: left;
    }
  }

  .bar {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 5px;
    background-color: transparent;

    .fill {
      transition: 333ms width ease-in-out;
      height: 100%;
      background-color: #00aaaa;
      width: 50%;
    }
  }
}

@keyframes pulse {
  0% {
    opacity: .5;
    transition: opacity .25s ease-in-out;
  }
  100% {
    opacity: 0;
    transition: opacity .25s ease-in-out;
  }
}
</style>
