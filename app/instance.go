package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"games.launch.launcher/config"
	ll "games.launch.launcher/logger"
	"github.com/gofrs/uuid"
	"io"
	"net"
	"net/url"
)

var gameConnPool = make(map[string]*net.Conn)

// Start the main instance of the application.
func (l *Launcher) StartFirstInstance() {
	// Start a TCP listener on the launcher designated port.
	listener, err := net.Listen("tcp", "127.0.0.1:"+config.LauncherPort)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("error starting listener: %v\n", err))
		return
	}

	// Close the listener when the function exits.
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error closing listener: %v\n", err))
		}
	}(listener)

	// Listen for subsequent instance connections, function runs as a goroutine, so it will not block.
	for {
		// Accept incoming connections.
		conn, err := listener.Accept()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error accepting connection: %v\n", err))
			continue
		}

		// Handle the connection with a subsequent instance in a goroutine.
		go l.handleSubsequentLauncherInstanceConnection(conn)
	}
}

// Handle a connection from a subsequent instance of the application.
func (l *Launcher) handleSubsequentLauncherInstanceConnection(conn net.Conn) {
	// Close the connection when the function exits.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error closing connection: %v\n", err))
		}
	}(conn)

	// Read the parameters from the subsequent launcher connection.
	data, err := io.ReadAll(conn)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("Error reading data: %v\n", err))
		return
	}

	// Convert the deep link to a string.
	deepLink := string(data)
	ll.Logger.Print(fmt.Sprintf("Received deep link: %s\n", deepLink))

	// Process the deep link.
	l.processDeepLink(deepLink)
}

// Do a deep link processing task.
func (l *Launcher) processDeepLink(deepLink string) {
	// Process the deep link as needed.
	ll.Logger.Print(fmt.Sprintf("Processing deep link: %s\n", deepLink))

	// Parse the deep link from the URL.
	parsedUrl, err := url.Parse(deepLink)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("Error parsing deep link: %v\n", err))
		return
	}

	// Get the app id from the deep link.
	queryParams := parsedUrl.Query()
	appId := queryParams.Get("appId")
	if appId == "" {
		ll.Logger.Error(fmt.Sprintf("Error parsing app id: %s\n", appId))
		return
	}

	// Check if the game client is already running.
	if gameConn, ok := gameConnPool[appId]; ok {
		_, err := (*gameConn).Write([]byte(deepLink))
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error sending deep link: %v\n", err))
			return
		}
	} else {
		appUuid := uuid.FromStringOrNil(appId)

		if appUuid.IsNil() {
			ll.Logger.Error(fmt.Sprintf("Error parsing app id: %s\n", appId))
		} else {
			err := l.LaunchApp(appUuid)
			if err != nil {
				ll.Logger.Error(fmt.Sprintf("Error launching app: %v\n", err))
				return
			}
		}
	}
}

// Start the game client listening TCP server.
func (l *Launcher) StartGameClientListener() {

	// Start a TCP listener on the designated port for game client connections
	listener, err := net.Listen("tcp", "127.0.0.1:"+config.GamePort)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("Error starting listener: %v\n", err))
		return
	}

	// Close the listener when the function exits.
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("Error closing listener: %v\n", err))
		}
	}(listener)

	for {
		// Accept incoming connections from the game client.
		conn, err := listener.Accept()
		if err != nil {
			ll.Logger.Error(fmt.Sprintln("Error accepting game client connection:", err))
			continue
		}

		// Handle the connection with the game client in a goroutine.
		go l.handleGameClientConnection(conn)
	}
}

// Handle a connection from the game client.
func (l *Launcher) handleGameClientConnection(conn net.Conn) {
	// Close the connection and delete it from the connection pool when the goroutine exits.
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing game client connection:", err)
		}
	}(conn)

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message from game client:", err)
		return
	}

	var messageStruct struct {
		AppId string `json:"appId"`
	}
	err = json.Unmarshal([]byte(message), &messageStruct)
	if err != nil {
		fmt.Println("Error parsing message from game client:", err)
		return
	}

	appId := messageStruct.AppId
	gameConnPool[appId] = &conn
	ll.Logger.Print(fmt.Sprintf("Game client connected for app id: %s\n", appId))
}
