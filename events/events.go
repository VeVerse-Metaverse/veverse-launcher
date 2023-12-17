// Package events contains all the events that are sent from the launcher to the UI.
package events

const (
	LauncherMetadata         = "launcher-metadata"          // launcher metadata received from the server
	LauncherUpdateAvailable  = "launcher-update-available"  // update available for launcher, proceed with update
	LauncherUpdateProgress   = "launcher-update-progress"   // update is in progress, user waiting for download to finish
	LauncherUpdateFailed     = "launcher-update-failed"     // update failed, but the launcher is still usable, and the user can retry or ignore the update
	LauncherUpdateDownloaded = "launcher-update-downloaded" // update downloaded, proceed with update
	LauncherReady            = "launcher-ready"             // launcher is ready to be used
	LauncherApps             = "launcher-apps"              // launcher apps received from the server
	LauncherApp              = "launcher-app"               // launcher app received from the server
	AppUpdateAvailable       = "app-update-available"       // app update available, can update or ignore
	AppUpdateProgress        = "app-update-progress"        // app update is in progress, user waiting for download to finish
	AppUpdateExtracting      = "app-update-extracting"      // app update archive downloaded, extracting files
	AppUpdateFailed          = "app-update-failed"          // app update failed, and the user can retry or ignore the update
	AppUpdateCompleted       = "app-update-completed"       // app update completed
)
