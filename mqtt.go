package main

import proto "github.com/huin/mqtt"

func broadcast(topic, message string) error {

	mqttClient.Publish(&proto.Publish{
		Header:    proto.Header{Retain: false},
		TopicName: topic,
		Payload:   proto.BytesPayload([]byte(message)),
	})

	return nil
}
