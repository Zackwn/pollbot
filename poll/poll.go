package poll

import (
	"github.com/bwmarrin/discordgo"
)

var emojis = [10]string{
	"1\ufe0f\u20e3", // 1
	"2\ufe0f\u20e3", // 2
	"3\ufe0f\u20e3", // 3
	"4\ufe0f\u20e3", // 4
	"5\ufe0f\u20e3", // 5
	"6\ufe0f\u20e3", // 6
	"7\ufe0f\u20e3", // 7
	"8\ufe0f\u20e3", // 8
	"9\ufe0f\u20e3", // 9
	"\U0001f51f",    // 10
}

type Poll struct {
	Question string
	Answers  []string
}

func (poll *Poll) Start(embed *discordgo.MessageEmbed, session *discordgo.Session, message *discordgo.Message) {
	// send poll embed
	embedMsg, _ := session.ChannelMessageSendEmbed(message.ChannelID, embed)

	// add reactions
	for i := range poll.Answers {
		session.MessageReactionAdd(embedMsg.ChannelID, embedMsg.ID, emojis[i])
	}
}

func (poll *Poll) BuildEmbed() *discordgo.MessageEmbed {
	answers := ""

	for i, answer := range poll.Answers {
		answers += "\n" + emojis[i] + "    " + answer + "\n"
	}

	return &discordgo.MessageEmbed{
		Title:       poll.Question,
		Description: answers,
	}
}

func NewPoll(text []rune) *Poll {
	poll := new(Poll)
	i := 0
	gotquestion := false
	for i < len(text) {
		char := text[i]
		if char == '"' {
			txtpart := ""
			j := i + 1
			for j < len(text) && text[j] != '"' {
				txtpart += string(text[j])
				j++
			}
			i = j + 1
			if len(txtpart) != 0 {
				if gotquestion {
					poll.Answers = append(poll.Answers, txtpart)
				} else {
					poll.Question = txtpart
					gotquestion = true
				}
			}
		} else {
			i++
		}
	}
	return poll
}
