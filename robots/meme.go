package robots

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type MemeBot struct {
}

func init() {
	w := &MemeBot{}
	RegisterRobot("meme", w)
}

func (w MemeBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go w.DeferredAction(p)
	return ""
}

func (w MemeBot) DeferredAction(p *Payload) {
	args := strings.Split(strings.TrimSpace(p.Text), ";")

	var meme Attachment
	var text string
	var user string
	var icon_url string

	if len(args) >= 3 {
		if Config.Tokens[p.TeamDomain].APIToken != "" {
			response, err := http.Get("https://slack.com/api/users.info?token=" + Config.Tokens[p.TeamDomain].APIToken + "&user=" + p.UserID)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			var resp UserResp
			err = json.Unmarshal(body, &resp)
			if err != nil {
				fmt.Println(err)
				return
			}

			user = resp.User.Name
			icon_url = resp.User.Profile.Image48
		} else {
			user = "Meme Bot"
		}

		text = "http://memegen.link/" + strings.TrimSpace(args[0]) + "/" + strings.Replace(strings.TrimSpace(args[1]), " ", "-", -1) + "/" + strings.Replace(strings.TrimSpace(args[2]), " ", "-", -1)
		meme.ImageURL = text + ".jpg"
	} else {
		user = "Meme Bot"
		meme.ImageURL = ""
		text = "Error: I need 3 arguments separated by semicolons " + p.UserName + "\nExample `/meme doge; such memes; much wow`"
	}

	response := &IncomingWebhook{
		Domain:   p.TeamDomain,
		Channel:  p.ChannelID,
		Username: user,
		Text:     text,
		Attachments: []Attachment{
			meme,
		},
		IconURL:     icon_url,
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}

	response.Send()
}

func (w MemeBot) Description() (description string) {
	return "Meme bot!\n\tUsage: /meme Template; Top Text; Bottom Text\n\tExpected Response: Dank Memes"
}
