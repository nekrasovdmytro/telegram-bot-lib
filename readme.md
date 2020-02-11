# Introduction

**Telegram bot lib** is called to make telegram bot development as easy as possible. The idea behind the lib to 
provide and interface to build flow of communication with user.

The basic element of telegram' bot is **`ability`**, simple saying, that what your bot can do - search for cheapest flight,
book hotel, control your smart home, give information about price for petrol and so on.

Each ability has name, description and **`taskFlow`**. Task flow is second very important thing in the lib.
A bot you build must have some ability and ability must have sequence of tasks what should be done. This I call taskFlow.

Let's say your bot can add note to some storage, Postgre database for example. Your ability will be called _/addnote_, and taskFlow will have two steps: 
1. Ask user to type his note (or send voice, photo, location, whatever)
2. After user clicked button _send_, his note must be stored and he should get message back like "Your note added"

Also each bot has two default commands (you don't need to define them, they are working out of the box) `/start` - starts the bot, prints list of abilities what bot can do, and
`/stop` - stops execution of current flow. 

Each TaskFlow has map with tasks. From code perspective it looks so
```go
type TaskFlow map[int]Task

type Task interface {
	Execute(input Input) *TaskResult
}
```

Each Task has method _Execute(input Input) *TaskResult_, inside of that method can be whatever you want.

# Code example

bot.go
```go
package bot

import (
	"github.com/nekrasovdmytro/telegram-bot-lib"
	"github.com/nekrasovdmytro/notebot/abilities"
)

func NewMyNoteBot(token, port, publicURL, redisUrl string) (*defaultTravelBot, error) {
	bl, err := telegrabotlib.NewBot(token, port, publicURL, redisUrl)

    //Here we set abilities - what this bot can do
	bl.SetAbilities(telegrabotlib.AbilityMap{
		"/addnote": {
			Short:       "note",
			Name:        "Add note",
			Flow:        abilities.NewNoteTaskFlow(),
			Description: "Some description will be printed out when bot will be started",
		},
	})

	return &defaultTravelBot{
		bot:      bl,
	}, err
}

type defaultTravelBot struct {
	bot      *telegrabotlib.BasicBot
}

func (d defaultTravelBot) Start() {
	d.bot.Start()
}

```

abilities/addnote.go
```go
package abilities

import (
	"github.com/nekrasovdmytro/telegram-bot-lib"
)

func NewNoteTaskFlow() telegrabotlib.TaskFlow {
	return telegrabotlib.TaskFlow{
		telegrabotlib.FirstTask: &telegrabotlib.SampleTask{
			Do: func(input telegrabotlib.Input) *telegrabotlib.TaskResult {

				return telegrabotlib.NewTaskResult([]*telegrabotlib.SingleResult{
					telegrabotlib.NewSingleResult(telegrabotlib.TEXT, "Hello"), //message one
					telegrabotlib.NewSingleResult(telegrabotlib.TEXT, "Add note"), //message two
				}, telegrabotlib.Step{Index: 10})
			},
		},
		10: &telegrabotlib.SampleTask{
			Do: func(input telegrabotlib.Input) *telegrabotlib.TaskResult {
				noteInput := input.(*telegrabotlib.TextInput)

				//store to some where - database, google note, whatever

				return telegrabotlib.NewTaskResult([]*telegrabotlib.SingleResult{
					telegrabotlib.NewSingleResult(telegrabotlib.TEXT, "Ok. I got you note"),
					telegrabotlib.NewSingleResult(telegrabotlib.TEXT, "Here is the link to your note - link"),
				}, telegrabotlib.Step{Index: telegrabotlib.LastStepIndex})
			},
		},
	}
}

```

main.go
```go
import (
    "github.com/nekrasovdmytro/telegrambot/bot"
    "log"
    "os"
)

func main() {
    var (
        port      = os.Getenv("PORT")
        publicURL = os.Getenv("PUBLIC_URL")
        redisURL = os.Getenv("REDIS_URL") 
        token     = os.Getenv("TOKEN") 
    )

    b, err := bot.NewMyNoteBot(token, port, publicURL, redisURL)
    if err != nil {
        log.Fatal(err)
        return
    }
    b.Start()
}
```

# Conclusion
If you have some question, feel free to contact me. For collaboration use this website http://nekrasov.one
