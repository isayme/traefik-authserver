import ReactDOM from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { Toaster } from './components/ui/toaster.tsx'
import './globals.css'
import './index.css'
import router from './router.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <>
    <RouterProvider router={router} />
    <Toaster />
  </>,
)
