
import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

const Home: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect to menu page with default restaurant ID
    navigate('/menu/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d', { replace: true });
  }, [navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <p className="text-xl">Redirecting to menu...</p>
    </div>
  );
};

export default Home;
