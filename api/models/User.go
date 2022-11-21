package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:225;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Creates a new hash for a given user password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Given a hashed password, this function will verify it
// against a provided user password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) beforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

// Initalises a new user and is called before Validate()
func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Looks at the requested action and validates user data
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("nickname required")
		}
		if u.Password == "" {
			return errors.New("password required")
		}
		if u.Email == "" {
			return errors.New("email required")
		}
		// checkmail will perform email validation
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	case "login":
		if u.Password == "" {
			return errors.New("password required")
		}
		if u.Email == "" {
			return errors.New("email required")
		}
		// checkmail will perform email validation
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("nickname required")
		}
		if u.Password == "" {
			return errors.New("password required")
		}
		if u.Email == "" {
			return errors.New("email required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// if there's an error, return an empty user with the error
	if err := db.Debug().Create(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error; err != nil {
		return &[]User{}, err
	}

	return &users, nil
}

func (u *User) FindUserById(db *gorm.DB, uid uint32) (*User, error) {
	if err := db.Debug().Model(User{}).Where("id= ?", uid).Take(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &User{}, errors.New("user not found")
		}
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	// Hashing password
	if err := u.beforeSave(); err != nil {
		log.Fatal(err)
	}

	// update the user
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	if err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error; err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}
