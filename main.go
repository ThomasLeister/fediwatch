package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/oschwald/geoip2-golang"
	"google.golang.org/protobuf/proto"

	fediwatchProto "thomas-leister.de/fediwatch/fediwatchProto"
)

/*
 * Const section
 */
const (
	WsMessageBufferSize = 9999
)

// Version string var
var versionString string = "0.0.0"

/*
 * Define ConnDir enum type (In, Out)
 */

type ConnDir int64

const (
	In ConnDir = iota
	Out
)

func (connDir ConnDir) String() string {
	if connDir == In {
		return "in"
	} else if connDir == Out {
		return "out"
	} else {
		return "(invalid)"
	}
}

type Config struct {
	HomeLocation  []float32 `toml:"homeLocation"`
	HttpPort      int16     `toml:"httpPort"`
	WebsocketPort int16     `toml:"websocketPort"`
	RedisHost     string    `toml:"redisHost"`
	RedisPort     int16     `toml:"redisPort"`
	DatabasePath  string    `toml:"databasePath"`
	WebsocketUrl  string    `toml:"websocketUrl"`
}

type Coordinates struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

type WebSettings struct {
	HomeLocation Coordinates `json:"homeLocation"`
	WebsocketUrl string      `json:"websocketUrl"`
}

/*
 * Global app vars
 */
var (
	config   Config
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
	broadcast    = make(chan []byte, WsMessageBufferSize)
	quit         = make(chan bool)
)

/*
 * Observe websocket connections and log strange disconnects
 */
func checkIfWsClosed(conn *websocket.Conn, connClosed chan bool) {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		connClosed <- true
	}()

	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		_, _, err := conn.ReadMessage()

		if err != nil {
			// If Connection is closed, we will receive an error here
			// We only want to log Strange errors, but not simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
	}
}

func handleWsConnection(w http.ResponseWriter, r *http.Request) {
	connClosed := make(chan bool)

	// Upgrade HTTP connection to Websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// Add new client to list
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	fmt.Println("Created a new Websocket connection")

	// In another Goroutine: Check if conn was closed by constantly trying to read from conn
	go checkIfWsClosed(conn, connClosed)

	// Wait for new message in broadcast channel or
	// Close command in connClosed channel
	for {
		select {
		case <-connClosed:
			// "websocket close" received. remove conn.
			fmt.Println("Received WS close!")
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
		case message := <-broadcast:
			// Receive a message and distribute it across all
			// websocket clients. Not send this to "this" / own WS client, but all,
			// since a broadcast message is lost once it has been withdrawn from the broadcast channel.
			for client := range clients {
				clientsMutex.Lock()
				err := client.WriteMessage(websocket.BinaryMessage, message)
				clientsMutex.Unlock()
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func handleWebSettings(w http.ResponseWriter, r *http.Request) {
	// Create a response struct with the home location
	response := WebSettings{
		HomeLocation: Coordinates{Lat: config.HomeLocation[0], Long: config.HomeLocation[1]},
		WebsocketUrl: config.WebsocketUrl,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and write it to the response writer
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*
 * Resolve an IP address to coordinates using the GeoIP database
 */
func resolveToCoordinates(ipAddress string, databasePath string) (float64, float64, error) {
	db, err := geoip2.Open(databasePath)
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()

	// Parse IP address
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return 0, 0, fmt.Errorf("Invalid IP address: %s\n", ipAddress)
	}

	// Get IP info
	record, err := db.City(ip)
	if err != nil {
		return 0, 0, err
	}

	// Extract GPS coordinates
	latitude := record.Location.Latitude
	longitude := record.Location.Longitude

	return latitude, longitude, nil
}

/*
 * Resolve an hostname to IP address(es) using a DNS lookup
 */
func resolveToIp(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		return "", err
	}

	// Select a random IP address
	randomIndex := rand.Intn(len(ips))
	randomIP := ips[randomIndex]

	return randomIP.String(), err
}

func handleRedisMessages(redisChannel <-chan *redis.Message, connDir ConnDir) {
	var protoConnDir fediwatchProto.Connection_Direction

	for {
		select {
		case msg := <-redisChannel:
			randIP, err := resolveToIp(msg.Payload)
			if err != nil {
				continue
			}

			latitude, longitude, err := resolveToCoordinates(randIP, config.DatabasePath)
			if err != nil {
				fmt.Printf("Failed to resolve IP address to location: %v\n", err)
				continue
			}

			// Output on console
			if connDir == In {
				protoConnDir = fediwatchProto.Connection_IN
				fmt.Printf("IN connection to: %s\t %s\t Latitude %.4f, Longitude %.4f\n", msg.Payload, randIP, latitude, longitude)
			} else {
				protoConnDir = fediwatchProto.Connection_OUT
				fmt.Printf("OUT connection to: %s\t %s\t Latitude %.4f, Longitude %.4f\n", msg.Payload, randIP, latitude, longitude)
			}

			// Feed into websocket broadcaster
			if len(broadcast) != cap(broadcast) {
				myConn := &fediwatchProto.Connection{
					Dir: protoConnDir,
					Lat: float32(latitude),
					Lng: float32(longitude),
				}

				broadcastBytes, err := proto.Marshal(myConn)
				if err != nil {
					fmt.Println("Encoding Connection object failed!")
					continue
				}

				// Send Protobuf bytes to websocket
				broadcast <- broadcastBytes
			}
		case <-time.After(time.Second * 5):
			fmt.Printf("Waiting for new %s connection\n", connDir)
		}
	}
}

func startStaticFilesServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	serverAddr := fmt.Sprintf("localhost:%d", config.HttpPort)
	log.Printf("HTTP server listening at %s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Failed to start the HTTP server: ", err)
	}
}

func startWebsocketServer() {
	http.HandleFunc("/ws", handleWsConnection)
	http.HandleFunc("/settings", handleWebSettings)

	serverAddr := fmt.Sprintf("localhost:%d", config.WebsocketPort)
	fmt.Printf("WebSocket server is listening at %s\n", serverAddr)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("Failed to start the Websocket server: ", err)
	}
}

func main() {
	// Greeting
	fmt.Printf("Starting FediWatch %s\n", versionString)

	// Read config file
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		log.Fatal("Could not find config file")
		os.Exit(1)
	}

	if _, err := toml.DecodeFile(f, &config); err != nil {
		log.Fatalf("Could not parse config file: %s", err)
		os.Exit(1)
	}

	// Configure Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()
	ctx := context.Background()

	// Create channels for inbound and outbound connections
	pubsub_inbound := redisClient.Subscribe(ctx, "activitypub-inbound-hosts")
	pubsub_outbound := redisClient.Subscribe(ctx, "activitypub-outbound-hosts")

	// Check both Redis pubsub channels for new messages
	ch_inbound := pubsub_inbound.Channel(redis.WithChannelSize(WsMessageBufferSize))
	ch_outbound := pubsub_outbound.Channel(redis.WithChannelSize(WsMessageBufferSize))

	fmt.Println("Waiting for Redis messages ...")

	go handleRedisMessages(ch_inbound, In)
	go handleRedisMessages(ch_outbound, Out)

	// Start static files server and websocket server...
	go startStaticFilesServer()
	go startWebsocketServer()

	<-quit
}
