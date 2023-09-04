package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int64     `db:"id,omitempty"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	Email           string    `db:"email"`
	EmailIsVerified bool      `db:"verified_email"`
	PictureURL      string    `db:"picture_url"`
	Locale          string    `db:"locale"`
	Active          int       `db:"user_active"`
	GoogleID        string    `db:"google_id"`
	StripeID        string    `db:"stripe_id"`
	Password        string    `db:"password"`
	Token           Token     `db:"-"` //Doesn't exist in the database
}

// Override the table name
func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll(upper up.Session) ([]*User, error) {
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (u *User) GetByEmail(email string, upper up.Session) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email =": email})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at DESC")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}
	user.Token = token

	return &user, nil
}

func (u *User) GetByGoogleID(googleID string, upper up.Session) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"google_id =": googleID})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at DESC")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	user.Token = token
	return &user, nil
}

// No token returned.
// Server-Side Only
func (u *User) GetByStripeID(stripeID string, upper up.Session) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"stripe_id =": stripeID})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) Get(id int, upper up.Session) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at DESC")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}
	user.Token = token

	return &user, nil
}

func (u *User) Update(user User, upper up.Session) error {
	user.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(user.ID)
	err := res.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id int, upper up.Session) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(user *User, upper up.Session) (int, error) {
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(newHash)
	collection := upper.Collection(u.Table())
	res, err := collection.Insert(user)
	if err != nil {
		return 0, err
	}

	user.ID = res.ID().(int64)

	return getInsertID(res.ID()), nil
}

func (u *User) ResetPassword(id int, password string, upper up.Session) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user, err := u.Get(id, upper)
	if err != nil {
		return err
	}

	u.Password = string(newHash)
	err = user.Update(*u, upper)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string, upper up.Session) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (u *User) CheckForRememberToken(id int, token string, upper up.Session) bool {
	var rememberToken RememberToken
	rt := RememberToken{}
	collection := upper.Collection(rt.Table())
	cond := up.Cond{"user_id =": id, "remember_token =": token}
	res := collection.Find(cond)
	err := res.One(&rememberToken)
	return err == nil
}
