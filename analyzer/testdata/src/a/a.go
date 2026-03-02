package a

import sl "log/slog"

func f() {
	token := "APS..."
	log := sl.Default()

	sl.Info("Server started")  // want "lowercase"
	log.Info("Server started") // want "lowercase"

	log.Info("token: " + token) // want "private info"
	log.Info("сервер запущен")  // want "must be in English"
	log.Info("server started!") // want "special symbols"
}
