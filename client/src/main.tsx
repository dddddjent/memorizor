import ReactDOM from 'react-dom/client'
import { createBrowserRouter, Outlet, RouterProvider } from 'react-router-dom'

const router = createBrowserRouter([
	{
		path: '/',
		element: (
			<div>
				111
				<Outlet />
			</div>
		),
		children: [
			{
				path: 'sub',
				element: <div>2222</div>,
			},
		],
	},
	{
		path: '/sss',
		element: 123,
	},
])

ReactDOM.createRoot(document.getElementById('root')!).render(
	<RouterProvider router={router} />,
)
