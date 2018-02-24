package bots

import (
	"net/http"
	"net/url"
	"strings"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type A3rtProvider struct {
}

func (p *A3rtProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) BotResult {
	r := BotResult{}

	var message string
	switch m.Type {
	case "text":
		var payloadText models.PayloadText
		json.Unmarshal(m.Payload, &payloadText)
		message = payloadText.Text
	case "image":
		message = "画像を受信しました"
	default:
		message = "メッセージを受信しました"
	}

	var cred models.A3rtCredencial
	json.Unmarshal(c, &cred)

	values := url.Values{}
	values.Set("apikey", cred.ApiKey)
	values.Add("query", message)
	req, err := http.NewRequest(
		"POST",
		"https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	var res models.A3RTResponse

	log.Printf("%#v\n", resp)
	//json.NewDecoder(resp.Body).Decode(res)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	log.Printf("%#v\n", string(body))
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("=================================")
	log.Printf("%#v\n", res)

	//if len(res.Results) > 0 {
	var textPayload utils.JSONText
	// A3RT
	err = json.Unmarshal([]byte("{\"text\": \""+res.Results[0].Reply+"\"}"), &textPayload)
	post := &models.Message{
		RoomId:  m.RoomId,
		UserId:  b.UserId,
		Type:    "text",
		Payload: textPayload,
	}
	posts := make([]*models.Message, 0)
	posts = append(posts, post)
	messages := &models.Messages{
		Messages: posts,
	}
	r.Messages = messages

	return r
}