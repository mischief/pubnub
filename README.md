pubnub
======

Go Bindings for the PubNub 3.3 REST API

Usage
=====

See pubnub_test.go for example use

Tests
=====

To test the api, use 'go test' in the root of the repository

```
mischief@omicron:~/code/go/src/github.com/mischief/pubnub$ go test -v
=== RUN TestUUID
--- PASS: TestUUID (0.00 seconds)
        pubnub_test.go:15: UUID: f61f59ea-69f7-ebb3-2c56-ef6bf3bbca6c
=== RUN TestPubNubTime
--- PASS: TestPubNubTime (0.09 seconds)
        pubnub_test.go:31: time response: 13650506022429176.000000
=== RUN TestPubNub
--- PASS: TestPubNub (0.53 seconds)
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506024157291"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506024482727"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506024858886"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506025050535"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506025524422"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506025886937"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506026283836"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506026497017"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506027071734"}
        pubnub_test.go:80: Publishing "Hello, World"
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
        pubnub_test.go:87: Publish response: []interface {}{1, "Sent", "13650506027503029"}
        pubnub_test.go:62: Subscriber got a message: []interface {}{"Hello, World"}
PASS
ok      github.com/mischief/pubnub      0.635s
```
