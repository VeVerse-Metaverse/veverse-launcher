<template>
  <nav class="library-navigation" v-if="launcher?.apps?.entities?.length > 0">
    <ul>
      <li>
        <a @click="router.go(-1)">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2"
               stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </a>
      </li>
      <li>
        <router-link to="/library">
          <span>Library</span>
        </router-link>
      </li>
    </ul>
  </nav>
  <main class="app" v-bind:style="{'background-image': `url(${image})`}">
    <nav class="app-navigation" v-if="app?.sdk">
      <ul>
        <li class="app-navigation__item">
          <router-link to="overview">Overview</router-link>
        </li>
        <li class="app-navigation__item">
          <router-link to="changelog">SDK</router-link>
        </li>
      </ul>
    </nav>
    <router-view v-slot="{ Component, route }">
      <transition :name="route.meta.transition || 'fade'">
        <div :key="route.name" class="app-body">
          <component :is="Component" :key="route.fullPath"/>
        </div>
      </transition>
    </router-view>
  </main>
</template>

<script lang="ts">
import {useRoute, useRouter} from "vue-router";
import {model} from "../../wailsjs/go/models";
import {onMounted, ref} from "vue";
import {GetAppMetadata, GetLauncherMetadata, LaunchApp} from "../../wailsjs/go/app/Launcher";
import AppV2 = model.AppV2;
import ReleaseV2 = model.ReleaseV2;
import LauncherV2 = model.LauncherV2;

/**
 * @description App component.
 */
export default {
  name: "App",
  setup() {
    /**
     * @description Router instance. Used for the navigate back button.
     */
    const router = useRouter();

    /**
     * @description Route instance. Used to get the app id.
     */
    const route = useRoute();

    /**
     *  @description App id extracted from the current route.
     */
    const id = route.params.id;

    /**
     * @description Launcher metadata. Used to detect number of apps.
     */
    const launcher = ref<LauncherV2>({} as LauncherV2);

    /**
     * @description App metadata.
     */
    const app = ref({} as AppV2);

    /**
     * @description App release metadata.
     */
    const release = ref({} as ReleaseV2);

    /**
     * @description App background image.
     */
    const image = ref("");

    /**
     * @description Formats the release printing the version and the name.
     * @param release Release to format.
     */
    const formatRelease = (release: ReleaseV2 | undefined) => {
      if (release) {
        return release.version + " - " + release.name;
      }
      return "";
    }

    /**
     * @description Launches the app.
     */
    const launchApp = async () => {
      try {
        await LaunchApp(app.value.id);
      } catch (e) {
        console.error(e);
      }
    };

    /**
     * @description Fetches the app metadata.
     */
    onMounted(async () => {
      try {
        // Retrieve the launcher metadata.
        launcher.value = await GetLauncherMetadata();

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

        // Extract the background image from list of files of the app.
        image.value = app.value.files.entities.find((file: any) => file.type === "launcher-app-background-image")?.url ?? "";

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
    });

    return {launcher, id, app, release, image, launchApp, formatRelease, router};
  }
}
</script>

<style lang="scss" scoped>
nav {
  font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  display: flex;
  flex-direction: row;

  &.app-navigation {
    font-size: 1.25rem;
    font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;

    ul {
      margin-left: 0.5rem;

      li a {
        text-outline: #000 1px 1px;
        transition: 333ms border-bottom ease-in-out, 333ms color ease-in-out, 333ms text-shadow ease-in-out;
        border-bottom: 1px solid transparent;
        color: #eee;
        text-shadow: -1px -1px 0 #000, 1px -1px 0 #000, -1px 1px 0 #000, 1px 1px 0 #000;

        &.active {
          border-bottom: 1px solid #fff;
          color: #fff;
        }

        &:hover {
          color: #fff;
          border-bottom: 1px solid #eee;
        }
      }
    }
  }

  ul {
    display: flex;
    flex-direction: row;
    list-style: none;
    padding: 0;
    margin: 1rem 0;

    li {
      display: flex;
      padding: 0 0.5rem;
      flex-direction: column;
      justify-content: center;

      a {
        cursor: pointer;
        display: flex;
        color: #777;
        transition: 333ms color ease-in-out;
        text-decoration: none;

        &:hover {
          color: #ccc;
        }

        svg {
          width: 1.5rem;
          height: 1.5rem;
        }
      }
    }
  }
}

main.app {
  font-family: "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
  font-weight: 400;
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 1rem;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-color: #000;

  .app-body {
    display: flex;
    flex-direction: row;
    flex: 1;
    align-content: flex-start;
  }

}
</style>