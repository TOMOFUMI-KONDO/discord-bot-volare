package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

type Request struct {
	AwsLogs struct {
		Data string `json:"data"`
	} `json:"awslogs"`
}

type LogData struct {
	LogEvents []struct {
		Message string `json:"message"`
	} `json:"logEvents"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	token := "Bot " + os.Getenv("TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	log.Println(request.Body)
	var req Request
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		log.Println("Error decoding reqBody", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	byteData, err := base64.StdEncoding.DecodeString(req.AwsLogs.Data)
	var logData LogData
	if err := json.Unmarshal(byteData, &logData); err != nil {
		log.Println("Error decoding logData", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	if err != nil {
		log.Println("Error decoding request", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	dg, err := discordgo.New(token)
	if err != nil {
		log.Println("Error creating Discord")
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	if err := dg.Open(); err != nil {
		log.Println("Error opening connection", err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}
	defer dg.Close()

	sendLog(dg, channelID, logData.LogEvents[0].Message)

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
