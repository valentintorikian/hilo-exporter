package thermostats

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/valentintorikian/hilo-client-go/hilo"
	"strconv"
)

type Collector struct {
	hiloClient                   *hilo.Hilo
	thermostatCurrentTemperature *prometheus.Desc
	thermostatTargetTemperature  *prometheus.Desc
	thermostatCurrentPower       *prometheus.Desc
	thermostatHeating            *prometheus.Desc
	thermostatHumidity           *prometheus.Desc
}

func (e *Collector) Collect(ch chan<- prometheus.Metric) {
	locations, err := e.hiloClient.Locations()
	if err != nil {
		log.WithError(err).Errorf("couldn't fetch hilo locations")
		return
	}
	for _, location := range locations {
		devices, err := e.hiloClient.Devices(location)
		if err != nil {
			log.WithError(err).Errorf("couldn't fetch hilo devices")
			return
		}
		for _, device := range devices {
			if device.Type != "Thermostat" {
				continue
			}
			attributes, err := e.hiloClient.DeviceAttributes(device)
			if err != nil {
				log.WithError(err).Errorf("couldn't fetch hilo device's attributes")
				return
			}
			for attributeName, attribute := range attributes {
				switch attributeName {
				case "humidity":
					ch <- prometheus.MustNewConstMetric(e.thermostatHumidity,
						prometheus.GaugeValue,
						attribute.Value,
						strconv.Itoa(location.Id),
						location.Name,
						strconv.Itoa(device.Id),
						device.Name,
					)
				case "heating":
					ch <- prometheus.MustNewConstMetric(e.thermostatHeating,
						prometheus.GaugeValue,
						attribute.Value,
						strconv.Itoa(location.Id),
						location.Name,
						strconv.Itoa(device.Id),
						device.Name,
					)
				case "targetTemperature":
					ch <- prometheus.MustNewConstMetric(e.thermostatTargetTemperature,
						prometheus.GaugeValue,
						attribute.Value,
						strconv.Itoa(location.Id),
						location.Name,
						strconv.Itoa(device.Id),
						device.Name,
					)
				case "currentTemperature":
					ch <- prometheus.MustNewConstMetric(e.thermostatCurrentTemperature,
						prometheus.GaugeValue,
						attribute.Value,
						strconv.Itoa(location.Id),
						location.Name,
						strconv.Itoa(device.Id),
						device.Name,
					)
				case "power":
					ch <- prometheus.MustNewConstMetric(e.thermostatCurrentPower,
						prometheus.GaugeValue,
						attribute.Value,
						strconv.Itoa(location.Id),
						location.Name,
						strconv.Itoa(device.Id),
						device.Name,
					)
				}

			}
		}
	}
}

func (e *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.thermostatCurrentTemperature
	ch <- e.thermostatTargetTemperature
	ch <- e.thermostatCurrentPower
	ch <- e.thermostatHumidity
	ch <- e.thermostatHeating
}

func NewCollector(hiloClient *hilo.Hilo) *Collector {
	e := &Collector{
		hiloClient: hiloClient,
		thermostatCurrentTemperature: prometheus.NewDesc(
			"hilo_thermostat_current_temperature_celsius",
			"Thermostat current temperature",
			[]string{"location_id", "location_name", "device_id", "device_name"},
			nil,
		),
		thermostatTargetTemperature: prometheus.NewDesc(
			"hilo_thermostat_target_temperature_celsius",
			"Thermostat target temperature",
			[]string{"location_id", "location_name", "device_id", "device_name"},
			nil,
		),
		thermostatCurrentPower: prometheus.NewDesc(
			"hilo_thermostat_power_watt",
			"Thermostat current power draw in watts",
			[]string{"location_id", "location_name", "device_id", "device_name"},
			nil,
		),
		thermostatHeating: prometheus.NewDesc(
			"hilo_thermostat_heating_percent",
			"Thermostat current heating strength",
			[]string{"location_id", "location_name", "device_id", "device_name"},
			nil,
		),
		thermostatHumidity: prometheus.NewDesc(
			"hilo_thermostat_humidity_percent",
			"Thermostat current humidity",
			[]string{"location_id", "location_name", "device_id", "device_name"},
			nil,
		),
	}
	return e
}
