package utils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestGetOutBoundIP(t *testing.T) {
	t.Log(GetOutBoundIP())
	NetworkInterfaceName = ""
	nets, e := net.Interfaces()
	assert.NoError(t, e)
	t.Log(nets)
	assert.NotZero(t, nets)
	for _, i := range nets {
		t.Log(i.Name)
	}
}
