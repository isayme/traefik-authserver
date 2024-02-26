import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { login } from '@/utils/account'
import { useState } from 'react'
import { redirect } from 'react-router-dom'
import { useToast } from './ui/use-toast'

export default function LoginForm() {
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
        redirect('/me')
      })
      .catch((err) => {
        toast({
          title: 'login fail',
          description: err.message,
        })
      })
  }

  return (
    <Card>
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
            id='username'
            placeholder='username'
            required
            type='username'
            value={username}
            onChange={handleUsernameChange}
          />
        </div>
        <div className='space-y-2'>
          <Label htmlFor='password'>Password</Label>
          <Input
            id='password'
            placeholder='password'
            required
            type='password'
            value={password}
            onChange={handlePasswordChange}
          />
        </div>
      </CardContent>
      <CardFooter>
        <Button className='w-full' onClick={handleLogin}>
          Login
        </Button>
      </CardFooter>
    </Card>
  )
}
