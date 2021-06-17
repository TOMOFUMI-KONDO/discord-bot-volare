package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

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

func handler(ctx context.Context, request Request) (string, error) {
	token := "Bot " + os.Getenv("TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	log.Println(request.AwsLogs.Data)

	byteData, err := base64.StdEncoding.DecodeString(request.AwsLogs.Data)
	var logData LogData
	if err := json.Unmarshal(byteData, &logData); err != nil {
		return handleErr("Error decoding logData", err)
	}

	if err != nil {
		return handleErr("Error decoding request", err)
	}

	dg, err := discordgo.New(token)
	if err != nil {
		return handleErr("Error creating Discord", err)
	}

	if err := dg.Open(); err != nil {
		return handleErr("Error opening connection", err)
	}
	defer dg.Close()

	sendLog(dg, channelID, logData.LogEvents[0].Message)

	return "success", nil
}

func sendLog(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func handleErr(msg string, err error) (string, error) {
	log.Println("[ERROR] " + msg + ":" + err.Error())
	return "ERROR", err
}

func main() {
	lambda.Start(handler)
}
