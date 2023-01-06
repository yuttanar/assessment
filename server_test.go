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

func TestCheckPortPatternMatch(t *testing.T) {
	assert.True(t, checkPortPatternMatch(":12565"), "they should be true")
	assert.True(t, checkPortPatternMatch(":2565"), "they should be true")
	assert.False(t, checkPortPatternMatch(":256500"), "they should be false")
	assert.False(t, checkPortPatternMatch("2565"), "they should be false")
	assert.False(t, checkPortPatternMatch(":32f5"), "they should be false")
}
