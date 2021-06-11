package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/justincremer/discord-bot/pkg/commands"
	"github.com/justincremer/discord-bot/pkg/logger"
)

var (
	Token string
	BotID string
	t0    time.Time
	err   error
)

func init() {
	Token = ""
	t0 = time.Now()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	logger.Must("Error creating discord session.", err)
	logger.WriteInfo("Session Created.")

	u, err := dg.User("@me")
	logger.Must("A problem occurred while obtaining account details.", err)

	BotID = u.ID

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)

	err = dg.Open()
	logger.Must("A problem occurred while opening a connection.", err)

	logger.WriteInfo("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s == nil || m == nil {
		return
	}

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "" {
		return
	}

	if m.Content[0] == '!' && strings.Count(m.Content, "!") < 2 {
		commands.ExecuteCommand(s, m.Message, t0)
		return
	}
}

func messageReactionAdd(s *discordgo.Session, reactMsg *discordgo.MessageReactionAdd) {
	_, err := s.ChannelMessage(reactMsg.ChannelID, reactMsg.MessageID)
	logger.Must("A problem occurred while getting a message.", err)
}
