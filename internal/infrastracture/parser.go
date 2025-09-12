package infrastracture

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"timetable-homework-tgbot/internal/domain/lesson"

	"github.com/PuerkitoBio/goquery"
)

func ParseLessonsStudent(url string) []lesson.LessonStudent {
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
				lessons = append(lessons, lesson.NewLessonStudent(subject, lessonType, tutor, startTime, weekday, room, week))
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

func ParseTeachers() []lesson.Teacher {
	const page = "https://table.nsu.ru/teacher"
	teachers := make([]lesson.Teacher, 0)
	resp, err := http.Get(page)
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
		fullURL := abs(page, href)
		teachers = append(teachers, lesson.NewTeacher(shortName, fullName, fullURL))
	})
	return teachers
}

func ParseLessonsTeacher(url string) []lesson.LessonTeacher {
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
				lessons = append(lessons, lesson.NewLessonTeacher(subject, lessonType, groups, startTime, weekday, room, week))
			})
		}
	})
	return lessons
}
