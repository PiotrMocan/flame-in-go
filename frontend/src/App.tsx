import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from '@/pages/Login';
import Dashboard from '@/pages/Dashboard';
import Layout from '@/pages/Layout';
import Companies from '@/pages/Companies';
import Users from '@/pages/Users';
import Customers from '@/pages/Customers';
import Funnels from '@/pages/Funnels';
import { AuthProvider } from '@/context/AuthContext';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="companies" element={<Companies />} />
            <Route path="users" element={<Users />} />
            <Route path="customers" element={<Customers />} />
            <Route path="funnels" element={<Funnels />} />
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;