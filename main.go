package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"message-expander/events"
);

func Env_load() {
    err := godotenv.Load()
    if err != nil {
       println("Error loading .env file") 
	   os.Exit(1);
    }
}

func main() {
	Env_load();

	dg, err := discordgo.New("Bot "+ os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.IntentsGuildMessages;
	
	dg.AddHandler(events.MessageCreate);
	
	err = dg.Open();
	if 	err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Logged in as", dg.State.User.Username);
	sc := make(chan os.Signal, 1);
	signal.Notify(sc, os.Interrupt, os.Kill);
	
	<-sc;
	dg.Close();
}
