package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"time"

	"go.uber.org/zap"
)

// AddLogger logs request/response pair
func AddLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := GetReqID(ctx)

		// Prepare fields to log
		var scheme string
		if ctx.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		proto := ctx.Request.Proto
		method := ctx.Request.Method
		remoteAddr := ctx.Request.RemoteAddr
		userAgent := ctx.Request.UserAgent()
		uri := strings.Join([]string{scheme, "://", ctx.Request.Host, ctx.Request.RequestURI}, "")

		// Log HTTP request
		logger.Debug("request started",
			zap.String("request-id", id),
			zap.String("http-scheme", scheme),
			zap.String("http-proto", proto),
			zap.String("http-method", method),
			zap.String("remote-addr", remoteAddr),
			zap.String("user-agent", userAgent),
			zap.String("uri", uri),
		)

		t1 := time.Now()

		ctx.Next()

		// Log HTTP response
		logger.Info("request completed",
			zap.String("request-id", id),
			zap.String("http-scheme", scheme),
			zap.String("http-proto", proto),
			zap.String("http-method", method),
			zap.String("remote-addr", remoteAddr),
			zap.String("user-agent", userAgent),
			zap.String("uri", uri),
			zap.Float64("elapsed-ms", float64(time.Since(t1).Nanoseconds())/1000000.0),
		)
	}
}
