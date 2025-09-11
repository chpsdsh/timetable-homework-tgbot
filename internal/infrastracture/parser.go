package infrastracture

import (
	"net/http"
	"strings"
	"timetable-homework-tgbot/internal/domain/lesson"

	"github.com/PuerkitoBio/goquery"
)

func ParseLessonsStudent(url string) []lesson.LessonStudent {
	lessons := make([]lesson.LessonStudent, 0)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		startTime := strings.TrimSpace(tr.Find("td").First().Text())
		tr.Find("div.cell").Each(func(i int, cell *goquery.Selection) {
			lessonType := strings.TrimSpace(cell.Find("span.type").First().Text())
			subject := strings.TrimSpace(cell.Find("div.subject").First().Text())
			room := strings.TrimSpace(cell.Find("div.room a").First().Text())
			tutor := strings.TrimSpace(cell.Find("a.tutor").First().Text())
			lessons = append(lessons, lesson.NewLessonStudent(subject, lessonType, tutor, startTime, "ТУТ ДОЛЖЕН БЫТЬ ДЕНЬ НЕДЕЛИ", room))
		})
	})
	return lessons
}
