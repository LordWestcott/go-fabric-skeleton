'use client'
 
import { useSearchParams } from 'next/navigation'
import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Api } from 'sst/node/api'

export default function LoginSuccess() {
    const params = useSearchParams()
    const token = params.get('token')
    const router = useRouter()

    const getSession = async () => {
        const { url } = await (await fetch(`/api/backend-url`)).json()
        console.log(url)
        
        const res = await fetch(`${url}/session`, {
            method: 'GET',
            headers: {
                Authorization: `Bearer ${token}`
            }
        })

        const data = await res.json()
        console.log(data)
    
        //TODO ADD TO SESSION CONTEXT

        router.replace('/')
    }

    useEffect(() => {
        getSession()
    }, [token])

    return (
        <div>
            <h1>Login Success</h1>
            <p>Token: {token}</p>
        </div>
    )
}