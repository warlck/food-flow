import React from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from '@/components/ui/use-toast';

const Login: React.FC = () => {
  const navigate = useNavigate();

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    toast({
      title: "Login successful",
      description: "Welcome back to FoodFlow!",
    });
    navigate('/');
  };

  return (
    <Layout>
      <div className="container mx-auto px-4 py-8 max-w-md">
        <h1 className="text-3xl font-bold mb-6 text-center">Login to FoodFlow</h1>
        
        <Card>
          <CardHeader>
            <CardTitle className="text-center">Welcome Back</CardTitle>
            <CardDescription className="text-center">
              Sign in to your account to continue
            </CardDescription>
          </CardHeader>
          
          <CardContent>
            <form onSubmit={handleLogin} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input id="email" type="email" placeholder="your@email.com" required />
              </div>
              
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <Label htmlFor="password">Password</Label>
                  <Link 
                    to="/forgot-password" 
                    className="text-sm text-food-primary hover:underline"
                  >
                    Forgot password?
                  </Link>
                </div>
                <Input id="password" type="password" placeholder="••••••••" required />
              </div>
              
              <Button 
                type="submit" 
                className="w-full bg-food-primary hover:bg-food-accent"
              >
                Sign In
              </Button>
            </form>
          </CardContent>
          
          <CardFooter className="flex flex-col space-y-4">
            <div className="text-center w-full">
              <span className="text-sm text-gray-500">Don't have an account? </span>
              <Link to="/signup" className="text-sm text-food-primary hover:underline">
                Sign up
              </Link>
            </div>
            
            {/* Social Login Buttons */}
            <div className="space-y-2 w-full">
              <Button variant="outline" className="w-full">
                Continue with Google
              </Button>
              <Button variant="outline" className="w-full">
                Continue with Facebook
              </Button>
            </div>
          </CardFooter>
        </Card>
      </div>
    </Layout>
  );
};

export default Login;
