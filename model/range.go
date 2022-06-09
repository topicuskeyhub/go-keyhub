package model

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

const (
	HEADER_REQ_RANGE  = "Range"
	HEADER_COUNT_MODE = "topicus-Range-Count"
	HEADER_RESP_RANGE = "Content-Range"
	RANGE_PER_PAGE    = 100
)

var contentRangeRegex *regexp.Regexp

func NewRange() *Range {

	if contentRangeRegex == nil {
		contentRangeRegex, _ = regexp.Compile("items (\\d+)-(\\d+)/(\\d+)")
	}

	r := &Range{}
	r.Setup(RANGE_PER_PAGE)

	return r
}

type Range struct {
	start   int
	end     int
	total   int
	perpage int
}

func (r *Range) Setup(perpage int) {
	r.start = 0
	r.end = perpage - 1
	r.total = perpage
	r.perpage = perpage
}

func (r *Range) ParseResponse(response *http.Response) {

	if response == nil {
		// Reset to page 1
		r.Setup(r.perpage)
	}

	matches := contentRangeRegex.FindStringSubmatch(response.Header.Get(HEADER_RESP_RANGE))
	// matches[0]  = all, matches[1]  = start (0 based), matches[2]  = end (0 based), matches[3]  = total

	if len(matches) > 0 {
		r.start, _ = strconv.Atoi(matches[1])
		r.end, _ = strconv.Atoi(matches[2])
		r.total, _ = strconv.Atoi(matches[3])
	}

}

func (r *Range) NextPage() bool {

	if r.end+1 >= r.total {
		return false
	}

	r.start = r.end + 1
	r.end = r.start + r.perpage

	return true
}

func (r *Range) PreviousPage() bool {

	if r.start <= 0 {
		return false
	}

	r.end = r.start - 1
	r.start = r.end - r.perpage

	// If start gets lower than 0, reset to page 0 values
	if r.start < 0 {
		r.Setup(r.perpage)
	}

	return true
}

func (r *Range) GetRequestRangeHeader() (string, string) {

	return HEADER_REQ_RANGE, fmt.Sprintf("items=%d-%d", r.start, r.end)
}
func (r *Range) GetRequestModeHeader() (string, string) {

	return HEADER_COUNT_MODE, "Exact"
}
