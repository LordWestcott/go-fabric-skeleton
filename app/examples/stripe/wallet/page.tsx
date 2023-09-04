"use client"

import { useAuth } from "@/contexts/AuthContext"
import { useEffect, useState } from "react";
import { StripeElementsOptions, loadStripe } from '@stripe/stripe-js';
import { Elements } from '@stripe/react-stripe-js';
import CreditCard from "@/components/CreditCard";
import SaveNewCreditCard from "@/components/SaveNewCreditCard";

const stripe = loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY_DEV || '');

export default function CustomerWallet() {
    const auth = useAuth()
    const [wallet, setWallet] = useState<any>(null);
    const [setupIntent, setSetupIntent] = useState<any>(null);

    const createNewSetupIntent = async () => {
        const response = await fetch(`${auth.api}/stripe/wallet`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${auth.session?.token}`
            },
        });

        if (!response.ok) {
            console.warn("Error creating new setup intent");
            return;
        }

        const json = await response.json();
        if (!json.success) {
            console.warn("Error creating new setup intent", json.message);
            return;
        }

        const setupIntent = json.data;
        console.log("Setup intent", setupIntent)
        setSetupIntent(setupIntent);
    }


    const getWallet = async () => {
        const response = await fetch(`${auth?.api}/stripe/wallet`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${auth?.session?.token}`
            },
        });

        if (!response.ok) {
            console.warn("Error getting wallet");
            return;
        }

        const json = await response.json();
        if (!json.success) {
            console.warn("Error getting wallet", json.message);
            return;
        }

        const wallet = json.data;
        console.log("Wallet", wallet)
        setWallet(wallet);
    }




     useEffect(() => {
            getWallet();
    }, [auth?.session]);


     if(auth?.loading){
        return (
            <>
                <h1>Customer Wallet</h1>
                <p>Loading...</p>
            </>
        )
     }

     if (!auth?.session){
        return (
            <>
                <h1>Customer Wallet</h1>
                <p>You need to be logged in to see this.</p>
            </>
        )
     }

     if (!auth?.session?.stripe_id){
        return (
            <>
                <h1>Customer Wallet</h1>
                <p>You are logged in, but you need a stripe id for this.</p>
            </>
        )
     }

     const options: StripeElementsOptions = {
        clientSecret: setupIntent?.client_secret,
        // appearance,
      };


     return (
            <>
                <h1>Customer Wallet</h1>
                <p>Wallet</p>

                <h2>Save New Card</h2>
                { !setupIntent && (
                    <button onClick={createNewSetupIntent}>Attach a new card</button>
                )}
                { setupIntent && (
                    <Elements stripe={stripe} options={options}>
                        <SaveNewCreditCard getWallet={getWallet} setSetupIntent={setSetupIntent} setupIntent={setupIntent}/>
                    </Elements>
                )}

                <h2>Existing Cards</h2>
                <select>
                    {wallet && wallet.map((paymentSource: any) => {
                        if(paymentSource) return <CreditCard key={paymentSource.id} card={paymentSource.card} />
                    })}
                </select>
            </>
     )
}