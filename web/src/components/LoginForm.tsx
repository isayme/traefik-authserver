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

export default function Component() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className='text-2xl'>Login</CardTitle>
        <CardDescription>
          Enter your email below to login to your account.
        </CardDescription>
      </CardHeader>
      <CardContent className='space-y-4'>
        <div className='space-y-2'>
          <Label htmlFor='email'>Email</Label>
          <Input id='email' placeholder='m@example.com' required type='email' />
        </div>
        <div className='space-y-2'>
          <Label htmlFor='password'>Password</Label>
          <Input id='password' required type='password' />
        </div>
      </CardContent>
      <CardFooter>
        <Button className='w-full'>Sign in</Button>
      </CardFooter>
    </Card>
  )
}
