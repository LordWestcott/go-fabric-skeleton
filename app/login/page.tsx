"use client"

import SessionCard from "@/components/SessionCard"
import { useAuth } from "@/contexts/AuthContext"
import { useRouter } from "next/navigation"
export default function Login() {

    const router = useRouter()
    const auth = useAuth()

    const loginWithGoogle = async () => {
        auth.signInWithGoogle()
    }

    if(auth?.loading){
        return (
            <div>
                <h1>Loading</h1>
            </div>
        )
    }

    if(auth?.session){
        return (
            <>
                <SessionCard />
                <button onClick={auth.signOut}>Sign Out</button>
            </>
        )
    }

    return (
        <div>
            <button onClick={loginWithGoogle}>Login With Google</button>
        </div>
    )
}

