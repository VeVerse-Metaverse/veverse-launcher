<template>
  <section class="app-content">
    <div class="app-header">
      <div class="app-title">
        <h1>{{ app?.name }}</h1>
      </div>
      <div class="app-release">
        <h2>{{ formatReleaseVersion(release) }}</h2>
        <p>{{ formatReleaseDescription(release) }}</p>
      </div>
    </div>
    <div class="app-description">
      <p>{{ app?.description }}</p>
    </div>
    <div class="app-links">
      <ul>
        <li v-for="link in app?.links" :key="link.id">
          <span>{{ link.name }}:</span> <a v-bind:href="link.url">{{ link.url }}</a>
        </li>
      </ul>
    </div>
  </section>

  <aside class="app-controls">
    <div class="buttons">
      <a class="button effect-default" @click="installApp($event)"
         v-if="!status?.installed && !status?.installing">Install</a>
      <a class="button effect-default" @click="updateApp($event)"
         v-if="status?.installed && !status?.installing && status?.updateAvailable === UpdateAvailability.Available">Update</a>
      <a class="button effect-default" @click="launchApp($event)"
         v-if="status?.installed && !status?.installing">Launch</a>
      <a class="button effect-default" @click="deleteApp($event)"
         v-if="0 && status?.installed && !status?.installing">Delete</a>
    </div>
    <!--    <div class="info">
          <ul>
            <li>Size - 1.1 Gb</li>
          </ul>
        </div>-->
  </aside>
</template>

<script lang="ts">
import * as runtime from "../../wailsjs/runtime";
import {model} from "../../wailsjs/go/models";
import {onMounted, ref} from "vue";
import {
  CheckForAppUpdates,
  DeleteApp,
  GetAppMetadata,
  InstallApp,
  IsAppInstalled,
  LaunchApp,
  UpdateApp
} from "../../wailsjs/go/app/Launcher";
import {AppStatus, UpdateAvailability} from "../common";
import {events} from "../common/events";
import {useRoute} from "vue-router";
import ReleaseV2 = model.ReleaseV2;
import AppV2 = model.AppV2;

export default {
  name: "AppOverview",
  setup() {
    /**
     * @description Route instance. Used to get the app id.
     */
    const route = useRoute();

    /**
     * @description App id.
     */
    const id = route.params.id;

    /**
     * @description App metadata.
     */
    const app = ref({} as AppV2);

    /**
     * @description Latest app release metadata.
     */
    const release = ref({} as ReleaseV2);

    /**
     * @description The status of the app.
     */
    const status = ref({} as AppStatus)

    /**
     * @description Launches the app.
     */
    const launchApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

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
        console.error(e);
      }
    };

    /**
     * @description Installs the app.
     */
    const installApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      try {
        if (status.value.installing) {
          console.log("App is already installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (status.value.installed) {
          console.log("App is already installed.");
          return;
        }

        await InstallApp(app.value.id);
      } catch (e) {
        console.error(e);
      }
    };

    /**
     * @description Updates the app.
     */
    const updateApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      try {
        if (status.value.installing) {
          console.log("App is already installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (!status.value.installed) {
          console.log("App is not installed.");
          return;
        }

        await UpdateApp(app.value.id);
      } catch (e) {
        console.error(e);
      }
    };

    /**
     * @description Deletes the app.
     */
    const deleteApp = async (e: Event) => {
      e.preventDefault();
      e.stopPropagation();

      try {
        if (status.value.installing) {
          console.log("App is installing.");
          return;
        }

        await updateAppInstalledStatus();

        if (!status.value.installed) {
          console.log("App is not installed.");
          return;
        }

        await DeleteApp(app.value.id);

        await updateAppInstalledStatus();
      } catch (e) {
        console.error(e);
      }
    };

    /**
     * @description Formats the release version, including the release name if available.
     */
    const formatReleaseVersion = (release: ReleaseV2) => {
      if (release.version === "") {
        return "";
      } else {
        if (release.name) {
          return release.version + " - " + release.name;
        } else {
          return release.version;
        }
      }
    }

    /**
     * @description Formats the release date.
     */
    const formatReleaseDate = (release: ReleaseV2): string => {
      // noinspection SuspiciousTypeOfGuard
      if (typeof (release.createdAt) === 'string') {
        release.createdAt = new Date(release.createdAt);
      }
      return release.createdAt?.toLocaleString() || '';
    };

    /**
     * @description Formats the release description.
     */
    const formatReleaseDescription = (release: ReleaseV2): string => {
      if (release && release.description) {
        return `(${release?.description || ''}, ${formatReleaseDate(release)})`;
      } else {
        return "";
      }
    };

    /**
     * @description Updates the app installed status.
     */
    const updateAppInstalledStatus = async () => {
      try {
        status.value.installed = await IsAppInstalled(app.value.id);
        status.value.installing = false;
        status.value.updateAvailable = await CheckForAppUpdates(app.value.id);
        console.log("updateAppInstalledStatus update available:", status.value.updateAvailable);
      } catch (e) {
        status.value.installed = false;
        status.value.installing = false;
        status.value.updateAvailable = UpdateAvailability.Unknown;

        console.error('updateAppInstalledStatus:', e);
      }
    };

    /**
     * @description Updates the app installed status and subscribes to events.
     */
    onMounted(async () => {
      try {
        // Retrieve the app metadata.
        app.value = await GetAppMetadata(id);

        // Convert the dates to Date objects if they are strings.
        // noinspection SuspiciousTypeOfGuard
        if (typeof app.value.createdAt === 'string') {
          app.value.createdAt = new Date(app.value.createdAt);
        }
        // noinspection SuspiciousTypeOfGuard
        if (typeof app.value.updatedAt === 'string') {
          app.value.updatedAt = new Date(app.value.updatedAt);
        }

        // Retrieve the latest release from list of releases.
        if (app.value.releases?.entities && app.value.releases.entities.length > 0) {
          for (const release of app.value.releases.entities) {
            // Convert the dates to Date objects if they are strings.
            // noinspection SuspiciousTypeOfGuard
            if (typeof release.createdAt === 'string') {
              release.createdAt = new Date(release.createdAt);
            }
            // noinspection SuspiciousTypeOfGuard
            if (typeof release.updatedAt === 'string') {
              release.updatedAt = new Date(release.updatedAt);
            }
          }
          release.value = app.value.releases?.entities[0] ?? {} as ReleaseV2;
        }
      } catch (e) {
        console.error(e);
      }

      await updateAppInstalledStatus();

      runtime.EventsOn(events.AppUpdateProgress, (appMetadata: AppV2, current: number, total: number) => {
        if (appMetadata.id === app.value.id && current / total < 1) {
          status.value.installing = true;
        }
      });

      runtime.EventsOn(events.AppUpdateCompleted, async (appMetadata: AppV2) => {
        if (appMetadata.id === app.value.id) {
          await updateAppInstalledStatus();
        }
      });

      runtime.EventsOn(events.AppUpdateFailed, async (appMetadata: AppV2, error: string) => {
        console.warn("AppUpdateFailed:", error);
        if (appMetadata.id === app.value.id) {
          await updateAppInstalledStatus();
        }
      });
    });

    return {
      app,
      release,
      status,
      installApp,
      updateApp,
      deleteApp,
      launchApp,
      formatReleaseVersion,
      formatReleaseDate,
      formatReleaseDescription,
      UpdateAvailability
    };
  }
}
</script>

<style lang="scss" scoped>

.app-content {
  position: relative;
  flex: 1;
  text-align: left;
  margin: 1rem;
  padding: 2rem;
  background-image: linear-gradient(150deg, rgba(0 0 0 / 40%) 0%, rgba(0 0 0 / 70%) 100%);
  border-radius: 10px;
  backdrop-filter: blur(5px);
  box-shadow: 0 2px 5px 1px rgb(0 0 0 / 65%);

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    border-radius: 10px;
    padding: 2px;
    background: linear-gradient(
            135deg,
            rgba(255 0 255 / 50%) 0%,
            rgba(255 255 255 / 0%) 40%,
            rgba(255 255 255 / 0%) 60%,
            rgba(0, 255, 255, 0.5) 100%
    );
    //noinspection CssInvalidPropertyValue,CssInvalidFunction
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: destination-out;
    mask-composite: exclude;
    user-select: none;
    z-index: -1;
  }

  &::after {
    position: absolute;
    content: "";
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    border-radius: 20px;
    //noinspection CssUnknownTarget
    background-image: url(data:image/png;base64,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);
    background-repeat: repeat;
    background-size: 80px;
  }
}

.app-controls {
  box-sizing: border-box;
  flex: 0 0 10rem;
  margin: 1rem;
  padding: 1rem;

  .buttons {
    box-sizing: border-box;
    margin: auto;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    .button {
      box-sizing: border-box;
      margin: 0 auto 0.5rem;
    }
  }
}

.app-release {
  h2 {
    font-family: "Bebas Neue", "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
    font-size: 1.5rem;
    font-weight: 400;
    margin: 1rem auto 0;
  }

  p {
    font-size: 0.9rem;
    color: #ccc;
    font-weight: 400;
    margin: 0 auto 1rem;
  }
}

.app-title {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;

  h1 {
    font-family: "Bebas Neue", "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
    font-size: 2.5rem;
    font-weight: 400;
    color: #fff;
    margin: 0;
  }
}

.app-header {
  .buttons {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;

    .button {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0.5rem 1rem;
      border-radius: 0.25rem;
      font-size: 1rem;
      font-weight: 700;
      color: #fff;
      background-color: #000;
      border: 1px solid #fff;
      cursor: pointer;
      transition: 333ms background-color ease-in-out, 333ms color ease-in-out, 333ms border-color ease-in-out;

      &:hover {
        background-color: #fff;
        color: #000;
        border-color: #000;
      }
    }
  }
}
</style>