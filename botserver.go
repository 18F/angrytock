/*
Script for starting the bot server
*/

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/geramirez/tock-bot/bot"
	"github.com/robfig/cron"
)

func main() {

	bot := bot.InitBot()
	bot.StoreSlackUsers()
	bot.SlapLateUsers()

	go bot.BotherSlackUsers()

	// Slap tock users a couple times
	c := cron.New()
	c.AddFunc("@weekly", func() {
		bot.StoreSlackUsers()
	})
	c.AddFunc("0 0 0 * * 1", func() {
		bot.SlapLateUsers()
	})
	c.Start()

	// Start server
	log.Print("Starting server on port :" + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
