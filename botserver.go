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

	c := cron.New()
	c.AddFunc("*/5 * * * * *", func() {
		log.Println("Slapping Users")
		bot.SlapLateUsers()
	})
	c.AddFunc("*/20 * * * * *", func() {
		log.Println("Collecting Users")
		bot.StoreSlackUsers()
	})
	c.Start()

	// Start server
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
