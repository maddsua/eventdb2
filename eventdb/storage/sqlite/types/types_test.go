package types_test

import (
	"fmt"
	"testing"

	"github.com/maddsua/eventdb2/storage/model"
	"github.com/maddsua/eventdb2/storage/sqlite/types"
)

func TestLabels(t *testing.T) {

	labels := model.StringMap{
		"user_id": "12345",
		"sus":     "true",
		"0":       "1",
		"________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________": "false",
	}

	data, err := types.EncodeStringMap(labels)
	if err != nil {
		t.Fatalf("encoding failed: %v", err)
	}

	decoded, err := types.DecodeStringMap(types.NullBlobSlice(data))
	if err != nil {
		t.Fatalf("decoding failed: %v", err)
	}

	var expectResult = func(key string, expect string) error {

		if val := decoded[key]; val != expect {
			return fmt.Errorf("%s invalid; expected: '%s'; got: '%s'", key, expect, val)
		}

		return nil
	}

	tests := []error{
		expectResult("user_id", "12345"),
		expectResult("sus", "true"),
		expectResult("0", "1"),
	}

	for _, err := range tests {
		if err != nil {
			t.Fatal("decoder failed", err)
		}
	}
}
