package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeFL(t *testing.T) {
	s := "FL022020204150504C45545620202020"
	s1, e := decode_fl(s)
	assert.Equal(t, "   APPLETV    ", s1)
	assert.Equal(t, nil, e)
}

func TestDecodeFL2(t *testing.T) {
	s := "FL004D2E564F4C20202D31382E356442"
	s1, e := decode_fl(s)
	assert.Equal(t, "M.VOL  -18.5dB", s1)
	assert.Equal(t, e, nil)
}

func TestDecodeMessage(t *testing.T) {
	f, e := decode_message("LM010f")
	assert.Equal(t, nil, e)
	assert.Equal(t, "Mode is WIDE SURR MUSIC", f)

	_, e = decode_message("foobar")
	assert.NotNil(t, e)
}
