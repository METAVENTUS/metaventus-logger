package logger

import (
	"github.com/joho/godotenv"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger configure le logger zerolog
func InitLogger() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msgf("Pas de fichier .env trouvé: %v", err)
	}

	zerolog.TimeFieldFormat = time.RFC3339

	var multiWriter io.Writer
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: true}

	// Exemple de configuration pour un environnement de production
	env := os.Getenv("SERVICE_MODE")
	switch env {
	case "release":
		multiWriter = consoleWriter
	default:
		// Créez ou ouvrez le fichier log pour le mode développement
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Erreur lors de l'ouverture du fichier log")
		}
		multiWriter = io.MultiWriter(consoleWriter, file)
	}

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()

	// Définir le niveau global de log à partir de la variable d'environnement LOG_LEVEL
	logLevel := getLogLevel()
	zerolog.SetGlobalLevel(logLevel)

	log.Info().Msgf("Logger initialisé avec le niveau %s", logLevel.String())
}

// getLogLevel retourne le niveau de log à partir de la variable d'environnement LOG_LEVEL
func getLogLevel() zerolog.Level {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}
