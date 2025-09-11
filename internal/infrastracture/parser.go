package infrastracture

import (
	"net/http"
	"strings"
	"timetable-homework-tgbot/internal/domain/lesson"

	"github.com/PuerkitoBio/goquery"
)

func ParseLessons(url string) []lesson.Lesson {
	lessons := make([]lesson.Lesson, 0)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		startTime := strings.TrimSpace(tr.Find("td").First().Text())
		tr.Find("div.cell").Each(func(i int, cell *goquery.Selection) {
			lessonType := strings.TrimSpace(cell.Find("span.type").First().Text())
			subject := strings.TrimSpace(cell.Find("div.subject").First().Text())
			room := strings.TrimSpace(cell.Find("div.room a").First().Text())
			tutor := strings.TrimSpace(cell.Find("a.tutor").First().Text())
			lessons = append(lessons, lesson.NewLesson(subject, lessonType, tutor, startTime, room))
		})
	})
	return lessons
}
