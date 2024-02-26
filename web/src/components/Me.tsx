import { IMe, logout, me } from '@/utils/account'
import { Loader2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Label } from './ui/label'
import { toast } from './ui/use-toast'

export default function LoginForm() {
  const navigate = useNavigate()

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [currentUser, setCurrentUser] = useState<IMe | null>(null)

  useEffect(function () {
    me()
      .then((user) => {
        setCurrentUser(user)
        console.log(user)
        setLoading(false)
      })
      .catch((err) => {
        console.error(err)
        setError(`error: ${err.message}`)
        setLoading(false)
      })
  }, [])

  function handleLogout() {
    logout()
      .then(() => {
        console.log('logout ok')
        navigate('/login')
      })
      .catch((err) => {
        toast({
          title: 'logout fail',
          description: err.message,
        })
      })
  }

  if (loading) {
    return <Loader2 className='mr-2 h-4 w-4 animate-spin' />
  }

  if (currentUser) {
    return (
      <Label>
        login username: {currentUser.username},{' '}
        <span
          className='text-blue-500 hover:cursor-pointer'
          onClick={handleLogout}
        >
          click here
        </span>{' '}
        to logout.
      </Label>
    )
  }

  return (
    <>
      {error || 'unexpect error'}, go{' '}
      <Link to={`/login`} className='text-blue-500'>
        login
      </Link>
      .
    </>
  )
}
