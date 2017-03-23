package line

import (
	"testing"

	"code.olipicus.com/trueselect_coupon/config"
)

func TestGetCoupon(t *testing.T) {
	config := config.LoadConfig("../config.json")
	channelSecret, _ := config.GetString("channel_secret")
	channelToken, _ := config.GetString("channel_token")

	app, err := NewLineApp(channelSecret, channelToken, config)
	if err != nil {
		t.Fatal("Create LineApp Error: ", err.Error())
	}

	coupon, err := app.getCoupon("LINE_USER_ID")

	if err != nil {
		t.Fatal("Get Coupon Error:", err.Error())
	}

	if coupon == "" {
		t.Fatal("Coupon is empty")
	}
}
