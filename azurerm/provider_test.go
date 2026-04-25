package azurerm

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
)

func TestFixNetworkInterfaceAcceleratedNetworking(t *testing.T) {
	const networkInterfaceID = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic-1"

	t.Run("overrides state value from azure api cache", func(t *testing.T) {
		a := &azurerm{
			networkInterfaceAcceleratedNetworking: map[string]bool{
				strings.ToLower(networkInterfaceID): true,
			},
		}

		input := cty.ObjectVal(map[string]cty.Value{
			"id":                            cty.StringVal(networkInterfaceID),
			"enable_accelerated_networking": cty.BoolVal(false),
			"name":                          cty.StringVal("nic-1"),
		})

		got, err := a.fixNetworkInterfaceAcceleratedNetworking(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var enabled bool
		if err := gocty.FromCtyValue(got.GetAttr("enable_accelerated_networking"), &enabled); err != nil {
			t.Fatalf("failed to decode bool value: %v", err)
		}

		if !enabled {
			t.Fatalf("expected enable_accelerated_networking=true, got false")
		}
	})

	t.Run("adds attribute when missing in terraform state", func(t *testing.T) {
		a := &azurerm{
			networkInterfaceAcceleratedNetworking: map[string]bool{
				strings.ToLower(networkInterfaceID): true,
			},
		}

		input := cty.ObjectVal(map[string]cty.Value{
			"id":   cty.StringVal(networkInterfaceID),
			"name": cty.StringVal("nic-1"),
		})

		got, err := a.fixNetworkInterfaceAcceleratedNetworking(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !got.Type().HasAttribute("enable_accelerated_networking") {
			t.Fatalf("expected enable_accelerated_networking attribute to be present")
		}

		var enabled bool
		if err := gocty.FromCtyValue(got.GetAttr("enable_accelerated_networking"), &enabled); err != nil {
			t.Fatalf("failed to decode bool value: %v", err)
		}

		if !enabled {
			t.Fatalf("expected enable_accelerated_networking=true, got false")
		}
	})

	t.Run("keeps state unchanged when id is not cached", func(t *testing.T) {
		a := &azurerm{
			networkInterfaceAcceleratedNetworking: map[string]bool{},
		}

		input := cty.ObjectVal(map[string]cty.Value{
			"id":                            cty.StringVal(networkInterfaceID),
			"enable_accelerated_networking": cty.BoolVal(false),
		})

		got, err := a.fixNetworkInterfaceAcceleratedNetworking(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !got.RawEquals(input) {
			t.Fatalf("expected state to remain unchanged when NIC id is not cached")
		}
	})
}
