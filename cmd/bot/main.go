package main

import (
	"fmt"
	infrastracture "timetable-homework-tgbot/internal/infrastracture/parser"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	log.Println(".env not found")
	//}
	//
	//if err := app.Run(); err != nil {
	//	log.Fatal(err)
	//}
	groups := infrastracture.ParseGroups("https://table.nsu.ru/faculty/ggf")
	fmt.Println(groups)
}
