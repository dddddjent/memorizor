import ReactDOM from 'react-dom/client'
import { createBrowserRouter, redirect, RouterProvider } from 'react-router-dom'
import SignIn from './app/signin'
import { generateURL } from './util'
import config from './config'
import SignUp from './app/signup'
import Dashboard from './app/dashboard'
import Profile from './app/profile'
import ProfileInfo from './app/profile_info'
import ProfileAccount from './app/profile_account'

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
	},
])

ReactDOM.createRoot(document.getElementById('root')!).render(
	<RouterProvider router={router} />,
)
