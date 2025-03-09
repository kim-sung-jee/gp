package mapstructure

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type ServiceConfig struct {
	Name string `json:"name"`
}

func ParseJson() {
	jsonData := `{"name": "UserService"}`

	var config ServiceConfig
	json.Unmarshal([]byte(jsonData), &config)
	fmt.Println([]byte(jsonData))
	fmt.Println("Decoded Name:", config.Name)

	encoded, _ := json.Marshal(config)
	fmt.Println(encoded)
	fmt.Println("Encoded JSON:", string(encoded))

	json.Unmarshal([]byte(encoded), &config)
	encoded2, _ := json.Marshal(config)
	fmt.Println(encoded2)
}

type ServiceMapStructure struct {
	Name string `json:"name" json:"hi" mapstructure:"name"`
}

func ParseMapStructure() {
	data := map[string]interface{}{
		"name": "UserService",
	}

	data2 := `{"name": "UserService", "hi": "Hello"}`

	// 맵 데이터를 구조체로 변환
	var mapService ServiceMapStructure
	mapstructure.Decode(data, &mapService)

	fmt.Println("Decoded Name:", mapService.Name)

	var serviceConfig ServiceMapStructure
	json.Unmarshal([]byte(data2), &serviceConfig)
	fmt.Println("Decoded Name:", serviceConfig.Name)
	fmt.Println("Decoded Config:", serviceConfig.Name)
	fmt.Println("Decoded Config:", serviceConfig)

}
