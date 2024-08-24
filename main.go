package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	IPAddress string `json:"ip_address"`
}

var users []User

func main() {
	db := ConnectToDB()
	if db.DB == nil {
		log.Println("Error in connecting to DB")
		return
	}
	ReadJSON()
	db.Insert(users)
}

func ReadJSON() {
	file, err := os.Open("data.json")
	if err != nil {
		log.Println("Error opening JSON file: ", err)
		return
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading JSON file: ", err)
		return
	}

	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		log.Println("Error unmarshalling JSON: ", err)
		return
	}
}
