package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	pb "github.com/iron-ore-26k/ore-pb-gen"
	"google.golang.org/grpc"
	"encoding/binary"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var GUILD_ID = "" // guild ID later set by env variables
var CHANNEL_ID = "" // channel ID later set by env variables

var dg *discordgo.Session

/* Enum to store song list, must match pb */
const (
	Wreck_of_the_ef int = 1
	Fortunate_son = 2
)

var song_files = [...]string{"unknown","wreck_of_the_ef.dca","fortunate_son.dca"}

var token string
var buffer = make([][]byte,0) // Buffer to store chosen song

var (
	port = flag.Int("port", 50051, "Server , port")
)

type server struct {
	pb.UnimplementedOreServiceServer
}

func (s *server) PlaySong(ctx context.Context, in *pb.PlaySongRequest) (*pb.PlaySongResponse, error) {
	log.Printf("Recieved: %v", in.GetSongToPlay())
	if(in.GetSongToPlay() != 0){
		if(in.GetSongToPlay() == pb.Song_SONG_WRECK_OF_EDMUND_FITZGERALD){
			err := loadSound(1) // load chosen song if it's not unknown
			if(err!=nil){
				fmt.Println(err)
			}
		}else{
			err := loadSound(2) // load chosen song if it's not unknown
			if(err!=nil){
				fmt.Println(err)
			}
		}
		
		err := playSound(dg, GUILD_ID, CHANNEL_ID)
		if(err!=nil){
			fmt.Println(err)
		}
	} else {
		fmt.Println("Error unknown song")
	}
	return &pb.PlaySongResponse{}, nil
}

func (s *server) Pause(ctx context.Context, in *pb.PauseRequest) (*pb.PauseResponse, error) {
	log.Printf("Recieved: pause request")
	return &pb.PauseResponse{}, nil
}

func main() {
	token = os.Getenv("BOT_TOKEN") // Get out bot toekn from the enviroment variables
	fmt.Println(token)
	GUILD_ID = os.Getenv("GUILD_ID")
	CHANNEL_ID = os.Getenv("CHANNEL_ID")

	// Load the sound file.
	err := loadSound(Wreck_of_the_ef) // load the wreck of the ef as default
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		fmt.Println("Please copy $GOPATH/src/github.com/bwmarrin/examples/airhorn/airhorn.dca to this directory.")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register guildCreate as a callback for the guildCreate events.
	// dg.AddHandler(guildCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	// dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOreServiceServer(s, &server{})
	log.Printf("server listening at %v, %v", lis.Addr(), lis.Addr().Network())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	dg.Close()

}


/* ----------- From example ----------- */

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "World of Warships")
	fmt.Println("Iron ore app is ready")
}



// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(songIndex int) error {

	file, err := os.Open("songs/" + song_files[songIndex])
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}