package line

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

//LineApp : struct
type LineApp struct {
	bot *linebot.Client
}

//NewLineApp : new line app
func NewLineApp(channelSecret string, channelToken string) (*LineApp, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineApp{
		bot: bot,
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

	log.Println(profile.UserID)

	app.replyText(event.ReplyToken, "หวัดดี")
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
