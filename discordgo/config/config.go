package config

import "os"

type Config struct 
{
	DiscordToken string
	DBDSN string
}

func LoadConfig() Config {
	return Config{
		
	DiscordToken: os.Getenv("DISCORD_TOKEN"),
	DBDSN: os.Getenv("DB_DSN"),
    }
}