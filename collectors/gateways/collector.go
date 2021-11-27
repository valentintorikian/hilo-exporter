package gateways

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/valentintorikian/hilo-client-go/hilo"
	"strconv"
)

type Collector struct {
	hiloClient                    *hilo.Hilo
	gatewayFirmwareVersion        *prometheus.Desc
	gatewayOnlineStatus           *prometheus.Desc
	gatewayZigbeePairingActivated *prometheus.Desc
	gatewayZigbeeChannel          *prometheus.Desc
}

func (e *Collector) Collect(ch chan<- prometheus.Metric) {
	locations, err := e.hiloClient.Locations()
	if err != nil {
		log.WithError(err).Errorf("couldn't fetch hilo locations")
	}
	for _, location := range locations {
		gateways, err := e.hiloClient.Gateways(location)
		if err != nil {
			log.WithError(err).Errorf("couldn't fetch gateway infos")
			return
		}
		for _, gateway := range gateways {
			var onlineStatus = 0
			if gateway.OnlineStatus == "Online" {
				onlineStatus = 1
			}
			var pairingActivated = 0
			if gateway.ZigBeePairingActivated {
				pairingActivated = 1
			}

			ch <- prometheus.MustNewConstMetric(e.gatewayOnlineStatus,
				prometheus.GaugeValue,
				float64(onlineStatus),
				strconv.Itoa(location.Id),
				location.Name,
			)
			ch <- prometheus.MustNewConstMetric(e.gatewayZigbeeChannel,
				prometheus.GaugeValue,
				float64(gateway.ZigBeeChannel),
				strconv.Itoa(location.Id),
				location.Name,
			)
			ch <- prometheus.MustNewConstMetric(e.gatewayZigbeePairingActivated,
				prometheus.GaugeValue,
				float64(pairingActivated),
				strconv.Itoa(location.Id),
				location.Name,
			)
			ch <- prometheus.MustNewConstMetric(e.gatewayFirmwareVersion,
				prometheus.GaugeValue,
				1,
				strconv.Itoa(location.Id),
				location.Name,
				gateway.FirmwareVersion,
			)
		}
	}
}

func (e *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.gatewayZigbeePairingActivated
	ch <- e.gatewayOnlineStatus
	ch <- e.gatewayFirmwareVersion
	ch <- e.gatewayZigbeeChannel
}

func NewCollector(hiloClient *hilo.Hilo) *Collector {
	e := &Collector{
		hiloClient: hiloClient,
		gatewayZigbeePairingActivated: prometheus.NewDesc(
			"hilo_gateway_peering_activated",
			"",
			[]string{"location_id", "location_name"},
			nil,
		),
		gatewayOnlineStatus: prometheus.NewDesc(
			"hilo_gateway_online",
			"Gateway online status",
			[]string{"location_id", "location_name"},
			nil,
		),
		gatewayFirmwareVersion: prometheus.NewDesc(
			"hilo_gateway_firmware_version",
			"",
			[]string{"location_id", "location_name", "firmware_version"},
			nil,
		),
		gatewayZigbeeChannel: prometheus.NewDesc(
			"hilo_gateway_zigbee_channel",
			"",
			[]string{"location_id", "location_name"},
			nil,
		),
	}
	return e
}
