import ReactDOM from 'react-dom/client'
import { createBrowserRouter, redirect, RouterProvider } from 'react-router-dom'
import SignIn from './app/signin'
import { generateURL } from './util'
import config from './config'

const router = createBrowserRouter([
    {
        path: '/',
        loader: () => {
            return redirect(generateURL(config.host, "/dashboard"))
        },
    },
    {
        path: '/signin',
        element: <SignIn />,
    },
    {
        path: '/dashboard',
        element: <div>111</div>,
    },
])

ReactDOM.createRoot(document.getElementById('root')!).render(
    <RouterProvider router={router} />,
)
