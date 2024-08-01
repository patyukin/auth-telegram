package router

import (
	"auth-telegram/internal/handler"
	"auth-telegram/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func Init(h *handler.Handler) http.Handler {
	prometheus.MustRegister(metrics.IncomingTraffic)

	r := http.ServeMux{}

	r.Handle("POST /sign-up", h.CORS(h.LogUser(http.HandlerFunc(h.SignUpHandler))))
	r.Handle("POST /resend-code", h.CORS(h.LogUser(http.HandlerFunc(h.ResendCodeHandler))))
	r.Handle("POST /sign-in", h.CORS(h.LogUser(http.HandlerFunc(h.SignInHandler))))
	r.Handle("POST /sign-in-verify", h.CORS(h.LogUser(http.HandlerFunc(h.SignInVerifyHandler))))

	r.Handle("GET /user-info/{id}", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.GetUserInfoHandler)))))
	r.Handle("PUT /user-info/{id}", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.SetUserInfoHandler)))))
	r.Handle("GET /change-password", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.ChangePasswordHandler)))))

	// admin handlers
	r.Handle("GET /admin/user-info", h.CORS(h.Auth(h.IsAdmin(h.LogUser(http.HandlerFunc(h.AdminGetUserInfoHandler))))))
	r.Handle("GET /v2/user-info", h.CORS(h.Auth(h.LogUser(http.HandlerFunc(h.GetUserInfoV2Handler)))))
	r.Handle("GET /metrics", promhttp.Handler())

	return &r
}
