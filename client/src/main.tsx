import ReactDOM from 'react-dom/client'
import { createBrowserRouter, redirect, RouterProvider } from 'react-router-dom'
import SignIn from './app/account/signin'
import { generateURL } from './util'
import config from './config'
import SignUp from './app/account/signup'
import Dashboard from './app/dashboard/dashboard'
import Profile from './app/account/profile'
import ProfileInfo from './app/account/profile_info'
import ProfileAccount from './app/account/profile_account'
import AllWords from './app/dashboard/all_words/all_words'
import { dashboardLoader } from './app/dashboard/dashboard_loader'

const router = createBrowserRouter([
    {
        path: '/',
        loader: () => {
            return redirect(generateURL(config.host, '/dashboard'))
        },
    },
    {
        path: '/signin',
        element: <SignIn />,
    },
    {
        path: '/signup',
        element: <SignUp />,
    },
    {
        path: '/profile',
        element: <Profile />,
        children: [
            {
                index: true,
                element: <ProfileInfo />,
            },
            {
                path: '/profile/info',
                element: <ProfileInfo />,
            },
            {
                path: '/profile/account',
                element: <ProfileAccount />,
            },
        ],
    },
    {
        path: '/dashboard',
        element: <Dashboard />,
        loader: dashboardLoader,
        children: [
            {
                index: true,
                element: <div>today</div>,
            },
            {
                path: '/dashboard/today',
                element: <div>today</div>,
            },
            {
                path: '/dashboard/all',
                element: <AllWords />,
            },
        ],
    },
])

ReactDOM.createRoot(document.getElementById('root')!).render(
    <RouterProvider router={router} />,
)
