package main

import (
	"os"
	"net/url"
	"net/http"
	"strings"
	"fmt"
	"setupconfig"
		"github.com/gin-gonic/gin/json"
	"log"
	"strconv"
)

func setMessage() (string, string,string,string) {
	config := setupconfig.ReadWriteConfig("ini.config")

	token := config["TOKEN"]
	messsage_file := config["MESSAGE_FILE"]
	first_name:=config["FIRST_NAME"]
	last_name:=config["LAST_NAME"]

	file, err := os.Open(messsage_file)
	if err != nil {
		panic("ERROR : message not found")
	}

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	read_file := make([]byte, stat.Size())

	_, err = file.Read(read_file)
	if err != nil {
		panic("ERROR : cant read file")

	}
	msg := string(read_file)

	return msg, token,first_name, last_name
}

func sendMessage(chatID string) {

	msg, token,_ ,_:= setMessage()
	uri := "https://api.telegram.org/bot" + token + "/sendMessage"
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", msg)
	apiClient := &http.Client{}

	resp, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	resp.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp.Header.Add("Authorization", "auth_token="+token)

	getResp, _ := apiClient.Do(resp)
	fmt.Println(getResp.Status)

}

func getChatID() string {

	var dataJson struct {
		OK     bool `json:"ok"`
		RESULT []struct {
			UPDATE_ID int `json:"update_id"`
			MESSAGE   struct {
				MESSAGE_ID int `json:"message_id"`
				FROM       struct {
					ID int `json:"id"`
					IS_BOT bool `json:"is_bot"`
					FIRST_NAME string `json:"first_name"`
					LAST_NAME string `json:"last_name"`
					LANGUEAGE_CODE string `json:"langueage_code"`
				} `json:"from"`
				CHAT struct {
					ID        int    `json:"id"`
					FIRSTNAME string `json:"first_name"`
					LASTNAME  string `json:"last_name"`
					TYPE      string `json:"type"`
				} `json:"chat"`
				DATE int `json:"date"`
				TEXT string `json:"text"`
			} `json:"message"`
		} `json:"result"`
	}

	_, token, first_id, last_id:= setMessage()
	url := "https://api.telegram.org/bot" + token + "/getUpdates"
	nurl := fmt.Sprintf(url)

	// Build the request
	req, err := http.NewRequest("GET", nurl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)

	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp_msg, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp_msg.Body.Close()

	if err := json.NewDecoder(resp_msg.Body).Decode(&dataJson); err != nil {
		log.Println(err)
	}

var chatid string
	for index,_:=range dataJson.RESULT{
		first_name:=dataJson.RESULT[index].MESSAGE.CHAT.FIRSTNAME
		last_name:=dataJson.RESULT[index].MESSAGE.CHAT.LASTNAME
		if first_name==first_id && last_name==last_id{
			chatid = strconv.Itoa(dataJson.RESULT[index].MESSAGE.CHAT.ID)
		}
	}

	return chatid
}

func main() {

	sendMessage(getChatID())

}
