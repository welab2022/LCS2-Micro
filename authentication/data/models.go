package data

import (
	"context"
	"database/sql"
	"errors"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/nicored/avatar"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name,omitempty"`
	LastName         string    `json:"last_name,omitempty"`
	Password         string    `json:"-"`
	Active           int       `json:"active"`
	LastLogin        time.Time `json:"last_login"`
	PasswordChangeAt time.Time `json:"password_changed_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ? ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// ? ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

// GET the data directory on postgres (/var/lib/postgres/data)
func (u *User) GetDataDirPath() string {
	var data_path string

	query := "show data_directory;"
	err := db.QueryRow(query).Scan(&data_path)
	if err != nil {
		return ""
	}

	return data_path
}

// GetAll returns a slice of all users, sorted by last name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, last_login, password_changed_at, created_at, updated_at
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("QueryContext err: %s", err)
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.LastLogin,
			&user.PasswordChangeAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail returns one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, last_login, password_changed_at, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.LastLogin,
		&user.PasswordChangeAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("%s does not exist!", email)
		return nil, err
	}

	return &user, nil
}

// GetOne returns one user by id
func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, last_login, password_changed_at, created_at, updated_at from users where id = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.LastLogin,
		&user.PasswordChangeAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set
		email = $1,
		first_name = $2,
		last_name = $3,
		user_active = $4,
		updated_at = $5
		where id = $6
	`

	_, err := db.ExecContext(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateAvatar(avatarHex []byte, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	log.Printf("email: %s", email)

	stmt := `update users set
		avatar = $1::bytea
		where email = $2
	`
	log.Printf("stmt: %s", stmt)

	_, err := db.ExecContext(ctx, stmt,
		avatarHex,
		email,
	)

	if err != nil {
		log.Printf("err: %s", err.Error())
		return err
	}

	return nil
}

func (u *User) GetAvatar(email string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select avatar from users where email = $1`

	var buf []byte

	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(&buf)
	if err != nil {
		return nil, err
	}

	return buf, nil

}

// Delete deletes one user from the database, by User.ID
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID deletes one user from the database, by ID
func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GenerateCircleAvatar(initial string) ([]byte, error) {
	size := 200
	newAvatar, err := avatar.NewAvatarFromInitials([]byte(initial), &avatar.InitialsOptions{
		FontPath:  "/app/Arial.ttf",           // Required
		Size:      size,                       // default 300
		NInitials: 2,                          // default 1 - If 0, the whole text will be printed
		TextColor: color.White,                // Default White
		BgColor:   color.RGBA{0, 0, 255, 255}, // Default color.RGBA{215, 0, 255, 255} (purple)
	})
	if err != nil {
		log.Printf("Generate Avatar error: %s", err)
		return nil, err
	}

	// square, _ := newAvatar.Square()
	// squareFile, _ := os.Create("./output/square_john_smith_initials.png")
	// defer squareFile.Close()
	// squareFile.Write(square)

	round, err := newAvatar.Circle()
	roundFile, _ := os.Create("./initial_circle.png")
	defer roundFile.Close()
	roundFile.Write(round)

	return round, err
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var fullname = user.FirstName + user.LastName
	avatar, err := u.GenerateCircleAvatar(fullname)
	if err != nil {
		log.Printf("Generate avatar failed, err: %s", err)
		return 0, err
	}

	stmt := `insert into users (email, first_name, last_name, password, avatar, user_active, last_login, password_changed_at, created_at, updated_at)
	values ($1, $2, $3, $4, $5::bytea, $6, $7, $8, $9, $10) returning id`

	row := db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		avatar,
		user.Active,
		user.LastLogin,
		user.PasswordChangeAt,
		time.Now(),
		time.Now(),
	)
	var newID int
	err = row.Scan(&newID)

	if err != nil {
		log.Printf("err: %s", err)
		log.Printf("return id: %d", newID)
		return newID, err
	}
	log.Printf("return id: %d", newID)

	return newID, nil
}

// ResetPassword is the method we will use to change a user's password.
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	var password_changed_at = time.Now()

	stmt := `update users set password = $1, password_changed_at = $2 where id = $3`
	_, err = db.ExecContext(ctx, stmt, hashedPassword, password_changed_at, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) LastLoginUpdate(_time time.Time, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update users set last_login = $1 where email = $2`
	_, err := db.ExecContext(ctx, stmt, _time, email)
	if err != nil {
		log.Printf("LastLoginUpdate: err: %s", err)
		return err
	}

	return nil
}

// PasswordMatches uses Go's bcrypt package to compare a user supplied password
// with the hash we have stored for a given user in the database. If the password
// and hash match, we return true; otherwise, we return false.
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
