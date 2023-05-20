package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	received := []int{}

	n.Handle("echo", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]interface{}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back.
		body["type"] = "echo_ok"

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("generate", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]interface{}{
			"type": "generate_ok",
			"id":   uuid.New(), // maybe kinda cheating but whatever
		})
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var bm struct {
			Message int `json:"message"`
		}
		if err := json.Unmarshal(msg.Body, &bm); err != nil {
			return err
		}
		received = append(received, bm.Message)
		return n.Reply(msg, map[string]interface{}{
			"type": "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]interface{}{
			"type":     "read_ok",
			"messages": received,
		})
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]interface{}{
			"type": "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
