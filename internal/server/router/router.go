package router

import (
	_ "auth-telegram/docs"
	"auth-telegram/internal/handler"
	"auth-telegram/internal/metrics"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/pprof"
)

// Init godoc
// @title Auth API
// @version 1.0
// @description Auth API for microservices
// @host http://0.0.0.0:1234
// @BasePath /
func Init(h *handler.Handler, srvAddress string) http.Handler {
	prometheus.MustRegister(metrics.IncomingTraffic)

	r := http.ServeMux{}

	r.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://0.0.0.0%s/swagger/doc.json", srvAddress)),
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

	// pprof
	r.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/handle/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("goroutine"))

	return &r
}
