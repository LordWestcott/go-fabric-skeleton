"use client"

import { useAuth } from "@/contexts/AuthContext";

export default function SessionCard() {
    const auth = useAuth();

    let response = (<></>);

    if(!auth){
        response = (
        <>
            <h1>No Auth</h1>
        </>
        )
    }

    if(!auth.session){
        response = (
            <>
                <h1>Auth But No Session</h1>
            </>
        )
    } else {
        response = (
        <>
            <h1>Session</h1>
            <p>{auth.session.firstName}</p>
            <p>{auth.session.lastName}</p>
            <p>{auth.session.email}</p>
            <p>{auth.session.api}</p>
        </>
        )
    }

    return (
        <>
            <div className="bg-white shadow overflow-hidden sm:rounded-lg text-black">
                <div className="px-4 py-5 sm:px-6">
                    {response}
                </div>
            </div>
        </>       
    )
}