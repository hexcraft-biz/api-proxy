package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
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

func TokenIntrospection(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		token := parts[1]

		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Client-Id")
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Client-Scope")

		// Admin API : POST /oauth2/introspect
		// Content-Type: application/x-www-form-urlencoded
		resp, err := http.PostForm(
			cfg.OAuth2AdminHost+"/oauth2/introspect",
			url.Values{"token": {token}},
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			introspect := hydraIntrospect{}
			json.NewDecoder(resp.Body).Decode(&introspect)

			if *introspect.Active == false {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			} else {
				if introspect.ClientID != nil && introspect.Scope != nil {
					c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Client-Id", *introspect.ClientID)
					c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Client-Scope", *introspect.Scope)
				}
			}
		} else if resp.StatusCode == http.StatusUnauthorized {
			err := hydraError{}
			json.NewDecoder(resp.Body).Decode(&err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.ErrorDescription})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
			return
		}

		// TODO Cache
	}
}

func Userinfo(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// Public API : POST /userinfo
		client := &http.Client{}

		userinfoUrl := cfg.OAuth2PublicHost + "/userinfo"
		req, err := http.NewRequest("GET", userinfoUrl, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Authenticated-User-Id")
		c.Request.Header.Del("X-" + cfg.OAuth2HeaderInfix + "-Authenticated-User-Email")

		req.Header.Set("Authorization", authHeader)
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			userinfo := hydraUserinfo{}
			json.NewDecoder(resp.Body).Decode(&userinfo)

			/*
				X-{infix}-Authenticated-User-Id: {pilgrimID}
				X-{infix}-Authenticated-User: {authentication_provider}:{media}:{identifier}
			*/

			if userinfo.UserID != nil && userinfo.UserIdentifier != nil && userinfo.UserIdentifierMedia != nil && userinfo.AuthenticationProvider != nil {
				c.Request.Header.Set("X-"+cfg.OAuth2HeaderInfix+"-Authenticated-User-Id", *userinfo.UserID)
				c.Request.Header.Set(
					"X-"+cfg.OAuth2HeaderInfix+"-Authenticated-User",
					*userinfo.AuthenticationProvider+":"+*userinfo.UserIdentifierMedia+":"+*userinfo.UserIdentifier,
				)
			}
		} else if resp.StatusCode == http.StatusUnauthorized {
			err := hydraError{}
			json.NewDecoder(resp.Body).Decode(&err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.ErrorDescription})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}
}
