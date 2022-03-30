package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
)

func TokenIntrospection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		// Admin API : POST /oauth2/introspect
		token := parts[1]
		adminURL, err := url.Parse(os.Getenv("OAUTH2_HOST"))
		if err != nil {
			log.Fatal(err)
		}

		hydraAdmin := client.NewHTTPClientWithConfig(
			nil,
			&client.TransportConfig{
				Schemes:  []string{adminURL.Scheme},
				Host:     adminURL.Host,
				BasePath: adminURL.Path,
			})

		paramsObj := admin.NewIntrospectOAuth2TokenParams()
		paramsObj.SetToken(token)

		res, hydraErr := hydraAdmin.Admin.IntrospectOAuth2Token(paramsObj)

		if hydraErr != nil {
			switch e := hydraErr.(type) {
			case *admin.IntrospectOAuth2TokenUnauthorized:
				fmt.Println(e)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
				return
			case *admin.IntrospectOAuth2TokenInternalServerError:
				fmt.Println(e)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
				return
			}
		}

		if *res.Payload.Active != true {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
		}

		// TODO Cache
	}
}
