package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func setLogOutputToFile(logPath string) {
	var err error
	var f *os.File
	f, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)

	log.SetOutput(f)
}

// Self-updater for the launcher.
func main() {
	var err error

	/** Arguments:
	 * 1. Path to the new launcher executable in the temp directory.
	 * 2. Path to the old launcher executable in the installation directory.
	 * 3. Enable logging (optional, enabled if argument is not empty).
	 */
	src := os.Args[1]
	dst := os.Args[2]
	logEnabled := os.Args[3]

	//region Logging
	setLogOutputToFile("updater.log")
	//endregion

	if logEnabled == "true" {
		log.Printf("removing %q", filepath.ToSlash(dst))
	}
	_, err = os.Stat(dst)
	if err == nil {
		var success bool

		// Retry several times in case the launcher is still running.
		for i := 0; i < 42; i++ {
			if logEnabled == "true" {
				log.Printf("attempt %d", i+1)
			}

			err := os.Remove(dst)
			if err == nil {
				if logEnabled == "true" {
					log.Printf("removed %q", filepath.ToSlash(dst))
				}
				success = true
				break
			} else {
				if logEnabled == "true" {
					log.Printf("failed to remove %q: %v", filepath.ToSlash(dst), err)
				}
			}

			time.Sleep(1 * time.Second)
		}

		if !success {
			log.Fatalf("failed to remove file %s: %v\n", filepath.ToSlash(dst), err)
		}
	} else if !os.IsNotExist(err) {
		log.Fatalf("failed to check if file exists: %v", err)
	}
	if logEnabled == "true" {
		log.Printf("renaming %q to %q", filepath.ToSlash(src), filepath.ToSlash(dst))
	}
	_, err = os.Stat(src)
	if err == nil {
		err := os.Rename(src, dst)
		if err != nil {
			log.Fatalf("failed to rename file %s to %s: %v", filepath.ToSlash(src), filepath.ToSlash(dst), err)
		}
	} else {
		if os.IsNotExist(err) {
			if logEnabled == "true" {
				log.Printf("file %q does not exist", filepath.ToSlash(src))
			}
		} else {
			log.Fatalf("failed to check if file exists: %v", err)
		}
	}

	path, err := filepath.Abs(dst)
	if err != nil {
		log.Fatalf("failed to get absolute path of %s: %v", filepath.ToSlash(dst), err)
	}

	// Start the new launcher.
	if logEnabled == "true" {
		log.Printf("starting %q", filepath.ToSlash(path))
	}
	cmd := exec.Command(path)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("failed to start process %s: %v", filepath.ToSlash(path), err)
	}
}
