package lesson

type LessonRoom struct {
	subject    string
	lessonType string
	tutor      string
	startTime  string
	weekday    string
	groups     []string
	week       string
}

func NewLessonRoom(subject string, lessonType string, tutor string, startTime string, weekday string, groups []string, week string) *LessonRoom {
	return &LessonRoom{subject: subject, lessonType: lessonType, tutor: tutor, startTime: startTime, weekday: weekday, groups: groups, week: week}
}

func (l LessonRoom) GetSubject() string {
	return l.subject
}

func (l LessonRoom) GetLessonType() string {
	return l.lessonType
}

func (l LessonRoom) GetTutor() string {
	return l.tutor
}

func (l LessonRoom) GetStartTime() string {
	return l.startTime
}

func (l LessonRoom) GetWeekday() string {
	return l.weekday
}

func (l LessonRoom) GetGroups() []string {
	return l.groups
}

func (l LessonRoom) GetWeek() string {
	return l.week
}
