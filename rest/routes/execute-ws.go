package routes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/florianwoelki/kira/internal/pool"
	"github.com/florianwoelki/kira/pkg"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type wsEvent struct {
	Event    string   `json:"event" binding:"required"`
	Language string   `json:"language,omitempty"`
	Content  string   `json:"content,omitempty"`
	Stdin    []string `json:"stdin,omitempty"`
}

type wsResponse struct {
	Type      string `json:"type"`
	RunOutput string `json:"runOutput"`
	Time      int64  `json:"time"`
	Error     string `json:"error"`
}

func ExecuteWs(c echo.Context, rceEngine *pkg.RceEngine) error {
	// Initialize later.
	var executionInformation pool.ExecutionInformation

	// Create a pipe channel to communicate with the stream.
	pipeChannel := pkg.PipeChannel{
		Data:                 make(chan pool.StreamOutput),
		Terminate:            make(chan bool),
		ExecutionInformation: make(chan pool.ExecutionInformation),
	}

	// Websocket connection does not get closed automatically.
	websocket.Handler(func(ws *websocket.Conn) {
		for {
			// Receive and parse send JSON data from the client.
			data := wsEvent{}
			err := websocket.JSON.Receive(ws, &data)
			if err != nil {
				fmt.Println("receiving error:", err)
				return
			}

			switch data.Event {
			case "execute":
				// Execute the code of the client.
				go rceEngine.DispatchStream(pool.WorkData{
					Lang:        data.Language,
					Code:        data.Content,
					Stdin:       data.Stdin,
					Tests:       []pool.TestResult{},
					BypassCache: true,
				}, pipeChannel)

				go func() {
				Executor:
					for {
						select {
						case execInformation := <-pipeChannel.ExecutionInformation:
							executionInformation = execInformation
						case output := <-pipeChannel.Data:
							response := wsResponse{
								Type:      "output",
								RunOutput: output.Output,
								Time:      output.Time,
								Error:     output.Error,
							}

							// Send the result of the code back to the client.
							err = websocket.JSON.Send(ws, response)
							if err != nil {
								fmt.Println("sending error:", err)
								break Executor
							}

							logResponse(data, response)
						case <-pipeChannel.Terminate:
							response := wsResponse{
								Type:      "terminate",
								RunOutput: "",
							}
							err = websocket.JSON.Send(ws, response)
							if err != nil {
								fmt.Println("sending error:", err)
								break Executor
							}

							rceEngine.CleanUp(executionInformation.User, executionInformation.TempDirName)

							// Reset the execution information and pool channel.
							executionInformation = pool.ExecutionInformation{}
							pipeChannel = pkg.PipeChannel{
								Data:                 make(chan pool.StreamOutput),
								Terminate:            make(chan bool),
								ExecutionInformation: make(chan pool.ExecutionInformation),
							}

							logResponse(data, response)
							break Executor
						}
					}
				}()
			case "terminate":
				if executionInformation == (pool.ExecutionInformation{}) {
					continue
				}

				if err != nil {
					// TODO: Maybe remove in the future, there is currently no other way to check if
					// the connection was closed by the client.
					if strings.Contains(err.Error(), "use of closed network connection") {
						continue
					}

					fmt.Println("receiving error:", err)
					continue
				}

				pipeChannel.Terminate <- true
			}
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

// logResponse takes in the socket data as a request and the to be logged response for
// that request. It will log to the specific logger with `pkg.Logger`.
func logResponse(request wsEvent, response wsResponse) {
	dataBytes, err := json.Marshal(request)
	if err != nil {
		fmt.Println("marshalling error:", err)
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("marshalling error:", err)
		return
	}

	pkg.Logger.Info(
		"ws-request",
		zap.String("requestBody", string(dataBytes)),
		zap.String("responseBody", string(responseBytes)),
	)
}
