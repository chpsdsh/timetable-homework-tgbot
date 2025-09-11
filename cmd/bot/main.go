package main

import (
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture"
)

func main() {
	lessons := infrastracture.Parse("https://table.nsu.ru/group/25204")
	fmt.Println(lessons)
}
