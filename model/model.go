// Package model contains the data models used by the launcher (some of them are deprecated in favor of the shared package).
package model

import (
	"dev.hackerman.me/artheon/veverse-shared/model"
	"encoding/gob"
	"fmt"
	ll "games.launch.launcher/logger"
	"github.com/google/uuid"
	"os"
	"time"
)

type Status struct {
	Downloading     bool    `json:"downloading"`     // is download in Progress, used to prevent multiple downloads, hides the launch and download buttons
	Progress        float64 `json:"progress"`        // download Progress in percent, used for the Progress bar
	UpdateAvailable bool    `json:"updateAvailable"` // is update available, used to show the update button and hide the launch button
	ShowButtons     bool    `json:"showButtons"`     // is launch or download buttons visible
	NextVersion     string  `json:"nextVersion"`     // the next version to update to
	Message         string  `json:"message"`         // status message, used to show the user what is happening
	ShowLoader      bool    `json:"showLoader"`      // is the status bar loader visible
}

type FileHeaders struct {
	Id      uuid.UUID `json:"id,omitempty"`
	Url     string    `json:"url,omitempty"`
	Size    int64     `json:"size,omitempty"`
	ETag    string    `json:"eTag,omitempty"`
	ModTime time.Time `json:"modTime,omitempty"`
}

func (r *FileHeaders) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to create file: %v\n", err))
		return fmt.Errorf("failed to create file: %v", err)
	}

	defer func(f *os.File) {
		err1 := f.Close()
		if err1 != nil {
			ll.Logger.Error(fmt.Sprintf("failed to close file: %v\n", err1))
		}
	}(f)

	enc := gob.NewEncoder(f)
	err = enc.Encode(&r)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to encode file: %v\n", err))
		return fmt.Errorf("failed to encode file: %v", err)
	}

	return nil
}

func (r *FileHeaders) LoadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to open file: %v\n", err))
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer func(f *os.File) {
		err1 := f.Close()
		if err1 != nil {
			ll.Logger.Error(fmt.Sprintf("failed to close file: %v\n", err1))
		}
	}(f)

	dec := gob.NewDecoder(f)
	err = dec.Decode(&r)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decode file: %v\n", err))
		return fmt.Errorf("failed to decode file: %v", err)
	}

	return nil
}

type FileWork struct {
	model.File
	*FileHeaders
}
