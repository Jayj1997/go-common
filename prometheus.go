package common

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusBoot(port int) error {
	http.Handle("/metrics", promhttp.Handler())

	// 启动web服务
	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
	if err != nil {
		return err
	}

	return nil
}
