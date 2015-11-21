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
	log.Println("Collecting Users")
	bot.StoreSlackUsers()
	c := cron.New()
	c.AddFunc("*/30 * * * * *", func() {
		log.Println("Slapping Users")
		bot.SlapLateUsers()
	})
	c.Start()

	// Start server
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
