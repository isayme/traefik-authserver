import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Github } from 'lucide-react'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { login } from '@/utils/account'
import { useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useToast } from './ui/use-toast'

export default function LoginForm() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { toast } = useToast()

  const [username, setUsername] = useState<string>('')
  const [password, setPassword] = useState<string>('')

  function handleUsernameChange(e: React.ChangeEvent<HTMLInputElement>) {
    e.preventDefault()
    setUsername(e.target.value)
  }

  function handlePasswordChange(e: React.ChangeEvent<HTMLInputElement>) {
    e.preventDefault()
    setPassword(e.target.value)
  }

  function handleLogin() {
    login(username, password)
      .then(() => {
        console.log(`login ok`)
        const nextUrl = searchParams.get('next_url')
        if (nextUrl) {
          window.location.href = nextUrl
          return
        }

        navigate('/me')
      })
      .catch((err) => {
        toast({
          title: 'login fail',
          description: err.message,
        })
      })
  }

  function handleLoginWithGithub() {
    let to = '/oauth/github/login'

    const nextUrl = searchParams.get('next_url')
    if (nextUrl) {
      to = `${to}?next_url=${nextUrl}`
    }

    window.location.href = to
  }

  function handleEnter(event: React.KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      handleLogin()
    }
  }

  return (
    <Card className='border-0'>
      <CardHeader>
        <CardTitle className='text-2xl'>Login</CardTitle>
        <CardDescription>
          Enter your username below to login to your account.
        </CardDescription>
      </CardHeader>
      <CardContent className='space-y-4'>
        <div className='space-y-2'>
          <Label htmlFor='username'>Username</Label>
          <Input
            className='focus-visible:ring-1 focus-visible:ring-offset-1 focus-visible:ring-blue-500 '
            id='username'
            placeholder='username'
            autoFocus={true}
            required
            type='username'
            value={username}
            onChange={handleUsernameChange}
            onKeyDown={handleEnter}
          />
        </div>
        <div className='space-y-2'>
          <Label htmlFor='password'>Password</Label>
          <Input
            id='password'
            className='focus-visible:ring-1 focus-visible:ring-offset-1 focus-visible:ring-blue-500 '
            placeholder='password'
            required
            type='password'
            value={password}
            onChange={handlePasswordChange}
            onKeyDown={handleEnter}
          />
        </div>
      </CardContent>
      <CardFooter>
        <Button
          className='w-full focus-visible:ring-1 focus-visible:ring-offset-1'
          onClick={handleLogin}
        >
          Login
        </Button>
      </CardFooter>

      <div className='relative'>
        <div className='absolute inset-0 flex items-center '>
          <Separator className='w-full opacity-80' />
        </div>
        <div className='relative flex justify-center text-xs uppercase'>
          <span className='bg-background px-2 text-muted-foreground'>OR</span>
        </div>
      </div>

      <div className='flex  justify-center items-center pt-3 pb-4'>
        <Button
          variant='outline'
          className='flex'
          type='button'
          onClick={handleLoginWithGithub}
        >
          <Github className='mr-2 h-4 w-4' />
          Login With GitHub
        </Button>
      </div>
    </Card>
  )
}
