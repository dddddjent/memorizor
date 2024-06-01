import '../style/util.css'
import '../style/signin.css'
import { useState } from 'react'
import axios from 'axios'
import config from '../config'
import { asAxiosError, generateURL } from '../util'
import { useNavigate } from 'react-router-dom'

type SignInResponse = {
    token_pair: {
        access_token: string
        refresh_token: {
            token_string: string
        }
    }
}

const SignIn = function() {
    const [userName, setUserName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [errorMessage, setErrorMessage] = useState('')
    const [useEmail, setUseEmail] = useState(false)
    const navigate = useNavigate()

    const submit = async () => {
        if (password.length == 0) {
            setErrorMessage('Password is missing')
            return
        }
        if (password.length < 6 || password.length > 10) {
            setErrorMessage(
                'Password should be greater or equal to 6 and less than or equal to 10 characters long.',
            )
            return
        }
        if (userName.length === 0 && !useEmail) {
            setErrorMessage('User name is missing.')
            return
        }
        if (userName.length > 30) {
            setErrorMessage('User name should be less or equal to 30 characters.')
            return
        }
        if (email.length === 0 && useEmail) {
            setErrorMessage('Email is missing.')
            return
        }

        try {
            const { data } = await axios.post<SignInResponse>(
                generateURL(config.api.account, '/signin'),
                {
                    user_name: userName,
                    email: email,
                    password: password,
                },
            )
            localStorage.setItem('access_token', data.token_pair.access_token)
            localStorage.setItem(
                'refresh_token',
                data.token_pair.refresh_token.token_string,
            )
            setErrorMessage('')
            return navigate('/dashboard')
        } catch (error) {
            asAxiosError(error, (error) => {
                console.log(error.response)
                switch (error.response?.status) {
                    case 400: {
                        setErrorMessage('Bad request.')
                        break
                    }
                    case 401: {
                        const id = useEmail ? 'Email' : 'User name'
                        setErrorMessage(`${id} or password is not correct.`)
                        break
                    }
                    case 500: {
                        setErrorMessage('Server is not responding. Please wait.')
                        break
                    }
                }
            })
        }
    }

    return (
        <div className='flex-center'>
            <div className='signin-root'>
                <div className='signin-title'>Sign In To Your Account</div>

                <div className='signin-form'>
                    {!useEmail && (
                        <div className='signin-row'>
                            <label className='signin-label'>User Name</label>
                            <input
                                type='text'
                                placeholder='Your user name'
                                value={userName}
                                onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
                                    setUserName(e.target.value)
                                }}
                                className='signin-input'
                            ></input>
                        </div>
                    )}
                    {useEmail && (
                        <div className='signin-row'>
                            <label className='signin-label'>Email</label>
                            <input
                                type='text'
                                placeholder='Your email'
                                value={email}
                                onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
                                    setEmail(e.target.value)
                                }}
                                className='signin-input'
                            ></input>
                        </div>
                    )}
                    <div className='signin-row'>
                        <label className='signin-label'>Password</label>
                        <input
                            type='password'
                            placeholder='Your password'
                            value={password}
                            onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
                                setPassword(e.target.value)
                            }}
                            className='signin-input'
                        ></input>
                    </div>
                </div>
                <div className='signin-switch-signin-method'>
                    <button
                        className='signin-switch-signin-button'
                        onClick={() => {
                            if (useEmail) {
                                setEmail('')
                            } else {
                                setUserName('')
                            }
                            setUseEmail(!useEmail)
                        }}
                    >
                        {useEmail ? 'Use user name' : 'Use email'}
                    </button>
                </div>
                <div className='signin-error-message'>{errorMessage}</div>
                <button className='signin-submit' onClick={submit}>
                    Sign In
                </button>
                <div className='signin-to-signup'>
                    <button
                        className='signin-to-signup-button'
                        onClick={() => navigate('/signup')}
                    >
                        No account? Goto sign up
                    </button>
                </div>
            </div>
        </div>
    )
}

export default SignIn
