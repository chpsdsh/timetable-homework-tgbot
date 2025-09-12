package lesson

type LessonTeacher struct {
	subject    string
	lessonType string
	groups     []string
	startTime  string
	weekday    string
	room       string
	week       string
}

func NewLessonTeacher(subject string, lessonType string, groups []string, startTime string, weekday string, room string, week string) LessonTeacher {
	return LessonTeacher{subject: subject, lessonType: lessonType, groups: groups, startTime: startTime, weekday: weekday, room: room, week: week}
}

func (l LessonTeacher) GetSubject() string {
	return l.subject
}

func (l LessonTeacher) GetLessonType() string {
	return l.lessonType
}

func (l LessonTeacher) GetTutor() []string {
	return l.groups
}

func (l LessonTeacher) GetStartTime() string {
	return l.startTime
}

func (l LessonTeacher) GetWeekday() string {
	return l.weekday
}

func (l LessonTeacher) GetRoom() string {
	return l.room
}
