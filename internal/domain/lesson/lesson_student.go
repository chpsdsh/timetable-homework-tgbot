package lesson

type LessonStudent struct {
	subject    string
	lessonType string
	tutor      string
	startTime  string
	weekday    string
	room       string
	week       string
}

func NewLessonStudent(subject string, lessonType string, tutor string, startTime string, weekday string, room string, week string) LessonStudent {
	return LessonStudent{subject: subject, lessonType: lessonType, tutor: tutor, startTime: startTime, weekday: weekday, room: room, week: week}
}

func (l LessonStudent) GetSubject() string {
	return l.subject
}

func (l LessonStudent) GetLessonType() string {
	return l.lessonType
}

func (l LessonStudent) GetTutor() string {
	return l.tutor
}

func (l LessonStudent) GetStartTime() string {
	return l.startTime
}

func (l LessonStudent) GetWeekday() string {
	return l.weekday
}

func (l LessonStudent) GetRoom() string {
	return l.room
}

func (l LessonStudent) GetWeek() string {
	return l.week
}
