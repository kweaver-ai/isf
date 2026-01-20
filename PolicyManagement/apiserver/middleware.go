package apiserver

import (
	"context"
	"net/http"
	"strings"

	"policy_mgnt/utils/errors"

	"policy_mgnt/utils/gocommon/api"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func oauth2Middleware(o api.OAuth2) gin.HandlerFunc {
	return func(c *gin.Context) {
		oauthOn := viper.GetBool("oauth_on")
		if !oauthOn {
			return
		}
		token, err := extraBearerToken(c.Request)
		if err != nil {
			// TODO: 记日志
			errorResponse(c, errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token invalid"}))
			return
		}
		if token == "" {
			// TODO: 记日志
			errorResponse(c, errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token empty"}))
			return
		}
		result, err := o.Introspection(context.Background(), token, nil)
		if err != nil {
			// TODO: 记日志
			errorResponse(c, errors.ErrInternalServerErrorPublic(&api.ErrorInfo{Cause: "introspection failed"}))
			return
		}

		if !result.Active {
			errorResponse(c, errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token does not active"}))
			return
		}

		if result.ClientID != result.Subject {
			userID := result.Subject
			c.Set("userid", userID)
			c.Set("ip", result.Extra["login_ip"])
		}
	}

}

func extraBearerToken(req *http.Request) (string, error) {
	hdr := req.Header.Get("Authorization")
	if hdr == "" {
		return "", nil
	}

	// Example: Bearer xxxx
	th := strings.SplitN(hdr, " ", 2)
	if len(th) != 2 {
		err := errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token invalid"})
		return "invalid", err
	}
	if strings.ToLower(th[0]) != "bearer" {
		err := errors.ErrUnauthorization(&api.ErrorInfo{Cause: "access_token invalid"})
		return "invalid", err
	}
	return th[1], nil
}
