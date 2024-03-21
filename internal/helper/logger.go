package helper

import (
	"encoding/json"
	"fmt"
)

func LogPretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b) + "\n")
}
