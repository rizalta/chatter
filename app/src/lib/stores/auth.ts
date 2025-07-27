import { PUBLIC_API_URL, PUBLIC_JWT_KEY } from "$env/static/public"
import { writable } from "svelte/store";
import { browser } from "$app/environment";
import { importJWK, jwtVerify } from "jose";

interface User {
  id: string,
  username: string,
  exp: number,
  iat: number,
}

export interface AuthState {
  isAuthenticated: boolean,
  user: User | null,
  token: string | null,
  loading: boolean,
}

export interface LoginCredentials {
  username: string,
  password: string,
}

interface LoginResponse {
  success: boolean,
  error?: string,
}

const jwk = JSON.parse(PUBLIC_JWT_KEY)

const verifyToken = async (token: string): Promise<User | null> => {
  try {
    const key = await importJWK(jwk, "RS256")
    const { payload } = await jwtVerify<User>(token, key)
    return payload
  } catch {
    return null
  }
}

const createAuthStore = () => {
  const initialState: AuthState = {
    isAuthenticated: false,
    user: null,
    token: null,
    loading: true,
  }

  const unauthenticatedState: AuthState = {
    isAuthenticated: false,
    user: null,
    token: null,
    loading: false,
  }

  const { subscribe, set } = writable<AuthState>(initialState)

  return {
    subscribe,

    init: async () => {
      if (!browser) return

      const token = localStorage.getItem("jwt_token")
      if (token) {
        const payload = await verifyToken(token)
        if (payload) {
          set({
            isAuthenticated: true,
            user: payload,
            token: token,
            loading: false,
          })
        } else {
          localStorage.removeItem("jwt_token")
          set(unauthenticatedState)
        }
      } else {
        localStorage.removeItem("jwt_token")
        set(unauthenticatedState)
      }
    },

    login: async (credentials: LoginCredentials): Promise<LoginResponse> => {
      try {
        const response = await fetch(`${PUBLIC_API_URL}/user/login`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(credentials),
        })

        if (!response.ok) {
          const err = await response.json()
          throw err
        }

        const { token }: { token: string } = await response.json()
        if (!token) {
          throw new Error("No token recieved")
        }

        const payload = await verifyToken(token)
        if (!payload) {
          throw new Error("Invalid token recieved")
        }

        set({
          isAuthenticated: true,
          user: payload,
          token: token,
          loading: false,
        })

        return { success: true }
      } catch (error) {
        set(unauthenticatedState)

        const errMessage = error instanceof Error ? error.message : "Unknown error occured"
        return { success: false, error: errMessage }
      }
    },

    logout: (): void => {
      localStorage.removeItem("jwt_token")
      set(unauthenticatedState)
    },
  }
}

export const auth = createAuthStore()
