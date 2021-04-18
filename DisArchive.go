package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	//"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func downloadImage(url, fileName string) error {
	response, err := http.Get(url)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Recieved NON-OK response code on url " + url)
	}
	defer response.Body.Close()

	file, err := os.Create("c:\\Users\\Alonzo\\Programming\\DisArchive\\DisArchive\\images\\" + fileName)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)

	if err != nil {
		return err
	}

	return nil
}

//put the last element of the 50 slice into a global variable then loop it in the message create event !start
var (
	lastChatID string
	medium     int
)

func archive(s *discordgo.Session, lastChatID, channelID string) string {

	//index the first and last message, make the last message first to keep going 100 back

	searchRange := 100
	message, err := s.ChannelMessages(channelID, searchRange, lastChatID, "", "")
	if err != nil {
		log.Println(err)
		return ""
	}

	//c:\Users\Alonzo\Programming\DisArchive\DisArchive\images
	//var messageID []int
	//gets last element in messages unique ID and makes it global
	log.Println("LAST CHAT ID: " + lastChatID)
	lastChatID = message[len(message)-1].ID
	//look for last message in range, and go another 100 back

	for _, content := range message {
		medium++
		log.Println(medium)

		if len(content.Attachments) != 0 {

			for _, foo := range content.Attachments {

				fileType := strings.SplitAfter(foo.Filename, ".")
				fileName := foo.ID + "." + fileType[1]
				//create your own folder for images and place the path below
				//only creates if file does not exist, file use unique IDs names so it should not make duplicates
				if _, err := os.Stat("c:\\Users\\Alonzo\\Programming\\DisArchive\\DisArchive\\images\\" + fileName); os.IsNotExist(err) {

					log.Println("Creating file " + fileName)
					err := downloadImage(foo.URL, fileName)
					if err != nil {
						log.Println(err)
					}

				}

			}

		}

	}
	archive(s, lastChatID, channelID)
	return lastChatID
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, "!start") {
		medium = 0
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		args := strings.SplitAfter(m.Content, " ")
		err := archive(s, args[1], args[2])
		if err == "" {
			log.Println(err)
			log.Println("Error reading chat history")
		}
	}
}
func main() {
	err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/isHeBoosted/killerkeys.env")
	if err != nil {
		log.Fatal(err)
	}
	dkey := os.Getenv("DisKey")
	dg, err := discordgo.New("Bot " + dkey)

	//log.Println(reflect.TypeOf(dg))
	if err != nil {
		fmt.Println(err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	dg.State.MaxMessageCount = 50
	discordgo.NewState()

	err1 := dg.Open()

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	fmt.Println("CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer dg.Close()
	//messageCreate()
}
