package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/her"
)

type hydraIntrospect struct {
	Active    *bool     `json:"active"`
	Scope     *string   `json:"scope"`
	ClientID  *string   `json:"client_id"`
	Sub       *string   `json:"sub"`
	Exp       *int      `json:"exp"`
	Iat       *int      `json:"iat"`
	Nbf       *int      `json:"nbf"`
	Aud       *[]string `json:"aud"`
	Iss       *string   `json:"iss"`
	TokenType *string   `json:"token_type"`
	TokenUse  *string   `json:"token_use"`
}

type hydraUserinfo struct {
	Aud                    *[]string `json:"aud"`
	Iat                    *int      `json:"iat"`
	Iss                    *string   `json:"iss"`
	Acr                    *string   `json:"acr,omitempty"`
	AuthTime               *int      `json:"auth_time,omitempty"`
	Rat                    *int      `json:"rat,omitempty"`
	Sub                    *string   `json:"sub,omitempty"`
	UserID                 *string   `json:"user_id,omitempty"`
	UserIdentifier         *string   `json:"user_identifier,omitempty"`
	UserIdentifierMedia    *string   `json:"user_identifier_media,omitempty"`
	AuthenticationProvider *string   `json:"authentication_provider,omitempty"`
}

type hydraError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func parseHydraIntrospect(oauth2AdminRootUrl, tokenstring string) (*hydraIntrospect, her.Error) {
	parts := strings.SplitN(tokenstring, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, her.NewErrorWithMessage(http.StatusUnauthorized, "Invalid token", nil)
	}

	resp, err := http.PostForm(oauth2AdminRootUrl+"/oauth2/introspect", url.Values{"token": {parts[1]}})
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		introspect := new(hydraIntrospect)
		json.NewDecoder(resp.Body).Decode(introspect)
		if *introspect.Active == false {
			return nil, her.ErrUnauthorized
		}

		return introspect, nil

	case http.StatusUnauthorized:
		err := new(hydraError)
		json.NewDecoder(resp.Body).Decode(err)
		return nil, her.New(http.StatusUnauthorized, err)
	}

	return nil, her.NewErrorWithMessage(http.StatusInternalServerError, "Unknown introspect response", nil)
}

func TokenIntrospection(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Client-Id")
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Client-Scope")

		if accessToken, err := parseHydraIntrospect(cfg.OAuth2AdminHost, c.GetHeader("Authorization")); err != nil {
			c.AbortWithStatusJSON(err.HttpR())
		} else {
			if accessToken.ClientID != nil {
				c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Client-Id", *accessToken.ClientID)
			}

			if accessToken.Scope != nil {
				c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Client-Scope", *accessToken.Scope)
			}
		}

		// TODO Cache
	}
}

// ================================================================
func parseHydraUserinfo(oauth2PublicRootUrl, tokenstring string) (*hydraUserinfo, her.Error) {
	if tokenstring == "" {
		return nil, her.NewErrorWithMessage(http.StatusUnauthorized, "Invalid token", nil)
	}

	client := &http.Client{}

	userinfoUrl := oauth2PublicRootUrl + "/userinfo"
	req, err := http.NewRequest("GET", userinfoUrl, nil)
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	}

	req.Header.Set("Authorization", tokenstring)
	resp, err := client.Do(req)
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		userinfo := new(hydraUserinfo)
		json.NewDecoder(resp.Body).Decode(userinfo)
		return userinfo, nil

	case http.StatusUnauthorized:
		err := new(hydraError)
		json.NewDecoder(resp.Body).Decode(err)
		return nil, her.New(http.StatusUnauthorized, err)
	}

	return nil, her.NewErrorWithMessage(http.StatusInternalServerError, "Unknown userinfo response", nil)
}

func Userinfo(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Authenticated-User-Id")
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Authenticated-User")

		if userinfo, err := parseHydraUserinfo(cfg.OAuth2PublicHost, c.GetHeader("Authorization")); err != nil {
			c.AbortWithStatusJSON(err.HttpR())
		} else {

			if userinfo.UserID != nil {
				c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Authenticated-User-Id", *userinfo.UserID)
			}

			authenticationProvider, media, identifier := "", "", ""
			if userinfo.AuthenticationProvider != nil {
				authenticationProvider = *userinfo.AuthenticationProvider
			}

			if userinfo.UserIdentifierMedia != nil {
				media = *userinfo.UserIdentifierMedia
			}

			if userinfo.UserIdentifier != nil {
				identifier = *userinfo.UserIdentifier
			}

			c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Authenticated-User", authenticationProvider+":"+media+":"+identifier)
		}
	}
}
