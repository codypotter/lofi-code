package hcaptcha

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type HCaptcha interface {
	VerifyHCaptcha(remoteip, token string) (bool, error)
}

func New(env string, secret string) HCaptcha {
	if env == "development" {
		return &mockHCaptcha{}
	}
	return &hCaptcha{
		hcaptchaSecret: secret,
	}
}

type hCaptcha struct {
	hcaptchaSecret string
}

type siteverifyResponse struct {
	Success bool `json:"success"`
}

func (h *hCaptcha) VerifyHCaptcha(remoteip, token string) (bool, error) {
	resp, err := http.PostForm("https://hcaptcha.com/siteverify", url.Values{
		"secret":   {h.hcaptchaSecret},
		"response": {token},
		"remoteip": {remoteip},
	})
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result siteverifyResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}
	return result.Success, nil
}

type mockHCaptcha struct{}

func (m *mockHCaptcha) VerifyHCaptcha(remoteip, token string) (bool, error) {
	log.Info().Str("remoteip", remoteip).Str("token", token).Msg("Mock hCaptcha verification")
	return true, nil
}
