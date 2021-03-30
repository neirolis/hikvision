package hikvision

import (
	"os"
	"testing"
)

var (
	addr = os.Getenv("TEST_HIKVISION_ADDR")
	user = os.Getenv("TEST_HIKVISION_USER")
	pass = os.Getenv("TEST_HIKVISION_PASS")

	c *Client
)

func TestNew(t *testing.T) {
	if addr == "" || user == "" || pass == "" {
		t.Skip("TEST_HIKVISION_ADDR, TEST_HIKVISION_USER, TEST_HIKVISION_PASS environment variable should be specified")
	}

	var err error
	c, err = New(addr, user, pass)
	if err != nil {
		t.Error(err)
	}
}

func TestPTZrelative(t *testing.T) {
	if c == nil {
		t.Skip("no client")
	}

	err := c.PTZrelative("1", PTZrelativeData{127, 127, 0})
	if err != nil {
		t.Error(err)
	}
}

func TestPTZabsolute(t *testing.T) {
	if c == nil {
		t.Skip("no client")
	}

	err := c.PTZabsolute("1", PTZabsoluteData{100, 1500, 0})
	if err != nil {
		t.Error(err)
	}
}

func TestPTZstatus(t *testing.T) {
	if c == nil {
		t.Skip("no client")
	}

	status, err := c.PTZstatus("1")
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", status)
}
