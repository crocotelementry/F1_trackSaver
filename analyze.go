package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/F1_trackSaver/structs"
	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"
)

var (
	position_struct structs.Position_struct

	right_max_x = 0
	right_max_y = 0
	right_max_z = 0
	left_max_x  = 0
	left_max_y  = 0
	left_max_z  = 0

	right_min_x = 0
	right_min_y = 0
	right_min_z = 0
	left_min_x  = 0
	left_min_y  = 0
	left_min_z  = 0
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

func find_min_and_max(redis_conn redis.Conn, string side, int incrementing_number) error {

  position_struct := new(structs.Position_struct)

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








func track_analyzer(redis_conn redis.Conn, int left_side_incrementing_number, int right_side_incrementing_number) {

  // set the min and max to the first position data
  err = set_min_and_max(redis_conn)

  // find the min and max for the right side of the track
  err = find_min_and_max(redis_conn, "right", right_side_incrementing_number)

  // find the min and max for the left side of the track
  err = find_min_and_max(redis_conn, "left", left_side_incrementing_number)




  // By now we should have our min and max values for both sides of the track, lets tell the user
  // that information!
  fmt.Println("Track general information:")
  fmt.Println("")
  fmt.Println("Minimum values:")
  fmt.Println("    X:       right:", right_max_x, "   left:", left_max_x)
  fmt.Println("    Y:       right:", right_max_y, "   left:", left_max_y)
  fmt.Println("    Z:       right:", right_max_z, "   left:", left_max_z)
  fmt.Println("")
  fmt.Println("Maximum values:")
  fmt.Println("    X:       right:", right_min_x, "   left:", left_min_x)
  fmt.Println("    Y:       right:", right_min_y, "   left:", left_min_y)
  fmt.Println("    Z:       right:", right_min_z, "   left:", left_min_z)

}
