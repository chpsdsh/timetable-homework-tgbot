package lesson

type Lesson struct {
	subject    string
	lessonType string
	tutor      string
	startTime  string
	weekday    string
	room       string
}

func NewLesson(subject string, lessonType string, tutor string, startTime string, room string) Lesson {
	return Lesson{subject: subject, lessonType: lessonType, tutor: tutor, startTime: startTime, room: room}
}

func (l Lesson) Getsubject() string {
	return l.subject
}

func (l Lesson) GetLessonType() string {
	return l.lessonType
}

func (l Lesson) GetTutor() string {
	return l.tutor
}

func (l Lesson) GetStartTime() string {
	return l.startTime
}

func (l Lesson) GetWeekday() string {
	return l.weekday
}

func (l Lesson) GetRoom() string {
	return l.room
}
