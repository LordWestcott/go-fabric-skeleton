"use client";

import { useAuth } from "@/contexts/AuthContext"
import { loadStripe } from "@stripe/stripe-js";

export default function ExampleStripeExternalCheckout() {
    
    const auth = useAuth()

    const handleClick = async () => {
        const stripe = await loadStripe(
            process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY_DEV || ''
         )

        const request = {
            currency: 'usd',
            product: {
                name: 'Emily\'s Example Product',
                description: 'Example Product Description',
                images: [
                    "https://plus.unsplash.com/premium_photo-1673125287084-e90996bad505?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=987&q=80"
                ],
                price: 2001,
            },
            quantity: 1,
            success_url: `${auth.host}/examples/stripe/redirects/payment-success?session_id={CHECKOUT_SESSION_ID}`,
            cancel_url: `${auth.host}/examples/stripe/redirects/payment-cancel?session_id={CHECKOUT_SESSION_ID}`,
            payment_method_types: ['card'],
        } 
        
        const response = await fetch(`${auth.api}/stripe/create-external-checkout-session`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(request),
        });

        if (!response.ok) {
            console.warn("Error creating stripe checkout session");
            return;
        }

        const json = await response.json();
        if (!json.success) {
            console.warn("Error creating stripe checkout session", json.message);
            return;
        }

        const sessionId = json.data.session_id;

        const result = await stripe?.redirectToCheckout({
            sessionId,
        });

        if (result?.error) {
            console.warn("Error redirecting to checkout", result?.error);
        }
    }
    
    return (
        <div>
            <h1>Example Stripe External Checkout</h1>
            <button onClick={handleClick}>Checkout</button>
        </div>
    )
}