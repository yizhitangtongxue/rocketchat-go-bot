package main

import (
	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
	"log"
	"net/url"
	"strings"
	"time"
)

// rocket.chat server host info
var (
	host      = "localhost:3000"
	debugMode = true
)

// bot account info
var (
	botEmail    = "AiBot@local.com"
	botPassword = "123456"
)

// bot login info ，only when login succeed
var (
	id           = ""
	token        = ""
	tokenExpires = int64(0)
)

// channel or room info
var (
	channelOrRoomId = "GENERAL"
)

// when someone use ‘@’ mention bot
var (
	namePrefix = "@AiBot"
)

// realtime api = 实时聊天
// rest api = 控制服务器
// doc = https://developer.rocket.chat/reference/api/realtime-api

func main() {
	serverUrl := &url.URL{Host: host}
	client, err := realtime.NewClient(serverUrl, debugMode)
	if err != nil {
		log.Fatalln("获取客户端实例失败", err)
	}
	credentials := &models.UserCredentials{
		Email:    botEmail,
		Password: botPassword,
	}
	loginRes, err := client.Login(credentials)
	if err != nil {
		log.Fatalln("登录失败", err)
	}
	log.Printf("登录成功 Id: %s,Token: %s,TokenExpires: %d\n", loginRes.ID, loginRes.Token, loginRes.TokenExpires)
	id = loginRes.ID
	token = loginRes.Token
	tokenExpires = loginRes.TokenExpires

	var channel = &models.Channel{
		ID: channelOrRoomId,
	}

	var chanMsg = make(chan models.Message, 1)

	err = client.SubscribeToMessageStream(channel, chanMsg)

	if err != nil {
		log.Fatalln("订阅消息失败", err)
	}

	go func() {
		for {
			select {
			case msg := <-chanMsg:
				handleMessage(client, msg)
			}
		}
	}()

	for {
		time.Sleep(10 * time.Second)
	}
}

func handleMessage(client *realtime.Client, msg models.Message) {

	var m = msg.Msg

	if m == "" {
		return
	}
	log.Println("接收到消息:", m)

	if strings.HasPrefix(m, namePrefix) {
		msgArr := strings.Split(m, namePrefix)
		var realMsg = msgArr[1]

		if realMsg != "" {

			// send req to ollama local api
			msgRes, err := SendMessage("qwen:7b", realMsg)

			if err != nil {
				log.Fatalln("调用Ollama服务失败", err)
			}

			var msg = &models.Message{
				RoomID: channelOrRoomId,
				Msg:    msgRes.Response,
			}
			_, err = client.SendMessage(msg)
			if err != nil {
				log.Fatalln("消息发送失败", err)
			}
			log.Println("消息发送成功")
		}

	}
}
