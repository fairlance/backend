package importer

import (
	"fmt"
	"net/http"
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
	Timestamps          map[string]time.Time
	ImporterStarted     string
}

func newPage(r *http.Request) page {
	pageState := page{Message: "ok"}
	query := r.URL.Query()
	offset := 0
	if query.Get("offset") != "" {
		o, err := strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			pageState.Message = err.Error()
		}
		offset = int(o)
	}
	pageState.Offset = offset

	limit := 10
	if query.Get("limit") != "" {
		l, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			pageState.Message = err.Error()
		}
		limit = int(l)
	}
	pageState.Limit = limit

	pageState.Type = "jobs"
	if query.Get("type") != "" {
		pageState.Type = query.Get("type")
	}

	return pageState
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

func (p page) FormatTimeHuman(t time.Time) string {
	return humanDuration(time.Now().Sub(t)) + " ago"
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

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.)
func humanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours()); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*3 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%f years", d.Hours()/24/365)
}
