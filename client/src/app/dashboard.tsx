import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

const Dashboard = function() {
    const navigate = useNavigate()
    useEffect(() => {
        if (localStorage.getItem('refresh_token') === null) {
            navigate('/signin')
        }
    }, [navigate])
    return <button onClick={() => navigate('/profile')}>Profile</button>
}

export default Dashboard
