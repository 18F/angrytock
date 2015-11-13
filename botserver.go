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

	c := cron.New()
	c.AddFunc("*/5 * * * * *", func() {
		log.Println("Slapping Users")
		bot := bot.InitBot()
		bot.SlapLateUsers()
		bot.FetchSlackUsers()
	})
	c.Start()

	// Start server
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
