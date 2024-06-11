import { useEffect, useState } from 'react'
import '../../style/util.css'
import '../../style/user.css'
import anonymouImageURL from '../../assets/anonymous.png'
import { asAxiosError, generateURL, tryRequest } from '../../util'
import { useNavigate } from 'react-router-dom'
import axios, { AxiosResponse } from 'axios'
import config from '../../config'
import Crop from './crop'

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
		case 'Name': {
			if (info.Tag === 'gte') {
				setErrorMessage(`Name should be longer than ${info.Param} characters.`)
			} else {
				setErrorMessage(`Name should be less than ${info.Param} characters.`)
			}
			break
		}
	}
}

const ProfileInfo = function () {
	const [userInfo, setUserInfo] = useState({
		name: '',
		profileImageURL: '',
		website: '',
		bio: '',
	})
	const [isEditing, setIsEditing] = useState(false)
	const [userEditingInfo, setUserEditingInfo] = useState({
		name: '',
		website: '',
		bio: '',
	})
	const [errorMessage, setErrorMessage] = useState('')
	const [openCrop, setOpenCrop] = useState(false)
	const navigate = useNavigate()

	useEffect(() => {
		tryRequest(
			async () => {
				const {
					data: {
						user: { name, profile_image_url, website, bio },
					},
				} = await axios.get(generateURL(config.api.account, '/me'))
				setUserInfo({
					name: name,
					profileImageURL: profile_image_url,
					website: website,
					bio: bio,
				})
				localStorage.setItem('profile_image', profile_image_url)
			},
			(error) => console.log(error),
			navigate,
		)
	}, [navigate])

	const handleEdit = () => {
		setUserEditingInfo({
			name: userInfo.name,
			website: userInfo.website,
			bio: userInfo.bio,
		})
		setIsEditing(true)
	}

	const handleSave = () => {
		tryRequest(
			async () => {
				const {
					data: {
						user: { name, website, bio },
					},
				} = await axios.post(
					generateURL(config.api.account, '/me'),
					userEditingInfo,
				)
				setUserInfo({
					...userInfo,
					name: name,
					website: website,
					bio: bio,
				})
				setErrorMessage('')
				setIsEditing(false)
			},
			(error) => {
				asAxiosError(error, (error) => {
					if (error.response?.status === 400)
						extracBadRequest(error.response, setErrorMessage)
					else {
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
			{openCrop && <Crop onClose={() => setOpenCrop(false)} />}
			<div id='user-image-container'>
				<img
					id='user-image'
					src={image_url === '' ? anonymouImageURL : image_url}
				/>
				<button onClick={() => setOpenCrop(true)}>Upload</button>
			</div>
			<div id='user-info'>
				<div className='user-info-row'>
					<span className='user-info-label'>Name: </span>
					{!isEditing && (
						<span className='user-info-content'>{userInfo.name}</span>
					)}
					{isEditing && (
						<input
							className='user-info-content'
							value={userEditingInfo.name}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									name: e.target.value,
								})
							}}
							placeholder='New name'
						/>
					)}
				</div>
				<div className='user-info-row'>
					<span className='user-info-label'>Website: </span>
					{!isEditing && (
						<span className='user-info-content'>{userInfo.website}</span>
					)}
					{isEditing && (
						<input
							className='user-info-content'
							value={userEditingInfo.website}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									website: e.target.value,
								})
							}}
							placeholder='New website'
						/>
					)}
				</div>
				<div className='user-info-row'>
					<span className='user-info-label'>Bio: </span>
					{!isEditing && (
						<span className='user-info-content' style={{ textAlign: 'left' }}>
							{userInfo.bio}
						</span>
					)}
					{isEditing && (
						<input
							className='user-info-content'
							value={userEditingInfo.bio}
							onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
								setUserEditingInfo({
									...userEditingInfo,
									bio: e.target.value,
								})
							}}
							placeholder='New bio'
						/>
					)}
				</div>
			</div>
			<div>{errorMessage}</div>
			<div>
				{!isEditing && <button onClick={handleEdit}>Edit</button>}
				{isEditing && (
					<div>
						<button onClick={handleSave}>Save</button>
						<button onClick={() => setIsEditing(false)}>Cancel</button>
					</div>
				)}
			</div>
		</div>
	)
}

export default ProfileInfo
