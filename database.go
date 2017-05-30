package gorth

const (
	insertSQL = `INSERT INTO Users
	(FirstName, LastName, Email, Password, Photo, Status, Token, TokenTimestamp,
	LastLogin) VALUES
	(?,?,?,?,?,?,?,NOW(),?)`
	updateLastLoginSQL = "UPDATE Users SET LastLogin = NOW() WHERE ID = ?"
	updateStatusSQL    = "UPDATE Users SET Status = ? WHERE ID = ?"
)

func updateLastLogin(userID int64) error {
	_, err := db.Exec(
		updateLastLoginSQL,
		userID)

	return err
}

func insertUser(user *User) error {
	result, err := db.Exec(
		insertSQL,
		user.FirstName,
		user.LastLogin,
		user.Email,
		user.Password,
		user.Photo,
		user.Status,
		user.Token,
		user.TokenTimestamp)

	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()

	return err
}

func updateStatus(userID int64, status string) error {
	_, err := db.Exec(
		updateStatusSQL,
		status,
		userID)

	return err
}

// CreateUsersTable sets up the user table in the db
func CreateUsersTable() error {
	createUserTableSQL := `
    CREATE TABLE IF NOT EXISTS Users (
        ID int(11) NOT NULL AUTO_INCREMENT,
        FirstName varchar(255) NOT NULL DEFAULT "",
        LastName varchar(255) NOT NULL DEFAULT "",
        Email varchar(255) NOT NULL DEFAULT "",
        Password varchar(255) NOT NULL DEFAULT "",
        Photo varchar(255) NOT NULL DEFAULT "",
        Status varchar(50) NOT NULL DEFAULT "",
        Role varchar(50) NOT NULL DEFAULT "",
        Token varchar(100) NOT NULL DEFAULT "",
        TokenTimestamp timestamp DEFAULT 0 NOT NULL,
        LastLogin timestamp DEFAULT 0 NOT NULL,
        Created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (ID),
        KEY email (Email),
        KEY status (Status)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Users, created for gorth'`

	_, err := db.Exec(createUserTableSQL)
	return err
}
