import '../../style/util.css'
import '../../style/dashboard.css'
import anonymouImageURL from '../../assets/anonymous.png'
import { useEffect, useState } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { generateURL, tryRequest } from '../../util'
import axios from 'axios'
import config from '../../config'

const Dashboard = function () {
	const [page, setPage] = useState<'today' | 'all'>('today')
	const [, setProfileImgURL] = useState('')
	const navigate = useNavigate()

	useEffect(() => {
		if (localStorage.getItem('refresh_token') === null) {
			navigate('/signin')
		}
		tryRequest(
			async () => {
				const {
					data: {
						user: { profile_image_url },
					},
				} = await axios.get(generateURL(config.api.account, '/me'))
				setProfileImgURL(profile_image_url)
				localStorage.setItem('profile_image', profile_image_url)
			},
			(error) => console.log(error),
			navigate,
		)
	}, [navigate])

	const handlePageChange = () => {
		if (page == 'today') {
			setPage('all')
			navigate('/dashboard/all')
		} else {
			setPage('today')
			navigate('/dashboard/today')
		}
	}

	const image_url =
		localStorage.getItem('profile_image') === null
			? ''
			: (localStorage.getItem('profile_image') as string)

	return (
		<div id='profile-root'>
			<div id='titlebar'>
				<button id='titlebar-change-button' onClick={handlePageChange}>
					{page == 'today' ? 'All Words' : 'Today'}
				</button>
				<div className='absolute-center-child'>
					<span id='title'>Word Memeorizer</span>
				</div>
				<button id='titlebar-button' onClick={() => navigate('/profile')}>
					<img
						id='titlebar-img'
						src={image_url === '' ? anonymouImageURL : image_url}
						alt='Profile Image'
					/>
				</button>
			</div>
			<hr />
			<div id='board'>
				<Outlet />
			</div>
		</div>
	)
}

export default Dashboard
