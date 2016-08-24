package conf

import (
	"os"

	"github.com/gen1us2k/log"
	"github.com/urfave/cli"
)

// Version stores current service version
var (
	Version            string
	TelegramBotToken   string
	TelegramWebhookURL string
	HTTPBindAddr       string
	LogLevel           string
)

type Configuration struct {
	data *CinematicaConfig
	app  *cli.App
}

// NewConfigurator is constructor and creates a new copy of Configuration
func NewConfigurator() *Configuration {
	Version = "0.1dev"
	app := cli.NewApp()
	app.Name = "Cinematica 2 telegram"
	app.Usage = "Get permieres from cinematica.kg and show to telegram users"
	return &Configuration{
		data: &CinematicaConfig{},
		app:  app,
	}
}

func (c *Configuration) fillConfig() *CinematicaConfig {
	return &CinematicaConfig{
		TelegramBotToken:   TelegramBotToken,
		TelegramWebhookURL: TelegramWebhookURL,
		HTTPBindAddr:       HTTPBindAddr,
	}
}

// Run is wrapper around cli.App
func (c *Configuration) Run() error {
	c.app.Before = func(ctx *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}
	c.app.Flags = c.setupFlags()
	return c.app.Run(os.Args)
}

// App is public method for Configuration.app
func (c *Configuration) App() *cli.App {
	return c.app
}

func (c *Configuration) setupFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "http_bind_addr",
			Value:       ":8090",
			Usage:       "Set address to bind http server",
			EnvVar:      "HTTP_BIND_ADDR",
			Destination: &HTTPBindAddr,
		},
		cli.StringFlag{
			Name:        "telegram_bot_token",
			Value:       "",
			Usage:       "Set telegram bot access token",
			EnvVar:      "TELEGRAM_TOKEN",
			Destination: &TelegramBotToken,
		},
		cli.StringFlag{
			Name:        "telegram_web_hook_url",
			Value:       "",
			Usage:       "Set telegram bot webhook url",
			EnvVar:      "TELEGRAM_WEBHOOK_URL",
			Destination: &TelegramWebhookURL,
		},
		cli.StringFlag{
			Name:        "loglevel",
			Value:       "debug",
			Usage:       "set log level",
			Destination: &LogLevel,
			EnvVar:      "LOG_LEVEL",
		},
	}

}

// Get returns filled BillginConfig
func (c *Configuration) Get() *CinematicaConfig {
	c.data = c.fillConfig()
	return c.data
}
