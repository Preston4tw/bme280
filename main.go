package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bme280Temp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_temp",
	})
	bme280Pressure = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_pressure",
	})
	bme280Humidity = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_humidity",
	})
)

type bme280data struct {
	temp     float64
	pressure float64
	humidity float64
}

func getbme280Data() bme280data {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("1")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()

	// Read temperature from the sensor:
	var env physic.Env
	if err = dev.Sense(&env); err != nil {
		log.Fatal(err)
	}
	var data bme280data
	data.temp = float64(env.Temperature.Fahrenheit())
	data.pressure, _ = strconv.ParseFloat(strings.TrimSuffix(env.Pressure.String(), "kPa"), 64)
	data.humidity, _ = strconv.ParseFloat(strings.TrimSuffix(env.Humidity.String(), "%rH"), 64)
	bme280Temp.Set(data.temp)
	bme280Pressure.Set(data.pressure)
	bme280Humidity.Set(data.humidity)
	fmt.Printf("%v\n", data)
	return data
}

func main() {
	go func() {
		for {
			getbme280Data()
			time.Sleep(5 * time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
