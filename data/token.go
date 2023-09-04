package data

import (
	"crypto/rand"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lordwestcott/gofabric/jwt"
	up "github.com/upper/db/v4"
)

type Token struct {
	ID        int64     `db:"id,omitempty" json:"id"`
	UID       int64     `db:"user_id" json:"uid"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	FirstName string    `db:"first_name" json:"first_name"`
	Email     string    `db:"email" json:"email"`
	PlainText string    `db:"token" json:"token"`
	// Hash      []byte    `db:"token_hash" json:"-"`
	Expires time.Time `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) GetUserForToken(token string, upper up.Session) (*User, error) {
	var u User
	var tok Token

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token =": token})
	err := res.One(&tok)
	if err != nil {
		return nil, err
	}

	collection = upper.Collection(u.Table())
	res = collection.Find(up.Cond{"id =": tok.UID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	u.Token = tok

	return &u, nil
}

func (t *Token) GetTokensForUser(uid int, upper up.Session) ([]*Token, error) {
	var tokens []*Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id =": uid}).OrderBy("created_at DESC")
	err := res.All(&tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (t *Token) Get(id int, upper up.Session) (*Token, error) {
	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) GetByToken(plainText string, upper up.Session) (*Token, error) {
	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *Token) Delete(id int, upper up.Session) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) DeleteByToken(plainText string, upper up.Session) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": plainText})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) Insert(token Token, u User, upper up.Session) error {
	collection := upper.Collection(t.Table())

	//delete existing tokens
	res := collection.Find(up.Cond{"user_id": u.ID})
	err := res.Delete()
	if err != nil {
		return err
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}

func (t *Token) GenerateToken(userID int64, jwt string, ttl time.Duration) (*Token, error) {
	token := &Token{
		UID:     userID,
		Expires: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = jwt

	// token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	// hash := sha256.Sum256([]byte(token.PlainText))
	// token.Hash = hash[:] //Converts to array

	return token, nil
}

func (t *Token) AuthenticateTokenFromRequest(r *http.Request, upper up.Session) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("No Authorization header found")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("No Authorization header found")
	}

	token := headerParts[1]

	jwt := jwt.JWT{}
	jwt.Secret = []byte(os.Getenv("JWT_SECRET"))

	details, err := jwt.VerifyJWT(token)
	if err != nil {
		return nil, err
	}

	if !details.Valid {
		return nil, errors.New("Invalid token")
	}

	tkn, err := t.GetByToken(token, upper)
	if err != nil {
		return nil, err
	}

	if tkn.Expires.Before(time.Now()) {
		return nil, errors.New("Token has expired")
	}

	user, err := t.GetUserForToken(token, upper)
	if err != nil {
		return nil, errors.New("No matching user found")
	}

	return user, nil
}

func (t *Token) AuthenticateToken(token string, upper up.Session) (*User, error) {
	if len(token) != 26 {
		return nil, errors.New("Invalid Token")
	}

	tkn, err := t.GetByToken(token, upper)
	if err != nil {
		return nil, err
	}

	if tkn.Expires.Before(time.Now()) {
		return nil, errors.New("Token has expired")
	}

	user, err := t.GetUserForToken(token, upper)
	if err != nil {
		return nil, errors.New("No matching user found")
	}

	return user, nil
}

func (t *Token) ValidToken(token string, upper up.Session) (bool, error) {
	user, err := t.GetUserForToken(token, upper)
	if err != nil {
		return false, errors.New("No matching user found")
	}

	if user.Token.PlainText == "" {
		return false, errors.New("No matching token found")
	}

	if user.Token.Expires.Before(time.Now()) {
		return false, errors.New("Token has expired")
	}

	return true, nil
}
