import { Outlet, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Button } from '@/components/ui/button';

export default function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <aside className="w-64 bg-white border-r">
        <div className="p-6">
          <h1 className="text-2xl font-bold text-primary">Flame CRM</h1>
          <p className="text-sm text-gray-500">Welcome, {user?.name}</p>
        </div>
        <nav className="mt-6 px-6 space-y-2">
            <Link to="/" className="block py-2 px-4 rounded hover:bg-gray-100 text-gray-700">Dashboard</Link>
            <Link to="/companies" className="block py-2 px-4 rounded hover:bg-gray-100 text-gray-700">Companies</Link>
            <Link to="/users" className="block py-2 px-4 rounded hover:bg-gray-100 text-gray-700">Users</Link>
            <Link to="/customers" className="block py-2 px-4 rounded hover:bg-gray-100 text-gray-700">Customers</Link>
        </nav>
        <div className="p-6 absolute bottom-0 w-64">
             <Button variant="outline" className="w-full" onClick={handleLogout}>Logout</Button>
        </div>
      </aside>

      <main className="flex-1 p-8">
        <Outlet />
      </main>
    </div>
  );
}