package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/crocotelementry/F1_GO/structs"
	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
)

var (
	header              structs.PacketHeader
	Motion_packet       structs.PacketMotionData
	Session_packet      structs.PacketSessionData
	Lap_packet          structs.PacketLapData
	redis_pool          = newPool() // newPool returns a pointer to a redis.Pool
	current_lap_number  = uint8(0)
	track_length        = uint16(0) // meters
	track_id            = int8(-2)  // since -1 is unknown and 0-21 are real ids, -2 is sufficiant
	total_packet_number = 0
	right_packet_number = 0
	left_packet_number  = 0

	addrs, _  = net.ResolveUDPAddr("udp", ":20777")
	sock, err = net.ListenUDP("udp", addrs)
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
		color.Green("Success\n")
	} else {
		color.Red("Error\n")
	}

	// Output: PONG
	return nil
}

// function that returns which side of the track the data is on
func which_side(current_lap_number) string {
	switch current_lap_number {
	case 2:
		return "right"
	case 4:
		return "left"
	default:
		log.Println("Asked which side for a lap number that was not 2 (right) or 4 (left)")
		return "error"
	}
}

// Function that returns the incrementing number for its corresponding side
func which_side_incrementing_number(current_lap_number) string {
	switch current_lap_number {
	case 2:
		return strconv.Itoa(right_packet_number)
	case 4:
		return strconv.Itoa(left_packet_number)
	default:
		log.Println("Asked incrementing number for its corresponding side for a lap number that was not 2 (right) or 4 (left)")
		return "error"
	}
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
	// incrementing_packet_number := 0
	// Create a reference point for our current lap. This is for adding things to catchup_packet number 2 stuff
	// current_lap_number := uint8(0)

	// fmt.Println("Track_saver started successfully")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("How to record track:")
	fmt.Println("     Lap 1:    Get feel for track")
	fmt.Println("     Lap 2:    Drive slowly around right side of track")
	fmt.Println("     Lap 3:    Buffer lap / prepare to switch sides")
	fmt.Println("     Lap 4:    Drive slowly around left side of track")
	fmt.Println("")
	fmt.Println("Once lap 4 is finished, track_saver will ask you if you are satisfied with your recording")
	fmt.Println("If yes is selected, track_saver will save the track then exit")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

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

			// We only want to record data when we are driving on the right during lap 2 or the left during lap 4
			if current_lap_number == 2 || current_lap_number == 4 {
				// If the packet we received is a motion_packet, read its binary into our motion_packet struct
				if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Motion_packet); err != nil {
					fmt.Println("binary.Read motion_packet failed:", err)
				}

				position_struct := structs.Position_struct{
					Frame_identifier: Motion_packet.M_header.M_frameIdentifier,
					Track_id:         track_id,
					Lap_number:       current_lap_number,
					Side:             which_side(current_lap_number),

					WorldPositionX: Motion_packet.M_carMotionData[Motion_packet.M_header.M_playerCarIndex].M_worldPositionX,
					WorldPositionY: Motion_packet.M_carMotionData[Motion_packet.M_header.M_playerCarIndex].M_worldPositionY,
					WorldPositionZ: Motion_packet.M_carMotionData[Motion_packet.M_header.M_playerCarIndex].M_worldPositionZ,
				}

				// Marshal the struct into json so we can save it in our redis database
				json_position_struct, err := json.Marshal(position_struct)
				if err != nil {
					fmt.Println(err)
				}

				if _, err := redis_conn.Do("SET", ("track_saver:" + which_side(current_lap_number) + ":" + which_side_incrementing_number(current_lap_number)), json_motion_packet); err != nil {

					fmt.Println("Adding track_saver data "+which_side(current_lap_number)+" to Redis database failed:", err)
					total_packet_number -= 1

					switch which_side(current_lap_number) {
					case "right":
						right_packet_number -= 1
					case "left":
						left_packet_number -= 1
					}
				}

				total_packet_number += 1

				switch which_side(current_lap_number) {
				case "right":
					right_packet_number += 1
				case "left":
					left_packet_number += 1
				}
			}

		case 1:
			// If the packet we received is the session_packet, read its binary into our session_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Session_packet); err != nil {
				fmt.Println("binary.Read session_packet failed:", err)
			}

			fmt.Println("Track length:", Session_packet.M_trackLength)

			// If this is our first session_packet received, make sure to save the track length
			if track_length == 0 {
				track_length = Session_packet.M_trackLength
			}

			// If this is our first session_packet received, make sure to save the M_trackId
			if track_id == -2 {
				track_id = Session_packet.M_trackId
			}

		case 2:
			// If the packet we received is the lap_packed, read its binary into our lap_packet struct
			if err := binary.Read(packet_bytes_reader, binary.LittleEndian, &Lap_packet); err != nil {
				fmt.Println("binary.Read Lap_packet failed:", err)
			}

			users_data := Lap_packet.M_lapData[Lap_packet.M_header.M_playerCarIndex]

			log.Println("Lap distance:", users_data.M_lapDistance)

			if current_lap_number == 0 {
				if users_data.M_currentLapNum != 1 {
					fmt.Println("")
					fmt.Println("Current Session alreday in progress!!!!")
					fmt.Println("Exiting Track_saver.....")
					fmt.Println("Please restart track_saver before starting a track session")
				}
			} else {

				if users_data.M_currentLapNum != current_lap_number {
					// When we get to a new lap
					switch users_data.M_currentLapNum {
					case 0:
						fmt.Println("Not sure if lap 0 is a thing? this is here just in case lol")

					case 1:
						fmt.Println("Lap 1: Get the feel for the track")
						fmt.Println("Next lap: You will be driving on the right side of the track at a slow speed")
						current_lap_number = 1

					case 2:
						fmt.Println("Lap 2: Drive slowly around right side of track")
						fmt.Println("Next lap: Buffer lap / prepare to switch sides")
						fmt.Println("")
						fmt.Println("Right side driving")
						fmt.Println("")
						fmt.Println("0%  |                                                                                                      | 100%")
						fmt.Printf("    ")
						current_lap_number = 2

					case 3:
						fmt.Println("Lap 3: Buffer lap / prepare to switch sides")
						fmt.Println("Next lap: You will be driving on the left side of the track at a slow speed")
						current_lap_number = 3

					case 4:
						fmt.Println("Lap 4: Drive slowly around left side of track")
						fmt.Println("")
						fmt.Println("Left side driving")
						fmt.Println("")
						fmt.Println("0%  |                                                                                                      | 100%")
						fmt.Printf("    ")
						current_lap_number = 4

					case 5:
						fmt.Println("Fourth Lap finished. Exiting track_saver.....")
						current_lap_number = 5
						// Do some saving stuff here and things :)
						redis_conn.Close()
						os.Exit(1)

					default:
						fmt.Println("")
						fmt.Println("")
						color.Red("Error\n")
						fmt.Println("Lap number not in range of 1-5 or some other error")
						fmt.Println("Exiting track_saver.....")
						redis_conn.Close()
						os.Exit(1)
					}
				} else {
					// When we are in middle of a lap
					switch users_data.M_currentLapNum {
					case 0:
						continue

					case 1:
						continue

					case 2:
						// right side driving
						progress := 100 * users_data.M_lapDistance / track_length

						if progress%2 == 0 {
							color.Green("|")
						}

					case 3:
						continue

					case 4:
						// left side driving
						progress := 100 * users_data.M_lapDistance / track_length

						if progress%2 == 0 {
							color.Green("|")
						}

					default:
						fmt.Println("")
						fmt.Println("")
						color.Red("Error\n")
						fmt.Println("Lap number not in range of 1-4 or some other error")
						fmt.Println("Exiting track_saver.....")
						redis_conn.Close()
						os.Exit(1)
					}
				}

			}

		default:
			continue
		}
	}
}
