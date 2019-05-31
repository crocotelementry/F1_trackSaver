package structs

type PacketHeader struct {
	M_packetFormat    uint16  // 2018
	M_packetVersion   uint8   // Version of this packet type, all start from 1
	M_packetId        uint8   // Identifier for the packet type, see below
	M_sessionUID      uint64  // Unique identifier for the session
	M_sessionTime     float32 // Session timestamp
	M_frameIdentifier uint32  // Identifier for the frame the data was retrieved on
	M_playerCarIndex  uint8   // Index of player's car in the array
}

// MOTION PACKET:
// The motion packet gives physics data for all the cars being driven.
// There is additional data for the car being driven with the goal of being able to drive a motion platform setup.
// Frequency: Rate as specified in menus
// Size: 1341 bytes
type CarMotionData struct {
	M_worldPositionX     float32 // World space X position
	M_worldPositionY     float32 // World space Y position
	M_worldPositionZ     float32 // World space Z position
	M_worldVelocityX     float32 // Velocity in world space X
	M_worldVelocityY     float32 // Velocity in world space Y
	M_worldVelocityZ     float32 // Velocity in world space Z
	M_worldForwardDirX   int16   // World space forward X direction (normalised)
	M_worldForwardDirY   int16   // World space forward Y direction (normalised)
	M_worldForwardDirZ   int16   // World space forward Z direction (normalised)
	M_worldRightDirX     int16   // World space right X direction (normalised)
	M_worldRightDirY     int16   // World space right Y direction (normalised)
	M_worldRightDirZ     int16   // World space right Z direction (normalised)
	M_gForceLateral      float32 // Lateral G-Force component
	M_gForceLongitudinal float32 // Longitudinal G-Force component
	M_gForceVertical     float32 // Vertical G-Force component
	M_yaw                float32 // Yaw angle in radians
	M_pitch              float32 // Pitch angle in radians
	M_roll               float32 // Roll angle in radians
}
type PacketMotionData struct {
	M_header PacketHeader // Header

	M_carMotionData [20]CarMotionData // Data for all cars on track

	// Extra player car ONLY data
	M_suspensionPosition     [4]float32 // Note: All wheel arrays have the following order:
	M_suspensionVelocity     [4]float32 // RL, RR, FL, FR
	M_suspensionAcceleration [4]float32 // RL, RR, FL, FR
	M_wheelSpeed             [4]float32 // Speed of each wheel
	M_wheelSlip              [4]float32 // Slip ratio for each wheel
	M_localVelocityX         float32    // Velocity in local space
	M_localVelocityY         float32    // Velocity in local space
	M_localVelocityZ         float32    // Velocity in local space
	M_angularVelocityX       float32    // Angular velocity x-component
	M_angularVelocityY       float32    // Angular velocity y-component
	M_angularVelocityZ       float32    // Angular velocity z-component
	M_angularAccelerationX   float32    // Angular velocity x-component
	M_angularAccelerationY   float32    // Angular velocity y-component
	M_angularAccelerationZ   float32    // Angular velocity z-component
	M_frontWheelsAngle       float32    // Current front wheels angle in radians
}

// SESSION PACKET:
// The session packet includes details about the current session in progress.
// Frequency: 2 per second
// Size: 147 bytes
type MarshalZone struct {
	M_zoneStart float32 // Fraction (0..1) of way through the lap the marshal zone starts
	M_zoneFlag  int8    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
}
type PacketSessionData struct {
	M_header PacketHeader // Header

	M_weather             uint8           // Weather - 0 = clear, 1 = light cloud, 2 = overcast, 3 = light rain, 4 = heavy rain, 5 = storm
	M_trackTemperature    int8            // Track temp. in degrees celsius
	M_airTemperature      int8            // Air temp. in degrees celsius
	M_totalLaps           uint8           // Total number of laps in this race
	M_trackLength         uint16          // Track length in metre
	M_sessionType         uint8           // 0 = unknown, 1 = P1, 2 = P2, 3 = P3, 4 = Short P, 5 = Q1, 6 = Q2, 7 = Q3, 8 = Short Q, 9 = OSQ, 10 = R, 11 = R2, 12 = Time Trial
	M_trackId             int8            // -1 for unknown, 0-21 for tracks, see appendix
	M_era                 uint8           // Era, 0 = modern, 1 = classic
	M_sessionTimeLeft     uint16          // Time left in session in seconds
	M_sessionDuration     uint16          // Session duration in seconds
	M_pitSpeedLimit       uint8           // Pit speed limit in kilometres per hour
	M_gamePaused          uint8           // Whether the game is paused
	M_isSpectating        uint8           // Whether the player is spectating
	M_spectatorCarIndex   uint8           // Index of the car being spectated
	M_sliProNativeSupport uint8           // SLI Pro support, 0 = inactive, 1 = active
	M_numMarshalZones     uint8           // Number of marshal zones to follow
	M_marshalZones        [21]MarshalZone // List of marshal zones – max 21
	M_safetyCarStatus     uint8           // 0 = no safety car, 1 = full safety car, 2 = virtual safety car
	M_networkGame         uint8           // 0 = offline, 1 = online
}

// LAP DATA PACKET:
// The lap data packet gives details of all the cars in the session.
// Frequency: Rate as specified in menus
// Size: 841 bytes
type LapData struct {
	M_lastLapTime       float32 // Last lap time in seconds
	M_currentLapTime    float32 // Current time around the lap in seconds
	M_bestLapTime       float32 // Best lap time of the session in seconds
	M_sector1Time       float32 // Sector 1 time in seconds
	M_sector2Time       float32 // Sector 2 time in seconds
	M_lapDistance       float32 // Distance vehicle is around current lap in metres – could be negative if line hasn’t been crossed yet
	M_totalDistance     float32 // Total distance travelled in session in metres – could be negative if line hasn’t been crossed yet
	M_safetyCarDelta    float32 // Delta in seconds for safety car
	M_carPosition       uint8   // Car race position
	M_currentLapNum     uint8   // Current lap number
	M_pitStatus         uint8   // 0 = none, 1 = pitting, 2 = in pit area
	M_sector            uint8   // 0 = sector1, 1 = sector2, 2 = sector3
	M_currentLapInvalid uint8   // Current lap invalid - 0 = valid, 1 = invalid
	M_penalties         uint8   // Accumulated time penalties in seconds to be added
	M_gridPosition      uint8   // Grid position the vehicle started the race in
	M_driverStatus      uint8   // Status of driver - 0 = in garage, 1 = flying lap, 2 = in lap, 3 = out lap, 4 = on track
	M_resultStatus      uint8   // Result status - 0 = invalid, 1 = inactive, 2 = active, 3 = finished, 4 = disqualified, 5 = not classified, 6 = retired
}
type PacketLapData struct {
	M_header PacketHeader // Header

	M_lapData [20]LapData // Lap data for all cars on track
}
