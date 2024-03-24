package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/realtime"
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
		ID: "GENERAL",
	}

	var chanMsg = make(chan models.Message)

	err = client.SubscribeToMessageStream(channel, chanMsg)

	if err != nil {
		log.Fatalln("订阅消息失败", err)
	}

	for {
		select {
		case msg := <-chanMsg:
			var m = msg.Msg
			log.Println("接收到消息:", m)

			if m == "" {
				continue
			}

			if strings.HasPrefix(m, "@AiBot") {
				msgArr := strings.Split(m, "@AiBot")
				var realMsg = msgArr[1]

				if realMsg != "" {
					var msg = &models.Message{
						RoomID: "GENERAL",
						Msg:    realMsg + "吗？",
					}
					_, err = client.SendMessage(msg)
					if err != nil {
						log.Fatalln("消息发送失败", err)
					}
					log.Println("消息发送成功")
				}

			}

		}
	}

}
