export const events = {
    // Launcher self-update, used in the SelfUpdate component.
    // New version is available, proceed with update.
    LauncherUpdateAvailable: "launcher-update-available",
    // Launcher self-update, used in the SelfUpdate and StatusBar components.
    // Update is in progress, user waiting for download to finish.
    // Payload: { progress: number, total: number }
    LauncherUpdateProgress: "launcher-update-progress",
    // Launcher self-update, used in the SelfUpdate component.
    // Update failed, but the launcher is still usable, and the user can retry or ignore the update.
    LauncherUpdateFailed: "launcher-update-failed",
    // Launcher self-update, used in the SelfUpdate component.
    // New version downloaded, proceed with update.
    LauncherUpdateDownloaded: "launcher-update-downloaded",
    // Launcher self-update, used in the SelfUpdate component.
    // Update is complete or no update required, launcher is ready to be used, open the app library.
    LauncherReady: "launcher-ready",
    // Application update, used in the StatusBar component.
    // Application update is in progress, user waiting for download to finish.
    // Payload: { id: string, progress: number, total: number }
    AppUpdateProgress: "app-update-progress",
    // Application update, used in the StatusBar component.
    // Update archive downloaded, extracting.
    AppUpdateExtracting: "app-update-extracting",
    // Application update, used in the StatusBar component.
    // Update failed, can retry or ignore.
    AppUpdateFailed: "app-update-failed",
    // Application update, used in the StatusBar component.
    // Application update completed and application is ready for launch.
    AppUpdateCompleted: "app-update-completed",
}

