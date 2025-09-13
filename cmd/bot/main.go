package main

import (
	"fmt"
	"timetable-homework-tgbot/internal/infrastracture"
)

func main() {
	lessons := infrastracture.ParseLessonsStudent("https://table.nsu.ru/group/23204")
	fmt.Println(lessons)
	teachers := infrastracture.ParseTeachers()
	fmt.Println(teachers)
	lessonsT := infrastracture.ParseLessonsTeacher("https://table.nsu.ru/teacher/1d3a63ba-083f-11e6-8153-000c29b4927a#slixdrrgyb")
	fmt.Println(lessonsT)
	rooms := infrastracture.ParseRooms()
	fmt.Println(rooms)
	lessonsR := infrastracture.ParseLessonsRoom("https://table.nsu.ru/room/3107#tltiyfkqde")
	fmt.Println(lessonsR)
}
