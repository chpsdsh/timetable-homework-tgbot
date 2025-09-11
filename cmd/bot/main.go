package main

import (
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture"
)

func main() {
	lessons := infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/25204")
	fmt.Println(lessons)
	return
}
