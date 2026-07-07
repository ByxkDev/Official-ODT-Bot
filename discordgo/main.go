package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" 

	"discordgo/config"
	"discordgo/utils"
	"discordgo/discord"
)

func main() {

	config.LoadEnvironment()
	cfg := config.LoadConfig()

	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	utils.SetDB(db)

	bot, err := discord.New(cfg.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Discord bot started")
	select {}
}