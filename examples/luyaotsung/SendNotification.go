package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func handleRequestError(resp *http.Response, err error) {
	if err != nil {
		if resp != nil {
			fmt.Printf("Request Failed:\n%+v\n", resp)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v\n", body)
		} else {
			fmt.Printf("Request failed, response is nil")
		}
		panic(err)
	}
}

func get_Room_ID(token string, name string) (string, error) {

	c := hipchat.NewClient(token)
	rooms, resp1, err := c.Room.List()
	handleRequestError(resp1, err)

	for _, room := range rooms.Items {
		if name == room.Name {
			fmt.Printf("%-25v%10v\n", room.Name, room.ID)
			fmt.Println("---")

			hooks, resp, err := c.Room.ListWebhooks(room.ID, nil)
			handleRequestError(resp, err)

			for _, webhook := range hooks.Webhooks {
				fmt.Printf("  %v %v\t%v\t%v\t%v\n", webhook.Name, webhook.ID, webhook.Event, webhook.URL, webhook.WebhookLinks.Links.Self)
			}
			fmt.Println("---")
			return strconv.Itoa(room.ID), nil
		}
	}
	return "", errors.New("Cannot find name of chart room , Please check again \n")
}

func send_Notify(token string, id string, message string, color string) {

	c := hipchat.NewClient(token)

	notifRq := &hipchat.NotificationRequest{Message: message, Color: color}
	resp2, err := c.Room.Notification(id, notifRq)

	if err != nil {
		fmt.Printf("Error during room notification %q\n", err)
		fmt.Printf("Server returns %+v\n", resp2)
	}
}

func main() {

	Access_Token := flag.String("token", "", "Access Token")
	Room_Name := flag.String("name", "", "Name of Chart Room")
	Notify_Message := flag.String("message", "Hey!! This is just a test message !!", "Message of Notification")
	Notify_Color := flag.String("color", "red", "Color of Notification")

	flag.Parse()

	var Room_ID string
	var err error

	if *Access_Token == "" || *Room_Name == "" {
		fmt.Println("Please Input Access Token and Room Name")
		os.Exit(-1)
	} else {
		Room_ID, err = get_Room_ID(*Access_Token, *Room_Name)

		if err != nil {
			fmt.Printf("Get Room ID Fail %+v \n", err)
		}

		fmt.Printf("Access Token %s \n", *Access_Token)
		fmt.Printf("Chat Room Name %s \n", *Room_Name)
		fmt.Printf("Chat Room ID %d \n", Room_ID)

	}

	send_Notify(*Access_Token, Room_ID, *Notify_Message, *Notify_Color)

	fmt.Printf("Notification sent ! \n")

}
