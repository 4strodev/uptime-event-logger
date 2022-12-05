package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	command *discordgo.ApplicationCommand = &discordgo.ApplicationCommand{
		Name:        "basic",
		Description: "A simple discord command",
	}
)

func main() {
	godotenv.Load()

	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		log.Fatal("Discord token not set!")
	}
	dg, err := discordgo.New(discordToken)
	if err != nil {
		log.Fatal("Cannot connect to discord:", err)
	}

	//dg.Identify.Intents |= discordgo.IntentMessageContent

	//dg.AddHandler(messageCreate)
	dg.AddHandler(commands)

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection", err)
	}

	dg.ApplicationCommandCreate(dg.State.User.ID, "", command)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-signalChannel
	secureShutdown(dg)
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// If the message sender id is the same as bot, means that this message was sent by our bot.
	if message.Author.ID == session.State.User.ID {
		return
	}

	messageAuthor := message.Author.Username
	log.Printf("%s Sent a new message\n", messageAuthor)
	response := fmt.Sprintf("@%s! You sent this message `%s`", messageAuthor, message.Content)
	_, err := session.ChannelMessageSend(message.ChannelID, response)
	if err != nil {
		log.Println(err)
	}
}

func commands(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	log.Println("Basic command")
	if interaction.ApplicationCommandData().Name == command.Name {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there! Congratulations, you just executed your first slash command",
			},
		})
	}
}

func secureShutdown(dg *discordgo.Session) {
	// Closing sockets. Also is a good idea close databse connections, etc.
	dg.Close()
}
