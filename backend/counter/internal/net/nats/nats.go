package nats

import (
	"github.com/nats-io/nats.go"
)

var Sub map[string]*nats.Subscription
var Channels = make(map[string]int)
var Enabled = false
var Nc *nats.Conn

func Enable() {
	var err error
	Nc, err = nats.Connect("nats://nats:4222")
	if err != nil {
		return
	}
	Enabled = true
}

// This brings me shame, but circular dependencies would have require minor rework
func Join(festivalCode string, broadcastTotal func(string) error) {
	if Channels[festivalCode] == 0 {
		Sub[festivalCode], _ = Nc.Subscribe(festivalCode, func(msg *nats.Msg) {
			broadcastTotal(festivalCode)
		})
	}
	Channels[festivalCode]++
}

func Leave(festivalCode string) {
	Channels[festivalCode]--
	if Channels[festivalCode] == 0 {
		Sub[festivalCode].Unsubscribe()
	}
}

func Update(festivalCode string) {
	if Channels[festivalCode] > 0 {
		Nc.Publish(festivalCode, []byte("ping"))
	}
}
