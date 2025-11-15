package lesson

type LessonType int

const (
	StudentLesson LessonType = iota
	TeacherLesson
	RoomLesson
)
