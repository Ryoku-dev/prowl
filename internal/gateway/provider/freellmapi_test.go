package provider

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFreeLLMAPIStreamChecksBeforeContent(t *testing.T) {
	for _, frame := range []string{
		`{"model":"wrong","choices":[{"delta":{"content":"do not leak"}}]}`,
		`{"choices":[{"delta":{"content":"do not leak"}}]}`,
		`{"error":{"message":"failed"}}`,
		`not json`,
	} {
		stream := newSSEStream(&http.Response{Body: io.NopCloser(strings.NewReader("data: " + frame + "\n\n"))}, func() {}, true, "wanted")
		stream.checkResponse = exactModelCheck
		chunk, err := stream.Recv()
		if err == nil || chunk != nil {
			t.Fatalf("unsafe frame was forwarded: %s", frame)
		}
		_ = stream.Close()
	}
	stream := newSSEStream(&http.Response{Body: io.NopCloser(strings.NewReader("data: {\"model\":\"wanted\",\"choices\":[{\"delta\":{\"content\":\"ok\"}}]}\n\ndata: {\"choices\":[],\"usage\":{}}\n\ndata: [DONE]\n\n"))}, func() {}, true, "wanted")
	stream.checkResponse = exactModelCheck
	defer stream.Close()
	for i := 0; i < 2; i++ {
		if _, err := stream.Recv(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := stream.Recv(); err != io.EOF {
		t.Fatalf("end = %v", err)
	}
}

func TestFreeLLMAPIEndpoints(t *testing.T) {
	r := NewRegistry()
	for platform, base := range map[string]string{
		"beatapi":        "https://api.beatapi.io/v1",
		"blaze":          "https://api.blazeapi.org/paid/v1",
		"blockrun":       "https://api.blockrun.ai/v1",
		"dreamprompting": "https://dreamprompting.com/api/v1",
		"inferbase":      "https://api.inferbase.ai/v1",
		"logfare":        "https://logfare.ai/v1",
		"parallel":       "https://api.parallel.ai/v1",
		"waterfall":      "https://api.getwaterfall.org/v1",
	} {
		t.Run(platform, func(t *testing.T) {
			p, ok := r.Resolve(platform, "")
			if !ok || p.BaseURL() != base || !Known(platform) {
				t.Fatalf("endpoint %s not registered at %s", platform, base)
			}
			if platform == "beatapi" || platform == "parallel" {
				if _, err := p.ChatCompletion(context.Background(), "test", &ChatRequest{}); !errors.Is(err, ErrChatUnsupported) {
					t.Fatalf("non-chat provider accepted chat: %v", err)
				}
				if _, err := p.StreamChatCompletion(context.Background(), "test", &ChatRequest{}); !errors.Is(err, ErrChatUnsupported) {
					t.Fatalf("non-chat provider accepted streaming chat: %v", err)
				}
			}
		})
	}
}

func TestFreeLLMAPIModelIdentity(t *testing.T) {
	for _, tc := range []struct {
		requested, returned string
		valid               bool
	}{
		{"groq/openai/gpt-oss-20b:free", "openai/gpt-oss-20b", true},
		{"codestral/codestral-latest", "codestral-latest", true},
		{"route/model:free", "route/model", true},
		{"route/model", "route/model", true},
		{"auto", "selected", true},
		{"route/model", "other", false},
		{"route/model", "", false},
	} {
		if err := dreamPromptingCheck(tc.requested, &ChatResponse{Model: tc.returned}); (err == nil) != tc.valid {
			t.Errorf("DreamPrompting %s -> %s: %v", tc.requested, tc.returned, err)
		}
	}
	for _, model := range []string{"", "other", "auto"} {
		if exactModelCheck("named", &ChatResponse{Model: model}) == nil {
			t.Errorf("accepted substitution %q", model)
		}
	}
	if err := exactModelCheck("named", &ChatResponse{Model: "named"}); err != nil {
		t.Fatal(err)
	}
}

func TestFreeLLMAPIValidation(t *testing.T) {
	for _, tc := range []struct {
		platform, path, body string
		status               int
		want                 KeyValidationStatus
	}{
		{"beatapi", "/models", `{"data":[]}`, 200, KeyValid},
		{"beatapi", "/models", `{}`, 200, KeyInconclusive},
		{"beatapi", "/models", `{"data":[],"error":{}}`, 200, KeyInconclusive},
		{"inferbase", "/credits", `{"balance_micros":0}`, 200, KeyValid},
		{"inferbase", "/credits", `{}`, 200, KeyInconclusive},
		{"blaze", "/usage", `{}`, 200, KeyValid},
		{"blockrun", "/account", `{}`, 401, KeyInvalid},
		{"waterfall", "/models", `{}`, 503, KeyInconclusive},
		{"parallel", "/search", `{}`, 422, KeyValid},
		{"parallel", "/search", `{}`, 401, KeyInvalid},
		{"parallel", "/search", `{}`, 500, KeyInconclusive},
	} {
		t.Run(tc.platform+tc.body+http.StatusText(tc.status), func(t *testing.T) {
			up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != tc.path {
					t.Errorf("path = %s, want %s", req.URL.Path, tc.path)
				}
				if tc.platform == "parallel" {
					body, _ := io.ReadAll(req.Body)
					if req.Method != "POST" || string(body) != "{}" || req.Header.Get("x-api-key") != "test" || req.Header.Get("Authorization") != "" {
						t.Errorf("invalid Parallel search probe")
					}
				} else {
					if req.Header.Get("Authorization") == "Bearer "+sentinelInvalidKey {
						w.WriteHeader(401)
						return
					}
					if req.Header.Get("Authorization") != "Bearer test" {
						t.Error("missing bearer")
					}
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer up.Close()
			c := compatOf(t, NewRegistry(), tc.platform)
			c.cfg.baseURL = up.URL
			c.cfg.validateURL = up.URL + tc.path
			if got := c.ValidateKey(context.Background(), "test"); got.Status != tc.want {
				t.Fatalf("validation = %+v, want %v", got, tc.want)
			}
		})
	}
}
