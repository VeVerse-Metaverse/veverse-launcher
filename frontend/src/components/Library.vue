<template>
  <main>
    <header>
      <h1>Library</h1>
    </header>
    <div class="library">
      <LibraryCard v-for="app in apps" :key="app.id" :app="app"/>
    </div>
  </main>
</template>

<script lang="ts">
import LibraryCard from "./LibraryCard.vue";
import {onMounted, ref} from "vue";
import {GetLauncherMetadata, IndexLauncherApps} from "../../wailsjs/go/app/Launcher";
import {model} from "../../wailsjs/go/models";
import {useRouter} from "vue-router";
import LauncherV2 = model.LauncherV2;
import AppV2 = model.AppV2;

export default {
  name: "Library",
  components: {LibraryCard},
  setup() {
    /**
     * @description Router instance. Used to navigate to the app page.
     */
    const router = useRouter();

    /**
     * @description Launcher metadata.
     */
    const metadata = ref<LauncherV2>({} as LauncherV2);

    /**
     * @description List of launcher apps, to be displayed in the library as cards.
     */
    const apps = ref<AppV2[]>([] as AppV2[]);

    /**
     * @description Offset for the apps list.
     */
    const offset = 0;

    /**
     * @description Limit for the apps list. Todo: implement pagination.
     */
    const limit = 100;

    /**
     * @description Fetch the launcher metadata and apps.
     */
    onMounted(async () => {
      try {
        metadata.value = await GetLauncherMetadata();
        apps.value = await IndexLauncherApps(offset, limit);

        // If only one app is available, redirect to the app page.
        if (apps.value.length === 1) {
          console.log(`Redirecting to app ${apps.value[0].id}...`);
          await router.push(`/apps/${apps.value[0].id}/overview`);
        }
      } catch (e) {
        console.error(e);
      }
    });

    return {metadata, apps};
  }
}
</script>

<style lang="scss" scoped>
main {
  header {
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-content: flex-start;
    margin: 1rem 2rem 0;

    h1 {
      text-align: left;
      font-family: "Bebas Neue", "Readex Pro", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
      font-size: 2.5rem;
      font-weight: 500;
      margin: 0;
    }
  }

  .library {
    box-sizing: border-box;
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 100%;
    padding: 2rem;
  }
}
</style>