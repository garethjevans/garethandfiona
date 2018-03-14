package rsvp

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS wedding DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE wedding;`,
	`CREATE TABLE IF NOT EXISTS rsvp (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		rsvp_id CHAR(36) NOT NULL,
		email VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS attendee (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		rsvp_id CHAR(36) NOT NULL,
        attending VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		dietry_requirements VARCHAR(255) NULL,
		wine VARCHAR(255) NULL,
		PRIMARY KEY (id)
	)`,
}

// mysqlDB persists Rsvps to a MySQL instance.
type mysqlDB struct {
	conn *sql.DB

	insert *sql.Stmt
	get    *sql.Stmt
	update *sql.Stmt
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
	if db.insert, err = conn.Prepare(insertStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare insert: %v", err)
	}
	if db.update, err = conn.Prepare(updateStatement); err != nil {
		return nil, fmt.Errorf("mysql: prepare update: %v", err)
	}

	return db, nil
}

func (db *mysqlDB) Close() {
	db.conn.Close()
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanRsvp(s rowScanner) (*Rsvp, error) {
	var (
		id            int64
		rsvp_id       sql.NullString
		email         sql.NullString
	)

	if err := s.Scan(&id, &rsvp_id, &email); err != nil {
		return nil, err
	}

	Rsvp := &Rsvp{
		ID:           id,
		RsvpID:       rsvp_id.String,
		Email:        email.String,
	}
	return Rsvp, nil
}

const getStatement = "SELECT * FROM rsvp WHERE rsvp_id = ?"

// GetRsvp retrieves a Rsvp by its ID.
func (db *mysqlDB) GetRsvp(id string) (*Rsvp, error) {
	Rsvp, err := scanRsvp(db.get.QueryRow(id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("mysql: could not find Rsvp with rsvp_id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("mysql: could not get Rsvp: %v", err)
	}
	return Rsvp, nil
}

const insertStatement = `
  INSERT INTO rsvp (
    rsvp_id, email
  ) VALUES (?, ?)`

// AddRsvp saves a given Rsvp, assigning it a new ID.
func (db *mysqlDB) AddRsvp(b *Rsvp) (id int64, err error) {
	r, err := execAffectingOneRow(db.insert, b.RsvpID, b.Email)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
	}
	return lastInsertID, nil
}

const updateStatement = `
  UPDATE rsvp
  SET email=?
  WHERE id = ? AND rsvp_id = ?`

// UpdateRsvp updates the entry for a given Rsvp.
func (db *mysqlDB) UpdateRsvp(b *Rsvp) error {
	if b.ID == 0 {
		return errors.New("mysql: Rsvp with unassigned ID passed into updateRsvp")
	}

	if b.RsvpID == "" {
		return errors.New("mysql: Rsvp with unassigned RsvpID passed into updateRsvp")
	}

	_, err := execAffectingOneRow(db.update, b.Email, b.ID, b.RsvpID)
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
