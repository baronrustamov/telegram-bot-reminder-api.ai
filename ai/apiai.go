package ai

import (
	"fmt"
	"log"

	apiai "github.com/meinside/api.ai-go"

	dbhelper "github.com/meinside/telegram-bot-reminder-api.ai/db"
)

// constants for api.ai
const (
	IntentNameMessage             = "message"
	IntentNameMessageConfirmedYes = "message-confirm-yes"
	IntentNameMessageConfirmedNo  = "message-confirm-no"

	ContextLifespan = 1
)

// setup agent
func SetupAgent(ai *apiai.Client, db *dbhelper.Database) {
	var existsMessage, existsConfirmYes, existsConfirmNo bool

	// check existence of intents
	if intents, err := ai.AllIntents(); err == nil {
		for _, intent := range intents {
			if intent.Name == IntentNameMessage {
				existsMessage = true
			} else if intent.Name == IntentNameMessageConfirmedYes {
				existsConfirmYes = true
			} else if intent.Name == IntentNameMessageConfirmedNo {
				existsConfirmNo = true
			}
		}
	}

	// create intents
	if existsMessage { // intent: message
		log.Printf("Intent '%s' already exists", IntentNameMessage)
	} else {
		createMessageIntent(ai, db)
	}
	if existsConfirmYes { // intent: message-confirm-yes
		log.Printf("Intent '%s' already exists", IntentNameMessageConfirmedYes)
	} else {
		createConfirmYesIntent(ai, db)
	}
	if existsConfirmNo { // intent: message-confirm-no
		log.Printf("Intent '%s' already exists", IntentNameMessageConfirmedNo)
	} else {
		createConfirmNoIntent(ai, db)
	}
}

func createMessageIntent(ai *apiai.Client, db *dbhelper.Database) {
	if res, err := ai.CreateIntent(apiai.IntentObject{
		Name:     IntentNameMessage,
		Auto:     true,
		Contexts: []string{}, // no input context
		UserSays: []apiai.UserSays{
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "숙제 빨리 해놓으라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "숙제 빨리 해놓으라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "내일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "내일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "6월 2일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "6월 2일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "09시 30분",
						Meta:  "@sys.time",
						Alias: "time",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "09시 30분",
						Meta:  "@sys.time",
						Alias: "time",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "5월 18일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text:  "밥 먹으라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "5월 18일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text:  "밥 먹으라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "12월 31일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text:  "23시",
						Meta:  "@sys.time",
						Alias: "time",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text:  "타종행사 보라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "보내줘",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text:  "12월 31일",
						Meta:  "@sys.date",
						Alias: "date",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text:  "23시",
						Meta:  "@sys.time",
						Alias: "time",
					},
					apiai.UserSaysData{
						Text: "에 ",
					},
					apiai.UserSaysData{
						Text:  "타종행사 보라고",
						Meta:  "@sys.any",
						Alias: "message",
					},
					apiai.UserSaysData{
						Text: " ",
					},
					apiai.UserSaysData{
						Text: "알려줘",
					},
				},
			},
		},
		Responses: []apiai.IntentResponse{
			apiai.IntentResponse{
				ResetContexts: false,
				AffectedContexts: []apiai.IntentAffectedContext{
					apiai.IntentAffectedContext{
						Name:     IntentNameMessage,
						Lifespan: ContextLifespan,
					},
				},
				Parameters: []apiai.IntentResponseParameter{
					apiai.IntentResponseParameter{
						Name:     "message",
						Value:    "$message",
						Required: true,
						DataType: "@sys.any",
						Prompts: []string{
							"무슨 메시지를 보내드릴까요?",
							"어떤 메시지를 보내드릴까요?",
						},
					},
					apiai.IntentResponseParameter{
						Name:     "date",
						Value:    "$date",
						Required: true,
						DataType: "@sys.date",
						Prompts: []string{
							"어느 날짜에 보내드릴까요?",
							"어느 날짜에 알려드릴까요?",
							"몇 월 몇 일에 보내드릴까요?",
							"몇 월 몇 일에 알려드릴까요?",
						},
					},
					apiai.IntentResponseParameter{
						Name:     "time",
						Value:    "$time",
						Required: true,
						DataType: "@sys.time",
						Prompts: []string{
							"어느 시간에 보내드릴까요?",
							"어느 시간에 알려드릴까요?",
							"몇 시에 보내드릴까요?",
							"몇 시에 알려드릴까요?",
						},
					},
				},
				Messages: []apiai.Message{
					apiai.TextResponseMessage("", []string{
						`$date $time에 "$message"라고 메시지를 보내드릴까요?`,
						`$date $time에 "$message"라고 알려드릴까요?`,
					}),
				},
			},
		},
		Priority: 500000,
	}); err != nil {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessage, err)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessage, err))
	} else if res.Status.Code != 200 {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessage, res.Status.ErrorDetails)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessage, res.Status.ErrorDetails))
	}
}

func createConfirmYesIntent(ai *apiai.Client, db *dbhelper.Database) {
	if res, err := ai.CreateIntent(apiai.IntentObject{
		Name:     IntentNameMessageConfirmedYes,
		Auto:     true,
		Contexts: []string{IntentNameMessage},
		UserSays: []apiai.UserSays{
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "그래",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "응",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "네",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "예",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "ㅇㅇ",
					},
				},
			},
		},
		Responses: []apiai.IntentResponse{
			apiai.IntentResponse{
				ResetContexts: true,
				Parameters: []apiai.IntentResponseParameter{
					apiai.IntentResponseParameter{
						Name:  "date",
						Value: "#message.date",
					},
					apiai.IntentResponseParameter{
						Name:  "time",
						Value: "#message.time",
					},
					apiai.IntentResponseParameter{
						Name:  "message",
						Value: "#message.message",
					},
				},
				Messages: []apiai.Message{
					apiai.TextResponseMessage("", []string{
						`$date $time에 "$message"를 보내드리겠습니다.`,
						`$date $time에 "$message"라고 알려드리겠습니다.`,
					}),
				},
			},
		},
		Priority: 500001,
	}); err != nil {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessageConfirmedYes, err)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessageConfirmedYes, err))
	} else if res.Status.Code != 200 {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessageConfirmedYes, res.Status.ErrorDetails)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessageConfirmedYes, res.Status.ErrorDetails))
	}
}

func createConfirmNoIntent(ai *apiai.Client, db *dbhelper.Database) {
	if res, err := ai.CreateIntent(apiai.IntentObject{
		Name:     IntentNameMessageConfirmedNo,
		Auto:     true,
		Contexts: []string{IntentNameMessage},
		UserSays: []apiai.UserSays{
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "아니",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "아니오",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "노노",
					},
				},
			},
			apiai.UserSays{
				Data: []apiai.UserSaysData{
					apiai.UserSaysData{
						Text: "ㄴㄴ",
					},
				},
			},
		},
		Responses: []apiai.IntentResponse{
			apiai.IntentResponse{
				ResetContexts: true,
				Messages: []apiai.Message{
					apiai.TextResponseMessage("", []string{
						"취소 되었습니다.",
					}),
				},
			},
		},
		Priority: 500001,
	}); err != nil {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessageConfirmedNo, err)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessageConfirmedNo, err))
	} else if res.Status.Code != 200 {
		log.Printf("*** Failed to create intent %s: %s", IntentNameMessageConfirmedNo, res.Status.ErrorDetails)

		db.LogError(fmt.Sprintf("failed to create intent %s: %s", IntentNameMessageConfirmedNo, res.Status.ErrorDetails))
	}
}
