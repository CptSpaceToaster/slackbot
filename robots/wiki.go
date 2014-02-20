package robots

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type WikiBot struct {
}

func init() {
	RegisterRobot("/wiki", func() (robot Robot) { return new(WikiBot) })
}

func (w WikiBot) Run(command *SlashCommand) (slashCommandImmediateReturn string) {
	go w.DeferredAction(command)
	return ""
}

func (w WikiBot) DeferredAction(command *SlashCommand) {
    text := strings.TrimSpace(command.Text)
    if text != "" {
        response := new(IncomingWebhook)
        response.Channel = command.Channel_ID
        response.Username = "Wiki Bot"
        response.Text = fmt.Sprintf("@%s: Searching google for wikis relating to: %s", command.User_Name, text)
        response.Icon_Emoji = ":ghost:"
        response.Unfurl_Links = 1
        response.Parse = "full"
        MakeIncomingWebhookCall(response)
        resp, err := http.Get(fmt.Sprintf("http://www.google.com/search?q=(site:en.wikipedia.org+OR+site:ja.wikipedia.org)+%s&btnI", url.QueryEscape(text)))                                   
        if err != nil {
            response.Text = fmt.Sprintf("@%s: %s", command.User_Name, "Error getting wikipedia link from google :(")
        } else {
            response.Text = fmt.Sprintf("@%s: %s", command.User_Name, resp.Request.URL.String())
        }
        MakeIncomingWebhookCall(response)
    }
}

func (w WikiBot) Description() (description string) {
	return "Wiki bot!\n\tUsage: /wiki\n\tExpected Response: @user: Link to wikipedia article!"
}