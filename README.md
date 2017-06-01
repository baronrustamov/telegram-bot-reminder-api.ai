# Reminder Bot (Telegram)

[api.ai](https://api.ai/)를 통해, 수신된 메시지를 예약된 시간에 보내주는

간단한 Telegram reminder bot.

[api.ai-go](https://github.com/meinside/api.ai-go)의 활용 예를 보이기 위해 만들었음.

## install

```bash
$ go get -d github.com/meinside/telegram-bot-reminder-api.ai
$ cd $GOPATH/src/github.com/meinside/telegram-bot-reminder-api.ai
```

## build

```bash
$ go build
```

## configure

샘플로 들어있는 config.json.sample을 config.json으로 복사, 고쳐서 사용

```bash
$ cp config.json.sample config.json
$ vi config.json
```

**telegram_api_token** 값은 본인의 telegram bot api token으로 교체,

**apiai_access_token** 값은 본인의 api.ai agent의 `developer access token`으로 교체.

## run

```bash
$ ./telegram-bot-reminder-api.ai
```

## license

MIT

