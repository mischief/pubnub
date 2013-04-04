package pubnub

import (
	"testing"
	"time"
)

func TestUUID(t *testing.T) {
	u, err := UUIDGen()

	if err != nil {
		t.Error(err)
	}

	t.Logf("UUID: %s", u)
}

func TestPubNubTime(t *testing.T) {
	pubnub := NewPubNub("demo", "demo", "", "", false)

	done := make(chan bool)

	pubnub.Time(func(msg []interface{}) bool {
		time, ok := msg[0].(float64)

		if !ok {
			t.Errorf("time response is not a float64")
		}

		if time != 0 {
			t.Logf("time response: %f", time)
		} else {
			t.Errorf("Invalid time response")
		}

		done <- true

		return false
	})

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("No response from time call")
	}
}

func TestPubNub(t *testing.T) {
	channel := "chat"
	message := "Hello, World"
	nmsgs := 10

	pubnub := NewPubNub("demo", "demo", "", "", false)

	done := make(chan bool, nmsgs)

	// subscribe to a PubNub channel in another goroutine, .Subscribe blocks
	go func() {

		// called when we get a message
		cb := func(msg []interface{}) bool {
			t.Logf("Subscriber got a message: %#v", msg)
			done <- true

			// return true to be called again, or false to stop the subscription
			return true
		}

		err := pubnub.Subscribe(channel, cb)
		if err != nil {
			t.Errorf("Subscribe error: %s", err)
		}
	}()

	// wait a moment..
	time.Sleep(100 * time.Millisecond)

	// publish some messages
	for i := 0; i < nmsgs; i++ {
		t.Logf("Publishing %q", message)
		resp, err := pubnub.Publish(channel, []interface{}{message})

		if err != nil {
			t.Errorf("Publish error: %s", err)
		}

		t.Logf("Publish response: %#v", resp)
	}

loop:
	for {

		select {
		// a message completed
		case <-done:
			nmsgs--

			if nmsgs == 0 {
				// success; all messages accounted for
				break loop
			}

			// failure; we timed out
		case <-time.After(1 * time.Second):
			// timeout, failure
			t.Errorf("Subscriber timed out on channel %s", channel)
			break loop
		}
	}
}
