# ODT-Discord-Bot

ODT Discord Bot is a lightweight and fast Discord bot built for our **GTA 1.27 Legacy Edition server**. It provides real-time server statistics and simple utility commands for players and staff inside Discord.

## Features

- Player statistics (`?playercount`)
- Help system (`?help`)
- Crew management system (`?crewhelp`)
- Crew status tracking (`?crewstatus`)
- Crew information lookup (`?crewinfo`)
- Crew creation system (`?crewcreate <name>`)
- Crew invite codes (`?crewcode`)
- Crew joining system (`?crewjoin <code>`)
- Crew leaving system (`?crewleave`)
- Crew deletion system (`?crewdelete`)
- Custom crew emblem support (`?crewemblem`)
- Crew listing system (`?crewlist`)
- Player information lookup (`?playerinfo <player>`)
- Server information (`?info`)
- Account verification (`?verify`)
- Important links (`?link`)
- Early supporter claiming (`?claim`)
- Staff management tools (`?staffhelp`)

## Requirements
- Go 1.18+
- DiscordGo library
  
```
cd Discord Go
go mod tidy
go run .
```
## Discord Permissions

Your bot must have the following permissions:

- Read Messages
- Send Messages
- Embed Links
- Manage Messages (optional, for deleting command messages)

Also enable in the Discord Developer Portal:

**Message Content Intent**

## Purpose

This bot is designed specifically for helping the server display live game statistics and improve player engagement through Discord.

## Notes

- Prefix-based commands (`?`)
- Lightweight and fast performance
- Easy to extend with new features
- Uses Discord embeds for clean output
