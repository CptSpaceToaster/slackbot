package robots

import (
	"bufio"
	"bytes"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type BabbleBot struct {
}

var dict = make(map[string][]string)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	order := 0
	context := [2]string{}

	w := &BabbleBot{}
	RegisterRobot("babble", w)

	regex, err := regexp.Compile("([\\w,’'\"-;]+)|([\\.?!])")
	if err != nil {
		log.Println("Error in regex")
		return // there was a problem with the regular expression.
	}

	fh, err := os.Open("dickens.txt")
	if err != nil {
		log.Println("It's not there")
		return
	}
	defer fh.Close() // f.Close will run when we're finished.

	word_count := 0
	f := bufio.NewReader(fh)
	buf := make([]byte, 1024)
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			//Usually an EOF
			break
		}

		line := regex.FindAllString(string(buf), -1)
		for _, word := range line {
			word_count++

			context[order%len(context)] = word
			order++

			var buffer bytes.Buffer

			for i := 0; i < len(context)-1; i++ {
				if context[(i+order)%len(context)] != "" {
					buffer.WriteString(context[(i+order)%len(context)])
					if i+2 != len(context) {
						buffer.WriteByte(' ')
					}
				}
			}
			//if word_count < 20 {
			//	log.Println("<" + buffer.String() + ">[" + word + "]")
			//}
			dict[buffer.String()] = append(dict[buffer.String()], word)
		}
	}
	log.Printf("Registerd %d words from dickens.txt", word_count)

	//s := regexp.MustCompile("\\s+").Split("This is a sentence", -1)
	//log.Println(scanner.Text())

	//split the text held by scanner by with regex by spaces
	//s := regexp.MustCompile("[^\\s]+").Split(scanner.Text())
	//log.Println(scanner.Text())

}

func (w BabbleBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go w.DeferredAction(p)
	return ""
}

func (w BabbleBot) DeferredAction(p *Payload) {
	order := 0
	context := [2]string{}

	var ret bytes.Buffer

	text := strings.TrimSpace(p.Text)
	//if text != "" {
	response := &IncomingWebhook{
		Domain:      p.TeamDomain,
		Channel:     p.ChannelID,
		Username:    "Babble Bot",
		IconEmoji:   ":ghost:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
	}
	/*
		resp, err := http.Get(fmt.Sprintf("http://www.google.com/search?q=(site:en.wikipedia.org+OR+site:ja.wikipedia.org)+%s&btnI", url.QueryEscape(text)))
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			message := fmt.Sprintf("ERROR: Non-200 Response from Google: %s", resp.Status)
			log.Println(message)
			response.Text = fmt.Sprintf("%s", message)
		} else if err != nil {
			response.Text = fmt.Sprintf("%s", "Error getting wikipedia link from google :(")
		} else {
			response.Text = fmt.Sprintf("%s", resp.Request.URL.String())
		}
		response.Send()
	*/
	order = len(context) - 1
	if text == "" {
		//TODO: arg parsing, and check to see if the thing is in the dictionary
	}
	context[0] = "The"
	//context[1] = "lion"
	context[1] = ""
	ret.WriteString(context[0])

	for {
		var buffer bytes.Buffer
		for i := 0; i < len(context)-1; i++ {
			if context[(i+order+1)%len(context)] != "" {
				buffer.WriteString(context[(i+order+1)%len(context)])
				if i+2 != len(context) {
					buffer.WriteByte(' ')
				}
			}
		}

		s := dict[buffer.String()][rand.Intn(len(dict[buffer.String()]))]
		context[order%len(context)] = s
		ret.WriteString(" " + s)
		if s[len(s)-1] == '.' {
			break
		}

		order++
	}
	response.Text = ret.String()
	response.Send()
	//}
}

func (w BabbleBot) Description() (description string) {
	return "Babble bot!\n\tUsage: !babble [topic]\n\tExpected Response: Poor English"
}
