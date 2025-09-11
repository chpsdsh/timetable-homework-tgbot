package lesson

type LessonTutor struct {
	subject    string
	lessonType string
	tutor      string
	startTime  string
	weekday    string
	room       string
}

func NewLesson(subject string, lessonType string, tutor string, startTime string, room string) LessonTutor {
	return LessonTutor{subject: subject, lessonType: lessonType, tutor: tutor, startTime: startTime, room: room}
}

func (l LessonTutor) GetSubject() string {
	return l.subject
}

func (l LessonTutor) GetLessonType() string {
	return l.lessonType
}

func (l LessonTutor) GetTutor() string {
	return l.tutor
}

func (l LessonTutor) GetStartTime() string {
	return l.startTime
}

func (l LessonTutor) GetWeekday() string {
	return l.weekday
}

func (l LessonTutor) GetRoom() string {
	return l.room
}
