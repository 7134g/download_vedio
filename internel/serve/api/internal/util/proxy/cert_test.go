package proxy

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLoadCert(t *testing.T) {
	if err := LoadCert(); err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, ca, nil)
	assert.NotEqual(t, private, nil)
}

func TestGenMITM(t *testing.T) {
	if err := GenMITM(); err != nil {
		t.Fatal(err)
	}

}

func TestNewCert(t *testing.T) {
	cd, pd, err := newCert("proxy", "location", time.Hour*24*365*10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(cd))
	t.Log(string(pd))
}
