"use client";
import { useAuth } from '@/contexts/AuthContext';
import Image from 'next/image'
import Link from 'next/link';

export default function Home() {
  const auth = useAuth();
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
<div>
      {!auth?.session && (
        <Link href="/login">Login</Link>
      )}

      {auth?.session && (
        <>
          <Link href="/login">Logout</Link>
          <br/>
          <br/>
          <h2>Stripe Examples</h2>
          <ul>
            <li><Link href="/examples/stripe/anon-payment-using-intent">Anonymous Payments Using Intents</Link></li>
            <li><Link href="/examples/stripe/external-checkout">External Checkout</Link></li>
            <li><Link href="/examples/stripe/wallet">Wallet</Link></li>
            <li><Link href="/examples/stripe/subscriptions">Subscriptions</Link></li>
          </ul>
        </>
      )}



      </div>
    </main>
  )
}
