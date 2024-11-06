package user

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNodeId(t *testing.T) {
	// Тестирование с корректным запросом
	req, err := http.NewRequest("GET", "/?nodeId=123", nil)
	if err != nil {
		t.Fatal(err)
	}
	nodeId, err := getNodeId(req.URL)
	assert.Nil(t, err)
	assert.Equal(t, 123, nodeId)

	// Тестирование с некорректным запросом (nodeId не является числом)
	req, err = http.NewRequest("GET", "/?nodeId=abc", nil)
	if err != nil {
		t.Fatal(err)
	}
	nodeId, err = getNodeId(req.URL)
	assert.NotNil(t, err)
	assert.Equal(t, 0, nodeId)

	// Тестирование с пустым запросом
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	nodeId, err = getNodeId(req.URL)
	assert.NotNil(t, err)
	assert.Equal(t, 0, nodeId)
}
