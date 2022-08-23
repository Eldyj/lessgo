package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/Goscord/goscord"
	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"github.com/Goscord/goscord/gateway"
)

var pref string = "~"
var bot *gateway.Session
var puts func(s string)
var msg *discord.Message

func runCMD(){
	var index string = getcmd(msg.Content)
	var indexes = [...]string{
		"help",
		"ping",
		"coin",
		"reverse",
		"profile",
		"avatar",
		"server",
		"say",
		"invalid command",
	}
	var descriptions = [...]string{
		"display this message",
		"check bot status",
		"bot flip the coin",
		"bot reverse your message",
		"show your profile",
		"show your avatar",
		"show server info",
		"bot say something",
		"invalid command",
	}
	var values = [...]func(){
		func() {
			var help string
			var i byte
			for ; i < byte(len(descriptions))-1; i++ {
				help += fmt.Sprintf("%s%s - %s\n", pref, indexes[i], descriptions[i])
			}
			puts(help)
		},

		func() {
			time1 := time.Now().UnixNano()
			puts(fmt.Sprintf("%d nano seconds", int(time.Now().UnixNano()-time1)))
		},

		func() {
			if rand.Int()%2 == 1 {
				puts("front side of coin")
			} else {
				puts("back side of coin")
			}
		},
		func() {
			yes := strings.Split(msg.Content, "~reverse ")[1]
			bruh := ""
			for i := len(yes) - 1; i >= 0; i-- {
				bruh += char_on(yes, byte(i))
			}
			puts(bruh)
		},
		func() {
			user := msg.Author

			if len(msg.Mentions) > 0 {
				user = msg.Mentions[0]
			}

			puts(fmt.Sprintf(
				"tag: <@%s>\n"+
					"discriminator: %s\n"+
					"id: %s\n",
				user.Id,
				user.Tag(),
				user.Id),
			)
		},
		func() {
			puts("")
		},
		func() {
			server, _ := bot.State().Guild(msg.GuildId)
			if server.Description == "" {
				server.Description = "no description"
			} else {
				server.Description += "\n"
			}
			if server.AfkChannelId == "" {
				server.AfkChannelId = "none"
			} else {
				server.AfkChannelId = "<#" + server.AfkChannelId + ">"
			}
			puts(fmt.Sprintf(
				"name: %s\n"+
					"id: %s\n"+
					"description: %s\n"+
					"channels: %d\n"+
					"emojis: %d\n"+
					"members: %d\n"+
					"owner: <@%s>\n"+
					"afk channel: %s",
				server.Name,
				server.Id,
				server.Description,
				len(server.Channels),
				len(server.Emojis),
				server.MemberCount,
				server.OwnerId,
				server.AfkChannelId,
			),
			)
		},

		func() {
			puts(msg.Content)
		},

		func() {
			puts(fmt.Sprintf("invalid command, `%shelp` for all commands", pref))
		},
	}
	var i byte
	for ; indexes[i] != index && i < byte(len(indexes))-1; i++ {
	}
	(values[i])()
}

func splitstr(str string, length byte) string {
	var result string = ""
	var i byte = 0
	for ; i < length; i++ {
		if char_on(str, i) != pref {
			result += char_on(str, i)
		}
	}
	return result
}

func char_on(str string, i byte) string {
	if byte(len(str)) > i {
		return strings.Split(str, "")[i]
	} else {
		return ""
	}
}

func getcmd(str string) string {
	var i byte
	var isValid func(string) bool = func(s string) bool {
		return s != "" && s != " " && s != "\t"
	}

	for ; isValid(char_on(str, i)); i++ {
	}

	return splitstr(str, i)
}

func main() {
	rand.NewSource(time.Now().UnixNano())
	conffile, _ := os.ReadFile("config.conf")
	conf := strings.Split(string(conffile), "\n")
	token := conf[0]
	status := discord.StatusType(conf[1])
	color, _ := strconv.Atoi(conf[2])

	bot = goscord.New(&gateway.Options{
		Token:   token,
		Intents: gateway.IntentGuildMessages + gateway.IntentGuilds + gateway.IntentGuildMembers,
	})

	bot.On("ready", func() {
		bot.SetStatus(status)
		bot.SetActivity(&discord.Activity{
			Name:    "~help",
			Type:    3,
			Details: "~help",
			State:   "my first go programm",
		})
	})

	bot.On("messageCreate", func(event *discord.Message) {
		msg = event
		if char_on(msg.Content, 0) == pref && msg.Author.Id != bot.Me().Id {
			e := embed.NewEmbedBuilder()
			e.SetColor(color)
			user := msg.Author

			if len(msg.Mentions) > 0 {
				user = msg.Mentions[0]
			}

			var webhook func(string) *embed.Embed = func(s1 string) *embed.Embed {
				e.SetTitle(getcmd(msg.Content))
				e.SetDescription(s1)

				switch getcmd(msg.Content) {
				case "server":
					server, _ := bot.State().Guild(msg.GuildId)
					e.SetThumbnail(fmt.Sprintf("http://cdn.discordapp.com/icons/%s/%s", server.Id, server.Icon))
				case "profile":
					e.SetThumbnail(user.AvatarURL())
				case "avatar":
					e.SetImage(user.AvatarURL())
				}
				return e.Embed()
			}
			puts = func(txt string) {
                                bot.Channel.ReplyMessage(msg.ChannelId, msg.Id, webhook(txt))
                        }
			runCMD()
		}
	})

	bot.Login()
	select {}
}
