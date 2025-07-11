package logger

import "github.com/sirupsen/logrus"

var L = logrus.New()

func Init(cfg Config) {
	L.SetLevel(cfg.Level)
	L.SetFormatter(&logrus.JSONFormatter{})
}