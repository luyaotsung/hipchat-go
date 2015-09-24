package main

import (
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"io/ioutil"
	"net/http"
)

var (
	Access_Token    = "tNgowr8imQKkISK3LBI1cHDVXmkjxUPcvlmktwen"
	Room_ID         = "1337031" //Webhook_Testing
	Check_Room_List = [...]string{"HC300000(All)", "Webhook_Testing"}
)

func main() {

	c := hipchat.NewClient(Access_Token)

	//Get All Room List
	rooms, resp1, err := c.Room.List()

	handleRequestError(resp1, err)

	for _, room := range rooms.Items {

		for _, Name := range Check_Room_List {

			if Name == room.Name {
				fmt.Printf("%-25v%10v\n", room.Name, room.ID)
				fmt.Println("---")

				hooks, resp, err := c.Room.ListWebhooks(room.ID, nil)
				handleRequestError(resp, err)

				for _, webhook := range hooks.Webhooks {
					fmt.Printf("  %v %v\t%v\t%v\t%v\n", webhook.Name, webhook.ID, webhook.Event, webhook.URL, webhook.WebhookLinks.Links.Self)

				}

				fmt.Println("---")
			}
		}

	}

	notifRq := &hipchat.NotificationRequest{Message: "<br> Send From Eli's Go Sample<br>", Color: "red", MessageFormat: "html"}

	resp2, err := c.Room.Notification(Room_ID, notifRq)

	//xxx := hipchat.CreateWebhookRequest

	//WBRequest := hipchat.CreateWebhookRequest{Name: "Create By Eli Blog", Event: "room_message", URL: "http://blog.min-jo.idv.tw:8081/"}
	//_, _, err = c.Room.CreateWebhook(Room_ID, &WBRequest)

	//c.Room.DeleteWebhook(Room_ID, "800703")

	if err != nil {
		fmt.Printf("Error during room notification %q\n", err)
		fmt.Printf("Server returns %+v\n", resp2)
		return
	}
	fmt.Printf("Notification sent ! to %v \n")

}

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
