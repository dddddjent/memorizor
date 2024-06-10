import '../../style/util.css'
import '../../style/signup.css'
import { useState } from 'react'
import axios, { AxiosResponse } from 'axios'
import config from '../../config'
import { asAxiosError, generateURL } from '../../util'
import { useNavigate } from 'react-router-dom'

type SignUpResponse = {
	token_pair: {
		access_token: string
		refresh_token: {
			token_string: string
		}
	}
}

type SignUpErrorData = {
	error: {
		type: string
		message: string
	}
	invalid_args: Array<{
		Field: string
		Param: string
		Tag: string
		Value: string
	}>
}

function extracBadRequest(
	response: AxiosResponse,
	setErrorMessage: React.Dispatch<React.SetStateAction<string>>,
) {
	const data = response.data as SignUpErrorData
	const invalid_args = data.invalid_args
	let missingFields = false
	for (const info of invalid_args) {
		if (info.Tag === 'required') {
			missingFields = true
			break
		}
	}
	if (missingFields === true) {
		setErrorMessage('Missing fields.')
		return
	}
	const info = invalid_args[0]
	switch (info.Field) {
		case 'UserName': {
			if (info.Tag === 'gte') {
				setErrorMessage(
					`User name should be longer than ${info.Param} characters.`,
				)
			} else {
				setErrorMessage(
					`User name should be less than ${info.Param} characters.`,
				)
			}
			break
		}
		case 'Email': {
			setErrorMessage('Invalid email format.')
			break
		}
		case 'Password': {
			if (info.Tag === 'gte') {
				setErrorMessage(
					`Password should be longer than ${info.Param} characters.`,
				)
			} else {
				setErrorMessage(
					`Password should be less than ${info.Param} characters.`,
				)
			}
			break
		}
	}
}

const SignUp = function () {
	const [userName, setUserName] = useState('')
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [errorMessage, setErrorMessage] = useState('')
	const navigate = useNavigate()

	const submit = async () => {
		try {
			const { data } = await axios.post<SignUpResponse>(
				generateURL(config.api.account, '/signup'),
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
			axios.defaults.headers.common[
				'Authorization'
			] = `Bearer ${localStorage.getItem('access_token')}`
			setErrorMessage('')
			return navigate('/dashboard')
		} catch (error) {
			asAxiosError(error, (error) => {
				console.log(error.response)
				switch (error.response?.status) {
					case 400: {
						extracBadRequest(error.response, setErrorMessage)
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
			<div className='signup-root'>
				<div className='signup-title'>Sign Up Your Account</div>

				<div
					className='signup-form'
					onKeyDownCapture={(e) => {
						if (e.key === 'Enter') {
							submit()
						}
					}}
				>
					<div className='signup-row'>
						<label className='signup-label'>User Name</label>
						<input
							type='text'
							placeholder='Your user name'
							value={userName}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserName(e.target.value)
							}}
							className='signup-input'
						></input>
					</div>
					<div className='signup-row'>
						<label className='signup-label'>Email</label>
						<input
							type='text'
							placeholder='Your email'
							value={email}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setEmail(e.target.value)
							}}
							className='signup-input'
						></input>
					</div>
					<div className='signup-row'>
						<label className='signup-label'>Password</label>
						<input
							type='password'
							placeholder='Your password'
							value={password}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setPassword(e.target.value)
							}}
							className='signup-input'
						></input>
					</div>
				</div>

				<div className='signup-error-message'>{errorMessage}</div>

				<button className='signup-submit' onClick={submit}>
					Sign In
				</button>
			</div>
		</div>
	)
}

export default SignUp
