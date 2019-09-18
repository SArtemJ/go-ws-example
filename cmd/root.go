package cmd

import (
	"fmt"
	"github.com/SArtemJ/wstest/messages"
	"github.com/SArtemJ/wstest/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strings"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "wssrv",
	Short: "wssrv",
	Long:  `wssrv`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRequiredParams()
		wsAppStart()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.go_ws_example.yaml)")
	rootCmd.PersistentFlags().String("ws.host", "localhost", "ws server host")
	rootCmd.PersistentFlags().Int("ws.port", 8099, "ws server port")
	viper.BindPFlag("ws.host", rootCmd.PersistentFlags().Lookup("ws.host"))
	viper.BindPFlag("ws.port", rootCmd.PersistentFlags().Lookup("ws.port"))

}

func initConfig() {
	viper.SetEnvPrefix("WSSRV")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("../.")
		viper.SetConfigName("go_ws_example")
	}

	err := viper.ReadInConfig()
	if err == nil {
		fmt.Print("Using config file:", viper.ConfigFileUsed())
	}
}

func wsAppStart() {
	addr, err := server.PreparedAddressPort()
	if err != nil {
		log.Fatal(err)
	}

	wspool := server.NewWsPool()
	messages.NewPool()

	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		mainHandler(ws, wspool)
	}))

	s := http.Server{Addr: ":" + addr, Handler: mux}
	log.Fatal(s.ListenAndServe())
}

func mainHandler(ws *websocket.Conn, wsp *server.WsPool) {
	go wsp.Start()

	wsp.NewClients <- ws
	for {
		var m messages.Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			wsp.StreamMessages <- messages.Message{err.Error()}
			wsp.DisconnectClient(ws)
			return
		}
		messages.SendersPool.Store(ws.RemoteAddr().String(), true)
		wsp.StreamMessages <- m
	}
}

func checkRequiredParams() {
	required := []string{
		"ws.host",
		"ws.port",
	}

	var check []string

	for _, item := range required {
		if viper.Get(item) == "" || viper.Get(item) == 0 {
			check = append(check, item)
		}
	}

	if len(check) != 0 {
		for _, item := range check {
			log.Fatalf("Missed required params - %s", item)
		}
	}
}
