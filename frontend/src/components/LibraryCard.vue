<template>
  <main class="card-wrapper">
    <div class="card" @click="navigateToAppOverview" v-bind:style="{'background-image': `url(${image})`}">
      <div class="card-image" v-bind:style="{'background-image': `url(${image})`}">
        <img v-bind:src="image" v-bind:alt="title"/>
      </div>
      <div class="card-header">
        <div class="card-title">
          <h3>{{ title }}</h3>
        </div>
        <div class="buttons">
          <a class="button effect-default" v-if="!status?.installed && !status?.installing" @click="installApp($event)">Install</a>
          <a class="button effect-default"
             v-if="status?.installed && !status?.installing && status?.updateAvailable === UpdateAvailability.Available"
             @click="updateApp($event)">Update</a>
          <a class="button effect-default" v-if="status?.installed&& !status?.installing" @click="launchApp($event)">Launch</a>
        </div>
      </div>
    </div>
  </main>
</template>

<script lang="ts">
import * as runtime from "../../wailsjs/runtime";
import {model} from "../../wailsjs/go/models";
import {onMounted, ref} from "vue";
import {useRouter} from "vue-router";
import {CheckForAppUpdates, InstallApp, IsAppInstalled, LaunchApp, UpdateApp} from "../../wailsjs/go/app/Launcher";
import {events} from "../common/events";
import {AppStatus, UpdateAvailability} from "../common";
import AppV2 = model.AppV2;

export default {
  props: {
    app: {
      type: Object,
      required: true
    }
  },
  setup(props: any) {
    /**
     * @description Router instance. Used to navigate to the application page when the card is clicked.
     */
    const router = useRouter();

    /**
     * @description The app metadata.
     */
    const app = ref({} as AppV2);

    /**
     * @description The status of the app.
     */
    const status = ref({} as AppStatus)

    /**
     * @description The image of the app.
     */
    const image = ref("");

    /**
     * @description The title of the app.
     */
    const title = ref("");

    /**
     * @description Launches the app.
     * @param e DOM event.
     */
    const launchApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      console.log("Launch app:", app.value.id);

      try {
        if (status.value.installing) {
          console.warn("App is installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (!status.value.installed) {
          console.warn("App is not installed.");
          return;
        }

        await LaunchApp(app.value.id);
      } catch (e) {
        console.error('LaunchApp:', e);
      }
    };

    /**
     * @description Installs the app.
     * @param e DOM event.
     */
    const installApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      try {
        if (status.value.installing) {
          console.warn("App is installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (status.value.installed) {
          console.warn("App is already installed.");
          return;
        }

        await InstallApp(app.value.id);
      } catch (e) {
        console.error('InstallApp:', e);
      }
    }

    /**
     * @description Updates the app.
     * @param e DOM event.
     */
    const updateApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      try {
        if (status.value.installing) {
          console.warn("App is installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (!status.value.updateAvailable) {
          console.warn("App is up to date.");
          return;
        }

        await UpdateApp(app.value.id);
      } catch (e) {
        console.error('UpdateApp:', e);
      }
    }

    /**
     * @description Navigates to the application overview page.
     */
    const navigateToAppOverview = () => {
      router.push({name: "AppOverview", params: {id: app.value.id}});
    };

    /**
     * @description Updates the app installed status.
     */
    const updateAppInstalledStatus = async () => {
      try {
        status.value.installed = await IsAppInstalled(app.value.id);
        status.value.installing = false;
        status.value.updateAvailable = await CheckForAppUpdates(app.value.id);

        console.log("card updateAppInstalledStatus update available:", status.value.updateAvailable);
      } catch (e) {
        status.value.installed = false;
        status.value.installing = false;
        status.value.updateAvailable = UpdateAvailability.Unknown;

        console.error('updateAppInstalledStatus:', e);
      }
    };

    /**
     * @description Track app installation progress and lock the installation button.
     */
    onMounted(async () => {
      app.value = props.app;
      title.value = app.value.name ?? "";
      image.value = app.value.files.entities.find((file: any) => file.type === "launcher-app-card-image")?.url ?? "";

      // Check if the app is installed.
      await updateAppInstalledStatus();

      // Track app installation progress and lock the installation button.
      runtime.EventsOn(events.AppUpdateProgress, (appMetadata: AppV2, current: number, total: number) => {
        if (appMetadata.id === app.value.id && current / total < 1) {
          status.value.installing = true;
        }
      });

      // Unlock the installation button when the app is installed.
      runtime.EventsOn(events.AppUpdateCompleted, async (appMetadata: AppV2) => {
        if (appMetadata.id === app.value.id) {
          await updateAppInstalledStatus();
        }
      });

      // Unlock the installation button when the app installation failed.
      runtime.EventsOn(events.AppUpdateFailed, async (appMetadata: AppV2, error: string) => {
        console.warn("AppUpdateFailed:", error);
        if (appMetadata.id === app.value.id) {
          await updateAppInstalledStatus();
        }
      })
    });

    return {
      status,
      image,
      title,
      navigateToAppOverview,
      installApp,
      updateApp,
      launchApp,
      UpdateAvailability
    };
  }
}
</script>

<style lang="scss" scoped>
main.card-wrapper {
  padding: 0 10px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  perspective: 2000px;
}

.card {
  display: flex;
  position: relative;
  flex-direction: column;
  width: 256px;
  height: 460px;
  border-radius: 0.25rem;
  overflow: hidden;
  box-shadow: 0 3px 7px 3px rgba(0, 0, 0, 0.25);
  transition: 333ms all ease-in-out;
  perspective: 2000px;
  background-position: center;
  background-origin: content-box;

  &::before {
    transition: inherit;
    position: absolute;
    top: 0;
    right: -50%;
    bottom: 0;
    left: -50%;
    content: "";
    opacity: 0.05;
    transform: rotate(30deg) translate(0, -80%);
    background: linear-gradient(
            to bottom,
            rgba(255, 255, 255, .2) 0%,
            rgba(255, 255, 255, .2) 90%,
            rgba(255, 255, 255, .4) 95%,
            rgba(255, 255, 255, 0) 100%
    );
  }

  &:hover {
    box-shadow: 0 10px 15px 1px rgba(0, 0, 0, 0.25);
    transform: rotateX(10deg) scale(1.05);

    &::before {
      opacity: 0.25;
      transform: rotate(30deg) translate(0, -60%);
    }
  }

  .card-image {
    flex: 1;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    overflow: hidden;

    img {
      display: block;
      width: 256px;
      height: 340px;
      object-fit: cover;
    }
  }

  .card-header {
    display: flex;
    box-sizing: border-box;
    flex-direction: column;
    width: 100%;
    height: auto;
    min-height: 142px; // to contain the card title and a button
    padding: 0.25rem;
    backdrop-filter: blur(15px);
    background: rgba(0, 0, 0, 0.75);

    .card-title {
      flex: 1;
      height: auto;
      padding: 0.25rem;
      color: #fff;

      h3 {
        font-weight: 400;
        font-size: 1.25rem;
        line-height: 1.5rem;
      }
    }

    .buttons {
      padding: 0.5rem;
    }
  }
}
</style>