<template>
  <main class="update">
    <div v-if="updateFailed">
      <p class="icon">
        <svg width="64px" height="64px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path d="M12,2A10,10,0,1,0,22,12,10,10,0,0,0,12,2Z"></path>
          <line x1="12" y1="8" x2="12" y2="12"></line>
          <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
      </p>
      <p class="message">{{ message }}</p>
      <p class="buttons">
        <a class="button effect-default" @click="updateFailedRetry()"><span>Retry</span></a>
        <a class="button effect-default" @click="updateFailedContinue()"><span>Continue</span></a>
      </p>
    </div>
    <div v-else>
      <p class="icon-download">

      </p>
      <p class="icon">
        <!--        <svg width="64px" height="64px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12,2A10,10,0,1,0,22,12,10,10,0,0,0,12,2Z"></path>
                  <polyline points="8 17 12 21 16 17"></polyline>
                  <line x1="12" y1="12" x2="12" y2="21"></line>
                </svg>-->
        <svg width="64px" height="64px" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path d="M12,2A10,10,0,1,0,22,12,10,10,0,0,0,12,2Z"></path>
          <polyline points="12 6 12 12 16 14"></polyline>
        </svg>
      </p>
      <p class="message">{{ message }}</p>
      <p class="progress" v-if="progress > 0 && progress < 1">{{ filterPercent(progress) }}</p>
    </div>
  </main>
</template>

<script lang="ts" setup>
import {onMounted, ref} from "vue";
import * as runtime from "../../wailsjs/runtime/runtime";
import {events} from "../common/events";
import {useRouter} from "vue-router";
import {UpdateLauncher} from "../../wailsjs/go/app/Launcher";
import {errors} from "../errors";

/**
 * @description Router instance. Used to navigate to the library page when the update is complete or failed.
 */
const router = useRouter();
/**
 * @description Progress of the update, in range [0, 1]. Used to display the progress percentage message.
 */
const progress = ref(0);
/**
 * @description Used to display the status message. Used both for progress and error messages.
 */
const message = ref("Checking for updates...");
/**
 * @description Whether the update failed. Toggles display of retry/continue buttons.
 */
const updateFailed = ref(false);

/**
 * @description Checks for updates on mount.
 */
onMounted(async () => {
  try {
    // Check for updates and update the launcher if available.
    await UpdateLauncher();
  } catch (e) {
    // No update available is not an error, just proceed to the library.
    if (e === errors.NoUpdateAvailable) {
      updateFailed.value = false;
      message.value = "No update available";
      await router.push("/library");
      return;
    }

    // Otherwise, display the error message.
    console.error(e);
    updateFailed.value = true;
    message.value = "Failed to check for updates";
  }
});

// Listen for events.LauncherUpdateProgress and update the progress and message accordingly.
runtime.EventsOn(events.LauncherUpdateProgress, (current: number, total: number) => {
  message.value = "Downloading update...";
  progress.value = Number.parseFloat((current / total).toFixed(2));
});

// Listen for events.LauncherUpdateAvailable and update the message accordingly.
runtime.EventsOn(events.LauncherUpdateAvailable, () => {
  message.value = "Update available";
});

// Listen for events.LauncherUpdateDownloaded and update the message accordingly.
runtime.EventsOn(events.LauncherUpdateDownloaded, () => {
  message.value = "Restarting...";
});

// Listen for events.LauncherUpdateFailed and update the updateFailed status and message accordingly.
runtime.EventsOn(events.LauncherUpdateFailed, () => {
  message.value = "Failed to update, please check your internet connection and try again later.";
  updateFailed.value = true;
});

// Listen for events.LauncherReady and navigate to the library page.
runtime.EventsOn(events.LauncherReady, () => {
  router.push("/library");
});

/**
 * @description Filters the progress value to a percentage string.
 * @param value Progress value in range [0, 1].
 */
const filterPercent = (value: number) => {
  if (value < 0) {
    return "0%"
  } else if (value > 1) {
    return "100%"
  }
  return `${Math.floor(value * 100)}%`
}

/**
 * @description Navigates to the library page when user clicks the continue button after update failed.
 */
const updateFailedContinue = () => {
  router.push("/library");
}

/**
 * @description Retries the update when user clicks the retry button after update failed.
 */
const updateFailedRetry = async () => {
  // Reset the progress and message.
  message.value = "Retrying...";
  updateFailed.value = false;
  progress.value = 0;

  try {
    // Check for updates and update the launcher if available.
    await UpdateLauncher();
  } catch (e) {
    // No update available is not an error, just proceed to the library.
    if (e == errors.NoUpdateAvailable) {
      await router.push("/library");
      return;
    }

    // Otherwise, display the error message.
    console.error(e);
    updateFailed.value = true;
    message.value = "Failed to check for updates";
  }
}
</script>

<style lang="scss" scoped>
main.update {
  display: flex;
  flex-direction: column;
  height: 100%;
  box-sizing: border-box;
  justify-content: center;
  align-items: center;
}

p.progress {
  font-family: "Bebas Neue", "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  font-size: 2.5rem;
  font-weight: 500;
  margin: 0;
  padding: 0;
}

p.message {
  font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  font-size: 1.0rem;
  font-weight: 500;
  margin: 2rem 0;
  padding: 0;
}

.icon svg {
  -webkit-animation: spin 3s linear infinite;
  -moz-animation: spin 3s linear infinite;
  animation: spin 3s linear infinite;

  path, polyline {
    fill: none;
    stroke: #ffffff;
    stroke-linecap: round;
    stroke-linejoin: round;
    stroke-width: 2;
  }
}

@-moz-keyframes spin {
  100% {
    -moz-transform: rotate(360deg);
  }
}

@-webkit-keyframes spin {
  100% {
    -webkit-transform: rotate(360deg);
  }
}

@keyframes spin {
  100% {
    -webkit-transform: rotate(360deg);
    transform: rotate(360deg);
  }
}
</style>