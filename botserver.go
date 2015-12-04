/*
Script for starting the bot server
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/18F/tock-bot/bot"
	"github.com/robfig/cron"
)

func main() {

	bot := bot.InitBot()
	bot.StoreSlackUsers()

	// Update the list of stored slack users weekly
	c := cron.New()
	c.AddFunc("@weekly", func() {
		bot.StoreSlackUsers()
	})
	c.Start()

	// Start go routine to listen to tock users
	go bot.ListenToSlackUsers()

	// Start server
	log.Print("Starting server on port :" + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
