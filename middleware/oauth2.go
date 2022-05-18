package middleware

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/api-proxy/config"
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
	Aud       *[]string `json:"aud"`
	Iat       *int      `json:"iat"`
	Iss       *string   `json:"iss"`
	Acr       *string   `json:"acr,omitempty"`
	AuthTime  *int      `json:"auth_time,omitempty"`
	Rat       *int      `json:"rat,omitempty"`
	Sub       *string   `json:"sub,omitempty"`
	UserID    *string   `json:"user_id,omitempty"`
	UserEmail *string   `json:"user_email,omitempty"`
}

type hydraError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func TokenIntrospection(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		token := parts[1]

		// Admin API : POST /oauth2/introspect
		// Content-Type: application/x-www-form-urlencoded
		resp, err := http.PostForm(
			cfg.Env.Oauth2AdminHost+"/oauth2/introspect",
			url.Values{"token": {token}},
		)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			introspect := hydraIntrospect{}
			json.NewDecoder(resp.Body).Decode(&introspect)

			if *introspect.Active == false {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			} else {
				if introspect.ClientID != nil && introspect.Scope != nil {
					ctx.Request.Header.Add("X-"+cfg.Env.Oauth2HeaderPrefix+"-Client-Id", *introspect.ClientID)
					ctx.Request.Header.Add("X-"+cfg.Env.Oauth2HeaderPrefix+"-Client-Scope", *introspect.Scope)
				}
			}
		} else if resp.StatusCode == http.StatusUnauthorized {
			err := hydraError{}
			json.NewDecoder(resp.Body).Decode(&err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.ErrorDescription})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
			return
		}

		// TODO Cache
	}
}

func Userinfo(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// Public API : POST /userinfo
		client := &http.Client{}

		userinfoUrl := cfg.Env.Oauth2PublicHost + "/userinfo"
		req, err := http.NewRequest("GET", userinfoUrl, nil)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		req.Header.Add("Authorization", authHeader)
		resp, err := client.Do(req)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			userinfo := hydraUserinfo{}
			json.NewDecoder(resp.Body).Decode(&userinfo)

			if userinfo.UserID != nil && userinfo.UserEmail != nil {
				ctx.Request.Header.Add("X-"+cfg.Env.Oauth2HeaderPrefix+"-Authenticated-User-Id", *userinfo.UserID)
				ctx.Request.Header.Add("X-"+cfg.Env.Oauth2HeaderPrefix+"-Authenticated-User-Email", *userinfo.UserEmail)
			}
		} else if resp.StatusCode == http.StatusUnauthorized {
			err := hydraError{}
			json.NewDecoder(resp.Body).Decode(&err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.ErrorDescription})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}
}

type oAuth2Scope struct {
	ID                 string `json:"id"`
	ResourceDomainName string `json:"resourceDomainName"`
	ResourceName       string `json:"resourceName"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type oAuth2Scopes []oAuth2Scope

type oAuth2ScopesError struct {
	Message string `json:"message"`
}

func VerifyScope(cfg *config.Config, internalHostname string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientScope := ctx.Request.Header.Get("X-" + cfg.Env.Oauth2HeaderPrefix + "-Client-Scope")

		if clientScope == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		clientScopes := strings.Split(clientScope, " ")

		// Scope API : GET /scopes/v1/scopes?resourceDomainName={internalHostname}&name=clientScope1|clientScope2
		client := &http.Client{}

		scopesUrl := cfg.Env.Oauth2ScopesHost + "/scopes/v1/scopes"
		req, err := http.NewRequest("GET", scopesUrl, nil)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}

		q := req.URL.Query()
		q.Add("resourceDomainName", internalHostname)
		q.Add("name", strings.Join(clientScopes[:], "|"))
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			scopes := oAuth2Scopes{}
			json.NewDecoder(resp.Body).Decode(&scopes)

			if len(scopes) == 0 {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
				return
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}
}
