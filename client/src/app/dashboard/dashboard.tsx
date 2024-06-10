import '../../style/util.css'
import '../../style/dashboard.css'
import anonymouImageURL from '../../assets/anonymous.png'
import { useEffect, useState } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { generateURL, tryRequest } from '../../util'
import axios from 'axios'
import config from '../../config'
import Detail, { DetailInterface } from './detail'

export type DetailState = Omit<DetailInterface, 'onClose'>

const Dashboard = function() {
    const [page, setPage] = useState<'today' | 'all'>('today')
    const [profileImgURL, setProfileImgURL] = useState('')
    const [detailParams, setDetailParams] = useState<DetailState>({
        show: true,
        word: { word: 'A', explanation: 'BBBBBB', url: 'www.a.com' },
        editable: { word: true, explanation: true, url: true },
        // editable: { word: false, explanation: false, url: false},
    })
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

    return (
        <div id='profile-root'>
            <div id='titlebar'>
                <button id='titlebar-change-button' onClick={handlePageChange}>
                    All Words
                </button>
                <div className='absolute-center-child'>
                    <span id='title'>Word Memeorizer</span>
                </div>
                <button id='titlebar-button' onClick={() => navigate('/profile')}>
                    <img
                        id='titlebar-img'
                        src={profileImgURL === '' ? anonymouImageURL : profileImgURL}
                        alt='Profile Image'
                    />
                </button>
            </div>
            <hr />
            <div id='board'>
                <Outlet />
            </div>
            {detailParams.show && (
                <Detail
                    show={detailParams.show}
                    word={detailParams.word}
                    editable={detailParams.editable}
                    onClose={() => setDetailParams({ ...detailParams, show: false })}
                />
            )}
        </div>
    )
}

export default Dashboard
