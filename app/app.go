package app

import (
	"context"
	"dev.hackerman.me/artheon/veverse-shared/executable"
	sm "dev.hackerman.me/artheon/veverse-shared/model"
	sp "dev.hackerman.me/artheon/veverse-shared/platform"
	"dev.hackerman.me/artheon/veverse-shared/unreal"
	"embed"
	"errors"
	"fmt"
	"games.launch.launcher/api"
	"games.launch.launcher/config"
	"games.launch.launcher/events"
	"games.launch.launcher/http"
	"games.launch.launcher/model"
	"games.launch.launcher/utils"
	"games.launch.launcher/version"
	"github.com/Masterminds/semver"
	"github.com/gofrs/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	goRuntime "runtime"
	"strings"
	"syscall"
	"time"
)

//go:embed updater/bin/updater.exe
var updater embed.FS

var ApplicationsDir = "apps"

var ErrorNoUpdateAvailable = errors.New("no update available")
var ErrorNoMetadata = errors.New("no metadata")
var ErrorNoReleases = errors.New("no releases")
var ErrorNoReleaseFiles = errors.New("no release files")
var ErrorNoReleaseFileSize = errors.New("no release file size")
var ErrorNoReleaseFileUrl = errors.New("no release file url")
var ErrorAppInstalled = errors.New("app installed")
var ErrorLauncherIsUpdating = errors.New("launcher is updating")
var ErrorAppIsUpdating = errors.New("app is updating")

type UpdateAvailability int

const (
	UpdateAvailabilityUpToDate UpdateAvailability = iota
	UpdateAvailabilityAvailable
	UpdateAvailabilityUnknown = -1
)

const RetryCount = 10

// Launcher struct
type Launcher struct {
	Ctx context.Context

	//region Persistent data

	GameConnPool       map[string]*net.Conn
	Metadata           *sm.LauncherV2
	UpdateAvailability UpdateAvailability
	IsUpdatingLauncher bool
	IsUpdatingApp      bool
	LastEvent          string

	Status model.Status `json:"status"` // the app status
	//endregion
}

// NewLauncher creates a new Launcher application struct
func NewLauncher() *Launcher {
	return &Launcher{
		Metadata:           nil,
		UpdateAvailability: UpdateAvailabilityUnknown,
		GameConnPool:       make(map[string]*net.Conn),
		Status: model.Status{
			Downloading:     false,
			Progress:        0,
			UpdateAvailable: false,
			ShowButtons:     false,
			NextVersion:     "",
		},
	}
}

// OnStartup is called when the app starts, requests the launcher metadata
func (l *Launcher) OnStartup(ctx context.Context) {
	l.Ctx = ctx

	// Start the first instance and listen for subsequent instance connections.
	go l.StartFirstInstance()

	// Start the game client listener for communication with the game client.
	go l.StartGameClientListener()

	if os.Args != nil && len(os.Args) > 1 {
		// Process the deep link if main instance was started with one.
		l.processDeepLink(os.Args[1])
	}

	//if err := l.UpdateLauncher(); err != nil && err != ErrorNoUpdateAvailable {
	//	runtime.LogErrorf(l.Ctx, "failed to update launcher: %w", err)
	//	l.EmitEvent(events.LauncherUpdateFailed)
	//	return
	//}
	//l.EmitEvent(events.LauncherReady)
}

// GetLauncherMetadata requests the app metadata from the backend
func (l *Launcher) GetLauncherMetadata() (*sm.LauncherV2, error) {
	runtime.LogInfof(l.Ctx, "GetLauncherMetadata")

	if config.LauncherId == "" {
		runtime.LogErrorf(l.Ctx, "launcher id is not set")
		return nil, fmt.Errorf("launcher id is not set")
	}
	launcherId := uuid.FromStringOrNil(config.LauncherId)

	var err error
	for i := 0; i < RetryCount; i++ {
		l.Metadata, err = api.GetLauncherMetadata(l.Ctx, launcherId)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get launcher metadata: %w", err)
		return nil, fmt.Errorf("failed to get launcher metadata: %w", err)
	}

	l.EmitEvent(events.LauncherMetadata, l.Metadata)

	return l.Metadata, nil
}

// SetLauncherUpdateStatus sets the update status of the launcher
func (l *Launcher) SetLauncherUpdateStatus(isUpdating bool, event string, args ...any) {
	l.IsUpdatingLauncher = isUpdating
	l.EmitEvent(event, args...)
}

// UpdateLauncher updates the launcher
func (l *Launcher) UpdateLauncher() error {
	runtime.LogWarningf(l.Ctx, "UpdateLauncher")

	if l.IsUpdatingLauncher {
		runtime.LogDebugf(l.Ctx, "launcher is already updating")
		return ErrorLauncherIsUpdating
	}

	_, err := l.CheckForUpdates()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to check for updates: %w", err)
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if l.UpdateAvailability != UpdateAvailabilityAvailable {
		// Delete the updater executable if it exists.

		// Get the path to the executable
		executablePath, err := os.Executable()
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to get current dir: %v", err)
			return fmt.Errorf("failed to get current dir: %v", err)
		}

		// Get the directory containing the executable
		dir := filepath.Dir(executablePath)

		// Cleanup updater executable, log file and temp dir.
		if //goland:noinspection GoBoolExpressions
		goRuntime.GOOS == "windows" {
			if _, err := os.Stat(filepath.Join(dir, "updater.exe")); err == nil {
				if err := os.Remove(filepath.Join(dir, "updater.exe")); err != nil {
					runtime.LogErrorf(l.Ctx, "failed to delete updater executable: %v", err)
					return fmt.Errorf("failed to delete updater executable: %w", err)
				}
			}
		} else {
			if _, err := os.Stat(filepath.Join(dir, "updater")); err == nil {
				if err := os.Remove(filepath.Join(dir, "updater")); err != nil {
					runtime.LogErrorf(l.Ctx, "failed to delete updater executable: %v", err)
					return fmt.Errorf("failed to delete updater executable: %w", err)
				}
			}
		}

		if logStat, err := os.Stat(filepath.Join(dir, "updater.log")); err == nil {
			if logStat.Size() == 0 {
				if err := os.Remove(filepath.Join(dir, "updater.log")); err != nil {
					runtime.LogErrorf(l.Ctx, "failed to delete updater log: %v", err)
					return fmt.Errorf("failed to delete updater log: %w", err)
				}
			}
		}

		if _, err := os.Stat(filepath.Join(dir, ".tmp")); err == nil {
			if err := os.RemoveAll(filepath.Join(dir, ".tmp")); err != nil {
				runtime.LogErrorf(l.Ctx, "failed to delete temp dir: %v", err)
				return fmt.Errorf("failed to delete temp dir: %w", err)
			}
		}

		return ErrorNoUpdateAvailable
	}

	if l.Metadata == nil {
		runtime.LogErrorf(l.Ctx, "launcher metadata is nil")
		return ErrorNoMetadata
	}

	var release sm.ReleaseV2
	if l.Metadata.Releases == nil {
		runtime.LogErrorf(l.Ctx, "launcher metadata releases is nil")
		return ErrorNoReleases
	}
	if len(l.Metadata.Releases.Entities) == 0 {
		runtime.LogErrorf(l.Ctx, "launcher metadata releases is empty")
		return ErrorNoReleases
	} else {
		release = l.Metadata.Releases.Entities[0]
	}

	releaseVersion, err := semver.NewVersion(release.Version)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to parse release version: %v", err)
		return fmt.Errorf("failed to parse release version: %w", err)
	}

	// print release data
	runtime.LogDebugf(l.Ctx, "release version: %s", releaseVersion.String())
	runtime.LogDebugf(l.Ctx, "release file number: %d", len(release.Files.Entities))

	var file *sm.File
	if release.Files == nil || len(release.Files.Entities) == 0 {
		runtime.LogErrorf(l.Ctx, "no release files")
		return ErrorNoReleaseFiles
	}

	// Find file with the correct platform and type.
	for _, f := range release.Files.Entities {
		if f.Platform == unreal.GetPlatformName() && f.Type == "launcher" {
			file = &f
			break
		}
	}

	if file == nil {
		runtime.LogErrorf(l.Ctx, "no release file for platform %s", unreal.GetPlatformName())
		return ErrorNoReleaseFiles
	}

	if file.Url == "" {
		runtime.LogErrorf(l.Ctx, "no release file url")
		return ErrorNoReleaseFileUrl
	}

	if file.Size == nil {
		runtime.LogErrorf(l.Ctx, "no release file size")
		return ErrorNoReleaseFileSize
	}
	fileSize := (uint64)(*file.Size)

	url := file.Url
	if url == "" || !strings.HasPrefix(url, "http") {
		runtime.LogErrorf(l.Ctx, "invalid url %s", url)
		return fmt.Errorf("invalid url %s", url)
	}

	l.SetLauncherUpdateStatus(true, events.LauncherUpdateProgress)

	downloadDir, err := getDownloadDir()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get download dir: %v", err)
		return fmt.Errorf("failed to get download dir: %w", err)
	}

	var fileName string
	if file.OriginalPath != nil {
		fileName = *file.OriginalPath
	} else {
		fileName = file.Id.String()
	}

	sourcePath := filepath.Join(downloadDir, fileName)

	counter := http.NewDownloadProgressTracker(fileSize, func(progress uint64, total uint64) {
		l.EmitEvent(events.LauncherUpdateProgress, progress, total)
	})
	err = http.DownloadFile(l.Ctx, sourcePath, url, counter)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to download file: %v", err)
		return fmt.Errorf("failed to download file: %w", err)
	}

	destinationPath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable: %v", err)
		return fmt.Errorf("failed to get executable: %w", err)
	}

	destinationDir := filepath.Dir(destinationPath)

	var updaterPath, updaterDestinationPath string
	if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "linux" {
		updaterPath = "updater/bin/updater"
		updaterDestinationPath = filepath.Join(destinationDir, "updater")
	} else if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "darwin" {
		updaterPath = "updater/bin/updater"
		updaterDestinationPath = filepath.Join(destinationDir, "updater")
	} else if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "windows" {
		updaterPath = "updater/bin/updater.exe"
		updaterDestinationPath = filepath.Join(destinationDir, "updater.exe")
	}
	// Read embedded updater.
	b, err := updater.ReadFile(updaterPath)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to read embedded updater: %v", err)
		return fmt.Errorf("failed to read embedded updater: %w", err)
	}
	// Write updater to disk.
	err = os.WriteFile(updaterDestinationPath, b, 0644)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to write updater to disk: %v", err)
		return fmt.Errorf("failed to create updater: %w", err)
	}

	// Start updater without logging (add third argument to the command to enable logging).
	var cmd *exec.Cmd
	if config.Logging == "true" {
		cmd = exec.Command(updaterDestinationPath, sourcePath, destinationPath, "true")
	} else {
		cmd = exec.Command(updaterDestinationPath, sourcePath, destinationPath)
	}

	// Hide updater window.
	if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	err = cmd.Start()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to start self update script: %v", err)
		return fmt.Errorf("failed to start self update script: %w", err)
	}

	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %v", err)
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	err = version.WriteVersion(executableDir, releaseVersion)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to set current launcher version: %v", err)
		return fmt.Errorf("failed to set current launcher version: %w", err)
	}

	l.EmitEvent(events.LauncherUpdateDownloaded)
	time.Sleep(1 * time.Second)

	os.Exit(0)

	return nil
}

// CheckForUpdates checks for updates
func (l *Launcher) CheckForUpdates() (UpdateAvailability, error) {
	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %v", err)
		return UpdateAvailabilityUnknown, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	currentVersion, err := version.ReadVersion(executableDir)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get current launcher version: %w", err)
		return UpdateAvailabilityUnknown, err
	}

	runtime.LogInfof(l.Ctx, "current version: %s", currentVersion)

	if l.Metadata == nil {
		_, err = l.GetLauncherMetadata()
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to get launcher metadata: %w", err)
			return UpdateAvailabilityUnknown, err
		}
	}

	var latestVersion *semver.Version
	for _, release := range l.Metadata.Releases.Entities {
		if latestVersion == nil {
			latestVersion, err = semver.NewVersion(release.Version)
			if err != nil {
				runtime.LogErrorf(l.Ctx, "failed to parse version: %w", err)
				return UpdateAvailabilityUnknown, err
			}
		} else {
			releaseVersion, err := semver.NewVersion(release.Version)
			if err != nil {
				runtime.LogErrorf(l.Ctx, "failed to parse semver: %w", err)
				return UpdateAvailabilityUnknown, fmt.Errorf("failed to parse semver: %w", err)
			}

			if latestVersion.GreaterThan(releaseVersion) {
				latestVersion = releaseVersion
			}
		}
	}

	if latestVersion == nil {
		runtime.LogErrorf(l.Ctx, "failed to find latest version")
		return UpdateAvailabilityUnknown, fmt.Errorf("failed to find latest version")
	}

	if latestVersion.GreaterThan(currentVersion) {
		l.UpdateAvailability = UpdateAvailabilityAvailable
	} else {
		l.UpdateAvailability = UpdateAvailabilityUpToDate
	}

	l.EmitEvent(events.LauncherUpdateAvailable, l.UpdateAvailability)

	return l.UpdateAvailability, nil
}

// IndexLauncherApps requests the apps from the launcher metadata from the backend
func (l *Launcher) IndexLauncherApps(offset int64, limit int64) ([]sm.AppV2, error) {
	runtime.LogInfof(l.Ctx, "IndexLauncherApps")

	if config.LauncherId == "" {
		runtime.LogErrorf(l.Ctx, "launcher id is not set")
		return nil, fmt.Errorf("launcher id is not set")
	}
	launcherId := uuid.FromStringOrNil(config.LauncherId)

	apps, err := api.IndexLauncherApps(l.Ctx, launcherId, offset, limit)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to index launcher apps: %v", err)
		return nil, fmt.Errorf("failed to index launcher apps: %w", err)
	}

	l.EmitEvent(events.LauncherApps, apps)

	return apps, nil
}

// GetAppMetadata requests the app metadata from the backend
func (l *Launcher) GetAppMetadata(id uuid.UUID) (*sm.AppV2, error) {
	app, err := api.GetAppMetadata(l.Ctx, id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app: %v", err)
		return nil, fmt.Errorf("failed to get app: %w", err)
	}

	l.EmitEvent(events.LauncherApp, app)

	return app, nil
}

// getAppExecutableById returns the executable path for the given app id located in the given directory
func (l *Launcher) getAppExecutableById(dir string, id uuid.UUID) (string, error) {
	path := filepath.Join(dir, id.String())
	if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "windows" {
		path = path + ".exe"
	}
	appExe, err := filepath.Abs(path)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get absolute path: %v", err)
		return "", err
	}
	fi, err := os.Stat(appExe)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get file info: %v", err)
		return "", err
	}

	if fi.IsDir() {
		runtime.LogErrorf(l.Ctx, "app executable is a directory")
		return "", fmt.Errorf("app executable is a directory")
	}

	return appExe, nil
}

// getAppExecutable returns the executable path for the given app located in the given directory
func (l *Launcher) getAppExecutableByName(dir string, name string) (string, error) {
	path := filepath.Join(dir, name)
	if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "windows" {
		path = path + ".exe"
	}
	appExe, err := filepath.Abs(path)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get absolute path: %v", err)
		return "", err
	}
	fi, err := os.Stat(appExe)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get file info: %v", err)
		return "", err
	}

	if fi.IsDir() {
		runtime.LogErrorf(l.Ctx, "app executable is a directory")
		return "", fmt.Errorf("app executable is a directory")
	}

	return appExe, nil
}

// getAppExecutable returns the executable path for the given app located in the given directory
func (l *Launcher) getAppExecutableByGenericName(dir string) (string, error) {
	path := filepath.Join(dir, "Metaverse")
	if //goland:noinspection GoBoolExpressions
	goRuntime.GOOS == "windows" {
		path = path + ".exe"
	}
	appExe, err := filepath.Abs(path)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get absolute path: %v", err)
		return "", err
	}
	fi, err := os.Stat(appExe)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get file info: %v", err)
		return "", err
	}

	if fi.IsDir() {
		runtime.LogErrorf(l.Ctx, "app executable is a directory")
		return "", fmt.Errorf("app executable is a directory")
	}

	return appExe, nil
}

// getAppExecutable returns the executable path for the given app located in the given directory
func (l *Launcher) getAppExecutable(id uuid.UUID, name string) (string, error) {
	var err error

	if l.Metadata == nil {
		_, err = l.GetLauncherMetadata()
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to get launcher metadata: %v", err)
			return "", err
		}
	}

	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %v", err)
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	applicationsDir := filepath.Join(executableDir, ApplicationsDir, id.String())

	_, err = os.Stat(applicationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			runtime.LogWarningf(l.Ctx, "applications directory does not exist: %s", applicationsDir)
			return "", err
		}
		runtime.LogErrorf(l.Ctx, "failed to stat applications directory: %s", applicationsDir)
		return "", fmt.Errorf("failed to stat app directory: %w", err)
	}

	var appPath string
	if !id.IsNil() {
		appPath, err = l.getAppExecutableById(applicationsDir, id)
		if err == nil {
			return appPath, nil
		} else {
			runtime.LogWarningf(l.Ctx, "failed to get app executable by id: %v", err)
		}
	}

	if name != "" {
		appPath, err = l.getAppExecutableByName(applicationsDir, name)
		if err == nil {
			return appPath, nil
		} else {
			runtime.LogWarningf(l.Ctx, "failed to get app executable by name: %v", err)
		}
	}

	appPath, err = l.getAppExecutableByGenericName(applicationsDir)
	if err == nil {
		return appPath, nil
	} else {
		runtime.LogWarningf(l.Ctx, "failed to get app executable by generic name: %v", err)
	}

	err = filepath.WalkDir(applicationsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to walk directory: %v", err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to get file info: %v", err)
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to open file: %v", err)
			return err
		}

		defer func() {
			err = f.Close()
			if err != nil {
				runtime.LogErrorf(l.Ctx, "failed to close file: %s", err)
			}
		}()

		var isExecutable bool
		isExecutable, err = executable.IsExecutable(f)
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to check if file is executable: %v", err)
			return err
		}

		if isExecutable {
			appPath = path
			return io.EOF
		}

		appPath = path

		return nil
	})

	if err != nil && err != io.EOF {
		runtime.LogErrorf(l.Ctx, "failed to walk directory: %v", err)
		return "", fmt.Errorf("failed to find app executable: %w", err)
	}

	if appPath == "" {
		runtime.LogErrorf(l.Ctx, "failed to find app executable")
		return "", fmt.Errorf("failed to find app executable")
	}

	appPath, err = filepath.Abs(appPath)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get absolute path: %v", err)
		return "", fmt.Errorf("failed to find app executable: %w", err)
	}

	return appPath, nil
}

// IsAppInstalled returns true if the app is installed and is ready to launch
func (l *Launcher) IsAppInstalled(id uuid.UUID) (bool, error) {
	app, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %v", err)
		return false, fmt.Errorf("failed to get app metadata: %w", err)
	}

	runtime.LogWarningf(l.Ctx, "check if app is installed: %+v", app)

	appExe, err := l.getAppExecutable(id, app.Name)
	if err != nil && !os.IsNotExist(err) {
		runtime.LogErrorf(l.Ctx, "failed to get app executable: %v", err)
		return false, fmt.Errorf("failed to get app executable: %w", err)
	}

	if appExe == "" {
		runtime.LogWarningf(l.Ctx, "app is not installed")
		return false, os.ErrNotExist
	}

	return true, nil
}

// CheckForAppUpdates checks for application updates
func (l *Launcher) CheckForAppUpdates(id uuid.UUID) (UpdateAvailability, error) {
	dir, err := l.getAppInstallationDir(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app installation dir: %v", err)
		return UpdateAvailabilityUnknown, fmt.Errorf("failed to get current dir: %v", err)
	}

	currentVersion, err := version.ReadVersion(dir)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get current app version: %w", err)
		return UpdateAvailabilityUnknown, err
	}

	runtime.LogInfof(l.Ctx, "current app version: %s", currentVersion)

	appMetadata, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %w", err)
		return UpdateAvailabilityUnknown, err
	}

	var latestVersion *semver.Version
	for _, release := range appMetadata.Releases.Entities {
		if latestVersion == nil {
			latestVersion, err = semver.NewVersion(release.Version)
			if err != nil {
				runtime.LogErrorf(l.Ctx, "failed to parse version: %w", err)
				return UpdateAvailabilityUnknown, err
			}
		} else {
			releaseVersion, err := semver.NewVersion(release.Version)
			if err != nil {
				runtime.LogErrorf(l.Ctx, "failed to parse semver: %w", err)
				return UpdateAvailabilityUnknown, fmt.Errorf("failed to parse semver: %w", err)
			}

			if latestVersion.GreaterThan(releaseVersion) {
				latestVersion = releaseVersion
			}
		}
	}

	if latestVersion == nil {
		runtime.LogErrorf(l.Ctx, "failed to find latest version")
		return UpdateAvailabilityUnknown, fmt.Errorf("failed to find latest version")
	}

	if latestVersion.GreaterThan(currentVersion) {
		l.UpdateAvailability = UpdateAvailabilityAvailable
	} else {
		l.UpdateAvailability = UpdateAvailabilityUpToDate
	}

	l.EmitEvent(events.AppUpdateAvailable, id, l.UpdateAvailability)

	return l.UpdateAvailability, nil
}

// InstallApp installs the app with the given id
func (l *Launcher) InstallApp(id uuid.UUID) error {
	if l.IsUpdatingApp {
		return ErrorAppIsUpdating
	}

	app, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %v", err)
		return err
	}

	if app == nil {
		runtime.LogErrorf(l.Ctx, "app %s not found", id)
		return fmt.Errorf("app %s not found", id)
	}

	if len(app.Releases.Entities) == 0 {
		runtime.LogErrorf(l.Ctx, "no releases found for app %s", app.Id)
		return fmt.Errorf("no releases found for app %s", app.Id)
	}

	runtime.LogInfof(l.Ctx, "installing app %s", app.Id)

	installed, err := l.IsAppInstalled(id)
	if err != nil && !os.IsNotExist(err) {
		runtime.LogErrorf(l.Ctx, "failed to check if app is installed: %s", err)
		return fmt.Errorf("failed to check if app is installed: %w", err)
	}

	if installed {
		runtime.LogWarningf(l.Ctx, "app is already installed")
		return ErrorAppInstalled
	} else {
		runtime.LogWarningf(l.Ctx, "app is not installed")
	}

	l.IsUpdatingApp = true

	release := app.Releases.Entities[0]
	if release.Archive {
		return l.installAppReleaseArchive(*app, release)
	} else {
		return l.installAppRelease(*app, release)
	}
}

// LaunchApp launches the app with the given id
func (l *Launcher) LaunchApp(id uuid.UUID) error {
	app, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %v", err)
		return fmt.Errorf("failed to get app metadata: %w", err)
	}

	appExe, err := l.getAppExecutable(id, app.Name)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app executable: %v", err)
		return fmt.Errorf("failed to get app executable: %w", err)
	}

	cmd := exec.Command(appExe)
	cmd.Dir = filepath.Dir(appExe)
	err = cmd.Start()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to start app: %v", err)
		return fmt.Errorf("failed to start app: %w", err)
	}

	return nil
}

func (l *Launcher) UpdateApp(id uuid.UUID) error {
	if l.IsUpdatingApp {
		return ErrorAppIsUpdating
	}

	app, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %v", err)
		return err
	}

	if app == nil {
		runtime.LogErrorf(l.Ctx, "app %s not found", id)
		return fmt.Errorf("app %s not found", id)
	}

	if len(app.Releases.Entities) == 0 {
		runtime.LogErrorf(l.Ctx, "no releases found for app %s", app.Id)
		return fmt.Errorf("no releases found for app %s", app.Id)
	}

	runtime.LogWarningf(l.Ctx, "updating app %s", app.Id)

	l.IsUpdatingApp = true

	release := app.Releases.Entities[0]
	if release.Archive {
		return l.installAppReleaseArchive(*app, release)
	} else {
		return l.installAppRelease(*app, release)
	}
}

func (l *Launcher) DeleteApp(id uuid.UUID) error {
	app, err := l.GetAppMetadata(id)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get app metadata: %v", err)
		return err
	}

	runtime.LogWarningf(l.Ctx, "deleting app %s", app.Id.String())

	installed, err := l.IsAppInstalled(id)
	if err != nil && !os.IsNotExist(err) {
		runtime.LogErrorf(l.Ctx, "failed to check if app is installed: %s", err)
		return fmt.Errorf("failed to check if app is installed: %w", err)
	}

	if !installed {
		runtime.LogWarningf(l.Ctx, "app is not installed")
		return fmt.Errorf("app is not installed")
	}

	appInstallationDir, err := l.getAppInstallationDir(id)
	if err != nil {
		runtime.LogWarningf(l.Ctx, "failed to get app installation directory: %s", err)
		return fmt.Errorf("failed to get app installation directory: %w", err)
	}

	err = os.RemoveAll(appInstallationDir)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to remove app directory: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateCompleted, app, "failed to remove app directory")
		return fmt.Errorf("failed to remove app directory: %w", err)
	}

	return nil
}

// EmitEvent emits an event to the frontend and stores the last event
func (l *Launcher) EmitEvent(event string, args ...any) {
	l.LastEvent = event
	runtime.EventsEmit(l.Ctx, event, args...)
}

// GetIsUpdatingApp returns if the launcher is currently updating any app
func (l *Launcher) GetIsUpdatingApp() bool {
	return l.IsUpdatingApp
}

// SetAppUpdateStatus sets the update status of the app
func (l *Launcher) SetAppUpdateStatus(isUpdating bool, event string, args ...any) {
	l.IsUpdatingApp = isUpdating
	l.EmitEvent(event, args...)
}

// GetLastEvent returns the last event that was emitted
func (l *Launcher) GetLastEvent() (string, error) {
	return l.LastEvent, nil
}

// OpenUrl opens the given url in the default browser
func (l *Launcher) OpenUrl(url string) error {
	if url == "" {
		runtime.LogErrorf(l.Ctx, "url is empty")
		return fmt.Errorf("url is empty")
	} else {
		err := sp.OpenUrl(url)
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to open url: %s", err)
			return fmt.Errorf("failed to open url: %w", err)
		}
	}

	return nil
}

// getAppInstallationDir returns the installation directory of the app with the given id
func (l *Launcher) getAppInstallationDir(id uuid.UUID) (string, error) {
	runtime.LogDebugf(l.Ctx, "getting working directory...")

	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %s", err)
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	runtime.LogDebugf(l.Ctx, "working directory: %s", executableDir)
	return filepath.Join(executableDir, ApplicationsDir, id.String()), nil
}

func (l *Launcher) installAppReleaseArchive(app sm.AppV2, release sm.ReleaseV2) error {
	runtime.LogDebugf(l.Ctx, "installing app release archive: %+v", release)

	id := app.Id

	runtime.LogDebugf(l.Ctx, "getting archive file...")
	var archive *sm.File
	for _, file := range release.Files.Entities {
		if file.Type == "release-archive" {
			runtime.LogDebugf(l.Ctx, "found archive file %s: %s", file.Id, file.Url)
			archive = &file
			break
		}
	}
	if archive == nil {
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "no archive file found")
		return fmt.Errorf("no archive file found")
	}
	runtime.LogDebugf(l.Ctx, "archive file found: %+v", archive)

	runtime.LogDebugf(l.Ctx, "getting working directory...")
	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "failed to get executable path")
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)
	runtime.LogDebugf(l.Ctx, "working directory: %s", executableDir)

	tempDownloadPath := filepath.Join(executableDir, ".tmp", id.String())
	runtime.LogDebugf(l.Ctx, "temp download path: %s", tempDownloadPath)
	appInstallationPath := filepath.Join(executableDir, ApplicationsDir, id.String())
	runtime.LogDebugf(l.Ctx, "app installation path: %s", appInstallationPath)

	counter := http.NewDownloadProgressTracker((uint64)(*archive.Size), func(progress uint64, total uint64) {
		l.EmitEvent(events.AppUpdateProgress, app, progress, total)
	})
	runtime.LogDebugf(l.Ctx, "downloading file to %s...", tempDownloadPath)
	err = http.DownloadFile(l.Ctx, tempDownloadPath, archive.Url, counter)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to download file: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "failed to download file")
		return fmt.Errorf("failed to download file: %w", err)
	}
	runtime.LogDebugf(l.Ctx, "downloaded file to %s", tempDownloadPath)

	l.SetAppUpdateStatus(false, events.AppUpdateExtracting, app)

	runtime.LogDebugf(l.Ctx, "extracting archive to %s...", appInstallationPath)
	err = utils.ExtractArchive(l.Ctx, tempDownloadPath, appInstallationPath)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to extract archive: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "failed to extract archive")
		return fmt.Errorf("failed to extract archive: %w", err)
	}
	runtime.LogDebugf(l.Ctx, "extracted archive to %s", appInstallationPath)

	runtime.LogDebugf(l.Ctx, "parsing release version: %s...", release.Version)
	v, err := semver.NewVersion(release.Version)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to parse release version: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "failed to parse release version")
		return fmt.Errorf("failed to parse release version: %w", err)
	}
	runtime.LogDebugf(l.Ctx, "parsed release version: %s", v.String())

	runtime.LogDebugf(l.Ctx, "writing version to %s...", appInstallationPath)
	err = version.WriteVersion(appInstallationPath, v)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to write version: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateFailed, app, "failed to write version")
		return fmt.Errorf("failed to write version: %w", err)
	}
	runtime.LogDebugf(l.Ctx, "wrote version to %s", appInstallationPath)

	runtime.LogDebugf(l.Ctx, "removing temporary download directory %s...", tempDownloadPath)
	err = os.RemoveAll(tempDownloadPath)
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to remove temporary download directory: %s", err)
		l.SetAppUpdateStatus(false, events.AppUpdateCompleted, app, "failed to remove temporary download directory")
		return fmt.Errorf("failed to remove temporary download directory: %w", err)
	}
	runtime.LogDebugf(l.Ctx, "removed temporary download directory %s", tempDownloadPath)

	l.SetAppUpdateStatus(false, events.AppUpdateCompleted, app)

	return nil
}

func (l *Launcher) installAppRelease(app sm.AppV2, release sm.ReleaseV2) error {
	runtime.LogDebugf(l.Ctx, "installing app release: %+v", release)

	id := app.Id

	var files []*sm.File
	for _, file := range release.Files.Entities {
		if file.Type == "release" {
			runtime.LogDebugf(l.Ctx, "found release file %s: %s", file.Id, file.Url)
			files = append(files, &file)
		}
	}

	if len(files) == 0 {
		runtime.LogErrorf(l.Ctx, "no release files found")
		return fmt.Errorf("no release files found")
	}

	// Get the path to the executable
	executablePath, err := os.Executable()
	if err != nil {
		runtime.LogErrorf(l.Ctx, "failed to get executable path: %s", err)
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	executableDir := filepath.Dir(executablePath)

	tempDownloadPath := filepath.Join(executableDir, ".tmp", id.String())
	appInstallationPath := filepath.Join(executableDir, ApplicationsDir, id.String())

	var totalProgress uint64 = 0
	var totalSize uint64 = 0

	// calculate total size for all files
	for _, file := range files {
		totalSize += uint64(*file.Size)
	}

	runtime.LogDebugf(l.Ctx, "total size: %d", totalSize)

	for _, file := range files {
		counter := http.NewDownloadProgressTracker(totalSize, func(progress uint64, total uint64) {
			// accumulate progress for all files and report it to the frontend as total progress
			totalProgress += progress
			l.EmitEvent(events.AppUpdateProgress, totalProgress, total)
		})
		// download next file
		err = http.DownloadFile(l.Ctx, tempDownloadPath, file.Url, counter)
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to download file: %s", err.Error())
		}
	}

	for _, file := range files {
		if file.OriginalPath == nil {
			runtime.LogErrorf(l.Ctx, "file %s has no original path", file.Id)
			continue
		}
		err = os.Rename(filepath.Join(tempDownloadPath, *file.OriginalPath), filepath.Join(appInstallationPath, *file.OriginalPath))
		if err != nil {
			runtime.LogErrorf(l.Ctx, "failed to move file: %s", err.Error())
		}
	}

	v, err := semver.NewVersion(release.Version)
	if err != nil {
		return fmt.Errorf("failed to parse release version: %w", err)
	}

	err = version.WriteVersion(appInstallationPath, v)
	if err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	err = os.RemoveAll(tempDownloadPath)
	if err != nil {
		return fmt.Errorf("failed to remove temporary download directory: %w", err)
	}

	return nil
}

// getDownloadDir returns the temporary download directory to store downloaded files
func getDownloadDir() (string, error) {
	cd, err := sp.GetCurrentDir()
	if err != nil {
		runtime.LogErrorf(nil, "failed to get current dir: %s", err.Error())
		return "", fmt.Errorf("failed to get current dir: %w", err)
	}

	return filepath.Join(cd, ".tmp"), nil
}
