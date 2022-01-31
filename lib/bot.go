package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/envfile"
	"github.com/yanzay/tbot/v2"
	"gorm.io/gorm"
)

func Log(s string, p ...interface{}) {
	log.Printf(s, p...)
}

var userpass string = ""

func H(hn string, h func(in *tbot.Message) (ret string, err error)) {
	c := bot.Client()

	bot.HandleMessage(hn, func(m *tbot.Message) {

		Log("Got msg: %#v", m)

		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)
		if m.From.IsBot {
			time.Sleep(time.Minute * 5)
			return
		}

		var ret string
		var err error
		switch {
		case hn == "/login.*" || hn == "/start":
			ret, err = h(m)
		default:
			sess := &Session{ChatID: m.Chat.ID}
			err := db.Where(sess).First(sess).Error
			if err == nil {
				ret, err = h(m)
			}
			if err == gorm.ErrRecordNotFound {
				ret = "Must login first - call /login <pass>"
				err = nil
			}

		}

		if err != nil {
			Log("Replying error: %#", err.Error())
			c.SendMessage(m.Chat.ID, "Error: "+err.Error())
			return
		}
		Log("Replying: %s", ret)
		c.SendMessage(m.Chat.ID, ret)
	})
}

var bot *tbot.Server

func Pub(ch string, msg string) (ret string, err error) {
	subs := make([]Sub, 0)
	err = db.Where(&Sub{Chan: ch}).Find(&subs).Error
	if err != nil {
		return
	}
	c := bot.Client()
	sb := strings.Builder{}
	for _, v := range subs {
		_, err = c.SendMessage(v.ChatID, fmt.Sprintf("<b>%s</b>: %s", ch, msg), tbot.OptParseModeHTML)
		if err != nil {
			sb.WriteString("* ")
			sb.WriteString(v.FirstName)
			sb.WriteString(" ")
			sb.WriteString(v.LastName)
			sb.WriteString(": ")
			sb.WriteString(err.Error())
			sb.WriteString("\n")
			Log(err.Error())
			err = nil
		} else {
			sb.WriteString("* ")
			sb.WriteString(v.FirstName)
			sb.WriteString(" ")
			sb.WriteString(v.LastName)
			sb.WriteString(": OK\n")
		}
	}
	ret = sb.String()
	return
}

func InitDB() error {
	db = &GDB{}
	err := db.Init()
	if err != nil {
		return err
	}

	err = db.AutoMigrates(&Sub{}, &Session{})
	if err != nil {
		return err
	}

	return nil
}

func InitBot() error {
	Log("Initiating Telegrambot - v 20220116A")
	var err error
	envfile.Load()
	bot = tbot.New(envfile.Must("TBOTKEY"))
	userpass = envfile.Must("USERPASS")

	Log("Setting up the pipes")

	H("/start", func(m *tbot.Message) (ret string, err error) {
		return "Hello 👽, please login.", nil
	})

	H("/login.*", func(m *tbot.Message) (ret string, err error) {
		parts := strings.Split(m.Text, " ")
		if len(parts) != 2 {
			ret = "usage: /login <password>"
			return
		}

		pass := parts[1]
		if pass == userpass {
			s := &Session{ChatID: m.Chat.ID}
			err = db.Save(s).Error
			ret = "Login ok"
		} else {
			time.Sleep(time.Second * 15)
		}

		return
	})

	H("/logout", func(m *tbot.Message) (ret string, err error) {
		s := &Session{ChatID: m.Chat.ID}
		err = db.Where(s).Delete(s).Error
		ret = "Logout ok"
		return

	})

	H("/pub.*", func(m *tbot.Message) (ret string, err error) {
		parts := strings.Split(m.Text, " ")
		if len(parts) < 3 {
			ret = "usage: /pub <channel name> message"
			return
		}
		ch := parts[1]
		msg := strings.Join(parts[2:], " ")
		return Pub(ch, msg)

	})

	H("/subbers.*", func(m *tbot.Message) (ret string, err error) {
		parts := strings.Split(m.Text, " ")
		if len(parts) != 2 {
			ret = "usage: /subbers <channel name>"
			return
		}
		ch := parts[1]

		subs := make([]Sub, 0)
		err = db.Where(&Sub{Chan: ch}).Find(&subs).Error
		if err != nil {
			return
		}

		sb := strings.Builder{}
		sb.WriteString("Subbers of [*")
		sb.WriteString(ch)
		sb.WriteString("*]: ")
		if len(subs) < 1 {
			sb.WriteString("NO ONE")
		}
		for _, v := range subs {
			sb.WriteString("\n * ")
			sb.WriteString(v.FirstName)
			sb.WriteString(" ")
			sb.WriteString(v.LastName)
			err = nil

		}
		ret = sb.String()
		return

	})

	H("/sub.*", func(m *tbot.Message) (ret string, err error) {
		chname := strings.Replace(m.Text, "/sub ", "", 1)
		chname = strings.TrimSpace(chname)
		s := &Sub{
			Chan:      chname,
			ChatID:    m.Chat.ID,
			Username:  m.From.Username,
			FirstName: m.From.FirstName,
			LastName:  m.From.LastName,
			JoinnedAt: time.Now(),
		}
		err = db.Save(s).Error
		return "Ok - subscribed " + chname, err

	})

	H("/unsub.*", func(m *tbot.Message) (ret string, err error) {
		chname := strings.Replace(m.Text, "/unsub ", "", 1)
		chname = strings.TrimSpace(chname)
		s := &Sub{
			Chan:   chname,
			ChatID: m.Chat.ID,
		}
		err = db.Where(s).Delete(s).Error
		return "Ok - subscribed " + chname, err
	})

	H("/mysubs", func(m *tbot.Message) (ret string, err error) {
		s := &Sub{
			ChatID: m.Chat.ID,
		}
		subs := make([]Sub, 0)
		err = db.Where(s).Find(&subs).Error
		if err != nil {
			return
		}

		sb := strings.Builder{}

		if len(subs) < 1 {
			sb.WriteString("No subs for you baby! 😘")
		}

		for _, v := range subs {
			sb.WriteString("* ")
			sb.WriteString(v.Chan)
			sb.WriteString("\n")
		}
		return sb.String(), err
	})

	H("/help", func(m *tbot.Message) (ret string, err error) {
		return `/sub to subscribe a channel
/unsub to unsubscribe
/mysubs to list user own subscriptions
/subbers to list all users subscribing a channel
/pub will send a message to a channel
`, nil
	})

	H(".*", func(m *tbot.Message) (ret string, err error) {
		ret = "Dont know what to do with: " + m.Text
		return
	})

	Log("Starting bot")
	err = bot.Start()
	if err != nil {
		Log(err.Error())
	}

	return nil
}

func InitHttp() error {

	apikey := envfile.Must("APIKEY")

	http.HandleFunc("/pub", func(rw http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != apikey {
			time.Sleep(time.Minute)
			return
		}
		pm := &PubMsg{}
		err := json.NewDecoder(r.Body).Decode(pm)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		ret, err := Pub(pm.Ch, pm.Msg)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Write([]byte(ret))

	})

	Log("Listening http")
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			Log("Error listening: %s", err.Error())
		}
	}()

	return nil
}

func Init() error {
	Log("Initiating Telegrambot - v 20220116A")
	envfile.Load()

	err := InitDB()
	if err != nil {
		return err
	}

	err = InitHttp()
	if err != nil {
		return err
	}
	err = InitBot()
	if err != nil {
		return err
	}
	return nil
}
