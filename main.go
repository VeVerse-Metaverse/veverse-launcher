package main

import (
	"embed"
	"fmt"
	"games.launch.launcher/app"
	"games.launch.launcher/config"
	ll "games.launch.launcher/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"net"
	"os"
	"path/filepath"
	goRuntime "runtime"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create a new logger instance.
	if config.Logging == "true" {
		// Get the path to the executable
		executablePath, err := os.Executable()
		if err != nil {
			log.Printf("Failed to get executable path: %v\n", err)
		}

		// Get the directory containing the executable
		executableDir := filepath.Dir(executablePath)

		// Create a log file in the same directory as the executable
		logFilePath := filepath.Join(executableDir, "launcher.log")

		// Try to create the log file to check permissions.
		f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)

		// If the file already exists, use it.
		if err != nil && os.IsExist(err) {
			ll.Logger = logger.NewFileLogger(logFilePath)
		} else {
			if err == nil {
				// Close the file if it was created so that the logger can open it.
				err := f.Close()
				if err != nil {
					log.Printf("Failed to close log file: %v\n", err)
				}
				ll.Logger = logger.NewFileLogger(logFilePath)
			} else {
				// Print error message if the file couldn't be created.
				log.Printf("Failed to create log file: %v\n", err)

				if goRuntime.GOOS == "windows" {
					// If the file doesn't exist, try to create the "LE7EL/logs" directory in the user's AppData directory.
					appData := os.Getenv("APPDATA")
					logsDir := filepath.Join(appData, "LE7EL", "logs")
					err = os.MkdirAll(logsDir, os.ModePerm)
					if err != nil {
						log.Printf("Failed to create logs directory: %v\n", err)
					} else {
						logFilePath = filepath.Join(logsDir, "launcher.log")

						// Try to crate or open the log file in the new location and check for write permissions.
						f, err = os.OpenFile(logFilePath, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0666)
						if err != nil && !os.IsExist(err) {
							log.Printf("Failed to create log file: %v, using default logger\n", err)
							ll.Logger = logger.NewDefaultLogger()
						} else {
							err = f.Close()
							if err != nil {
								log.Printf("Failed to close log file: %v\n", err)
								ll.Logger = logger.NewDefaultLogger()
							} else {
								// Use the file logger if the launcher has been built with file logging enabled and the log file was created successfully.
								ll.Logger = logger.NewFileLogger(logFilePath)
							}
						}
					}
				}
				// todo: macos and linux support
			}
		}
	} else {
		ll.Logger = logger.NewDefaultLogger()
	}

	var err error

	// Attempt to connect to the single instance port, failure means this is the first instance.
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+config.LauncherPort, time.Second)
	if err != nil {
		ll.Logger.Print(fmt.Sprintf("No other instance found, starting the main instance: %v\n", err))

		// Create a new launcher application instance.
		launcher := app.NewLauncher()

		//region Main application logic goes here

		// Run the launcher application.
		err = wails.Run(&options.App{
			Title:     "Launcher",
			Width:     1280,
			Height:    720,
			MinWidth:  800,
			MinHeight: 600,
			MaxWidth:  1280,
			MaxHeight: 720,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 15, G: 15, B: 15, A: 1},
			OnStartup:        launcher.OnStartup,
			Bind: []interface{}{
				launcher,
			},
			Frameless:  false,
			Fullscreen: false,
			Logger:     ll.Logger,
		})

		//endregion

		// Handle any errors that occurred during the application run.
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error running application: %v\n", err))
		}
	} else {
		// Start the subsequent instance and send the deep link to the first instance then exit.
		startSubsequentInstance(conn)
	}
}

// Start a subsequent instance of the application to send the deep link to the main instance.
func startSubsequentInstance(conn net.Conn) {
	// Check if a deep link was passed to the application.
	var deepLink string
	if len(os.Args) > 1 {
		deepLink = os.Args[1]
	}

	if deepLink == "" {
		ll.Logger.Print("No deep link found, exiting subsequent instance\n")
		os.Exit(0)
	}

	ll.Logger.Print(fmt.Sprintf("Starting subsequent instance with args: %v\n", deepLink))

	// Send the deep link to the first instance
	_, err := conn.Write([]byte(deepLink))
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("Error sending deep link: %v\n", err))
		return
	}

	err = conn.Close()
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("Error closing connection: %v\n", err))
		return
	}

	// Exit the subsequent instance after sending the deep link.
	os.Exit(0)
}
