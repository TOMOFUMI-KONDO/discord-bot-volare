package main

import (
	"fmt"
	"os"
	"log"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Log string `json:"log"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token := "Bot " + os.Getenv("TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	var req Request
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		panic(err)
	}

	dg, err := discordgo.New(token)
	if err != nil {
		fmt.Println("Error creating Discord")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	if err := dg.Open(); err != nil {
		fmt.Println("Error opening connection", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	defer dg.Close()

	sendLog(dg, channelID, req.Log)

	return events.APIGatewayProxyResponse{
		Body:       "Hello",
		StatusCode: 200,
	}, nil
}

func sendLog(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func main() {
	lambda.Start(handler)
}
