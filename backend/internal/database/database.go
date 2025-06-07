package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thyamix/festival-counter/internal/models"
)

var DB *sql.DB

func InitDB() {
	fmt.Println("Started init DB")
	_, err := os.Stat("./data")
	if os.IsNotExist(err) {
		os.Mkdir("./data", os.ModeDir|0755)
	}

	fmt.Println("Created /data file")

	if err != nil {
		log.Println(err)
	}

	DB, err = sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created .db file")

	err = initTables()
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}

	fmt.Println("Created tables for db file")

	fmt.Println("Database started")
	fmt.Println("Successfully init the DB")
}

func initTables() error {

	sqlBytes, err := os.ReadFile("static/init.sql")
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		return err
	}

	return nil
}

func GetTotal(festivalCode string) (int, int, error) {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = (SELECT id FROM festival WHERE code = ?)", festivalCode).Scan(&eventId)

	var total, maxGauge sql.NullInt64

	err := DB.QueryRow("SELECT SUM(value) FROM active WHERE event_id = ?", eventId).Scan(&total)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get sum of data: %v", err)
	}

	err = DB.QueryRow(`SELECT MAX(max) FROM max 
		WHERE time = (SELECT MAX(time) from max WHERE event_id = ?)
		AND event_id = ?`, eventId, eventId).Scan(&maxGauge)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get most recent max: %v", err)
	}

	if total.Valid {
		return int(total.Int64), int(maxGauge.Int64), nil
	}
	return 0, 0, nil
}

func AddValue(value int, festivalCode string) error {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = (SELECT id FROM festival WHERE code = ?)", festivalCode).Scan(&eventId)

	time := time.Now().Unix()
	_, err := DB.Exec("INSERT INTO active (value, time, event_id) VALUES (?, ?, ?)", value, time, eventId)

	if err != nil {
		return fmt.Errorf("failed to add %d to db error: %v", value, err)
	}

	return nil
}

func ChangeMax(newMax int) {
	var fid int
	DB.QueryRow("SELECT id FROM festival WHERE active = 1").Scan(&fid)
	time := time.Now().Unix()
	_, err := DB.Exec("INSERT INTO max (max, time, fid) VALUES (?, ?, ?)", newMax, time, fid)

	if err != nil {
		log.Printf("Failed to change max to %d in db:\nError: %v\n", newMax, err)
	}
}

func IsValidFestivalId(festivalId int) bool {
	var result int
	DB.QueryRow("SELECT COUNT(1) FROM festival WHERE id = ?", festivalId).Scan(&result)
	return (result != 1)
}

func IsValidEventId(eventId int) bool {
	var result int
	DB.QueryRow("SELECT COUNT(1) FROM event WHERE id = ?", eventId).Scan(&result)
	return (result != 1)
}

func GetFestival(festivalCode string) (*models.FestivalData, error) {
	var festival models.FestivalData
	err := DB.QueryRow("SELECT id, expires_at, created_at, pin, password, code FROM festival WHERE code = ?", festivalCode).Scan(&festival.Id, &festival.ExpiresAt, &festival.CreatedAt, &festival.Pin, &festival.PasswordHash, &festival.Code)
	if err != nil {
		return nil, fmt.Errorf("this one here: %v", err)
	}
	return &festival, nil
}

func Reset(festivalId int) int {
	var oldEventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = ?", festivalId).Scan(&oldEventId)
	time := time.Now().Unix()

	DB.Exec("UPDATE event SET active = 0, last_used_at = ? WHERE id = ?", time, oldEventId)

	DB.Exec("INSERT INTO event (created_at, last_used_at, festival_id, active) VALUES (?, ?, ?, ?)", time, time, festivalId, 1)

	archiveEvent(oldEventId)

	return oldEventId
}

func archiveInactive() {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1").Scan(&eventId)

	DB.Exec("INSERT INTO archive (value, time, event_id) SELECT value, time, event_id FROM active WHERE event_id != ?", eventId)
	DB.Exec("DELETE FROM active WHERE festival_id != ?", eventId)
}

func archiveEvent(eventId int) {
	DB.Exec("INSERT INTO archive (value, time, event_id) SELECT value, time, event_id FROM active WHERE event_id = ?", eventId)
	DB.Exec("DELETE FROM active WHERE event_id = ?", eventId)
}

func GetArchivedFestivalsIdsTimes() ([]int, []int) {
	var id, time int
	var ids, times []int

	rows, err := DB.Query("SELECT id, time FROM festival WHERE active = 0 ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &time)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
		times = append(times, time)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return ids, times
}

func getNumberOfEntries(id int) int {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM archive WHERE festival_id = ?", id).Scan(&count)
	return count
}

/*
Takes in the id of festival as well as the chunk (10k entries) number index 0
Return the up to 10k values from chunk or nil if out of range and bool to say if there is more
*/
func GetFestivalEventEntriesChunk(festivalId int, eventId int, chunk int, archived bool) ([][]string, bool) {
	const CHUNKSIZE = 10000
	count := getNumberOfEntries(festivalId)
	numberOfChunks := count / CHUNKSIZE
	if numberOfChunks < chunk {
		return nil, false
	}

	var offset = chunk * CHUNKSIZE

	var rows *sql.Rows
	var err error

	if archived {
		rows, err = DB.Query("SELECT value, time FROM archive WHERE event_id = ? ORDER BY id LIMIT ? OFFSET ?", eventId, CHUNKSIZE, offset)
	} else {
		rows, err = DB.Query("SELECT value, time FROM active WHERE event_id = ? ORDER BY id LIMIT ? OFFSET ?", eventId, CHUNKSIZE, offset)
	}
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var output [][]string
	idCount := CHUNKSIZE*chunk + 1

	for rows.Next() {
		var row []string
		var values [4]int

		values[0] = idCount
		idCount += 1

		if err = rows.Scan(&values[1], &values[2]); err != nil {
			log.Fatal(err)
		}
		if err = DB.QueryRow(`SELECT max FROM max WHERE id = (
			SELECT MAX(id) FROM max WHERE event_id = ? AND time <= ?)`, eventId, values[2]).Scan(&values[3]); err != nil {
			values[3] = 0
		}

		for i := range values {
			row = append(row, fmt.Sprintf("%v", values[i]))
		}

		output = append(output, row)
	}

	return output, (count - (chunk+1)*CHUNKSIZE) > 0
}

func CreateUser() (int, error) {
	var userId int64
	result, err := DB.Exec(`INSERT INTO user DEFAULT VALUES`)
	if err != nil {
		return -1, err
	}
	userId, err = result.LastInsertId()

	return int(userId), err
}

func CreateRefreshToken(token models.RefreshToken) error {
	_, err := DB.Exec("INSERT INTO refresh_token (token, expires_at, last_used_at, user_id, revoked) VALUES (? ,? ,? ,? ,? )", token.Token, token.ExpiresAt, token.LastUsedAt, token.UserId, token.Revoked)
	return err
}

func CreateAccessToken(token models.AccessToken) error {
	_, err := DB.Exec("INSERT INTO access_token (token, expires_at, user_id) VALUES (? ,? ,?)", token.Token, token.ExpiresAt, token.UserId)
	return err
}

func GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := DB.QueryRow(`
	SELECT id, user_id, token, expires_at, last_used_at, revoked FROM refresh_token WHERE token = ?`, token).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.LastUsedAt,
		&refreshToken.Revoked,
	)
	return &refreshToken, err
}

func GetAccessToken(token string) (*models.AccessToken, error) {
	var accessToken models.AccessToken
	err := DB.QueryRow("SELECT id, token, expires_at, user_id FROM access_token WHERE token = ?", token).Scan(
		&accessToken.Id,
		&accessToken.Token,
		&accessToken.ExpiresAt,
		&accessToken.UserId,
	)
	return &accessToken, err
}

/*
Update the Access token in the database, only expiresAt, lastUsedAt and revoked can be changed
*/
func UpdateRefreshToken(token models.RefreshToken) error {
	_, err := DB.Exec("UPDATE refresh_token SET expires_at = ?, last_used_at = ?, revoked = ? WHERE id = ?", token.ExpiresAt, token.LastUsedAt, token.Revoked, token.Id)
	return err
}

/*
Update the Access token in the database, only expiresAt can be changed
*/
func UpdateAccessToken(token models.AccessToken) error {
	_, err := DB.Exec("UPDATE refresh_token SET expires_at = ? WHERE id = ?", token.ExpiresAt, token.Id)
	return err
}

func CreateFestival(festival models.FestivalData) (int, error) {
	result, err := DB.Exec("INSERT INTO festival (code, password, pin, created_at, expires_at) VALUES (?, ?, ?, ?, ?)", festival.Code, festival.PasswordHash, festival.Pin, festival.CreatedAt, festival.ExpiresAt)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	_, err = DB.Exec("INSERT INTO event (created_at, last_used_at, festival_id, active) VALUES (?, ?, ?, ?)", festival.CreatedAt, festival.CreatedAt, id, 1)
	if err != nil {
		return -1, err
	}
	return int(id), err
}

func IsNewFestivalCode(code string) bool {
	var id int
	err := DB.QueryRow("SELECT id FROM festival WHERE code = ?", code).Scan(id)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}

func AddFestivalAccess(accessToken string, festival models.FestivalData) error {
	var id int
	err := DB.QueryRow("SELECT user_id FROM access_token WHERE token = ?", accessToken).Scan(&id)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO festival_access (festival_id, user_id, last_used_at) VALUES (? ,? ,? )", festival.Id, id, festival.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetEventTotalAndMax(festivalId int) (int, int, error) {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = ?", festivalId).Scan(&eventId)

	var total, maxGauge sql.NullInt64

	err := DB.QueryRow("SELECT SUM(value) FROM active WHERE event_id = ?", eventId).Scan(&total)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to get sum of data: %v", err)
	}

	err = DB.QueryRow(`SELECT max FROM max 
		WHERE time = (SELECT MAX(time) from max WHERE event_id = ?)`, festivalId).Scan(&maxGauge)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to get most recent max: %v", err)
	}

	if total.Valid {
		return int(total.Int64), int(maxGauge.Int64), nil
	}
	return 0, 0, nil
}

func GetFestivalAccess(userId int, festivalId int) (*models.FestivalAccess, error) {
	var access models.FestivalAccess
	err := DB.QueryRow("SELECT id, festival_id, last_used_at, user_id FROM festival_access WHERE user_id = ? AND festival_id = ?", userId, festivalId).Scan(&access.Id, &access.FestivalId, &access.LastUsedAt, &access.UserId)
	if err != nil {
		return nil, err
	}
	_, err = DB.Exec("UPDATE festival_access SET last_used_at = ? WHERE id = ?", time.Now().Unix(), access.Id)
	if err != nil {
		log.Println("failed to update festival access last used at: ", err)
	}
	return &access, nil
}

func GetArchivedEventsIdsTimes(festivalId int) ([]int, []int, error) {
	var id, time int
	var ids, times []int

	rows, err := DB.Query("SELECT id, last_used_at FROM event WHERE active = 0 AND festival_id = ? ORDER BY id DESC", festivalId)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &time)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
		times = append(times, time)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return ids, times, nil
}

func ChangeCurrentEventMax(festivalCode string, newMax int) {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = (SELECT id FROM festival WHERE code = ?)", festivalCode).Scan(&eventId)
	time := time.Now().Unix()
	_, err := DB.Exec("INSERT INTO max (max, time, event_id) VALUES (?, ?, ?)", newMax, time, eventId)

	if err != nil {
		log.Printf("Failed to change max to %d in db:\nError: %v\n", newMax, err)
	}
}

func GetActiveEventId(festivalCode string) int {
	var eventId int
	DB.QueryRow("SELECT id FROM event WHERE active = 1 AND festival_id = (SELECT id FROM festival WHERE code = ?)", festivalCode).Scan(&eventId)
	return eventId
}
