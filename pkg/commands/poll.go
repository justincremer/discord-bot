package commands

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/justincremer/discord-bot/pkg/logger"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type newPollRequest struct {
	Title    string   `json:"title"`
	Options  []string `json:"options"`
	Multi    bool     `json: "multi"`
	DupCheck string   `json: "dupcheck"`
	Captcha  bool     `json: "captcha"`
}

type newPollResponse struct {
	Id       int      `json: "id"`
	Title    string   `json: "title"`
	Options  []string `json: "options"`
	Multi    bool     `json: "multi"`
	DupCheck string   `json: "dupcheck"`
	Captcha  bool     `json: "captcha"`
}

// HandlePollCommand controls the !poll command.
// Given a topic, the function will send a message with said topic, along with 3 initial reactions to allow voting.
func HandlePollCommand(s *discordgo.Session, m *discordgo.Message, pollTopic string) {
	message, err := s.ChannelMessageSend(m.ChannelID, pollTopic)
	logger.Must("Failed to send message: ", err)

	go s.MessageReactionAdd(m.ChannelID, message.ID, "👍")
	go s.MessageReactionAdd(m.ChannelID, message.ID, "🍴")
	go s.MessageReactionAdd(m.ChannelID, message.ID, "👎")
}

// HandleStrawPollCommand controls the !strawpoll command
// !strawpoll {title} option1 option2 option3...
// Given a title and options, the command creates a post request to the strawpoll.me API and returns a link to
// the newly created poll
func HandleStrawPollCommand(s *discordgo.Session, m *discordgo.Message, pollTitle string, pollOptions []string) {
	request := &newPollRequest{
		Title:    pollTitle,
		Options:  pollOptions,
		DupCheck: "normal",
	}
	reqData, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", "https://www.strawpoll.me/api/v2/polls", bytes.NewBuffer(reqData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := client.Do(req)
	logger.Must("Poll client failed: ", err)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var newPoll newPollResponse
	err = json.Unmarshal(body, &newPoll)
	logger.Must("Failed to parse JSON: ", err)

	go s.ChannelMessageSend(m.ChannelID, "https://www.strawpoll.me/"+strconv.Itoa(newPoll.Id))
}
