package demo

import (
	sl "log/slog"
)

func Demo() {
	//password := "123"
	token := "APSP0DSAD#@DS"
	//apiKey := "1asdsd213asf"
	sl.Info("Server started") // алиас
	logger := sl.Default()
	logger.Info("Server started") // метод
	log := sl.Default()
	log.Info("a" + "b" + "c")
	log.Info("Server " + "started")
	log.Info("token: " + token)
	log.Info("сервер " + "запущен")

}
