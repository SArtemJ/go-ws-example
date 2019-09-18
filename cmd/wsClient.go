package cmd

import (
	"bufio"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/SArtemJ/wstest/messages"
	"github.com/SArtemJ/wstest/server"
	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
	"log"
	"os"
)

var callclient = &cobra.Command{
	Use:   "callclient",
	Short: "callclient",
	Long:  `callclient`,
	Run: func(cmd *cobra.Command, args []string) {
		StartClient()
	},
}

func init() {
	rootCmd.AddCommand(callclient)
}

func StartClient() {
	ws, err := connectToServer()
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var m messages.Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error read message: ", err.Error())
				break
			}
			fmt.Println("You got message: ", m.Data)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msgText := scanner.Text()
		if msgText == "" {
			continue
		}
		m := messages.Message{
			Data: msgText,
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Error sending message: ", err.Error())
			break
		}
	}
}

func connectToServer() (*websocket.Conn, error) {
	h, err := server.PreparedAddressHost()
	if err != nil {
		return nil, err
	} else {
		p, err := server.PreparedAddressPort()
		if err != nil {
			return nil, err
		} else {
			srvAddr := fmt.Sprintf("ws://%s:%s", h, p)
			return websocket.Dial(srvAddr, "", "http://"+randomdata.IpV4Address())
		}
	}
}
