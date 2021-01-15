package handler

import (
	"bytes"
	"encoding/json"
	"flag"
	"general_game/gmodel"
	"general_game/gutils"
	"git_webhooks/controller"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "teenpattionline.kz:8081", "http service address")

// FailBroadcast used to send server fault
func FailBroadcast() {
	allLsofers, err := controller.GetAllActive()
	log.Print("here1")
	if err != nil {
		log.Print(err.Error())
		return
	}

	for _, single := range allLsofers {
		log.Print("sending for 1")
		go LSOFSingle(single.ChatID, "Server is downed and restarted")
	}
}

// Broadcast used to braodcast message
func Broadcast(Message string) {
	allLsofers, err := controller.GetAllActive()
	if err != nil {
		log.Print(err.Error())
		return
	}

	for _, single := range allLsofers {
		go LSOFSingle(single.ChatID, Message)
	}
}

// LSOFBroadcast used to send lsof for subs
func LSOFBroadcast() {
	allLsofers, err := controller.GetAllLsof()

	if err != nil {
		log.Print(err.Error())
		return
	}

	//here we should do a script
	cmd, err := exec.Command("/bin/sh", "/home/danko/src/git_webhooks/countOfSocket.sh").Output()
	if err != nil {
		log.Print("error in executing bash for lsof " + err.Error())
	}

	// goroutines reading
	dat, err := ioutil.ReadFile("/home/danko/src/points_socket/goCount.dat")
	if err != nil {
		log.Print(err.Error())
	}

	for _, single := range allLsofers {
		go LSOFSingle(single.ChatID, "lsof: "+string(cmd)+" go: "+string(dat))
	}
}

// LSOFSingle used to send lsof for sub
func LSOFSingle(chatID int64, message string) {
	// Create the request body struct
	reqBody := &gmodel.SendMessageTelegram{
		ChatID: chatID,
		Text:   message,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Print(err.Error())
		return
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot"+gutils.TokenTelegramTeen+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Print("normal" + err.Error())
		return
	}
	log.Print("done")
	return
}

// LSOFUnicast used to unicast
func LSOFUnicast(chatID int64) {
	//here we should do a script
	cmd, err := exec.Command("/bin/sh", "/home/danko/src/git_webhooks/countOfSocket.sh").Output()
	if err != nil {
		log.Print("error in executing bash for lsof " + err.Error())
	}

	// goroutines reading
	dat, err := ioutil.ReadFile("/home/danko/src/points_socket/goCount.dat")
	if err != nil {
		log.Print(err.Error())
	}

	// Create the request body struct
	reqBody := &gmodel.SendMessageTelegram{
		ChatID: chatID,
		Text:   "lsof: " + string(cmd) + " go: " + string(dat),
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Print(err.Error())
		return
	}

	//chatStr := strconv.FormatInt(chatID, 10)
	// Send a post request with your token
	urlt := "https://api.telegram.org/bot" + gutils.TokenTelegramTeen + "/sendMessage"
	log.Print(urlt)
	res, err := http.Post(urlt, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("normalerror %d", res.StatusCode)
		return
	}
	return
}

// KillRoom used to kill room
func KillRoom(RoomID int) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/room_killer"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(RoomID)))
	if err != nil {
		log.Println("write:", err)
		return
	}
}
