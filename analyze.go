package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
  "math"

	"github.com/F1_trackSaver/structs"
	// "github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
)

var (
  compress_track_struct structs.Compress_track_struct

	right_max_x = float32(0)
	right_max_y = float32(0)
	right_max_z = float32(0)
	left_max_x  = float32(0)
	left_max_y  = float32(0)
	left_max_z  = float32(0)

	right_min_x = float32(0)
	right_min_y = float32(0)
	right_min_z = float32(0)
	left_min_x  = float32(0)
	left_min_y  = float32(0)
	left_min_z  = float32(0)
)

func set_min_and_max(redis_conn redis.Conn) error {
	left_data_initial_position_struct := new(structs.Position_struct)
	right_data_initial_position_struct := new(structs.Position_struct)

	right_side_data_initial, err := (redis_conn.Do("GET", "track_saver:right:0"))
	if err != nil {
		fmt.Println("", err)
		return err
	}

	left_side_data_initial, err := (redis_conn.Do("GET", "track_saver:left:0"))
	if err != nil {
		fmt.Println("", err)
		return err
	}

	err = json.Unmarshal(right_side_data_initial.([]byte), &left_data_initial_position_struct)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(left_side_data_initial.([]byte), &right_data_initial_position_struct)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Set the max and min values for x, y, z on the first save
	right_max_x = right_data_initial_position_struct.WorldPositionX
	right_max_y = right_data_initial_position_struct.WorldPositionY
	right_max_z = right_data_initial_position_struct.WorldPositionZ
	right_min_x = right_data_initial_position_struct.WorldPositionX
	right_min_y = right_data_initial_position_struct.WorldPositionY
	right_min_z = right_data_initial_position_struct.WorldPositionZ

	left_max_x = left_data_initial_position_struct.WorldPositionX
	left_max_y = left_data_initial_position_struct.WorldPositionY
	left_max_z = left_data_initial_position_struct.WorldPositionZ
	left_min_x = left_data_initial_position_struct.WorldPositionX
	left_min_y = left_data_initial_position_struct.WorldPositionY
	left_min_z = left_data_initial_position_struct.WorldPositionZ

	return nil
}

func find_min_and_max(redis_conn redis.Conn, side string, incrementing_number int) error {

	position_structs := new(structs.Position_struct)

	// We start at 1 since we already got the data from 0 for our initial setting for min and max
	for struct_number := 1; struct_number < incrementing_number; struct_number += 1 {

		position_data, err := (redis_conn.Do("GET", "track_saver"+":"+side+":"+strconv.Itoa(struct_number)))
		if err != nil {
			log.Println("Getting position data for ", side, " side unsuccesfull on struct_number ", struct_number, err)
			return err
		}

		err = json.Unmarshal(position_data.([]byte), &position_structs)
		if err != nil {
			log.Println(err)
			return err
		}

		switch side {
		case "right":
			// start comparing to our min and maxs to see if we found a new min or a new max

			// x
			if position_structs.WorldPositionX > right_max_x {
				right_max_x = position_structs.WorldPositionX
			} else {
				if position_structs.WorldPositionX < right_min_x {
					right_min_x = position_structs.WorldPositionX
				}
			}

			// y
			if position_structs.WorldPositionY > right_max_y {
				right_max_y = position_structs.WorldPositionY
			} else {
				if position_structs.WorldPositionY < right_min_y {
					right_min_y = position_structs.WorldPositionY
				}
			}

			// z
			if position_structs.WorldPositionZ > right_max_z {
				right_max_z = position_structs.WorldPositionZ
			} else {
				if position_structs.WorldPositionZ < right_min_z {
					right_min_z = position_structs.WorldPositionZ
				}
			}

		case "left":
			// start comparing to our min and maxs to see if we found a new min or a new max

			// x
			if position_structs.WorldPositionX > left_max_x {
				left_max_x = position_structs.WorldPositionX
			} else {
				if position_structs.WorldPositionX < left_min_x {
					left_min_x = position_structs.WorldPositionX
				}
			}

			// y
			if position_structs.WorldPositionY > left_max_y {
				left_max_y = position_structs.WorldPositionY
			} else {
				if position_structs.WorldPositionY < left_min_y {
					left_min_y = position_structs.WorldPositionY
				}
			}

			// z
			if position_structs.WorldPositionZ > left_max_z {
				left_max_z = position_structs.WorldPositionZ
			} else {
				if position_structs.WorldPositionZ < left_min_z {
					left_min_z = position_structs.WorldPositionZ
				}
			}

		}
		// end of switch

	}
	// end of for loop

	// If no errors, executes cleanly, return nil
	return nil
}



func compress_track(redis_conn redis.Conn, side string, incrementing_number int) error {

  position_structs := new(structs.Position_struct)
  distance_structs := new(structs.Distance_struct)

  current_distance := 0

  var shared_distance_positions_x []float32

  var shared_distance_positions_y []float32

  var shared_distance_positions_z []float32


  fmt.Println("Compressing track:", side)


  // We start at 1 since we already got the data from 0 for our initial setting for min and max
	for struct_number := 1; struct_number < incrementing_number; struct_number += 1 {

		position_data, err := (redis_conn.Do("GET", "track_saver"+":"+side+":"+strconv.Itoa(struct_number)))
		if err != nil {
			log.Println("Getting position data for ", side, " side unsuccesfull on struct_number ", struct_number, err)
			return err
		}

		err = json.Unmarshal(position_data.([]byte), &position_structs)
		if err != nil {
			log.Println(err)
			return err
		}


    distance_data, err := (redis_conn.Do("GET", "distance:" + strconv.FormatUint(uint64(position_structs.Frame_identifier), 10)))
    if err != nil {
      fmt.Println("Getting distance data for", side, "on lap number", current_lap_number, "from frame number:", position_structs.Frame_identifier, "from Redis database failed:", err)
      return err
    }

    err = json.Unmarshal(distance_data.([]byte), &distance_structs)
		if err != nil {
			log.Println(err)
			return err
		}


    // Now we should have the position data and its distance around the track for a specific frame of data
    rounded_distance := math.Round(float64(distance_structs.Distance))

    if int(rounded_distance) != current_distance {
      temp_x_value := float32(0)
      temp_y_value := float32(0)
      temp_z_value := float32(0)

      // average last distances values
      for _, value := range shared_distance_positions_x {
        temp_x_value += value
      }

      for _, value := range shared_distance_positions_y {
        temp_y_value += value
      }

      for _, value := range shared_distance_positions_z {
        temp_z_value += value
      }

      temp_x_value = temp_x_value / float32(len(shared_distance_positions_x))
      temp_y_value = temp_y_value / float32(len(shared_distance_positions_y))
      temp_z_value = temp_z_value / float32(len(shared_distance_positions_z))


      // add distance to redis
      compress_track_struct := structs.Compress_track_struct{
        Frame_identifier: position_structs.Frame_identifier,
        Side: side,
        Distance: int(rounded_distance),
        WorldPositionX: temp_x_value,
        WorldPositionY: temp_y_value,
        WorldPositionZ: temp_z_value,
      }

      // Marshal the struct into json so we can save it in our redis database
      json_compress_track_struct, err := json.Marshal(compress_track_struct)
      if err != nil {
        fmt.Println(err)
        return err
      }

      if _, err := redis_conn.Do("SET", ("compress_track:" + side + ":" + strconv.Itoa(int(rounded_distance))), json_compress_track_struct); err != nil {
        fmt.Println("Adding compress_track data from side", side, "on lap number", current_lap_number, "from distance:", rounded_distance, " to Redis database failed:", err)
        return err
      }


      // Now set our slices to nil
      shared_distance_positions_x = nil
      shared_distance_positions_y = nil
      shared_distance_positions_z = nil

      // set the current_distance to the new distance
      current_distance = int(rounded_distance)
    }

    shared_distance_positions_x = append(shared_distance_positions_x, position_structs.WorldPositionX)
    shared_distance_positions_y = append(shared_distance_positions_y, position_structs.WorldPositionY)
    shared_distance_positions_z = append(shared_distance_positions_z, position_structs.WorldPositionZ)


  }

  return nil
}



func print_compress_track_data(redis_conn redis.Conn) error{

  compress_track_struct_right := new(structs.Compress_track_struct)
  compress_track_struct_left := new(structs.Compress_track_struct)

  for distance_data_compress_interval := 0; distance_data_compress_interval < int(track_length); distance_data_compress_interval += 1 {

    distance_data_compress_right, err := redis_conn.Do("GET", ("compress_track:right:" + strconv.Itoa(distance_data_compress_interval)))
    if err != nil {
      fmt.Println("getting compress_track data from side", "right", "on distance:", distance_data_compress_interval, "from Redis database failed:", err)
      return err
    }

    distance_data_compress_left, err := redis_conn.Do("GET", ("compress_track:left:" + strconv.Itoa(distance_data_compress_interval)))
    if err != nil {
      fmt.Println("getting compress_track data from side", "left", "on distance:", distance_data_compress_interval, "from Redis database failed:", err)
      return err
    }

    err = json.Unmarshal(distance_data_compress_right.([]byte), &compress_track_struct_right)
    if err != nil {
      log.Println(err)
      return err
    }

    err = json.Unmarshal(distance_data_compress_left.([]byte), &compress_track_struct_left)
		if err != nil {
			log.Println(err)
			return err
		}

    fmt.Println("")
    fmt.Println("Distance:", distance_data_compress_interval)
    fmt.Println("left X:", compress_track_struct_left.WorldPositionX, "   right X:", compress_track_struct_right.WorldPositionX)
    fmt.Println("left Y:", compress_track_struct_left.WorldPositionY, "   right Y:", compress_track_struct_right.WorldPositionY)
    fmt.Println("left Z:", compress_track_struct_left.WorldPositionZ, "   right Z:", compress_track_struct_right.WorldPositionZ)
  }

  return nil
}




func analyse_track(redis_conn redis.Conn, right_side_incrementing_number int, left_side_incrementing_number int) {

	// set the min and max to the first position data
	err = set_min_and_max(redis_conn)
  if err != nil {
    fmt.Println(err)
  }

	// find the min and max for the right side of the track
	err = find_min_and_max(redis_conn, "right", right_side_incrementing_number)
  if err != nil {
    fmt.Println(err)
  }


	// find the min and max for the left side of the track
	err = find_min_and_max(redis_conn, "left", left_side_incrementing_number)
  if err != nil {
    fmt.Println(err)
  }


	// By now we should have our min and max values for both sides of the track, lets tell the user
	// that information!
	fmt.Println("Track general information:")
	fmt.Println("")
	fmt.Println("Maximum values:")
	fmt.Println("    X:       right:", right_max_x, "   left:", left_max_x)
	fmt.Println("    Y:       right:", right_max_y, "   left:", left_max_y)
	fmt.Println("    Z:       right:", right_max_z, "   left:", left_max_z)
	fmt.Println("")
	fmt.Println("Minimum values:")
	fmt.Println("    X:       right:", right_min_x, "   left:", left_min_x)
	fmt.Println("    Y:       right:", right_min_y, "   left:", left_min_y)
	fmt.Println("    Z:       right:", right_min_z, "   left:", left_min_z)


  fmt.Println("")
  fmt.Println("")

  err = compress_track(redis_conn, "right", right_side_incrementing_number)
  if err != nil {
    fmt.Println(err)
  }

  err = compress_track(redis_conn, "left", left_side_incrementing_number)
  if err != nil {
    fmt.Println(err)
  }

  // lets see what our compressed tracks data looks like
  err = print_compress_track_data(redis_conn)
  if err != nil {
    fmt.Println(err)
  }


}
