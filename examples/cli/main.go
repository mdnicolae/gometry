package main

import (
	"fmt"
	"github.com/mdnicolae/gometry"
)

func main() {
	//If identifier will not be provided, it will use the one marked as default in the gometry.json file
	gmy, err := gometry.Init()

	if err != nil {
		fmt.Printf("Failed to initialize gmy: %v\n", err)
		return
	}

	gmy.Info("This is an info message", map[string]interface{}{"user": "john_doe"})
	gmy.Debug("Debugging details", map[string]interface{}{"module": "auth"})
	gmy.Warning("This is a warning", map[string]interface{}{"disk_space": "low"})
	gmy.Error("An error occurred", map[string]interface{}{"error_code": 500})
	gmy.Info("No additional params")
	gmy.Critical("This is really bad", map[string]interface{}{"error_code": 500})

	gmyNoColor, err := gometry.Init("cli-no-color")
	if err != nil {
		fmt.Printf("Failed to initialize gmyNoColor: %v\n", err)
		return
	}

	gmyNoColor.Info("This is an info message", map[string]interface{}{"user": "john_doe"})
	gmyNoColor.Error("This is an info message", map[string]interface{}{"user": "john_doe"})
}
