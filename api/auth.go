package api

import (
	"bytes"
	"context"
	vContext "dev.hackerman.me/artheon/veverse-shared/context"
	"dev.hackerman.me/artheon/veverse-shared/model"
	"encoding/json"
	"fmt"
	glConfig "games.launch.launcher/config"
	glSession "games.launch.launcher/session"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
)

func Login(ctx context.Context, project string, email string, password string) (context.Context, error) {
	var (
		requestBody []byte
		err         error
	)

	requestBody, err = json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	if err != nil {
		runtime.LogErrorf(ctx, "failed to marshal login request body: %s", err.Error())
		return nil, err
	}

	url := fmt.Sprintf("%s/auth/login", glConfig.Api2Url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		runtime.LogErrorf(ctx, "failed to create login request: %s", err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to login to %s: %s", url, err.Error())
		return ctx, err
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			runtime.LogErrorf(ctx, "failed to close login response body: %s", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		runtime.LogErrorf(ctx, "failed to login to %s, status code: %d, error: %s", url, resp.StatusCode, err.Error())
		return ctx, fmt.Errorf("failed to login to %s, status code: %d, error: %s", url, resp.StatusCode, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		runtime.LogErrorf(ctx, "failed to read login response body: %s", err.Error())
		return ctx, err
	}

	var v model.Wrapper[string]
	if err = json.Unmarshal(body, &v); err != nil {
		runtime.LogErrorf(ctx, "failed to unmarshal login response body: %s", err.Error())
		return ctx, err
	}

	if v.Status == "error" {
		runtime.LogErrorf(ctx, "authentication error %d: %s\n", resp.StatusCode, v.Message)
		return ctx, fmt.Errorf("authentication error %d: %s\n", resp.StatusCode, v.Message)
	} else if v.Status == "ok" {
		err1 := glSession.SaveSession(project, v.Payload)
		if err1 != nil {
			runtime.LogErrorf(ctx, "failed to save session: %v", err1)
			return context.WithValue(ctx, vContext.Token, v.Payload), fmt.Errorf("failed to save session: %v", err1)
		}

		return context.WithValue(ctx, vContext.Token, v.Payload), nil
	}

	runtime.LogErrorf(ctx, "unknown authentication error: %s", v.Message)
	return ctx, fmt.Errorf(v.Message)
}
