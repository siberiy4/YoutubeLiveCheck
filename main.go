package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)
func main(){
	err:=godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	apiKey:=os.Getenv("YOUTUBE_API_KEY")
	fmt.Println(apiKey)

}
