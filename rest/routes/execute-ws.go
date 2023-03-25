package routes

import (
	"fmt"

	"github.com/florianwoelki/kira/pkg"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type socketData struct {
	Language string `json:"language"`
	Content  string `json:"content"`
}

type wsResponse struct {
	RunOutput string `json:"runOutput"`
}

func ExecuteWs(c echo.Context, rceEngine *pkg.RceEngine) error {
	websocket.Handler(func(ws *websocket.Conn) {
		dataCh := make(chan string)
		terminateCh := make(chan bool)
		defer ws.Close()

		// Receive and parse send JSON data from the client.
		data := socketData{}
		err := websocket.JSON.Receive(ws, &data)
		if err != nil {
			fmt.Println("receiving error:", err)
			return
		}

		// Execute the code of the client.
		go rceEngine.ExecuteWs(data.Content, data.Language, dataCh, terminateCh)

		for {
			select {
			case shouldTerminate := <-terminateCh:
				fmt.Println("terminate", shouldTerminate)
				return
			case output := <-dataCh:
				// Send the result of the code back to the client.
				err = websocket.JSON.Send(ws, wsResponse{
					RunOutput: output,
				})
				if err != nil {
					fmt.Println("sending error:", err)
					return
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
