package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPortIsBeingUsed(t *testing.T) {
	port := ":11223"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		t.Error("Failed to start listener:", err)
	}
	defer ln.Close()

	assert.False(t, checkPortIsBeingUsed(port), "they should be false")

}
