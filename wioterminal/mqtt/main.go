// This is a sensor station that uses a RTL8720DN running on the device UART2.
// It creates an MQTT connection that publishes a message every second
// to an MQTT broker.
//
// In other words:
// Your computer <--> USB-CDC <--> MCU <--> UART2 <--> RTL8720DN <--> Internet <--> MQTT broker.
//
// You must also install the Paho MQTT package to build this program:
//
//	go get -u github.com/eclipse/paho.mqtt.golang
//
// You can check that mqttpub/mqttsub is running successfully with the following command.
//
//	mosquitto_sub -h test.mosquitto.org -t sago35/tinygo/tx
//	mosquitto_pub -h test.mosquitto.org -t sago35/tinygo/rx -m "{"Temperature": 9.87, "Humidity": 54.32}"
package main

import (
	"device/sam"
	"fmt"
	"machine"
	"math/rand"
	"time"

	"github.com/sago35/tinygo-examples/wioterminal/initialize"
	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/rtl8720dn"
)

// You can override the setting with the init() in another source code.
// func init() {
//    ssid = "your-ssid"
//    password = "your-password"
//    debug = true
//    server = "tinygo.org"
// }

var (
	ssid     string
	password string
	server   string = "tcp://test.mosquitto.org:1883"
	debug           = false
)

var lastRequestTime time.Time
var conn net.Conn
var adaptor *rtl8720dn.RTL8720DN

func main() {
	err := run()
	for err != nil {
		fmt.Printf("error: %s\r\n", err.Error())
		time.Sleep(5 * time.Second)
	}
}

// change these to connect to a different UART or pins for the ESP8266/ESP32
var (
	cl      mqtt.Client
	topicTx = "sago35/tinygo/tx"
	topicRx = "sago35/tinygo/rx"
)

func subHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\r\n", msg.Payload())
}

func readTemperatureAndHumidity(sensor bme280.Device) (float64, float64) {
	temp, _ := sensor.ReadTemperature()
	hum, _ := sensor.ReadHumidity()
	return float64(temp) / 1000, float64(hum) / 100
}

func run() error {
	_, err := initialize.Wifi(ssid, password, 10*time.Second)
	if err != nil {
		return err
	}

	// Enable 3V3 output
	machine.OUTPUT_CTR_3V3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.OUTPUT_CTR_3V3.Low()

	// Initialize I2C0 using BCM0 (SDA) and BCM1 (SCL)
	// I2C SERCOM4 : SCL0_PIN (PA12) + SDA0_PIN (PA13)
	i2c := &machine.I2C{Bus: sam.SERCOM4_I2CM, SERCOM: 4}
	i2c.Configure(machine.I2CConfig{SCL: machine.SCL0_PIN, SDA: machine.SDA0_PIN})

	// Initialize BME280
	sensor := bme280.New(i2c)
	sensor.Configure()

	rand.Seed(time.Now().UnixNano())

	opts := mqtt.NewClientOptions()
	opts.AddBroker(server).SetClientID("tinygo-client-" + randomString(10))

	println("Connecting to MQTT broker at", server)
	cl = mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		failMessage(token.Error().Error())
	}

	// subscribe
	token := cl.Subscribe(topicRx, 0, subHandler)
	token.Wait()
	if token.Error() != nil {
		failMessage(token.Error().Error())
	}

	go publishing(sensor)

	select {}

	// Right now this code is never reached. Need a way to trigger it...
	println("Disconnecting MQTT...")
	cl.Disconnect(100)

	println("Done.")

	return nil
}

func publishing(sensor bme280.Device) {
	for {
		temp, hum := readTemperatureAndHumidity(sensor)
		fmt.Printf("%.2f °C %.2f %%\r\n", temp, hum)
		body := fmt.Sprintf(`{"Temperature": %.2f, "Humidity": %.2f}`, temp, hum)
		data := []byte(body)
		token := cl.Publish(topicTx, 0, false, data)
		token.Wait()
		if token.Error() != nil {
			println(token.Error().Error())
		}

		time.Sleep(2000 * time.Millisecond)
	}
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func failMessage(msg string) {
	for {
		println(msg)
		time.Sleep(1 * time.Second)
	}
}
