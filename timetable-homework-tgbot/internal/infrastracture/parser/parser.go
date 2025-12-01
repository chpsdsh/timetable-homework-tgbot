package httpparser

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"timetable-homework-tgbot/internal/domain/lesson"
	"timetable-homework-tgbot/internal/domain/urlselector"

	"github.com/PuerkitoBio/goquery"
)

const (
	teachersPage  = "https://table.nsu.ru/teacher"
	roomPage      = "https://table.nsu.ru/room"
	facultiesPage = "https://table.nsu.ru/faculties"
	nsuPage       = "https://table.nsu.ru"
)

type Parser struct {
}

var instance *Parser

var once sync.Once

func GetParser() *Parser {
	once.Do(func() {
		instance = &Parser{}
	})
	return instance
}

func (p *Parser) ParseLessonsStudent(groupUrl string) []lesson.LessonStudent {
	page := nsuPage + groupUrl
	resp, err := http.Get(page)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	table := doc.Find("table.time-table").First()
	if table == nil {
		log.Fatal("table not found")
	}

	days := make([]string, 0)

	table.Find("tr").First().Find("th").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		days = append(days, s.Text())
	})

	lessons := make([]lesson.LessonStudent, 0)

	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}

		tds := tr.Find("td")
		startTime := strings.TrimSpace(tds.Eq(0).Text())

		for col := 1; col < tds.Length(); col++ {
			weekdayIdx := col - 1
			if weekdayIdx < 0 || weekdayIdx >= len(days) {
				continue
			}
			weekday := days[weekdayIdx]
			td := tds.Eq(col)
			cells := td.Find("div.cell")
			cells.Each(func(i int, cell *goquery.Selection) {
				lessonType := strings.TrimSpace(cell.Find("span.type").First().Text())
				subject := strings.TrimSpace(cell.Find("div.subject").First().Text())
				room := strings.TrimSpace(cell.Find("div.room a").First().Text())
				tutor := strings.TrimSpace(cell.Find("a.tutor").First().Text())
				week := strings.TrimSpace(cell.Find("div.week").First().Text())
				lessons = append(lessons, *lesson.NewLessonStudent(subject, lessonType, tutor, startTime, weekday, room, week))
			})
		}
	})
	return lessons
}

func abs(base, href string) string {
	bu, _ := url.Parse(base)
	ru, _ := url.Parse(href)
	return bu.ResolveReference(ru).String()
}

func (p *Parser) ParseTeachers() []urlselector.Teacher {
	teachers := make([]urlselector.Teacher, 0)
	resp, err := http.Get(teachersPage)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	doc.Find("a.tutors_item").Each(func(i int, s *goquery.Selection) {

		shortName := strings.TrimSpace(s.Text())
		fullName, _ := s.Attr("title")
		href, _ := s.Attr("href")
		fullURL := abs(teachersPage, href)
		teachers = append(teachers, urlselector.Teacher{ShortName: shortName, FullName: fullName, FullURL: fullURL})
	})
	return teachers
}

func (p *Parser) ParseRooms() []urlselector.Room {
	rooms := make([]urlselector.Room, 0)
	resp, err := http.Get(roomPage)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("a.tutors_item").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		fullURL := abs(roomPage, href)
		rooms = append(rooms, urlselector.Room{Room: name, FllURL: fullURL})
	})
	return rooms
}

func (p *Parser) ParseFaculties() []urlselector.Faculty {
	faculties := make([]urlselector.Faculty, 0)
	resp, err := http.Get(facultiesPage)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("a.faculty").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		fullURL := abs(facultiesPage, href)
		faculties = append(faculties, urlselector.Faculty{FacultyName: name, FullUrl: fullURL})
	})
	return faculties
}

func (p *Parser) ParseGroups(facultyURL string) []urlselector.Group {
	groups := make([]urlselector.Group, 0)
	seen := make(map[string]struct{})

	resp, err := http.Get(facultyURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	doc.Find("a.group").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")

		key := name + "|" + href
		if _, ok := seen[key]; ok {
			// уже добавляли такую группу → пропускаем
			return
		}
		seen[key] = struct{}{}

		groups = append(groups, urlselector.Group{
			GroupName: name,
			GroupUrl:  href,
		})
	})

	return groups
}

func (p *Parser) ParseLessonsTeacher(url string) []lesson.LessonTeacher {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	table := doc.Find("table.time-table").First()
	if table == nil {
		log.Fatal("table not found")
	}

	days := make([]string, 0)

	table.Find("tr").First().Find("th").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		days = append(days, s.Text())
	})

	lessons := make([]lesson.LessonTeacher, 0)

	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}

		tds := tr.Find("td")
		startTime := strings.TrimSpace(tds.Eq(0).Text())

		for col := 1; col < tds.Length(); col++ {
			weekdayIdx := col - 1
			if weekdayIdx < 0 || weekdayIdx >= len(days) {
				continue
			}
			weekday := days[weekdayIdx]
			td := tds.Eq(col)
			cells := td.Find("div.cell")
			cells.Each(func(i int, cell *goquery.Selection) {
				lessonType := strings.TrimSpace(cell.Find("span.type").First().Text())
				subject := strings.TrimSpace(cell.Find("div.subject").First().Text())
				room := strings.TrimSpace(cell.Find("div.room a").First().Text())
				groups := make([]string, 0)
				cells.Find("div.groups").Find("a.group").Each(func(i int, group *goquery.Selection) {
					groups = append(groups, strings.TrimSpace(group.Text()))
				})
				week := strings.TrimSpace(cell.Find("div.week").First().Text())
				lessons = append(lessons, *lesson.NewLessonTeacher(subject, lessonType, groups, startTime, weekday, room, week))
			})
		}
	})
	return lessons
}

func (p *Parser) ParseLessonsRoom(url string) []lesson.LessonRoom {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	table := doc.Find("table.time-table").First()
	if table == nil {
		log.Fatal("table not found")
	}

	days := make([]string, 0)

	table.Find("tr").First().Find("th").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		days = append(days, s.Text())
	})

	lessons := make([]lesson.LessonRoom, 0)

	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}

		tds := tr.Find("td")
		startTime := strings.TrimSpace(tds.Eq(0).Text())

		for col := 1; col < tds.Length(); col++ {
			weekdayIdx := col - 1
			if weekdayIdx < 0 || weekdayIdx >= len(days) {
				continue
			}
			weekday := days[weekdayIdx]
			td := tds.Eq(col)
			cells := td.Find("div.cell")
			cells.Each(func(i int, cell *goquery.Selection) {
				lessonType := strings.TrimSpace(cell.Find("span.type").First().Text())
				subject := strings.TrimSpace(cell.Find("div.subject").First().Text())
				teacher := strings.TrimSpace(cell.Find("a.tutor").First().Text())
				groups := make([]string, 0)
				cells.Find("div.groups").Find("a.group").Each(func(i int, group *goquery.Selection) {
					groups = append(groups, strings.TrimSpace(group.Text()))
				})
				week := strings.TrimSpace(cell.Find("div.week").First().Text())
				lessons = append(lessons, *lesson.NewLessonRoom(subject, lessonType, teacher, startTime, weekday, groups, week))
			})
		}
	})
	return lessons
}
