CREATE
OR REPLACE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();

RETURN NEW;

END;

$$ LANGUAGE plpgsql;

drop table if exists users cascade;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    google_id character varying(255) UNIQUE,
    email character varying(255) NOT NULL UNIQUE,
    verified_email boolean NOT NULL DEFAULT false,
    picture_url character varying(255),
    locale character varying(50),
    stripe_id character varying(255),
    password character varying(60) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists remember_tokens;

CREATE TABLE remember_tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    remember_token character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON remember_tokens FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists tokens;

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    first_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    -- token_hash bytea NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp without time zone NOT NULL
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON tokens FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists products;

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    is_active boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    stripe_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255),
    is_recurring boolean NOT NULL DEFAULT false
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON products FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists prices;

CREATE TABLE prices (
    id SERIAL PRIMARY KEY,
    is_active boolean NOT NULL DEFAULT false,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    product_id integer NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    stripe_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255),
    billing_period character varying(255) NOT NULL,
    amount integer NOT NULL
);

drop table if exists subscriptions;

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    product_id integer NOT NULL REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,
    price_id integer NOT NULL REFERENCES prices(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stripe_subscription_id character varying(255) NOT NULL,
    stripe_customer_id character varying(255) NOT NULL,
    stripe_price_id character varying(255) NOT NULL,
    status character varying(255) NOT NULL
);

CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON subscriptions FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();


CREATE TRIGGER set_timestamp BEFORE
UPDATE
    ON prices FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
