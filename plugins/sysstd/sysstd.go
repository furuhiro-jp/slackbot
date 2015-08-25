package sysstd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kyokomi/slackbot/plugins"
)

type Plugin struct {
	debug bool
}

func (p Plugin) SetTimezone(tmText string) {
	reTimeZone, ok := timeZone[strings.ToUpper(tmText)]
	if ok {
		// TODO: かなり強引なので注意...BotContextでTimezone持つべき
		time.Local = &reTimeZone
	}
}

func (p Plugin) ExecuteCommand(args ...string) string {
	switch {
	case dateCommand.Contains(args[0]):
		return time.Now().Format(defaultTimeFormat)
	case setTimezoneCommand.Contains(args[0]):
		if len(args) >= 2 {
			p.SetTimezone(args[1])
		}
		return fmt.Sprintf("%#v", *time.Local)
	}
	return "`command error`"
}

func (r Plugin) CheckMessage(event plugins.BotEvent, message string) (bool, string) {
	if r.debug {
		log.Printf("message   [%s]\n", message)
		log.Printf("botLinkID [%s]\n", event.BotLinkID())
		log.Printf("botName   [%s]\n", event.BotName())
		log.Printf("botID     [%s]\n", event.BotID())
	}

	var cmdArgs []string
	if strings.HasPrefix(message, event.BotLinkID()) {
		cmdArgs = strings.Fields(message[len(event.BotLinkID()):])
	} else if strings.HasPrefix(message, event.BotName()) {
		cmdArgs = strings.Fields(message[len(event.BotName()):])
	} else if strings.HasPrefix(message, event.BotID()) {
		cmdArgs = strings.Fields(message[len(event.BotID()):])
	} else {
		return false, message
	}

	if r.debug {
		log.Println(cmdArgs)
	}

	for cmdKey, cmd := range commandList {
		if cmd.Contains(cmdArgs[0]) {
			cmdArgs[0] = cmdKey
			return true, strings.Join(cmdArgs, " ")
		}
	}

	return false, message
}

func (r Plugin) DoAction(event plugins.BotEvent, message string) bool {
	cmdArgs := strings.Fields(message)
	if _, ok := commandList[cmdArgs[0]]; !ok {
		return true
	}
	event.Reply(r.ExecuteCommand(cmdArgs...))
	return false // next ok
}

var _ plugins.BotMessagePlugin = (*Plugin)(nil)