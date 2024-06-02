import { useEffect, useState } from 'react'
import '../style/util.css'
import '../style/user.css'
import anonymouImageURL from '../assets/anonymous.png'
import { generateURL, tryRequest } from '../util'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../config'

const ProfileInfo = function () {
	const [userInfo, setUserInfo] = useState({
		name: '',
		profileImageURL: '',
		website: '',
		bio: '',
	})
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
			},
			(error) => console.log(error),
			navigate,
		)
	}, [navigate])

	return (
		<div className='flex-center'>
			<div id='user-image-container'>
				<img
					id='user-image'
					src={
						userInfo.profileImageURL === ''
							? anonymouImageURL
							: userInfo.profileImageURL
					}
					alt='Profile Image'
				/>
			</div>
			<div id='user-info'>
				<div className='user-info-row'>
					<span className='user-info-label'>Name: </span>
					<span className='user-info-content'>{userInfo.name}</span>
				</div>
				<div className='user-info-row'>
					<span className='user-info-label'>Website: </span>
					<span className='user-info-content'>{userInfo.website}</span>
				</div>
				<div className='user-info-row'>
					<span className='user-info-label'>Bio: </span>
					<span className='user-info-content' style={{ textAlign: 'left' }}>
						{userInfo.bio}
					</span>
				</div>
			</div>
		</div>
	)
}

export default ProfileInfo
