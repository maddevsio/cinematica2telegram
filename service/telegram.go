package service

import (
	"fmt"
	"strings"

	"github.com/gen1us2k/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"gopkg.in/telegram-bot-api.v4"
)

type TelegramService struct {
	BaseService

	logger       log.Logger
	cinema       *Cinematica
	e            *echo.Echo
	bot          *tgbotapi.BotAPI
	updateChan   chan tgbotapi.Update
	botCommands  map[string]func(*tgbotapi.Message)
	botChatTexts map[string]func(*tgbotapi.Message)
}

func (ts *TelegramService) Name() string {
	return "telegram_service"
}

func (ts *TelegramService) Init(cinema *Cinematica) error {
	ts.cinema = cinema
	ts.logger = log.NewLogger(ts.Name())
	bot, err := tgbotapi.NewBotAPI(ts.cinema.Config().TelegramBotToken)

	if err != nil {
		return err
	}
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(ts.cinema.Config().TelegramWebhookURL))
	if err != nil {
		return err
	}

	ts.bot = bot
	ts.updateChan = make(chan tgbotapi.Update, 100)
	ts.e = echo.New()
	ts.e.POST("/", ts.handleBotRequests)
	ts.e.GET(fmt.Sprintf("/%s", bot.Token), ts.handleBotRequests)
	ts.e.POST(fmt.Sprintf("/%s", bot.Token), ts.handleBotRequests)
	ts.botCommands = make(map[string]func(*tgbotapi.Message))
	ts.botChatTexts = make(map[string]func(*tgbotapi.Message))
	ts.botCommands["/start"] = ts.onStartCommand
	ts.botCommands["/help"] = ts.onHelpCommand

	return nil
}

func (ts *TelegramService) Run() error {
	ts.cinema.waitGroup.Add(1)
	go ts.e.Run(standard.New(ts.cinema.Config().HTTPBindAddr))
	ts.handleRun()

	return nil
}

func (ts *TelegramService) handleBotRequests(c echo.Context) error {
	var update tgbotapi.Update
	if err := c.Bind(&update); err != nil {
		return err
	}
	ts.updateChan <- update
	return nil
}

func (ts *TelegramService) handleRun() {
	for update := range ts.updateChan {
		ts.logger.Infof("%+v\n", update)

		if update.Message != nil {
			ts.onMessage(update.Message)
		}
		if update.CallbackQuery != nil {
			ts.logger.Info("Answering")
			ts.bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
		}
		ts.logger.Infof("%+v\n", update.CallbackQuery)

	}
}

func (ts *TelegramService) onMessage(message *tgbotapi.Message) {
	messageText := strings.ToLower(message.Text)
	if message.IsCommand() {
		ts.logger.Info(messageText)
		if _, ok := ts.botCommands[messageText]; ok {
			ts.botCommands[messageText](message)
		}
		return
	}
	me, err := ts.bot.GetMe()
	if err != nil {
		ts.logger.Error("Error getting myself: %s", err)
	}
	if me.UserName != message.From.UserName {
		ts.logger.Info(messageText)
	}
}

func (ts *TelegramService) onStartCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет, меня зовут Ололоша и я могу сказать про")
	msg.ReplyMarkup = ts.getKeyboard()
	ts.bot.Send(msg)
}

func (ts *TelegramService) onHelpCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		message.Chat.ID,
		`
            /start - начать работу со мной
            /help - помощь
        `)
	msg.ReplyMarkup = ts.getKeyboard()
	ts.bot.Send(msg)
}

func (ts *TelegramService) onRemindMeCommand(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я пока-что не умею так делать")
	msg.ReplyToMessageID = message.MessageID
	ts.bot.Send(msg)
}

func (ts *TelegramService) getKeyboard() tgbotapi.InlineKeyboardMarkup {
	b := tgbotapi.NewInlineKeyboardButtonData("Ближайшие события", "nearest")
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(b))
}
