import { createBrowserRouter } from 'react-router-dom'
import LoginForm from './components/LoginForm'
import Me from './components/Me'
import ErrorPage from './error-page'

const router = createBrowserRouter([
  {
    path: '/',
    element: <Me />,
    errorElement: <ErrorPage />,
  },
  {
    path: '/me',
    element: <Me />,
  },
  {
    path: '/login',
    element: <LoginForm />,
  },
])

export default router
