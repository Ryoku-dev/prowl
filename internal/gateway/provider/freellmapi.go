package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Endpoint and wire deviations follow FreeLLMAPI's server/src/providers adapters.
func registerFreeLLMAPIEndpoints(r *Registry) {
	for _, cfg := range []providerCfg{
		{platform: "beatapi", name: "BeatAPI", baseURL: "https://api.beatapi.io/v1", noChat: true},
		{platform: "blaze", name: "BlazeAPI", baseURL: "https://api.blazeapi.org/paid/v1", validateURL: "https://api.blazeapi.org/paid/v1/usage"},
		{platform: "blockrun", name: "BlockRun", baseURL: "https://api.blockrun.ai/v1", validateURL: "https://api.blockrun.ai/v1/account"},
		{platform: "dreamprompting", name: "DreamPrompting", baseURL: "https://dreamprompting.com/api/v1", noStream: true, checkResponse: dreamPromptingCheck},
		{platform: "inferbase", name: "Inferbase", baseURL: "https://api.inferbase.ai/v1", validateURL: "https://api.inferbase.ai/api/v1/inference/credits", noStream: true},
		{platform: "logfare", name: "Logfare", baseURL: "https://logfare.ai/v1"},
		{platform: "parallel", name: "Parallel", baseURL: "https://api.parallel.ai/v1", noChat: true},
		{platform: "waterfall", name: "Waterfall", baseURL: "https://api.getwaterfall.org/v1"},
	} {
		cfg.validate = freeLLMAPIValidate
		if !cfg.noChat && cfg.checkResponse == nil {
			cfg.checkResponse = exactModelCheck
			cfg.verifyModelIdentity = true
		}
		// shortcut: buffer DreamPrompting's aliased identities and Inferbase's
		// receipt frames; add native streaming when their frame handling is needed.
		r.register(cfg)
	}
}

func exactModelCheck(requested string, resp *ChatResponse) error {
	if resp.Model == "" || resp.Model != requested {
		return &HTTPError{Status: 502, Message: "provider returned a different or missing model identity"}
	}
	return nil
}

func dreamPromptingCheck(requested string, resp *ChatResponse) error {
	stripped := strings.TrimSuffix(requested, ":free")
	_, unprefixed, hasPrefix := strings.Cut(stripped, "/")
	if requested == "auto" || (resp.Model != "" && (resp.Model == requested || resp.Model == stripped || (hasPrefix && resp.Model == unprefixed))) {
		return nil
	}
	return &HTTPError{Status: 502, Message: "DreamPrompting returned a different or missing model identity"}
}

func freeLLMAPIValidate(ctx context.Context, c *compat, apiKey string) KeyValidationResult {
	endpoint := c.cfg.validateURL
	if endpoint == "" {
		endpoint = c.BaseURL() + "/models"
	}
	method, headers := http.MethodGet, c.authMap(apiKey)
	var payload []byte
	if c.Platform() == "parallel" {
		// Only keyed REST search is supported, not anonymous MCP search.
		if apiKey == "no-key" || strings.TrimSpace(apiKey) == "" {
			return Inconclusive("Parallel requires an API key; anonymous MCP is not implemented")
		}
		endpoint, method, payload = c.BaseURL()+"/search", http.MethodPost, []byte("{}")
		headers = map[string]string{"x-api-key": apiKey, "Content-Type": "application/json"}
	}
	status, body, _, err := c.httpCall(ctx, method, endpoint, headers, payload, c.cfg.validateTimeout)
	if err != nil {
		return Inconclusive(err.Error())
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return classifyValidation(c.Name(), status, body)
	}
	if c.Platform() == "parallel" {
		if status == http.StatusUnprocessableEntity {
			return Valid()
		}
		return Inconclusive("Parallel empty search did not return the expected 422")
	}
	if status < 200 || status >= 300 {
		return Inconclusive(c.Name() + " credential probe failed")
	}
	if c.Platform() == "beatapi" {
		var models struct {
			Data  []json.RawMessage `json:"data"`
			Error json.RawMessage   `json:"error"`
		}
		if json.Unmarshal(body, &models) != nil || models.Data == nil || (len(models.Error) != 0 && string(models.Error) != "null") {
			return Inconclusive("BeatAPI returned an invalid model list")
		}
	}
	if c.Platform() == "inferbase" {
		var credits struct {
			Balance *float64 `json:"balance_micros"`
		}
		if json.Unmarshal(body, &credits) != nil || credits.Balance == nil {
			return Inconclusive("Inferbase returned an invalid credits response")
		}
	}
	if c.endpointIgnoresCredentials(ctx, method, endpoint) {
		return Inconclusive(c.Name() + " credential probe also accepts an invalid key")
	}
	return Valid()
}
