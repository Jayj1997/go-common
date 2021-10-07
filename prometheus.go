package common

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func PromethusBoot(port int) error {
	http.Handle("/metrics", promhttp.Handler())

	// 启动web服务
	err := http.ListenAndServe("0.0.0.0"+strconv.Itoa(port), nil)
	logrus.Info("promethus 启动成功")
	if err != nil {
		return err
	}

	return nil
}
