package main

import (
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture"
)

func main() {
	lessons := infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/25204")
	teachers := infrastracture.ParseTeachers("Абдула")
	fmt.Println(lessons)
	fmt.Println(teachers)
}
