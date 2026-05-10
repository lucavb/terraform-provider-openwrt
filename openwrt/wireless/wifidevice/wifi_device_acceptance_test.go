//go:build acceptance.test

package wifidevice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	client, providerBlock := runOpenWrtServerWithWireless(
		ctx,
		*dockerPool,
		t,
	)
	options := lucirpc.Options{
		"channel": lucirpc.String("6"),
		"hwmode":  lucirpc.String("11g"),
		"txpower": lucirpc.Integer(15),
		"type":    lucirpc.String("mac80211"),
	}
	ok, err := client.CreateSection(ctx, "wireless", "wifi-device", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_wireless_wifi_device" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "channel", "6"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "hwmode", "11g"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "txpower", "15"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "type", "mac80211"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	_, providerBlock := runOpenWrtServerWithWireless(
		ctx,
		*dockerPool,
		t,
	)

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_wireless_wifi_device" "testing" {
	channel = "6"
	hwmode = "11g"
	id = "testing"
	txpower = 15
	type = "mac80211"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "channel", "6"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "hwmode", "11g"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "txpower", "15"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "type", "mac80211"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_wireless_wifi_device.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_wireless_wifi_device" "testing" {
	band = "6g"
	channel = "11"
	hwmode = "11a"
	htmode = "HE80"
	id = "testing"
	txpower = 20
	type = "mac80211"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "band", "6g"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "channel", "11"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "hwmode", "11a"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "htmode", "HE80"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "txpower", "20"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "type", "mac80211"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
