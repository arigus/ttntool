// Copyright © 2016 Hylke Visser
// MIT Licensed - See LICENSE file

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/TheThingsNetwork/server-shared"
	"github.com/apex/log"
	"github.com/htdvisser/ttntool/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CLI Variables
var (
	followAll bool
	showRaw   bool
)

var mqttClient *MQTT.Client
var devices []string

var followCmd = &cobra.Command{
	Use:   "follow devAddr [devAddr [...]]",
	Short: "Follow the messages from devices",
	Long:  `Connects to The Things Network and prints out messages received from the specified devices.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 && !followAll {
			cmd.Help()
			return
		}

		// Set up MQTT Client
		setupMQTT()

		// Subscribe to device topics
		switch {
		case followAll:
			devices = []string{"+"}
			log.Debug("Will subscribe to all devices")
		case len(args) == 1:
			devices = args
			log.Debugf("Will subscribe to device %s", devices[0])
		default:
			devices = args
			log.Debugf("Will subscribe to %d devices", len(devices))
		}

		// Connect
		connectMQTT()

		// Keep running...
		for {
			time.Sleep(60 * time.Second)
		}

	},
}

func init() {
	RootCmd.AddCommand(followCmd)

	followCmd.Flags().BoolVar(&followAll, "all", false, "Follow all devices")
	followCmd.Flags().BoolVar(&showRaw, "raw", false, "Show raw data")
}

func setupMQTT() {
	broker := fmt.Sprintf("tcp://%s:1883", viper.GetString("broker"))
	opts := MQTT.NewClientOptions().AddBroker(broker)

	clientID := fmt.Sprintf("ttntool-%d", util.RandString(15))
	opts.SetClientID(clientID)

	opts.SetKeepAlive(20)

	opts.SetOnConnectHandler(func(client *MQTT.Client) {
		log.Info("Connected to The Things Network")
		subscribeToDevices()
	})

	opts.SetDefaultPublishHandler(handleMessage)

	opts.SetConnectionLostHandler(func(client *MQTT.Client, err error) {
		log.WithError(err).Error("Connection Lost. Reconnecting...")
	})

	mqttClient = MQTT.NewClient(opts)
}

func connectMQTT() {
	log.Infof("Connecting to The Things Network...")
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func subscribeToDevices() {
	log.Info("Subscribing to devices...")
	for _, devAddr := range devices {
		topic := fmt.Sprintf("nodes/%s/packets", devAddr)
		if token := mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			log.WithField("topic", topic).WithError(token.Error()).Fatal("Could not subscribe")
		}
	}
}

// messageHandler is called when a new message arrives
func handleMessage(client *MQTT.Client, msg MQTT.Message) {

	// Unmarshal JSON to RxPacket
	var packet shared.RxPacket
	err := json.Unmarshal(msg.Payload(), &packet)
	if err != nil {
		log.WithField("topic", msg.Topic()).WithError(err).Warn("Failed to unmarshal JSON.")
		return
	}

	// Decode payload
	data, err := base64.StdEncoding.DecodeString(packet.Data)
	if err != nil {
		log.WithField("topic", msg.Topic()).WithError(err).Warn("Failed to decode Payload.")
		return
	}

	ctx := log.WithFields(log.Fields{
		"devAddr": packet.NodeEui,
	})

	if showRaw {
		ctx = ctx.WithField("raw_data", fmt.Sprintf("%x", data))
	}

	// Check for unprintable characters
	unprintable, _ := regexp.Compile(`[^[:print:]]`)
	if unprintable.Match(data) {
		ctx.Debug("Received Message")
	} else {
		ctx.WithField("message", fmt.Sprintf("%s", data)).Info("Received Message")
	}

}
