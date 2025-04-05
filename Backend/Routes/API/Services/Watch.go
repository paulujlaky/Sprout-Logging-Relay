package Routes

import (
	"elucid503/SproutLoggingRelay/Misc"
	"elucid503/SproutLoggingRelay/Types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WatchLogs(GinContext *gin.Context) {
		
	// This will be a WebSocket route
	
	Upgrader := websocket.Upgrader{

		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	Conn, Err := Upgrader.Upgrade(GinContext.Writer, GinContext.Request, nil)

	if Err != nil {

		GinContext.JSON(500, Types.Response{

			Message: "Failed to upgrade connection",
			JSON:    nil,
		})

		return

	}
	
	Conn.WriteJSON(Types.Response{

		Message: "Connected to WebSocket",
		
	})

	// Add websocket to list of Connections
	
	Misc.AddWebSocket(Conn)
	
	// Ping/Pong Handling 
	
	Conn.SetCloseHandler(func(code int, text string) error {
				
		Misc.RemoveWebSocket(Conn)
		
		return nil
		
	})

	// Keep reading to detect disconnects
	
	go func() {
		
		defer func() {
			
			Misc.RemoveWebSocket(Conn)
			Conn.Close()
									
		}()

		for {
			
			_, MessageBytes, ReadErr := Conn.ReadMessage()
			
			if ReadErr != nil {
				
				break
				
			}

			// Check if the message is a ping

			if strings.Contains(strings.ToLower(string(MessageBytes)), "ping") {

				Conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
				continue

			}
			
		}
		
	}()
	
}		