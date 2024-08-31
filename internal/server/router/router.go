package router

import (
	"auth-telegram/internal/handler"
	"auth-telegram/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title Example API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /
func Init(h *handler.Handler) http.Handler {
	prometheus.MustRegister(metrics.IncomingTraffic)

	r := http.ServeMux{}

	r.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://0.0.0.0:1234/swagger/doc.json"),
	))
	// metrics
	r.Handle("GET /metrics", promhttp.Handler())

	// public handlers
	r.Handle("POST /sign-up", h.CORS(h.LogUser(http.HandlerFunc(h.SignUpHandler))))
	r.Handle("POST /v2/sign-up", h.CORS(h.LogUser(http.HandlerFunc(h.SignUpV2Handler))))
	r.Handle("POST /v3/sign-up", h.CORS(h.LogUser(http.HandlerFunc(h.SignUpV3Handler))))
	r.Handle("POST /resend-code", h.CORS(h.LogUser(http.HandlerFunc(h.ResendCodeHandler))))
	r.Handle("POST /sign-in", h.CORS(h.LogUser(http.HandlerFunc(h.SignInHandler))))
	r.Handle("POST /sign-in-verify", h.CORS(h.LogUser(http.HandlerFunc(h.SignInVerifyHandler))))
	r.Handle("POST /v2/sign-in", h.CORS(h.LogUser(http.HandlerFunc(h.SignInV2Handler))))
	r.Handle("POST /refresh", h.CORS(http.HandlerFunc(h.GenerateRefreshTokensHandler)))
	r.Handle("POST /v1/validate-token", h.CORS(http.HandlerFunc(h.ValidateTokenHandler)))

	// auth handlers
	r.Handle("GET /user-info/{id}", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.GetUserInfoHandler)))))
	r.Handle("PUT /user-info/{id}", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.SetUserInfoHandler)))))
	r.Handle("PUT /change-password/{id}", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.ChangePasswordHandler)))))

	// admin handlers
	r.Handle("GET /admin/user-info/{id}", h.CORS(h.Auth(h.IsAdmin(h.LogUser(http.HandlerFunc(h.AdminGetUserInfoHandler))))))
	r.Handle("GET /v2/user-info/{id}", h.CORS(h.Auth(h.IsAdmin(h.LogUser(http.HandlerFunc(h.GetUserInfoV2Handler))))))

	return &r
}
