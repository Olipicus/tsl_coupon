package line

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"code.olipicus.com/trueselect_coupon/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

//LineApp : struct
type LineApp struct {
	bot    *linebot.Client
	config *config.Table
}

//NewLineApp : new line app
func NewLineApp(channelSecret, channelToken string, config *config.Table) (*LineApp, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineApp{
		bot:    bot,
		config: config,
	}, nil
}

//CallbackHandler : for Line Webhook
func (app *LineApp) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
			log.Println("Invalid Signature")
			log.Println("X-Line-Signature: " + r.Header.Get("X-Line-Signature"))
		} else {
			w.WriteHeader(500)
			log.Println("Unknow error")
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeFollow {
			app.follow(event)
		}
	}
}

func (app *LineApp) follow(event *linebot.Event) {
	profile, err := app.bot.GetProfile(event.Source.UserID).Do()

	if err != nil {
		log.Fatal("Get Line Profile Error")
	}

	if coupon, err := app.getCoupon(profile.UserID); err == nil {
		msgTemplate, _ := app.config.GetString("coupon_message")
		msg := strings.Replace(msgTemplate, "{coupon}", coupon, -1)

		app.replyText(event.ReplyToken, msg)
	}

}

func (app *LineApp) getCoupon(userID string) (string, error) {
	url, _ := app.config.GetString("get_coupon_url")
	url = strings.Replace(url, "{user_id}", userID, -1)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result resultCoupon
	json.Unmarshal(body, &result)

	return result.CouponCode, nil
}

func (app *LineApp) replyText(replyToken, text string) error {
	if _, err := app.bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}
