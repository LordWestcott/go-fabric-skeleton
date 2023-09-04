package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type Product struct {
	ID          int64     `db:"id,omitempty" json:"id"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	StripeID    string    `db:"stripe_id" json:"stripe_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	IsRecurring bool      `db:"is_recurring" json:"is_recurring"`
}

func (p *Product) Table() string {
	return "products"
}

func (p *Product) GetAll(upper up.Session) ([]*Product, error) {
	collection := upper.Collection(p.Table())

	var all []*Product

	res := collection.Find().OrderBy("created_at DESC")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (p *Product) GetByStripeID(stripeID string, upper up.Session) (*Product, error) {
	collection := upper.Collection(p.Table())

	var product Product

	res := collection.Find(up.Cond{"stripe_id =": stripeID})
	err := res.One(&product)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	return &product, nil
}

func (p *Product) Get(id int, upper up.Session) (*Product, error) {
	collection := upper.Collection(p.Table())

	var product Product

	res := collection.Find(up.Cond{"id =": id})
	err := res.One(&product)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	return &product, nil
}

func (p *Product) Update(product Product, upper up.Session) error {
	product.UpdatedAt = time.Now()

	collection := upper.Collection(p.Table())
	res := collection.Find(product.ID)
	err := res.Update(&product)
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Delete(id int, upper up.Session) error {
	collection := upper.Collection(p.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

func (p *Product) Insert(product *Product, upper up.Session) (int, error) {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	collection := upper.Collection(p.Table())
	res, err := collection.Insert(product)
	if err != nil {
		return 0, err
	}

	product.ID = res.ID().(int64)

	return getInsertID(res.ID()), nil
}
