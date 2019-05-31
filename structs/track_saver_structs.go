package structs

type Position_struct struct {
  Frame_identifier uint32  // Identifier for the frame the data was retrieved on
  Track_id         int8    // -1 for unknown, 0-21 for tracks, see appendix
  Lap_number       uint8   // Current lap number
  Side             string

	WorldPositionX float32 // World space X position
	WorldPositionY float32 // World space Y position
	WorldPositionZ float32 // World space Z position
}

type Distance_struct struct {
  Frame_identifier uint32  // Identifier for the frame the data was retrieved on

  Distance         float32 // Distance vehicle is around current lap in metres – could be negative if line hasn’t been crossed yet
}
