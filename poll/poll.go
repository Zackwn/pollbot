package poll

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	emojis = [10]string{
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
	checkmarkemoji = "\u2705"
)

type PollError struct {
	message string
}

func (err PollError) Error() string {
	return err.message
}

func newPollError(message string) PollError {
	return PollError{message: message}
}

type Poll struct {
	Question string
	Answers  []string

	authorID string
}

func (poll *Poll) Winner(session *discordgo.Session, msg *discordgo.Message) (*discordgo.MessageReactions, int) {
	msg, _ = session.ChannelMessage(msg.ChannelID, msg.ID)

	var winner *discordgo.MessageReactions
	var winnerIndex int

	i := 0
	for i < len(emojis) && i < len(msg.Reactions) {
		reaction := msg.Reactions[i]
		if reaction.Emoji.Name == emojis[i] {
			if winner != nil {
				if reaction.Count > winner.Count {
					winner = reaction
					winnerIndex = i
				}
			} else {
				winner = reaction
				winnerIndex = i
			}
		}
		i++
	}
	return winner, winnerIndex
}

func (poll *Poll) Start(embed *discordgo.MessageEmbed, session *discordgo.Session, message *discordgo.Message) {
	// send poll embed
	embedMsg, _ := session.ChannelMessageSendEmbed(message.ChannelID, embed)

	// add reactions
	for i := range poll.Answers {
		session.MessageReactionAdd(embedMsg.ChannelID, embedMsg.ID, emojis[i])
	}
	session.MessageReactionAdd(embedMsg.ChannelID, embedMsg.ID, checkmarkemoji)

	var removeReactionHandler func()
	removeReactionHandler = session.AddHandler(func(s *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		fmt.Println(reaction.UserID, poll.authorID)
		if reaction.UserID == poll.authorID && reaction.Emoji.Name == checkmarkemoji {
			winner, winnerIndex := poll.Winner(session, embedMsg)

			content := "🕗 Time's up, the result is: `" + winner.Emoji.Name + " " + poll.Answers[winnerIndex] + "` 🥳"
			session.ChannelMessageSend(reaction.ChannelID, content)
			removeReactionHandler()
		}
	})
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

func NewPoll(authorID string, text []rune) (*Poll, error) {
	poll := new(Poll)
	poll.authorID = authorID
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
	if len(poll.Answers) > 10 {
		return nil, newPollError("error: There can only be up to 10 asnwers")
	}
	return poll, nil
}
