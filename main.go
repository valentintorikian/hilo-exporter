package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/valentintorikian/hilo-client-go/hilo"
	"github.com/valentintorikian/hilo-exporter/collectors/gateways"
	"github.com/valentintorikian/hilo-exporter/collectors/thermostats"
	"net/http"
)

const (
	metricsRoute = "/metrics"
)

var (
	listenAddress = flag.String("web.listen-address", ":6464", "Address to listen on for web interface and telemetry.")
	hiloUsername  = flag.String("hilo.username", "", "Hilo username")
	hiloPassword  = flag.String("hilo.password", "", "Hilo password")
)

func main() {
	flag.Parse()
	hiloClient := hilo.NewHilo(*hiloUsername, *hiloPassword)
	hiloThermostatCollector := thermostats.NewCollector(hiloClient)
	hiloGatewayCollector := gateways.NewCollector(hiloClient)
	prometheus.MustRegister(hiloThermostatCollector)
	prometheus.MustRegister(hiloGatewayCollector)

	http.Handle(metricsRoute, promhttp.Handler())

	s := &http.Server{
		Addr: *listenAddress,
	}

	log.Fatal(s.ListenAndServe())
}
