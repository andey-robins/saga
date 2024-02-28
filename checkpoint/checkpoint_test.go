package checkpoint

import (
	"os"
	"testing"
)

func TestSave(t *testing.T) {
	cname := os.TempDir() + "/test.json"
	data := "test data"

	Save(cname, data)

	var s string

	Load(cname, &s)

	if s != data {
		t.Errorf("Expected %s, got %s", data, s)
	}
}
