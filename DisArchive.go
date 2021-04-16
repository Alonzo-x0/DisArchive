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

	x, err := io.Copy(file, response.Body)
	log.Println(x)
	log.Println("saved file " + fileName)
	if err != nil {
		return err
	}

	return nil
}

func archive(s *discordgo.Session) {
	//810621872844570645
	message, err := s.ChannelMessages("396535726855946240", 50, "816035932687564830", "", "")
	if err != nil {
		log.Println(err)
		return
	}
	//c:\Users\Alonzo\Programming\DisArchive\DisArchive\images
	//var messageID []int
	for _, content := range message {
		if len(content.Attachments) != 0 {
			for _, foo := range content.Attachments {
				fileType := strings.SplitAfter(foo.Filename, ".")
				err := downloadImage(foo.URL, foo.ID+"."+fileType[1])
				if err != nil {
					log.Println(err)
				}
			}

		}

	}
	log.Println("done")
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, "!start") {
		archive(s)
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
