package Routes

import (
	"elucid503/SproutLoggingRelay/Misc"
	"elucid503/SproutLoggingRelay/Types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WatchServices(GinContext *gin.Context) {
		
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
	
	Conn.SetPongHandler(func(appData string) error {
				
		Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		
		return nil
		
	})

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
			
			_, _, ReadErr := Conn.ReadMessage()
			
			if ReadErr != nil {
				
				break
				
			}
			
		}
		
	}()
	
}		