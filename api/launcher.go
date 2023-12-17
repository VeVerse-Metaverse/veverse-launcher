// Package api contains all the API calls to the LE7EL XR API.
package api

import (
	"context"
	sm "dev.hackerman.me/artheon/veverse-shared/model"
	vUnreal "dev.hackerman.me/artheon/veverse-shared/unreal"
	"encoding/json"
	"fmt"
	"games.launch.launcher/config"
	"github.com/gofrs/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
)

// GetLauncherMetadata returns the launcher metadata for the given launcher id.
func GetLauncherMetadata(ctx context.Context, id uuid.UUID) (*sm.LauncherV2, error) {
	if id.IsNil() {
		if config.LauncherId == "" {
			runtime.LogErrorf(ctx, "launcher id is not set")
			return nil, fmt.Errorf("launcher id is not set")
		} else {
			id = uuid.FromStringOrNil(config.LauncherId)
		}
	}

	url := fmt.Sprintf("%s/launchers/public/%s?platform=%s", config.Api2Url, id, vUnreal.GetPlatformName())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create a HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create a HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to send a HTTP GET request: %v", err)
		return nil, fmt.Errorf("failed to send a HTTP GET request: %w", err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			runtime.LogErrorf(ctx, "error closing http response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			runtime.LogErrorf(ctx, "failed to read response body: %v", err)
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		runtime.LogErrorf(ctx, "failed to get launcher metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
		return nil, fmt.Errorf("failed to get launcher metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var v sm.Wrapper[sm.LauncherV2]
	if err = json.Unmarshal(body, &v); err != nil {
		runtime.LogErrorf(ctx, "failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &v.Payload, err
}

func RequestLauncherReleaseMetadata(ctx context.Context, offset int64, limit int64) ([]sm.ReleaseV2, error) {
	if config.LauncherId == "" {
		runtime.LogErrorf(ctx, "launcher id is not set")
		return nil, fmt.Errorf("launcher id is not set")
	}

	url := fmt.Sprintf("%s/launchers/public/%s/releases?platform=%s&offset=%d&limit=%d", config.Api2Url, config.LauncherId, vUnreal.GetPlatformName(), offset, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create a HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create a HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to send a HTTP GET request: %v", err)
		return nil, fmt.Errorf("failed to send a HTTP GET request: %w", err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			runtime.LogErrorf(ctx, "error closing http response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			runtime.LogErrorf(ctx, "failed to read response body: %v", err)
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		runtime.LogErrorf(ctx, "failed to get launcher releases from %s, status code: %d, content: %s", url, resp.StatusCode, body)
		return nil, fmt.Errorf("failed to get launcher releases from %s, status code: %d, content: %s", url, resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var v sm.Wrapper[[]sm.ReleaseV2]
	if err = json.Unmarshal(body, &v); err != nil {
		runtime.LogErrorf(ctx, "failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return v.Payload, err
}

// IndexLauncherApps returns a list of apps for the given launcher id.
func IndexLauncherApps(ctx context.Context, id uuid.UUID, offset int64, limit int64) ([]sm.AppV2, error) {
	if id.IsNil() {
		if config.LauncherId == "" {
			runtime.LogErrorf(ctx, "launcher id is not set")
			return nil, fmt.Errorf("launcher id is not set")
		} else {
			id = uuid.FromStringOrNil(config.LauncherId)
		}
	}

	url := fmt.Sprintf("%s/launchers/public/%s/apps?platform=%s&offset=%d&limit=%d", config.Api2Url, id, vUnreal.GetPlatformName(), offset, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create a HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create a HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to send a HTTP GET request: %v", err)
		return nil, fmt.Errorf("failed to send a HTTP GET request: %w", err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			runtime.LogErrorf(ctx, "error closing http response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			runtime.LogErrorf(ctx, "failed to read response body: %v", err)
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		runtime.LogErrorf(ctx, "failed to get launcher apps metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
		return nil, fmt.Errorf("failed to get launcher apps metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var v sm.Wrapper[[]sm.AppV2]
	if err = json.Unmarshal(body, &v); err != nil {
		runtime.LogErrorf(ctx, "failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return v.Payload, err
}

// GetAppMetadata returns the app metadata for the given app id.
func GetAppMetadata(ctx context.Context, id uuid.UUID) (*sm.AppV2, error) {
	if id.IsNil() {
		runtime.LogErrorf(ctx, "app id is not set")
		return nil, fmt.Errorf("app id is not set")
	}

	url := fmt.Sprintf("%s/apps/public/%s?platform=%s", config.Api2Url, id, vUnreal.GetPlatformName())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create a HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create a HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to send a HTTP GET request: %v", err)
		return nil, fmt.Errorf("failed to send a HTTP GET request: %w", err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			runtime.LogErrorf(ctx, "error closing http response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			runtime.LogErrorf(ctx, "failed to read response body: %v", err)
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		runtime.LogErrorf(ctx, "failed to get app metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
		return nil, fmt.Errorf("failed to get launcher metadata from %s, status code: %d, content: %s", url, resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var v sm.Wrapper[sm.AppV2]
	if err = json.Unmarshal(body, &v); err != nil {
		runtime.LogErrorf(ctx, "failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &v.Payload, err
}
