package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/drawbridge/config"
	"github.com/hexcraft-biz/drawbridge/constants"
	"github.com/hexcraft-biz/her"
	"github.com/hexcraft-biz/xuuid"
)

func Dogmas(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requesterId *xuuid.UUID
		if uID := c.Request.Header.Get("X-" + cfg.OAuth2HeaderInfix + "-Authenticated-User-Id"); uID != "" {
			if authUserID, err := xuuid.Parse(uID); err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
				return
			} else {
				requesterId = &authUserID
			}
		}

		var scheme string
		if scheme = c.Request.Header.Get("X-Forwarded-Proto"); scheme == "" {
			scheme = "https"
		}

		if resultDestination, err := cfg.Dogmas.CanAccess(
			c.Request.Header.Get("X-"+cfg.OAuth2HeaderInfix+"-Client-Scope"),
			c.Request.Method,
			scheme+"://"+c.Request.Host+c.Request.URL.Path,
			requesterId,
		); err != nil {
			c.AbortWithStatusJSON(her.NewError(http.StatusForbidden, err, nil).HttpR())
			return
		} else {

			c.Set(constants.MiddlewareKeyProxyRoute, resultDestination)
		}
	}
}
