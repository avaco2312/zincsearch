package main

import (
	"fmt"
	"log"
	"prueba/zinc"
)

func main() {
	res, err := zinc.Query("conference", 0, 1)
	if err != nil {
		log.Fatal(err)
	}
	for _, email := range res.Hits.Hits {
		fmt.Println(email.Source.Id)
		fmt.Println(email.Source.Subject)
		fmt.Println(email.Source.From)
		fmt.Println(email.Source.To)
		fmt.Println(email.Source.Content)
	}
}
