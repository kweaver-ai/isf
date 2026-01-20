package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type OAuth2 interface {
	Introspection(ctx context.Context, token string, scopes []string) (IntrospectionResult, error)
}

type Config struct {
	TokenEndpoint, IntrospectEndpoint, InternalAccountEndpoint string
}

type oauth2 struct {
	httpClient Client
	config     Config
	trace      Tracer
}

var (
	oOnce sync.Once
	ur    OAuth2
)

func NewOAuth2() OAuth2 {
	oOnce.Do(func() {
		// TODO: 支持自定义httpclient配置，比如连接池，长连接等
		config := Config{
			IntrospectEndpoint: getIntrospectEndpoint(),
		}

		ur = &oauth2{
			httpClient: NewHttpClient(),
			config:     config,
			trace:      NewARTrace(),
		}
	})

	return ur
}

func getHydraAdminURL() url.URL {
	schema := os.Getenv("HYDRA_ADMIN_PROTOCOL")
	host := os.Getenv("HYDRA_ADMIN_HOST")
	port := os.Getenv("HYDRA_ADMIN_PORT")
	url := url.URL{
		Scheme: schema,
		Host:   fmt.Sprintf("%v:%v", host, port),
	}
	return url
}

func getIntrospectEndpoint() string {
	url := getHydraAdminURL()
	url.Path = "/admin/oauth2/introspect"
	return url.String()
}

type IntrospectionResult struct {
	Active            bool                   `json:"active"`                       // Active is a boolean indicator of whether or not the presented token is currently active.
	Audience          []string               `json:"aud,omitempty"`                // Audience contains a list of the token's intended audiences.
	ClientID          string                 `json:"client_id,omitempty"`          // ClientID is aclient identifier for the OAuth 2.0 client that requested this token.
	ExpiresAt         int64                  `json:"exp,omitempty"`                // Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire.
	Extra             map[string]interface{} `json:"ext,omitempty"`                // Extra is arbitrary data set by the session.
	IssuedAt          int64                  `json:"iat,omitempty"`                // Issued at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued.
	IssuerURL         string                 `json:"iss,omitempty"`                // IssuerURL is a string representing the issuer of this token
	NotBefore         int64                  `json:"nbf,omitempty"`                // NotBefore is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token is not to be used before.
	ObfuscatedSubject string                 `json:"obfuscated_subject,omitempty"` // ObfuscatedSubject is set when the subject identifier algorithm was set to "pairwise" during authorization. It is the `sub` value of the ID Token that was issued.
	Scope             string                 `json:"scope,omitempty"`              // Scope is a JSON string containing a space-separated list of scopes associated with this token.
	Subject           string                 `json:"sub,omitempty"`                // Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token.
	TokenUse          string                 `json:"token_use,omitempty"`          // TokenUse is the introspected token's purpose, for example `access_token` or `refresh_token`.
	TokenType         string                 `json:"token_type,omitempty"`         // TokenType is the introspected token's type, always is `Bearer`.
	Username          string                 `json:"username,omitempty"`           // Username is a human-readable identifier for the resource owner who authorized this token.
}

func (o *oauth2) Introspection(ctx context.Context, token string, scopes []string) (result IntrospectionResult, err error) {
	var resp *http.Response
	o.trace.SetInternalSpanName("内省")
	ctx, span := o.trace.AddInternalTrace(ctx)
	defer func() {
		o.trace.TelemetrySpanEnd(span, err)
	}()

	data := url.Values{"token": []string{token}}
	if len(scopes) > 0 {
		data["scope"] = []string{strings.Join(scopes, " ")}
	}
	req, err := http.NewRequest("POST", o.config.IntrospectEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = o.httpClient.Do(ctx, req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Introspection failed, status code is %d", resp.StatusCode)
		return
	}

	defer resp.Body.Close()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return
}
