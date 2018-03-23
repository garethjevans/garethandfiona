package rsvp

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"github.com/go-sql-driver/mysql"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS wedding DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE wedding;`,
	`CREATE TABLE IF NOT EXISTS rsvp (
		id        INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		rsvp_id   CHAR(36) NOT NULL,
		rsvp_date DATETIME NULL,
		email     VARCHAR(255) NOT NULL,
		name      VARCHAR(255) NOT NULL,
		comments  VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS guests (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		rsvp_id CHAR(36) NOT NULL,
        attending TINYINT(1) NOT NULL,
		name VARCHAR(255) NOT NULL,
		comments VARCHAR(255) NULL,
		PRIMARY KEY (id)
	)`,
}

// mysqlDB persists Rsvps to a MySQL instance.
type mysqlDB struct {
	conn        *sql.DB

	get         *sql.Stmt
	getGuests   *sql.Stmt
	update      *sql.Stmt
	updateGuest *sql.Stmt
}

// Ensure mysqlDB conforms to the WeddingDatabase interface.
var _ WeddingDatabase = &mysqlDB{}

type MySQLConfig struct {
	Username, Password string
	Host string
	Port int
	UnixSocket string
}

// dataStoreName returns a connection string suitable for sql.Open.
func (c MySQLConfig) dataStoreName(databaseName string) string {
	var cred string
	// [username[:password]@]
	if c.Username != "" {
		cred = c.Username
		if c.Password != "" {
			cred = cred + ":" + c.Password
		}
		cred = cred + "@"
	}

	if c.UnixSocket != "" {
		return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, databaseName)
	}
	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, c.Host, c.Port, databaseName)
}

// newMySQLDB creates a new WeddingDatabase backed by a given MySQL server.
func NewMySQLDB(config MySQLConfig) (WeddingDatabase, error) {
	if err := config.ensureTableExists(); err != nil {
		return nil, err
	}

	conn, err := sql.Open("mysql", config.dataStoreName("wedding"))
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
	}

	db := &mysqlDB{
		conn: conn,
	}

	if db.get, err = conn.Prepare(getStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare get: %v", err)
	}
	if db.getGuests, err = conn.Prepare(getGuestsStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare getGuests: %v", err)
	}
	if db.update, err = conn.Prepare(updateStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare update: %v", err)
	}
	if db.updateGuest, err = conn.Prepare(updateGuestStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare updateGuest: %v", err)
	}

	return db, nil
}

func (db *mysqlDB) Close() {
	db.conn.Close()
}

func (db *mysqlDB) Exec(statement string) (sql.Result, error) {
	return db.conn.Exec(statement)
}

func (db *mysqlDB) DB() *sql.DB {
	return db.conn
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanRsvp(s rowScanner) (*Rsvp, error) {
	var (
		id            int64
		rsvp_id       sql.NullString
		email         sql.NullString
		name          sql.NullString
		comments      sql.NullString
	)

	if err := s.Scan(&id, &rsvp_id, &email, &name, &comments); err != nil {
		return nil, err
	}

	Rsvp := &Rsvp{
		ID:           id,
		RsvpID:       rsvp_id.String,
		Email:        email.String,
		Name:         name.String,
		Comments:     comments.String,
	}
	return Rsvp, nil
}

func scanGuest(s rowScanner) (*Guest, error) {
	var (
		id            int64
		rsvp_id       sql.NullString
		name          sql.NullString
		attending     sql.NullBool
		comments      sql.NullString
	)

	if err := s.Scan(&id, &rsvp_id, &name, &attending, &comments); err != nil {
		return nil, err
	}

	Guest := &Guest{
		ID:           id,
		RsvpID:       rsvp_id.String,
		Name:         name.String,
		Attending:    attending.Bool,
		Comments:     comments.String,
	}
	return Guest, nil
}

const getStatement = `SELECT id,rsvp_id,email,name,comments FROM rsvp WHERE rsvp_id = ?`
const getGuestsStatement = `SELECT id,rsvp_id,name,attending,comments FROM guests WHERE rsvp_id = ?`

// GetRsvp retrieves a Rsvp by its ID.
func (db *mysqlDB) GetRsvp(id string) (*Rsvp, error) {
	log.Printf("GetRsvp(%s)", id)
	Rsvp, err := scanRsvp(db.get.QueryRow(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mysql: could not find Rsvp with rsvp_id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get Rsvp: %v", err)
	}

	guests, err := db.getGuestsByRsvpId(id)
	if err != nil {
		log.Printf("got error reading guests - %v", err)
		return nil, err
	}

	log.Printf("Got Guests %s", guests)

	Rsvp.Guests = guests
	return Rsvp, nil
}

// GetGuestsByRsvp retrieves a list of guests by its rsvp id.
func (db *mysqlDB) getGuestsByRsvpId(id string) ([]*Guest, error) {
	log.Printf("getGuestsByRsvpId(%s)", id)
	rows, err := db.getGuests.Query(id)

	if err != nil {
		return nil, fmt.Errorf("mysql: could not get Guests: %v", err)
	}

	var guests []*Guest
	for rows.Next() {
		guest, err := scanGuest(rows)
		if err != nil {
			return nil, fmt.Errorf("mysql: could not read row: %v", err)
		}

		log.Printf("converted %s", guest)
		guests = append(guests, guest)
	}

	return guests, nil
}

const updateStatement = `UPDATE rsvp SET email=?, name=?, comments=? WHERE id = ? AND rsvp_id = ?`
const updateGuestStatement = `UPDATE guests SET attending = ?, name = ?, comments = ? WHERE id = ? AND rsvp_id = ?`

// UpdateRsvp updates the entry for a given Rsvp.
func (db *mysqlDB) UpdateRsvp(b *Rsvp) error {
	log.Printf("UpdateRsvp(%s)", b)
	if b.ID == 0 {
		return errors.New("mysql: Rsvp with unassigned ID passed into updateRsvp")
	}

	if b.RsvpID == "" {
		return errors.New("mysql: Rsvp with unassigned RsvpID passed into updateRsvp")
	}

	_, err := execAffectingOneRow(db.update, b.Email, b.Name, b.Comments, b.ID, b.RsvpID)

    for _, guest := range b.Guests {
		db.updateGuestByGuest(guest)
	}

	return err
}

func (db *mysqlDB) updateGuestByGuest(b *Guest) error {
	log.Printf("updateGuestByGuest(%s)", b)
	if b.ID == 0 {
		return errors.New("mysql: Guest with unassigned ID passed into updateGuest")
	}

	if b.RsvpID == "" {
		return errors.New("mysql: Guest with unassigned RsvpID passed into updateGuest")
	}

	_, err := execAffectingOneRow(db.updateGuest,b.Attending, b.Name, b.Comments, b.ID, b.RsvpID)
	return err
}

// ensureTableExists checks the table exists. If not, it creates it.
func (config MySQLConfig) ensureTableExists() error {
	conn, err := sql.Open("mysql", config.dataStoreName(""))
	if err != nil {
		return fmt.Errorf("mysql: could not get a connection: %v", err)
	}
	defer conn.Close()

	// Check the connection.
	if conn.Ping() == driver.ErrBadConn {
		return fmt.Errorf("mysql: could not connect to the database. " +
			"could be bad address, or this address is not whitelisted for access.")
	}

	if _, err := conn.Exec("USE wedding"); err != nil {
		// MySQL error 1049 is "database does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
			return createTable(conn)
		}
	}

	if _, err := conn.Exec("DESCRIBE rsvp"); err != nil {
		// MySQL error 1146 is "table does not exist"
		if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
			return createTable(conn)
		}
		// Unknown error.
		return fmt.Errorf("mysql: could not connect to the database: %v", err)
	}
	return nil
}

// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	r, err := stmt.Exec(args...)
	if err != nil {
		return r, fmt.Errorf("mysql: could not execute statement: %v", err)
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
	} else if rowsAffected != 1 {
		return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
	}
	return r, nil
}
