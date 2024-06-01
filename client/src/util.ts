import axios, { AxiosError } from "axios"

export function generateURL(api: string, path: string): string {
    return api + path
}

export function asAxiosError(err: unknown, process: (error: AxiosError<any, any>) => void) {
    if (axios.isAxiosError(err)) {
        process(err as AxiosError<any, any>)
    }
    else {
        console.log('unexpected error: ', err)
    }
}
