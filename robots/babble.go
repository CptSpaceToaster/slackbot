package robots

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type BabbleBot struct {
}

type BabbleConfig struct {
	Source string `json:"source"`
}

// Don't worry about this slice of Dictionaries that are keyed with strings, and reference slices of strings
var Dict = make(map[string][]string)

// Local Bot config
var BotConfig = new(BabbleConfig)

//var Dicts = make([]map[string][]string)

// The regexp we use
var WordRegex = regexp.MustCompile("([\\w’';,]+)|([\\.?!\"])")
var PuncRegex = regexp.MustCompile("[\\.,?!]")
var CommaRegex = regexp.MustCompile(",")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "babble.json")

	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("ERROR: Error opening babblebot's config: %s", err)
			return
		}
		err = json.Unmarshal(config, BotConfig)
		if err != nil {
			log.Printf("ERROR: Error parsing babblebot's config: %s", err)
			return
		}
	} else {
		log.Printf("WARNING: Could not find configuration file babble.json in %s", *ConfigDirectory)
	}

	index := 0
	context := [2]string{}

	w := &BabbleBot{}
	RegisterRobot("babble", w)

	// TODO: Ability to read more than one file (avoid symlinks)
	fh, err := os.Open(filepath.Join(*ConfigDirectory, BotConfig.Source))
	if err != nil {
		log.Println("Bot source file not found")
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

		line := WordRegex.FindAllString(string(buf), -1)
		for _, word := range line {
			word_count++

			context[index%len(context)] = word
			index++

			var buffer bytes.Buffer

			for i := 0; i < len(context)-1; i++ {
				if context[(i+index)%len(context)] != "" {
					buffer.WriteString(context[(i+index)%len(context)])
					if i+2 != len(context) {
						buffer.WriteByte(' ')
					}
				}
			}
			//if word_count < 20 {
			//	log.Println("<" + buffer.String() + ">[" + word + "]")
			//}
			Dict[buffer.String()] = append(Dict[buffer.String()], word)
		}
	}
	log.Printf("Registerd %d words from %s", word_count, BotConfig.Source)

	//s := regexp.MustCompile("\\s+").Split("This is a sentence", -1)
	//log.Println(scanner.Text())

	//split the text held by scanner by with WordRegex by spaces
	//s := regexp.MustCompile("[^\\s]+").Split(scanner.Text())
	//log.Println(scanner.Text())

}

func (w BabbleBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go w.DeferredAction(p)
	return ""
}

func (w BabbleBot) DeferredAction(p *Payload) {
	//order := 0
	index := 0
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
	index = len(context) - 1
	if text == "" {
		//TODO: arg parsing, and check to see if the thing is in the Dictionary
	}
	context[0] = "The"
	//context[1] = "lion"
	context[1] = ""
	ret.WriteString(context[0])

	// TODO: These have to be set if the user types in a quote
	// Keeps track of whether or not we have an open quote.  Nested quotes are not possible atm
	quote := false
	quote_flag := false
	punc_flag := false
	last_word := ""

	for {
		var buffer bytes.Buffer
		for i := 0; i < len(context)-1; i++ {
			if context[(i+index+1)%len(context)] != "" {
				buffer.WriteString(context[(i+index+1)%len(context)])
				if i+2 != len(context) {
					buffer.WriteByte(' ')
				}
			}
		}

		// Obtain word
		s := Dict[buffer.String()][rand.Intn(len(Dict[buffer.String()]))]

		// Immediatly exit if we chose another puncuation mark when the punctuation flag was set
		if PuncRegex.MatchString(s) && punc_flag == true {
			continue
		} else {
			punc_flag = true
		}

		//TODO: Make this a case statement, and handle quotes
		//TODO: Quotes can follow commas, and they look weird if they close
		switch s {
		case "\"":
			// If we JUST started a quote
			if quote_flag {
				continue
			}

			if quote {
				// If the last word was a comma, we actually DON'T want to put a quote there
				if CommaRegex.MatchString(last_word) {
					continue
				} else {
					ret.WriteString(s)
				}
			} else {
				// If we open a quote, set the quote flag
				ret.WriteString(" " + s)
				quote_flag = true
			}
			// Toggle whether we are in a quote or not. TODO: Quotes can't be nested
			quote = !quote
			// Set the puncuation flag
			punc_flag = true
		case ".":
			ret.WriteString(s)
		case "?":
			ret.WriteString(s)
		case "!":
			ret.WriteString(s)
		default:
			if quote_flag {
				// If the quote_flag is set, then don't print a leading space, and consume the flag
				ret.WriteString(s)
				quote_flag = false
			} else {
				ret.WriteString(" " + s)
			}
			// Consume the puncuation flag
			punc_flag = false
		}

		// Accept the new word
		context[index%len(context)] = s

		// End sequence
		if s == "." {
			if quote {
				// Close a dangeling quote
				ret.WriteString("\"")
			}
			break
		}
		last_word = s
		index++
	}
	response.Text = ret.String()
	response.Send()
	//}
}

func (w BabbleBot) Description() (description string) {
	return "Babble bot!\n\tUsage: !babble [topic]\n\tExpected Response: Poor English"
}
