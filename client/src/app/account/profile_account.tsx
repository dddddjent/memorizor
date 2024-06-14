import { useEffect, useState } from 'react'
import '../../style/util.css'
import '../../style/user.css'
import anonymouImageURL from '../../assets/anonymous.png'
import { asAxiosError, generateURL, tryRequest } from '../../util'
import { useNavigate } from 'react-router-dom'
import axios, { AxiosResponse } from 'axios'
import config from '../../config'

type UpdateErrorData = {
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
	const data = response.data as UpdateErrorData
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

const ProfileAccount = function () {
	const [userInfo, setUserInfo] = useState({
		userName: '',
		profileImageURL: '',
		email: '',
		password: '',
	})
	const [isEditing, setIsEditing] = useState(false)
	const [userEditingInfo, setUserEditingInfo] = useState({
		user_name: '',
		email: '',
		password: '',
	})
	const [errorMessage, setErrorMessage] = useState('')
	const navigate = useNavigate()

	useEffect(() => {
		tryRequest(
			async () => {
				const {
					data: {
						user: { user_name, profile_image_url, email },
					},
				} = await axios.get(generateURL(config.api.account, '/me'))
				setUserInfo({
					userName: user_name,
					profileImageURL: profile_image_url,
					email: email,
					password: '******',
				})
				localStorage.setItem('profile_image', profile_image_url)
			},
			(error) => console.log(error),
			navigate,
		)
	}, [navigate])

	const handleEdit = () => {
		setUserEditingInfo({
			user_name: userInfo.userName,
			email: userInfo.email,
			password: '',
		})
		setIsEditing(true)
	}

	const handleSave = () => {
		tryRequest(
			async () => {
				const {
					data: {
						user: { user_name, email },
					},
				} = await axios.post(
					generateURL(config.api.account, '/me'),
					userEditingInfo.password === ''
						? {
								user_name: userEditingInfo.user_name,
								email: userEditingInfo.email,
								// eslint-disable-next-line no-mixed-spaces-and-tabs
						  }
						: userEditingInfo,
				)
				setUserInfo({
					...userInfo,
					userName: user_name,
					email: email,
					password: userEditingInfo.password,
				})
				setErrorMessage('')
				setIsEditing(false)
			},
			(error) => {
				asAxiosError(error, (error) => {
					if (error.response?.status === 400)
						extracBadRequest(error.response, setErrorMessage)
					else if (error.response?.status === 409) {
						const response = error.response?.data as {
							error: {
								type: string
								message: string
							}
						}
						if (response.error.message.search('Resource Email') === 0)
							setErrorMessage('Email already exists')
						else if (response.error.message.search('Resource User') === 0)
							setErrorMessage('User name already exists')
						else setErrorMessage('Conflicted')
						console.log(error)
					} else {
						setErrorMessage('Failed to save. Please retry.')
						console.log(error)
					}
				})
			},
			navigate,
		)
	}

	const image_url =
		localStorage.getItem('profile_image') === null
			? ''
			: (localStorage.getItem('profile_image') as string)

	return (
		<div className='flex-center'>
			<div id='user-image-container'>
				<img
					id='user-image'
					src={image_url === '' ? anonymouImageURL : image_url}
				/>
			</div>
			<div id='user-info'>
				<div className='user-info-row'>
					<span className='user-info-label'>User name: </span>
					{!isEditing && (
						<span className='user-info-content'>{userInfo.userName}</span>
					)}
					{isEditing && (
						<input
							className='user-info-content'
							value={userEditingInfo.user_name}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									user_name: e.target.value,
								})
							}}
							placeholder='New user name'
						/>
					)}
				</div>
				<div className='user-info-row'>
					<span className='user-info-label'>Email: </span>
					{!isEditing && (
						<span className='user-info-content'>{userInfo.email}</span>
					)}
					{isEditing && (
						<input
							className='user-info-content'
							value={userEditingInfo.email}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									email: e.target.value,
								})
							}}
							placeholder='New email'
						/>
					)}
				</div>
				{isEditing && (
					<div className='user-info-row'>
						<span className='user-info-label'>Password: </span>
						<input
							type='password'
							className='user-info-content'
							value={userEditingInfo.password}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									password: e.target.value,
								})
							}}
							placeholder='New password or leave it empty'
						/>
					</div>
				)}
			</div>
			<div>{errorMessage}</div>
			<div
				style={{
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center',
					width: '100%',
				}}
			>
				{!isEditing && (
					<button
						onClick={handleEdit}
						style={{
							width: '3rem',
						}}
					>
						Edit
					</button>
				)}
				{isEditing && (
					<div>
						<button
							onClick={handleSave}
							style={{
								marginRight: '1rem',
								width: '3.5rem',
							}}
						>
							Save
						</button>
						<button
							onClick={() => setIsEditing(false)}
							style={{
								width: '3.5rem',
							}}
						>
							Cancel
						</button>
					</div>
				)}
			</div>
		</div>
	)
}

export default ProfileAccount
