package main

import (
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture"
)

func main() {

	//teachers := infrastracture.ParseTeachers("Абдула")
	//infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/23204")
	lessons := infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/23204")
	fmt.Println(lessons)
	//fmt.Println(teachers)
}
