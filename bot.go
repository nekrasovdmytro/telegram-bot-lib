package telegrabotlib

import (
    tb "gopkg.in/tucnak/telebot.v2"
    "log"
    "strconv"
    "strings"
)

type Bot interface {
    Abilities() AbilityMap
    Execute(userID string, r Recipient, executable Executable, input Input) *TaskResult
    SendMessage(r Recipient, what interface{})
    Start()
}

type Recipient interface {
    Recipient() string
}

type Message interface {
    Send(userId string, what interface{})
}

func NewBot(token, port, publicURL string) (*BasicBot, error) {
    webHook := &tb.Webhook{
        Listen:   ":" + port,
        Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
    }

    pref := tb.Settings{
        Token:  token,
        Poller: webHook,
    }

    telBot, err := tb.NewBot(pref)

    return &BasicBot{
        TelegramBot:telBot,
        UfManager: NewUserFlowManager(NewUserFlowSession(NewRedisSession())),
        specialCommands: map[string]struct{}{
            "/start": {},
            "/stop": {},
        },
    }, err
}

type BasicBot struct {
    TelegramBot     *tb.Bot
    UfManager       *UserFlowManager
    abilities       AbilityMap
    specialCommands map[string]struct{}
}

func (d BasicBot) SetAbilities(a AbilityMap) {
    d.abilities = a
}

func (d BasicBot) init() {
    d.TelegramBot.Handle("/start", func(m *tb.Message) {
        d.RenderStartFrame(m)
    })

    d.TelegramBot.Handle("/stop", func(m *tb.Message) {
        d.SendMessage(m.Sender, "Stopping")
        d.RenderStartFrame(m)
        d.UfManager.uSession.FinishFlow(strconv.Itoa(m.Sender.ID))
    })

    d.TelegramBot.Handle(tb.OnText, func(m *tb.Message) {
        if _, ok := d.specialCommands[m.Text]; ok {
            return
        }

        d.UfManager.ExecuteFlow(d, strconv.Itoa(m.Sender.ID), m.Sender, &TextInput{Text: m.Text})
    })
}

func (d BasicBot) Start() {
    log.Printf("Authorized on account %s", d.TelegramBot.Me.FirstName)
    d.init()
    d.TelegramBot.Start()
}

func (d BasicBot) Abilities() AbilityMap {
    return d.abilities
}

func (d BasicBot) Execute(userID string, r Recipient, executable Executable, input Input) *TaskResult {
    res := executable(input)

    log.Print("Executed data:")
    log.Print(res)

    for _, sr := range res.Result() {
        switch sr.Type {
        case LOCATION:
            d.SendMessage(r, sr.Result().(*tb.Location))
        default:
            d.SendMessage(r, sr.Result())
        }

    }

    return res
}

func (d BasicBot) SendMessage(r Recipient, what interface{}) {
    if _, err := d.TelegramBot.Send(r, what, tb.ModeHTML); err != nil {
        log.Print(err)
    }
}

func (d BasicBot) RenderStartFrame(m *tb.Message) {
    for key, a := range d.abilities  {
        d.SendMessage(m.Sender, key + " " + a.Name + "\n" + a.Description)
    }

    var inlineKeys [][]tb.ReplyButton
    for k := range d.abilities {
        inlineBtn := tb.ReplyButton{
            Text: k,
        }
        inlineKeys = append(inlineKeys, []tb.ReplyButton{inlineBtn})

        d.TelegramBot.Handle(&inlineBtn, func(c *tb.Callback) {
            d.UfManager.ExecuteFlow(d, strconv.Itoa(c.Sender.ID), c.Sender, &TextInput{Text: k})
        })
    }

    text := strings.Join(
        []string{
            "Hey ",
            m.Sender.FirstName,
            ", ",
            "please use menu buttons to interact with me",
        },
        "",
    )

    if _, err := d.TelegramBot.Send(m.Sender,text , &tb.ReplyMarkup{ReplyKeyboard: inlineKeys, ReplyKeyboardRemove: true, ResizeReplyKeyboard: true}); err != nil {
        log.Print(err)
    }
}