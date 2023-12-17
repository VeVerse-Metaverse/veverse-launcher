//go:build debug

package app

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"time"
)

// debugPrintMessageToFile prints a string message to a file
func (l *Launcher) debugPrintMessageToFile(message string, useRandomFileName bool) {
	var filename = "debug.log"

	if useRandomFileName {
		fileid, err := uuid.NewV4()
		if err != nil {
			logrus.Errorf("failed to generate file id: %v", err)
			return
		}
		filename = fmt.Sprintf("debug-%s.log", fileid)
	}

	// open file and write error
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Errorf("failed to open error log file: %v", message)
		return
	}
	defer f.Close()

	runtime.LogErrorf(l.Ctx, "%v", f)

	f.WriteString(fmt.Sprintf("%s: %s", time.Now().Format(time.RFC3339), message))

	options := runtime.MessageDialogOptions{
		Type:    "error",
		Title:   "error",
		Message: fmt.Sprintf("%s", message),
		Buttons: []string{"ok"},
	}
	_, err = runtime.MessageDialog(l.Ctx, options)
	if err != nil {
		return
	}
}
