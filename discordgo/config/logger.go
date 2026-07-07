package config

import "log"

func Info(msg string) {
  log.Println("[INFO]", msg)
}

func Error(err error) {
  log.Println("[ERROR]", err)
}