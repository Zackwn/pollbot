// Discord ASCII Game of Life
package gol

import (
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

const width = 9
const height = 9

type display [height][width]uint8

var pixel []string = []string{
	":white_large_square:",
	":green_square:",
}

var (
	dead        uint8 = 0
	alive       uint8 = 1
	randomState       = func() uint8 {
		return uint8(rand.Intn(2))
	}
)

const maxFrames = 15

func setup(display *display) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < height; i++ {
		display[i] = [width]uint8{}
		for j := 0; j < width; j++ {
			display[i][j] = randomState()
		}
	}
}

func nearbyLives(display *display, i, j int) int {
	count := 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if (x+i >= 0 && x+i < width) && (y+j >= 0 && y+j < height) {
				if display[y+j][x+i] == alive {
					count++
				}
			}
		}
	}
	return count
}

func draw(display *display, session *discordgo.Session, msg *discordgo.Message) {
	content := ""
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			content += pixel[display[i][j]]
		}
		content += "\n"
	}
	_, err := session.ChannelMessageEdit(msg.ChannelID, msg.ID, content)
	if err != nil {
		log.Println(err)
		return
	}
}

func next(display *display) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			count := nearbyLives(display, i, j)
			if display[i][j] == alive {
				// dead by lonineless or overpopulation
				if count < 2 || count > 3 {
					display[i][j] = dead
				}
			} else {
				// new life
				if count == 3 {
					display[i][j] = alive
				}
			}
		}
	}
}

func Run(session *discordgo.Session, message *discordgo.MessageCreate) {
	msg, err := session.ChannelMessageSend(message.ChannelID, "Starting...")
	if err != nil {
		log.Println(err)
		return
	}
	session.ChannelMessageDelete(message.ChannelID, message.ID)

	display := display([height][width]uint8{})
	var frames int

	setup(&display)
	draw(&display, session, msg)
	for {
		next(&display)
		draw(&display, session, msg)

		if frames > maxFrames {
			break
		}
		frames++

		// Slowdown animation to avoid spamming discord api too much.
		// Sometimes it takes a little longer because updating the message might be slow.
		time.Sleep(time.Second / 2)
	}

	// timeout before deleting message
	time.Sleep(2 * time.Second)
	session.ChannelMessageDelete(msg.ChannelID, msg.ID)
}
