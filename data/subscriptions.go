package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// These are active subscriptions that we are earning on.
// These are bought recurring products.
type Subscription struct {
	ID                   int64     `db:"id,omitempty"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
	UID                  int64     `db:"user_id"`
	ProductID            int64     `db:"product_id"`
	PriceID              int64     `db:"price_id"`
	StripeSubscriptionID string    `db:"stripe_subscription_id"`
	StripeCustomerID     string    `db:"stripe_customer_id"`
	StripePriceID        string    `db:"stripe_price_id"`
	Status               string    `db:"status"`
}

// Override the table name
func (s *Subscription) Table() string {
	return "subscriptions"
}

func (s *Subscription) GetAll(upper up.Session) ([]*Subscription, error) {
	collection := upper.Collection(s.Table())

	var all []*Subscription

	res := collection.Find().OrderBy("created_at DESC")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (s *Subscription) GetByUserID(userID string, upper up.Session) (*Subscription, error) {
	var sub Subscription
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"user_id =": userID})
	err := res.One(&sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *Subscription) GetByStripeSubscriptionID(stripeSubscriptionID string, upper up.Session) (*Subscription, error) {
	var sub Subscription
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"stripe_subscription_id =": stripeSubscriptionID})
	err := res.One(&sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *Subscription) Get(id int64, upper up.Session) (*Subscription, error) {
	var sub Subscription
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *Subscription) Update(sub Subscription, upper up.Session) error {
	sub.UpdatedAt = time.Now()
	collection := upper.Collection(s.Table())
	res := collection.Find(sub.ID)
	err := res.Update(&sub)
	if err != nil {
		return err
	}
	return nil
}

func (s *Subscription) Delete(id int, upper up.Session) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *Subscription) Insert(sub *Subscription, upper up.Session) (int, error) {
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res, err := collection.Insert(sub)
	if err != nil {
		return 0, err
	}

	sub.ID = res.ID().(int64)

	return getInsertID(res.ID()), nil
}

func (s *Subscription) UpdateStatus(id int64, status string, upper up.Session) error {
	sub, err := s.Get(id, upper)
	if err != nil {
		return err
	}

	sub.Status = status
	err = sub.Update(*sub, upper)
	if err != nil {
		return err
	}

	return nil
}
