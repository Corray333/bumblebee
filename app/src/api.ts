import axios, { type InternalAxiosRequestConfig } from 'axios'
import { isAxiosError } from 'axios'



interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
    useAuth?: boolean
}

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
})

api.interceptors.request.use((config: CustomAxiosRequestConfig) => {
    if (config.useAuth) {
        const token = import.meta.env.VITE_API_TOKEN
        if (token) {
        config.headers['Authorization'] = token
        }
    }
    return config
})

export {api}