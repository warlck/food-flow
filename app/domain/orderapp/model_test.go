package orderapp_test

import (
	"encoding/json"
	"testing"

	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

func TestToAppOrderEmptySelectionsEncodeAsArrays(t *testing.T) {
	t.Parallel()

	appOrder := orderapp.ToAppOrder(orderbus.Order{
		Items: []orderbus.OrderItem{{}},
	})

	if appOrder.Items[0].Modifiers == nil {
		t.Fatal("modifiers must be a non-nil empty slice")
	}
	if appOrder.Items[0].Addons == nil {
		t.Fatal("addons must be a non-nil empty slice")
	}

	data, _, err := appOrder.Encode()
	if err != nil {
		t.Fatalf("encode order: %v", err)
	}

	var payload struct {
		Items []struct {
			Modifiers []json.RawMessage `json:"modifiers"`
			Addons    []json.RawMessage `json:"addons"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("decode order: %v", err)
	}
	if payload.Items[0].Modifiers == nil || len(payload.Items[0].Modifiers) != 0 {
		t.Fatalf("modifiers must encode as [], got %s", data)
	}
	if payload.Items[0].Addons == nil || len(payload.Items[0].Addons) != 0 {
		t.Fatalf("addons must encode as [], got %s", data)
	}
}
