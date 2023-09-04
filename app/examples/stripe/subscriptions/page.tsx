"use client";
import { Session, useAuth } from "@/contexts/AuthContext";
import { useEffect, useState } from "react";
import { Stripe, StripeElementsOptions, loadStripe } from "@stripe/stripe-js";
import {
  CardElement,
  Elements,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";

type Product = {
  id: number;
  name: string;
  description: string;
  is_recurring: boolean;
  prices: Price[];
};

type Price = {
  id: number;
  amount: number;
  name: string;
  description: string;
  billing_period: string;
};

const stripe = loadStripe(
  process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY ||
    process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY_DEV ||
    ""
);

const getUserStripeSubcriptionsData = async (
  session: Session,
  api: string,
  setData: (data: any) => void
) => {
  if (!session || !session.token) return;

  const response = await fetch(`${api}/stripe/subscriptions`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${session.token}`,
    },
  });

  if (!response.ok) {
    console.warn("Error getting user stripe subscriptions data");
    return;
  }

  const json = await response.json();
  if (!json.success) {
    console.warn("Error getting user stripe subscriptions data");
    console.error(json.message);
    return;
  }

  console.log("Subscription Data: ", json.data);

  setData(json.data);
};

const getAvailableProducts = async (
  session: Session,
  api: string,
  setData: (data: any) => void
) => {
  const response = await fetch(`${api}/products`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${session.token}`,
    },
  });

  if (!response.ok) {
    console.warn("Error getting available products");
    return;
  }

  const json = await response.json();

  if (!json.success) {
    console.warn("Error getting available products");
    console.error(json.message);
    return;
  }

  console.log("Product Data: ", json);

  setData(json.data.products as Product[]);
};

export default function Subscriptions() {
  const auth = useAuth();
  if (auth?.loading) {
    return (
      <>
        <h1>Subscriptions</h1>
        <p>Loading...</p>
      </>
    );
  }

  if (!auth?.session) {
    return (
      <>
        <h1>Subscriptions</h1>
        <p>You need to be logged in to see this.</p>
      </>
    );
  }

  if (!auth?.session?.stripe_id) {
    return (
      <>
        <h1>Subscriptions</h1>
        <p>You are logged in, but you need a stripe id for this.</p>
      </>
    );
  }

  return (
    <div>
      <h1>Subscriptions</h1>
      <Elements stripe={stripe}>
        <SubscribeToPlan session={auth.session} />
      </Elements>
    </div>
  );
}

function SubscribeToPlan({ session }: { session: Session }) {
  const auth = useAuth();
  const stripe = useStripe();
  const elements = useElements();

  const [plan, setPlan] = useState<number | null>(null);
  const [subscriptions, setSubscriptions] = useState<>([]);
  const [loading, setLoading] = useState(false);
  const [avaliableProducts, setAvaliableProducts] = useState<Product[]>([]);

  useEffect(() => {
    if (!session || !session.token || !auth.api) return;
    getAvailableProducts(session, auth.api, setAvaliableProducts);
    getUserStripeSubcriptionsData(session, auth.api, setSubscriptions);
  }, [session, auth]);

  const cancelSubscription = async (subscriptionId: number) => {
    if (!session || !session.token || !auth.api) return;
    setLoading(true);
    const request = {
      subscription_id: subscriptionId,
    };
    const response = await fetch(`${auth.api}/stripe/cancel-subscription`, {
      method: "PATCH",
      headers: {
        Authorization: `Bearer ${session.token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      console.warn("Error cancelling subscription");
      return;
    }

    const json = await response.json();
    if (!json.success) {
      console.warn("Error cancelling subscription");
      console.error(json.message);
      return;
    }

    await getUserStripeSubcriptionsData(session, auth.api!, setSubscriptions);
    setLoading(false);
  };

  const handleSubmit = async (event: any) => {
    try {
      event.preventDefault();

      if (!stripe || !elements) {
        if (!stripe) console.warn("Stripe.js has not loaded yet.");
        if (!elements) console.warn("Stripe Elements has not loaded yet.");
        return;
      }

      console.log("Creating payment method.", CardElement)

      const cardElement = elements.getElement(CardElement);

        if (!cardElement) {
            console.warn("Card Element not found.");
            return;
        }
        
        console.log("Got here")
        
        const { paymentMethod, error } = await stripe.createPaymentMethod({
            type: "card",
            card: cardElement,
        });
        
    console.log("Got here 2")
      console.log("Payment Method: ", paymentMethod);
      console.log("Payment Method Error: ", error);

      if (error) {
        console.warn(error.message);
        return;
      }

      console.log("Sending request to create subscription.", {
        price_id: plan,
        stripe_payment_method_id: paymentMethod?.id,
        payment_method_attached: false,
        MakeDefaultPaymentMethod: true,
      });

      const response = await fetch(`${auth.api}/stripe/create-subscription`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          price_id: plan,
          stripe_payment_method_id: paymentMethod?.id,
          payment_method_attached: false,
          MakeDefaultPaymentMethod: true,
        }),
      });

      console.log("Response: ", response);

      const res = await response.json();
      const {data, success, message} = res;
      const subscription = data;

      console.log("Subscription: ", subscription);


      const { latest_invoice } = subscription;

      console.log("latest invoice: ", latest_invoice);
      console.log("Payment Intent: ", latest_invoice.payment_intent);

      if (latest_invoice.payment_intent) {
        const { client_secret, status } = latest_invoice.payment_intent;

        console.log("Client Secret", client_secret);
        console.log("Status", status);

        if (status === "requires_action") {
          const { error: confirmationError } = await stripe.confirmCardPayment(
            client_secret
          );
          if (confirmationError) {
            console.error(confirmationError);
            alert("unable to confirm card.");
            return;
          }
        }

        alert("You are subscribed!");
        await getUserStripeSubcriptionsData(
          session,
          auth.api!,
          setSubscriptions
        );
      }

      setPlan(null);
    } catch (error) {
      console.error(error);
    }
  };

  if (loading) {
    return (
      <>
        <p>Loading...</p>
      </>
    );
  }

  console.log("Avaliable Products: ", avaliableProducts)

  return (
    <>
      <hr />
      <h2>Select a Product</h2>
      {avaliableProducts.map((product) => (
        <>
          <div key={product.id}>
            <h3>{product.name}</h3>
            <p>{product.description}</p>
            <h4>Prices</h4>
            {product.prices.map((price) => (
              <div key={price.id}>
                <p>{price.name}</p>
                <p>{price.description}</p>
                <p>{price.amount}</p>
                <button onClick={() => setPlan(price.id)}>Select</button>
              </div>
            ))}
          </div>
        </>
      ))}
      <hr />

      <p>
        Selected Plan: <strong>{plan}</strong>
      </p>

      <hr />

      <form onSubmit={handleSubmit} hidden={!plan}>
        <CardElement />
        <button type="submit" disabled={!stripe || loading}>
          Subscribe & Pay
        </button>
      </form>

      <div>
        <h3>Manage Current Subscriptions</h3>
        <div>
          {subscriptions.map((sub: any) => (
            <div key={sub.id}>
              <br />
              <div>
                <strong>{sub.id}.</strong>{" "}
                {!sub.cancel_at_period_end && (
                  <>
                    Next payment of{" "}
                    <strong>
                      {sub.items.data[0].plan.currency}{" "}
                      {(sub.items.data[0].plan.amount / 100).toFixed(2)}
                    </strong>{" "}
                    due{" "}
                    <strong>
                      {new Date(
                        sub.current_period_end * 1000
                      ).toLocaleDateString()}
                    </strong>
                    <br />
                    <button
                      onClick={() => cancelSubscription(sub.id)}
                      disabled={loading}
                    >
                      Cancel
                    </button>
                  </>
                )}
                {sub.cancel_at_period_end && (
                  <>
                    Has been cancelled. Will expire on{" "}
                    <strong>
                      {new Date(
                        sub.current_period_end * 1000
                      ).toLocaleDateString()}
                    </strong>
                    <br />
                  </>
                )}
              </div>
              <br />
            </div>
          ))}
        </div>
      </div>
    </>
  );
}
