package lib

import (
	"net/http"
	"strconv"
)

func GetPageAndLimitFromRequest(r *http.Request) (int, int, error) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return 0, 0, err
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return 0, 0, err
	}
	return page, limit, nil
}
