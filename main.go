package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	apiai "github.com/meinside/api.ai-go"
	bot "github.com/meinside/telegram-bot-go"

	aihelper "github.com/meinside/telegram-bot-reminder-api.ai/ai"
	dbhelper "github.com/meinside/telegram-bot-reminder-api.ai/db"
)

const (
	DbFilename     = "db.sqlite"
	ConfigFilename = "config.json"

	CommandStart         = "/start"
	CommandListReminders = "/list"
	CommandCancel        = "/cancel"
	CommandHelp          = "/help"

	MessageCancel           = "취소"
	MessageCommandCanceled  = "명령이 취소 되었습니다."
	MessageReminderCanceled = "알림이 취소 되었습니다."
	MessageTextNeeded       = "텍스트를 입력해 주세요."
	MessageError            = "오류가 발생했습니다."
	MessageNoReminders      = "예약된 알림이 없습니다."
	MessageSaveFailed       = "알림 저장을 실패 했습니다"
	MessageCancelWhat       = "어떤 알림을 취소하시겠습니까?"
	MessageTimeIsPastFormat = "2006.1.2 15:04는 이미 지난 시각입니다"
	MessageTimeParseError   = "시간이 올바르지 않습니다"
	MessageSendingBackFile  = "받은 파일을 다시 보내드립니다."
	MessageUsage            = `사용법:

* 사용 예:
"내일 저녁 9시에 뉴스 보라고 보내줘"
"12월 31일 오후 11시에 신년 타종행사 보라고 알려줘"

* 기타 명령어:
/list : 예약된 알림 조회
/cancel : 예약된 알림 취소
/help : 본 사용법 확인

* 문의:
https://github.com/meinside/telegram-bot-reminder-api.ai
`

	// messages for api.ai errors
	MessageApiAiErrorFormat         = "api.ai 오류: %s"
	MessageApiAiDetailedErrorFormat = "api.ai 오류: %s (%s)"
)

var telegram *bot.Bot
var ai *apiai.Client
var db *dbhelper.Database
var _location *time.Location

var _conf config
var _maxNumTries int
var _monitorIntervalSeconds int
var _telegramIntervalSeconds int
var _restrictUsers bool
var _allowedUserIds []string

var _isVerbose bool

type config struct {
	TelegramApiToken        string   `json:"telegram_api_token"`
	ApiaiAccessToken        string   `json:"apiai_access_token"`
	MonitorIntervalSeconds  int      `json:"monitor_interval_seconds"`
	TelegramIntervalSeconds int      `json:"telegram_interval_seconds"`
	MaxNumTries             int      `json:"max_num_tries"`
	RestrictUsers           bool     `json:"restrict_users,omitempty"`
	AllowedUserIds          []string `json:"allowed_user_ids"`
	IsVerbose               bool     `json:"is_verbose,omitempty"`
}

func openConfig() (conf config, err error) {
	if file, err := ioutil.ReadFile(ConfigFilename); err == nil {
		if err := json.Unmarshal(file, &conf); err == nil {
			return conf, nil
		} else {
			return config{}, err
		}
	} else {
		return config{}, err
	}
}

func init() {
	var err error
	if _conf, err = openConfig(); err != nil {
		panic(err)
	} else {
		if _conf.MonitorIntervalSeconds <= 0 {
			_conf.MonitorIntervalSeconds = 10
		}
		_monitorIntervalSeconds = _conf.MonitorIntervalSeconds

		if _conf.TelegramIntervalSeconds <= 0 {
			_conf.TelegramIntervalSeconds = 1
		}
		_telegramIntervalSeconds = _conf.TelegramIntervalSeconds

		if _conf.MaxNumTries < 0 {
			_conf.MaxNumTries = 10
		}
		_maxNumTries = _conf.MaxNumTries

		_restrictUsers = _conf.RestrictUsers
		_allowedUserIds = _conf.AllowedUserIds

		telegram = bot.NewClient(_conf.TelegramApiToken)
		telegram.Verbose = _conf.IsVerbose

		ai = apiai.NewClient(_conf.ApiaiAccessToken)
		ai.Verbose = _conf.IsVerbose

		db = dbhelper.OpenDb(DbFilename)

		_location, _ = time.LoadLocation("Local")
		_isVerbose = _conf.IsVerbose
	}
}

// check if given Telegram id is allowed or not
func isAllowedId(id string) bool {
	if _restrictUsers == false {
		return true
	}

	for _, v := range _allowedUserIds {
		if v == id {
			return true
		}
	}

	return false
}

func monitorQueue(monitor *time.Ticker, client *bot.Bot) {
	for {
		select {
		case <-monitor.C:
			processQueue(client)
		}
	}
}

func processQueue(client *bot.Bot) {
	queue := db.DeliverableQueueItems(_maxNumTries)

	if _isVerbose {
		log.Printf("Checking queue: %d items...", len(queue))
	}

	for _, q := range queue {
		go func(q dbhelper.QueueItem) {
			// send message
			message := fmt.Sprintf("%s", q.Message)
			options := map[string]interface{}{}
			if sent := client.SendMessage(q.ChatId, message, options); !sent.Ok {
				log.Printf("*** failed to send reminder: %s", *sent.Description)
			} else {
				// mark as delivered
				if !db.MarkQueueItemAsDelivered(q.ChatId, q.Id) {
					log.Printf("*** failed to mark chat id: %d, queue id: %d", q.ChatId, q.Id)
				}
			}

			// increase num tries
			if !db.IncreaseNumTries(q.ChatId, q.Id) {
				log.Printf("*** failed to increase num tries for chat id: %d, queue id: %d", q.ChatId, q.Id)
			}
		}(q)
	}
}

func processUpdate(b *bot.Bot, update bot.Update, err error) {
	if err == nil {
		if update.HasMessage() {
			username := *update.Message.From.Username

			if !isAllowedId(username) {
				log.Printf("*** Id not allowed: %s", username)

				return
			}

			chatId := update.Message.Chat.Id

			// 'is typing...'
			b.SendChatAction(chatId, bot.ChatActionTyping)

			message := ""
			options := map[string]interface{}{
				"reply_markup": bot.ReplyKeyboardMarkup{ // show keyboards
					Keyboard: [][]bot.KeyboardButton{
						[]bot.KeyboardButton{
							bot.KeyboardButton{
								Text: CommandListReminders,
							},
						},
						[]bot.KeyboardButton{
							bot.KeyboardButton{
								Text: CommandCancel,
							},
						},
						[]bot.KeyboardButton{
							bot.KeyboardButton{
								Text: CommandHelp,
							},
						},
					},
					ResizeKeyboard: true,
				},
			}

			if update.Message.HasText() { // text
				txt := *update.Message.Text

				if strings.HasPrefix(txt, CommandStart) { // /start
					message = MessageUsage
				} else if strings.HasPrefix(txt, CommandListReminders) {
					reminders := db.UndeliveredQueueItems(chatId)
					if len(reminders) > 0 {
						for _, r := range reminders {
							message += fmt.Sprintf("➤ %s (%s)\n", r.Message, r.FireOn.Format("2006.1.2 15:04"))
						}
					} else {
						message = MessageNoReminders
					}
				} else if strings.HasPrefix(txt, CommandCancel) {
					reminders := db.UndeliveredQueueItems(chatId)
					if len(reminders) > 0 {
						// inline keyboards
						keys := make(map[string]string)
						for _, r := range reminders {
							keys[fmt.Sprintf("➤ %s (%s)", r.Message, r.FireOn.Format("2006.1.2 15:04"))] = fmt.Sprintf("%s %d", CommandCancel, r.Id)
						}
						buttons := bot.NewInlineKeyboardButtonsAsRowsWithCallbackData(keys)

						// add a button for canceling command
						cancel := CommandCancel
						buttons = append(buttons, []bot.InlineKeyboardButton{
							bot.InlineKeyboardButton{
								Text:         MessageCancel,
								CallbackData: &cancel,
							},
						})

						// options
						options["reply_markup"] = bot.InlineKeyboardMarkup{
							InlineKeyboard: buttons,
						}

						message = MessageCancelWhat
					} else {
						message = MessageNoReminders
					}
				} else if strings.HasPrefix(txt, CommandHelp) {
					message = MessageUsage
				} else {
					// send query to api.ai
					if response, err := ai.QueryText(apiai.QueryRequest{
						Query:     []string{txt},
						SessionId: sessionIdFor(chatId),
						Language:  apiai.Korean,
					}); err == nil {
						if response.Status.ErrorType == apiai.Success {
							if response.Result.ActionIncomplete {
								message = response.Result.Fulfillment.Speech
							} else {
								message = processQueryResponse(chatId, response)
							}
						} else {
							message = fmt.Sprintf(MessageApiAiDetailedErrorFormat, response.Status.ErrorType, response.Status.ErrorDetails)
						}
					} else {
						message = fmt.Sprintf(MessageApiAiErrorFormat, err)
					}
				}
			} else {
				message = MessageTextNeeded
			}

			// send message
			if len(message) <= 0 {
				message = MessageError
			}
			if sent := b.SendMessage(chatId, message, options); !sent.Ok {
				log.Printf("*** failed to send message: %s", *sent.Description)
			}
		} else if update.HasCallbackQuery() {
			processCallbackQuery(b, update)
		}
	} else {
		log.Printf("*** error while receiving update (%s)", err.Error())
	}
}

// process incoming callback query
func processCallbackQuery(b *bot.Bot, update bot.Update) bool {
	// process result
	result := false

	query := *update.CallbackQuery
	txt := *query.Data

	var message string = MessageError
	if strings.HasPrefix(txt, CommandCancel) {
		if txt == CommandCancel {
			message = MessageCommandCanceled
		} else {
			cancelParam := strings.TrimSpace(strings.Replace(txt, CommandCancel, "", 1))
			if queueId, err := strconv.Atoi(cancelParam); err == nil {
				if db.DeleteQueueItem(query.Message.Chat.Id, int64(queueId)) {
					message = MessageReminderCanceled
				} else {
					log.Printf("*** Failed to delete reminder")
				}
			} else {
				log.Printf("*** Unprocessable callback query: %s", txt)
			}
		}
	} else {
		log.Printf("*** Unprocessable callback query: %s", txt)
	}

	// answer callback query
	if apiResult := b.AnswerCallbackQuery(query.Id, map[string]interface{}{"text": message}); apiResult.Ok {
		// edit message and remove inline keyboards
		options := map[string]interface{}{
			"chat_id":    query.Message.Chat.Id,
			"message_id": query.Message.MessageId,
		}
		if apiResult := b.EditMessageText(message, options); apiResult.Ok {
			result = true
		} else {
			log.Printf("*** Failed to edit message text: %s", *apiResult.Description)

			db.LogError(fmt.Sprintf("failed to edit message text: %s", *apiResult.Description))
		}
	} else {
		log.Printf("*** Failed to answer callback query: %+v", query)

		db.LogError(fmt.Sprintf("failed to answer callback query: %+v", query))
	}

	return result
}

func sessionIdFor(chatId int64) string {
	return fmt.Sprintf("ss_%d", chatId)
}

func processQueryResponse(chatId int64, response apiai.QueryResponse) string {
	var message string = response.Result.Fulfillment.Speech

	// if confirmed yes,
	if response.Result.Metadata.IntentName == aihelper.IntentNameMessageConfirmedYes {
		params := response.Result.Parameters

		// check params
		if msg, ok := params["message"]; ok {
			if dt, ok := params["date"]; ok {
				if tm, ok := params["time"]; ok {
					// parse date & time
					if when, err := time.ParseInLocation(
						"2006-01-02 15:04:05",
						fmt.Sprintf("%s %s", dt, tm),
						_location,
					); err == nil {
						if when.Unix() >= time.Now().Unix() {
							// save it to DB
							if !db.Enqueue(chatId, msg.(string), when) {
								message = MessageSaveFailed
							}
						} else {
							message = when.Format(MessageTimeIsPastFormat)
						}
					} else {
						message = MessageTimeParseError
					}
				}
			}
		}
	}

	return message
}

func main() {
	// get info about this bot
	if me := telegram.GetMe(); me.Ok {
		// delete webhook (getting updates will not work when wehbook is set up)
		if unhooked := telegram.DeleteWebhook(); unhooked.Ok {
			// monitor queue
			log.Printf("> Starting monitoring queue...")
			go monitorQueue(
				time.NewTicker(time.Duration(_monitorIntervalSeconds)*time.Second),
				telegram,
			)

			// setup api.ai agent
			log.Printf("> Setting up agent...")
			aihelper.SetupAgent(ai, db)

			// wait for new updates
			log.Printf("> Starting bot: @%s (%s)", *me.Result.Username, me.Result.FirstName)
			telegram.StartMonitoringUpdates(0, _telegramIntervalSeconds, processUpdate)
		} else {
			panic("failed to delete webhook")
		}
	} else {
		panic("failed to get info of the bot")
	}
}
