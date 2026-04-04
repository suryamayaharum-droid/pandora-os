package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("=== IOT INTEGRATION SYSTEM ===")
	
	devices := []string{
		"Smart Sensors",
		"Security Cameras",
		"Smart Lights",
		"Temperature Sensors",
		"Door Locks",
		"Smart TVs",
		"Voice Assistants",
		"Smart Appliances",
		"GPS Trackers",
		"Weather Stations",
	}
	
	fmt.Println("\nDISPOSITIVOS IOT:")
	for i, d := range devices {
		status := []string{"online", "standby", "error"}[rand.Intn(3)]
		icon := "O"
		if status == "online" { icon = "*" }
		fmt.Printf("  [%s] %d. %s (%s)\n", icon, i+1, d, status)
	}
	
	// MQTT topics
	fmt.Println("\nMQTT TOPICS:")
	topics := []string{
		"home/sensors/temperature",
		"home/sensors/motion",
		"home/lights/bedroom",
		"home/security/door",
		"home/+/status",
	}
	for _, t := range topics {
		fmt.Printf("  - %s\n", t)
	}
	
	// Data flow simulation
	fmt.Println("\nDADOS RECEBIDOS:")
	for i := 0; i < 5; i++ {
		val := 20 + rand.Float32()*10
		fmt.Printf("  -> temp: %.1f°C | hum: %d%% | motion: %v\n", 
			val, rand.Intn(100), rand.Float32() > 0.5)
	}
	
	fmt.Println("\nOK - IOT INTEGRATION PRONTA!")
}
