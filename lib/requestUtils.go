package lib

import (
	"database/sql"
	"log"
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
func ConnectPostgres(conStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println("Err open db connection:", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
