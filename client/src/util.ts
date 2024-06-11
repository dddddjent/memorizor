import axios, { AxiosError } from 'axios'
import { NavigateFunction } from 'react-router-dom'
import config from './config'
import { CSSProperties } from 'react'

export function generateURL(api: string, path: string): string {
    return api + path
}

export function asAxiosError(
    err: unknown,
    process: (error: AxiosError) => void,
) {
    if (axios.isAxiosError(err)) {
        process(err as AxiosError)
    } else {
        console.log('Unexpected error: ', err)
    }
}

export async function tryRequest(
    process: () => Promise<void>,
    processError: (err: unknown) => void,
    navigate: NavigateFunction,
) {
    axios.defaults.headers.common[
        'Authorization'
    ] = `Bearer ${localStorage.getItem('access_token')}`
    await process().catch((err: AxiosError) => {
        const responseData = err.response?.data as {
            error: {
                type: string
                message: string
            }
        }
        if (responseData.error.message.includes('expire')) {
             (async () => {
                try {
                    console.log('Begin to refresh')
                    const {
                        data: {
                            token_pair: {
                                access_token,
                                refresh_token: { token_string },
                            },
                        },
                    } = await axios.post(generateURL(config.api.account, '/token'), {
                        refresh_token: localStorage.getItem('refresh_token'),
                    })

                    localStorage.setItem('access_token', access_token)
                    localStorage.setItem('refresh_token', token_string)
                    axios.defaults.headers.common[
                        'Authorization'
                    ] = `Bearer ${localStorage.getItem('access_token')}`
                    console.log('Refreshed')
                    await process()
                } catch (error) {
                    asAxiosError(error, (error) => {
                        const responseData = error.response?.data as {
                            error: {
                                type: string
                                message: string
                            }
                        }
                        if (responseData.error.message.includes('expire')) {
                            console.log('Refresh token expired')
                            navigate('/signin')
                        } else {
                            console.log('Unexpected error')
                        }
                    })
                    processError(error)
                }
            })()
        } else {
            processError(err)
        }
    })
}

export const leftAlignedPosition = (percentage: number): CSSProperties => {
    return {
        position: 'absolute',
        left: percentage.toString() + '%',
        top: '50%',
        transform: 'translate(0%, -50%)',
    }
}
