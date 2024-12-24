package middlewares

import (
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

func CheckXRealIP(cidr string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			realIP := req.Header.Get("X-Real-IP")
			if realIP == `` {
				return ctx.JSON(http.StatusForbidden, map[string]string{"error": "cannot get X-Real-IP header value"})
			}
			_, subnet, err := net.ParseCIDR(cidr)
			if err != nil {
				return ctx.JSON(http.StatusForbidden, map[string]string{"error": "cannot parse CIDR"})
			}
			ip := net.ParseIP(realIP)
			if !subnet.Contains(ip) {
				return ctx.JSON(http.StatusForbidden, map[string]string{"error": "client ip is not in CIDR range"})
			}

			return next(ctx)
		}
	}
}
