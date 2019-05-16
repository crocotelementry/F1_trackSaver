package main

import (
  "bytes"
  "encoding/binary"

  "github.com/crocotelementry/F1_GO/structs"
  "github.com/gomodule/redigo/redis"
)

var (
	header                                 structs.PacketHeader
	Motion_packet                          structs.PacketMotionData
	Session_packet                         structs.PacketSessionData
  redis_ping_done                        = make(chan bool)
	redis_pool                             = newPool() // newPool returns a pointer to a redis.Pool
  incrementing_motion_packet_number      = 0
	incrementing_session_packet_number     = 0
)








// To establish connectivity in redigo, you need to create a redis.Pool object which is a pool of connections to Redis.
func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// ping tests connectivity for redis (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	// fmt.Println("PING Response = ", s)
	fmt.Print("Redis connection       ")

	if s == "PONG" {
		redis_ping_done <- true
		color.Green("Success")
	} else {
		redis_ping_done <- false
		color.Red("Error")
	}

	// Output: PONG
	return nil
}

func main() {

  // get a connection from the pool (redis.Conn)
	redis_conn := redis_pool.Get()
	// use defer to close the connection when the function completes

	c := make(chan os.Signal, 2)
	// When we close F1_go by using control-c, we will catch it, flush the redis database empty and then
	// close the connection to the redis database before closing F1_GO
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		redis_conn.Do("FlushAll")
		fmt.Println("\n")
		log.Println("               redis flushed")
		redis_conn.Close()
		os.Exit(1)
	}()

	defer func() {
		redis_conn.Close()
	}()

	// call Redis PING command to test connectivity
	err := ping(redis_conn)
	if err != nil {
		fmt.Println("Problem with connection to Redis database", err)
	}

	// Set number of SETs to redis database to zero
	incrementing_packet_number := 0
  // Create a reference point for our current lap. This is for adding things to catchup_packet number 2 stuff
	current_lap_number := uint8(0)




	for {
		buf := make([]byte, 1341)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("readfromudp error::: ", err)
		}

		// Set a new reader which we will use to cast into our structs.
		// This reader is for the header, which we determine what packet we have and what index our users car is in.
		// Bytes 3 in the udp packet will be the packet number and byte 20 will be the index of the users car.
		header_bytes_reader := bytes.NewReader(buf[0:21])
		packet_bytes_reader := bytes.NewReader(buf)

		// Read the binary of the udp packet header into our struct
		if err := binary.Read(header_bytes_reader, binary.LittleEndian, &header); err != nil {
			fmt.Println("binary.Read header failed:", err)
		}

		// Depending on which packet we have, which we find by looking at header.M_packetId
		// We use a switch statement to then read the whole binary udp packet into its associated struct

		switch header.M_packetId {
		case 0:
			// If the packet we received is a motion_packet, read its binary into our motion_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Motion_packet); err != nil {
				fmt.Println("binary.Read motion_packet failed:", err)
			}

      // Marshal the struct into json so we can save it in our redis database
			json_motion_packet, err := json.Marshal(Motion_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Motion_packet.M_header.M_sessionUID, 10) + ":0:" + strconv.Itoa(incrementing_motion_packet_number)), json_motion_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_motion_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_motion_packet_number += 1

    case 1:
			// If the packet we received is the session_packet, read its binary into our session_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Session_packet); err != nil {
				fmt.Println("binary.Read session_packet failed:", err)
			}

      // Marshal the struct into json so we can save it in our redis database
			json_session_packet, err := json.Marshal(Session_packet)
			if err != nil {
				fmt.Println(err)
			}

			if _, err := redis_conn.Do("SET", (strconv.FormatUint(Session_packet.M_header.M_sessionUID, 10) + ":1:" + strconv.Itoa(incrementing_session_packet_number)), json_session_packet); err != nil {
				fmt.Println("Adding json_motion_packet to Redis database failed:", err)
				incrementing_packet_number -= 1
				incrementing_session_packet_number -= 1
			}
			incrementing_packet_number += 1
			incrementing_session_packet_number += 1

		default:
			continue
		}
	}
}
