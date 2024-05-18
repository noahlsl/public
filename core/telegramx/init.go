package telegramx

import (
	json "github.com/bytedance/sonic"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (c *Cfg) NewClient() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		panic(err)
	}
	return bot
}

type TgMsg struct {
	Project string
	Env     string
	Game    string
	Period  interface{}
	Remark  interface{}
	Remind  []string
}

func NewTgMsg(p, n string) *TgMsg {
	return &TgMsg{
		Project: p,
		Game:    n,
		Period:  nil,
		Remark:  nil,
	}
}
func (m *TgMsg) WithEnv(v interface{}) *TgMsg {
	m.Env = "测试环境"
	if v == "main" {
		m.Env = "正式环境"
	}
	return m
}

// WithPeriod 关联人
func (m *TgMsg) WithPeriod(p interface{}) *TgMsg {
	m.Period = p
	return m
}

// WithRemark 消息期数
func (m *TgMsg) WithRemark(r interface{}) *TgMsg {
	m.Remark = r
	return m
}

// WithRemind @人
func (m *TgMsg) WithRemind(r ...string) *TgMsg {
	m.Remind = r
	return m
}

func (m *TgMsg) ToStr() string {
	marshal, _ := json.Marshal(m)
	return string(marshal)
}

func SendMsgStr(bot *tgbotapi.BotAPI, m string, id int64) {
	msg := tgbotapi.NewMessage(id, m)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func SendMsg(bot *tgbotapi.BotAPI, ms *TgMsg, id int64, remind bool) {
	var (
		names []string
	)

	// 获取管理员名称
	if remind {
		names = GetMembers(bot, id)
	}

	ms = ms.WithRemind(names...)
	msg := tgbotapi.NewMessage(id, ms.ToStr())

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func GetMembers(bot *tgbotapi.BotAPI, id int64) []string {
	// 获取群组成员
	config := tgbotapi.ChatAdministratorsConfig{ChatConfig: tgbotapi.ChatConfig{
		ChatID: id,
	}}
	members, err := bot.GetChatAdministrators(config)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// 提取成员的用户ID
	names := make([]string, 0)
	for _, member := range members {
		names = append(names, member.User.UserName)
	}

	return names
}
