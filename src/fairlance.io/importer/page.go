package importer

import (
	"strconv"
	"time"

	"fairlance.io/application"
)

type page struct {
	Message             string
	Entities            map[string]interface{}
	Document            map[string]interface{}
	Offset              int
	Limit               int
	TotalInDB           int
	TotalInSearchEngine int
	Type                string
}

func (p page) PrevPageLabel() string {
	if p.Offset >= p.Limit {
		return strconv.Itoa(p.Offset-p.Limit+1) + "-" + strconv.Itoa(p.Offset)
	}

	return ""
}

func (p page) NextPageLabel() string {
	if p.Offset+p.Limit < p.TotalInDB {
		return strconv.Itoa(p.Offset+p.Limit+1) + "-" + strconv.Itoa(p.Offset+(p.Limit*2))
	}

	return ""
}

func (p page) CurrentPageLabel() string {
	return strconv.Itoa(p.Offset+1) + "-" + strconv.Itoa(p.Offset+p.Limit)
}

func (p page) PrevPageURL() string {
	if p.Offset >= p.Limit {
		return "?offset=" + strconv.Itoa(p.Offset-p.Limit) + "&limit=" + strconv.Itoa(p.Limit)
	}

	return ""
}

func (p page) NextPageURL() string {
	if p.Offset+p.Limit < p.TotalInDB {
		return "?offset=" + strconv.Itoa(p.Offset+p.Limit) + "&limit=" + strconv.Itoa(p.Limit)
	}

	return ""
}

func (p page) CurrentPageURL() string {
	return "?offset=" + strconv.Itoa(p.Offset) + "&limit=" + strconv.Itoa(p.Limit)
}

func (p page) FormatTime(t time.Time) string {
	return t.Format(time.RFC822)
}

func (p page) GetName(doc interface{}) string {
	switch doc.(type) {
	case application.Freelancer:
		f := doc.(application.Freelancer)
		return f.FirstName + " " + f.LastName
	case application.Job:
		j := doc.(application.Job)
		return j.Name
	default:
		return ""
	}
}
