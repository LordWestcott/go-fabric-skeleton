"use client"

import { useRouter } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";

export type authContextType = { 
    session: Session | null; 
    loading: boolean; 
    token: string | null;
    getSession: () => Promise<void>; 
    signInWithGoogle: () => Promise<void>; 
    signOut: () => Promise<void>; 
    api: string | null;
    host: string | null;
}

export type Session = { //TODO: Move this to a types file
    email: string;
    firstName: string;
    lastName: string;
    picture_url: string;
    token: string;
    stripe_id: string;
    locale: string;
    created_at: Date;
    //TODO: Add is new user
};

type AuthContextProps = {
    children: React.ReactNode;
};

const authContextDefaultValues: authContextType = {
    session: null,
    loading: true,
    token: null,
    getSession: async () => { },
    signInWithGoogle: async () => { },
    signOut: async () => { },
    api: null,
    host: null,
}

export const AuthContext = createContext<authContextType>(authContextDefaultValues);

export function useAuth() {
    return useContext(AuthContext);
}

export function AuthProvider({ children }: AuthContextProps) {
    const [session, setSession] = useState<Session | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [token, setToken] = useState<string | null>(null);
    const [api, setApi] = useState<string | null>(null);
    const [host, setHost] = useState<string | null>(null);
    const router = useRouter();

    const getApi = async () => {
        let url = process.env.NEXT_PUBLIC_API_URL || process.env.NEXT_PUBLIC_API_URL_DEV || null;
        
        console.log("API URL", url);
        console.log("NEXT_PUBLIC_API_URL", process.env.NEXT_PUBLIC_API_URL);
        console.log("NEXT_PUBLIC_API_URL_DEV", process.env.NEXT_PUBLIC_API_URL_DEV);

        if(typeof url === "undefined"){
            url = null;
        }

        setApi(url);
        return url;
    };

    const getHost = async () => {
        let url = window.location.origin;
        setHost(url);
    };

    useEffect(() => {
        getHost();
        getApi();
    }, []);

    const getUserInfo = async (token: string | null) => {
        let backend = api;
        
        if (!backend) {
             backend = await getApi() || null;
        }
        
        if (!backend || !token) {
            //TODO may need to redirect to login.
            console.warn("No api or token", {api, token});
            return;
        }; 

        const fetchUserInfo = async () => {
            const response = await fetch(
                `${backend}/session`,
                {
                    method: "GET",
                    headers: {
                    Authorization: `Bearer ${token}`,
                    },
                }
                );
                if (!response.ok) {
                    console.warn("Error getting user info");
                    // signOut();
                    return;
                }
                const data = await response.json();
                return data;
        };

        try {
            console.log("Fetching user info");
            var data = await fetchUserInfo();
            console.log("User info", data);

            if (!data?.success) {
                console.warn("Error getting user info");
                signOut();
                //check if we are on the login page
                if (window.location.pathname !== "/login") {
                    router.replace("/login");
                }else{
                    setLoading(false);
                    console.log("Would have directed to login.")
                }
                return;
            }

            return data.data;
        } catch (error) {
            console.log("get user failed: ->>>", error)
        }
    };

    const signOut = async () => {
        localStorage.removeItem("session");
        setToken(null);
        setSession(null);
        router.replace("/login");
    };

    const getSession = async () => {
        const token = localStorage.getItem("session");
        setToken(token);
        if (token) {
            const user = await getUserInfo(token);
            if (user) setSession(user);
        }
        setLoading(false);
    };

    const signInWithGoogle = async () => {
        let url = await getApi();
        if (!url) {
            //TODO may need to redirect to login.
            console.warn("No api address");
            console.log(process.env);
            return;
        }; 
        window.location.href = `${url}/google-signin`;
    };

    useEffect(() => {
            getSession();
    }, []);

    useEffect(() => {
        const search = window.location.search;
        const params = new URLSearchParams(search);
        const token = params.get("token");
        if (token) {
            setToken(token);
            localStorage.setItem("session", token);
            window.location.replace(window.location.origin);
        }
    }, []);

    const value: authContextType = { //TODO: This MAY need memoization
        session,
        loading,
        token,
        getSession,
        signInWithGoogle,
        signOut,
        api,
        host,
    };

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}