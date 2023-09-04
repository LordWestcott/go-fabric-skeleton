package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type Price struct {
	ID            int64     `db:"id,omitempty" json:"id"`
	IsActive      bool      `db:"is_active" json:"is_active"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	ProductID     int64     `db:"product_id" json:"product_id"`
	StripeID      string    `db:"stripe_id" json:"stripe_id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	BillingPeriod string    `db:"billing_period" json:"billing_period"`
	Amount        int       `db:"amount" json:"amount"`
}

func (u *Price) Table() string {
	return "prices"
}

func (u *Price) GetAll(upper up.Session) ([]*Price, error) {
	collection := upper.Collection(u.Table())

	var all []*Price

	res := collection.Find().OrderBy("created_at DESC")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (p *Price) GetByStripeID(stripeID string, upper up.Session) (*Price, error) {
	collection := upper.Collection(p.Table())

	var price Price

	res := collection.Find(up.Cond{"stripe_id =": stripeID})
	err := res.One(&price)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	return &price, nil
}

func (p *Price) GetByProductID(productID int64, upper up.Session) ([]*Price, error) {
	collection := upper.Collection(p.Table())

	var all []*Price

	res := collection.Find(up.Cond{"product_id =": productID}).OrderBy("created_at DESC")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (p *Price) Get(id int, upper up.Session) (*Price, error) {
	collection := upper.Collection(p.Table())

	var price Price

	res := collection.Find(up.Cond{"id =": id})
	err := res.One(&price)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	return &price, nil
}

func (p *Price) Update(price Price, upper up.Session) error {
	price.UpdatedAt = time.Now()

	collection := upper.Collection(p.Table())
	res := collection.Find(price.ID)
	err := res.Update(&price)
	if err != nil {
		return err
	}

	return nil
}

func (p *Price) Delete(id int, upper up.Session) error {
	collection := upper.Collection(p.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (p *Price) Insert(price *Price, upper up.Session) (int, error) {
	price.CreatedAt = time.Now()
	price.UpdatedAt = time.Now()

	collection := upper.Collection(p.Table())
	res, err := collection.Insert(price)
	if err != nil {
		return 0, err
	}

	price.ID = res.ID().(int64)

	return getInsertID(res.ID()), nil
}
