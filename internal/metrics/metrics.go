package metrics

import "github.com/prometheus/client_golang/prometheus"

var IncomingTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "incoming_traffic",
	Help: "Incoming traffic to the application",
})

var SignUpV1RegisterTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sign_up_v1_register",
	Help: "Sign up v1 register",
})

var SignUpV2RegisterTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sign_up_v2_register",
	Help: "Sign up v2 register",
})

var SignUpV3RegisterTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sign_up_v3_register",
	Help: "Sign up v3 register",
})

var GetInfoUserTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "get_info_user",
	Help: "Get info user",
})

var SuccessfulAuthentications = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "successful_authentications",
	Help: "Successful authentications",
})

var TotalAuthentications = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "total_authentications",
	Help: "Total authentications",
})

var FailedAuthentications = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "failed_authentications",
	Help: "Failed authentications",
})
