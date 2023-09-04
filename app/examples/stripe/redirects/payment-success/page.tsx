"use client";

export function CheckoutSuccess_GetSessionID() {
    const url = window.location.href
    const sessionId = new URL(url).searchParams.get("session_id");
    return sessionId
}

export default function StripePaymentSuccessRedirect(){

    const sessionId = CheckoutSuccess_GetSessionID()
    //now you could do something with it, like show an invoice, and update the session in the database

    return (
        <div>
            <h1>Payment Succeeded, sessionid: {sessionId}</h1>
        </div>
    )
}