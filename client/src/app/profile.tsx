import { Outlet, useNavigate } from 'react-router-dom'
import '../style/util.css'
import '../style/profile.css'
import { useEffect, useState } from 'react'
import axios from 'axios'
import { generateURL, tryRequest } from '../util'
import config from '../config'

const Profile = function () {
	const navigate = useNavigate()
	const [isInfoPage, setIsInfoPage] = useState(true)

	useEffect(() => {
		if (localStorage.getItem('refresh_token') === null) {
			navigate('/signin')
		}
	}, [navigate])

	const handleSignOut = () => {
		tryRequest(
			async () => {
				await axios.post(generateURL(config.api.account, '/signout'), {
					refresh_token: localStorage.getItem('refresh_token'),
				})
				navigate('/signin')
				localStorage.clear()
			},
			(error) => {
				console.log(error)
				navigate('/signin')
			},
			navigate,
		)
	}

	return (
		<div id='profile'>
			<div id='profile-sidebar'>
				<button id='profile-back' onClick={() => navigate('/dashboard')}>
					Back
				</button>
				<button
					id='profile-information'
					className='profile-sidebar-button'
					style={{ backgroundColor: isInfoPage ? '#ffffff' : '#666666' }}
					onClick={() => {
						setIsInfoPage(true)
						navigate('/profile/info')
					}}
				>
					Information
				</button>
				<button
					id='profile-acccount'
					className='profile-sidebar-button'
					style={{ backgroundColor: !isInfoPage ? '#ffffff' : '#666666' }}
					onClick={() => {
						setIsInfoPage(false)
						navigate('/profile/account')
					}}
				>
					Account
				</button>
				<button id='profile-signout' onClick={handleSignOut}>
					Signout
				</button>
			</div>
			<div id='profile-detail'>
				<Outlet />
			</div>
		</div>
	)
}

export default Profile
