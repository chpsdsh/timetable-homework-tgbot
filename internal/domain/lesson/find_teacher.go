package lesson

type Teacher struct {
	shortName string
	fullName  string
	fullURL   string
}

func NewTeacher(shortName string, fullName string, fullURL string) Teacher {
	return Teacher{shortName: shortName, fullName: fullName, fullURL: fullURL}
}

func (t Teacher) GetShortName() string {
	return t.shortName
}

func (t Teacher) GetFullName() string {
	return t.fullName
}

func (t Teacher) GetFullURL() string {
	return t.fullURL
}
