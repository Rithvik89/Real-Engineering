package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RespBody struct {
	Message string `json:"message"`
	IsError bool   `json:"is_error"`
	Error   string `json:"error"`
}

func main() {
	// go object

	resp := RespBody{
		Message: "User Registered Successfully",
		IsError: false,
	}

	// using marshal

	jsonEncoded, err := json.Marshal(resp)
	if err != nil {
		error := fmt.Errorf("unable to encode with error, %s", err)
		fmt.Println(error)
	}

	fmt.Printf("Json Encoded byte slice from Marshal: %s \n", string(jsonEncoded))

	// using json encode

	json.NewEncoder(os.Stdout).Encode(resp)

}
