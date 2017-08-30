package messaging

import (
	lib "github.com/k8guard/k8guardlibs"
	msg "github.com/k8guard/k8guardlibs/messaging"
	"github.com/k8guard/k8guardlibs/messaging/types"
)

var MessageProducer types.MessageProducer

func InitBroker() {
	s, err := msg.CreateMessageProducer(
		types.MessageBrokerType(lib.Cfg.MessageBroker), types.DISCOVER_CLIENTID, lib.Cfg)
	if err != nil {
		lib.Log.Error("Creating Messagging Producer ", err)
		panic(err)
	}
	MessageProducer = s
}

func TestBrokerWithTestMessage() error {
	// Sending Test Data
	err := MessageProducer.SendData(types.TEST_MESSAGE, "Testing")
	if err != nil {
		lib.Log.Error("Error trying to send test data to broker ", err)
	}
	return err
}

func SendData(kind types.MessageType, name string, message interface{}) {
	lib.Log.Debugf("Sending %s: %v to broker", name, message)
	err := MessageProducer.SendData(kind, message)
	if err != nil {
		lib.Log.Error("Error trying to send message to broker ", err)
	}
}

func InitStatsHandler() {
	lib.Log.Debug("Initializing stats handler")
	if MessageProducer == nil {
		InitBroker()
	}
	MessageProducer.InitStatsHandler()
}

func CloseBroker() {
	lib.Log.Debug("Closing broker")
	MessageProducer.Close()
}
