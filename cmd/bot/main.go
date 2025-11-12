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
	lessons := infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/23204")
	fmt.Println(lessons)
}
