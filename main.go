package main

import (
	"encoding/json"
	"fmt"
	"general_game/gmodel"
	"log"
	"strconv"
	"strings"
	"teen_webhooks/controller"
	"teen_webhooks/handler"
	"time"

	"net/http"

	"gopkg.in/robfig/cron.v3"
)

const (
	githubPointsSocket   = "/webhooks"
	telegramPointsSocket = "/pointssocket"
	roomHook             = "/roomHook"
	faultHook            = "/faultsocket"
	// SertificateName name
	SertificateName string = "/etc/ssl/teenpattionline_kz.pem"

	// SertificateKey key
	SertificateKey string = "/etc/ssl/private/server.key"
)

func main() {

	telega := &http.Server{
		Addr:         ":8443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Cron scheduled push notification
	c := cron.New()
	c.AddFunc("*/15 * * * *", func() {
		handler.LSOFBroadcast()
	})

	c.Start()

	// Telegram webhook
	http.HandleFunc("/thook", func(w http.ResponseWriter, req *http.Request) {
		log.Print("it handling")
		// First, decode the JSON response body
		body := &gmodel.RequestMessageTelegram{}
		if err := json.NewDecoder(req.Body).Decode(body); err != nil {
			fmt.Println("could not decode request body", err)
			return
		}

		commandSended := strings.ToLower(body.Message.Text)
		log.Print(body.Message.Chat.ID)

		switch commandSended {
		case "/start":
			err := controller.SetStart(body.Message.Chat.ID)
			if err != nil {
				log.Print(err)
			}
			go handler.LSOFUnicast(body.Message.Chat.ID)
			break
		case "/stop":
			err := controller.SetStop(body.Message.Chat.ID)
			if err != nil {
				log.Print(err)
			}
			break
		case "/startlsof":
			err := controller.SetStart(body.Message.Chat.ID)
			if err != nil {
				log.Print(err)
			}
			go handler.LSOFUnicast(body.Message.Chat.ID)
			break
		case "/stoplsof":
			err := controller.SetLsofStop(body.Message.Chat.ID)
			if err != nil {
				log.Print(err)
			}
			break
		default:
			// Killer comand
			if strings.HasPrefix(commandSended, "kill:") {
				onlyNumberString := strings.TrimPrefix(commandSended, "kill:")
				ons, err := strconv.Atoi(onlyNumberString)
				if err != nil {
					log.Print(err.Error())
					return
				}
				handler.KillRoom(ons)
			}
			break
		}
	})

	// Room not responding
	http.HandleFunc(roomHook, func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("roomId")
		if roomID == "" {
			http.Error(w, "Get 'file' not specified in url.", http.StatusBadRequest)
			return
		}

		ons, _ := strconv.Atoi(roomID)

		//here we should do a script
		// _, err := exec.Command("/bin/sh", "/home/danko/src/points_socket/lastroom.sh", roomID).Output()
		// if err != nil {
		// 	//log.Print("lick fat girl pussy")
		// 	fmt.Println(err)
		// }

		go handler.Broadcast("Dead RoomID: " + roomID)
		handler.KillRoom(ons)
	})

	// FaultHook
	http.HandleFunc(faultHook, func(w http.ResponseWriter, r *http.Request) {
		go handler.FailBroadcast()
	})

	if err := telega.ListenAndServeTLS(SertificateName, SertificateKey); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
