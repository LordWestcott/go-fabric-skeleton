"use client";

import Head from 'next/head';
import React, { useState, useEffect } from 'react';
import { Elements } from '@stripe/react-stripe-js';
import { StripeElementsOptions, loadStripe } from '@stripe/stripe-js';
import CheckoutForm from '@/components/CheckoutForm';
import { useAuth } from '@/contexts/AuthContext';

const stripe = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY_DEV || '');

export default function Home() {
  const [clientSecret, setClientSecret] = useState('');
  const [paymentIntent, setPaymentIntent] = useState('');

  const auth = useAuth()

  console.log(auth.api);

  useEffect(() => {
    if (!auth?.api) {return}
    // Create PaymentIntent as soon as the page loads using our local API
    console.log("fetching from", `${auth.api}/stripe/create-anon-payment-intent`)
    fetch(`${auth.api}/stripe/create-anon-payment-intent`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        amount: 30000,
        currency: 'usd',
      }),
    })
      .then((res) => {
        console.log("res: ", res)
        res.json().then((json) => {
            console.log("json: ", json)
            setClientSecret(json.data.client_secret);
            setPaymentIntent(json.data.payment_intent);
            })
    })
  }, [auth]);

  const appearance = {
    theme: 'stripe',
    labels: 'floating',
  };
  const options: StripeElementsOptions = {
    clientSecret,
    // appearance,
  };

  return (
    <div>
      <Head>
        <title>Stripe Elements</title>
      </Head>
      <h1 className="text-2xl bold mb-4 mt-4 text-center">
        Accept payments with credit card
      </h1>
      <div className="flex justify-center m-auto w-1/2">
        <div className="w-1/3 flex justify-start items-center">
          {/* <img src="/stripe.svg" alt="Stripe" /> */}
        </div>
        <div className="w-1/3 flex justify-center items-center text-4xl bold">
          +
        </div>
        <div className="w-1/3 flex justify-end items-center">
          {/* <img src="/next.svg" alt="Stripe" /> */}
        </div>
      </div>
      {clientSecret && (
        <Elements options={options} stripe={stripe}>
          <CheckoutForm paymentIntent={paymentIntent} />
        </Elements>
      )}
    </div>
  );
}