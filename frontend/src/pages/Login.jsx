import React, { useState } from 'react';
import { API_URL } from '../config';
import { useToast } from '../components/UI/ToastNotification'; // Import useToast
// import { Google } from 'lucide-react'; // Original incorrect import

export default function Login({ onLoginSuccess }) {
  const [isRegistering, setIsRegistering] = useState(false);
  const showToast = useToast(); // Initialize useToast

  const handleLogin = async (e) => {
    e.preventDefault();
    const form = e.target;
    try {
        const res = await fetch(`${API_URL}/login`, {
          method: 'POST', headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({username: form.username.value, password: form.password.value})
        });
        if(res.ok) {
          const data = await res.json();
          localStorage.setItem('token', data.token);
          showToast("Login berhasil!", 'success');
          onLoginSuccess();
        } else { 
            const errorData = await res.json();
            showToast(`Login Gagal! ${errorData.error || 'Username atau password salah.'}`, 'error'); 
        }
    } catch (e) { showToast(`Error koneksi backend: ${e.message}`, 'error'); }
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    const form = e.target;
    try {
        const res = await fetch(`${API_URL}/register`, {
          method: 'POST', headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({username: form.username.value, password: form.password.value})
        });
        if(res.ok) {
          showToast("Registrasi berhasil! Silakan login.", 'success');
          setIsRegistering(false);
        } else { 
            const errorData = await res.json();
            showToast(`Registrasi Gagal! ${errorData.error || 'Username sudah terdaftar atau input tidak valid.'}`, 'error'); 
        }
    } catch (e) { showToast(`Error koneksi backend: ${e.message}`, 'error'); }
  };

  return (
    <div className="flex min-h-screen animate-fade-in">
        {/* Left Side: Visual & Quote */}
        <div 
          className="relative hidden md:flex md:w-1/2 bg-cover bg-center items-center justify-center p-8"
          style={{ backgroundImage: "url('/placeholder-marble-bg.jpg')" }} // Placeholder background image
        >
          <div className="absolute inset-0 bg-black bg-opacity-60"></div> {/* Overlay */}
          <div className="relative z-10 text-white text-center max-w-lg">
            <h2 className="text-4xl font-bold italic mb-4 leading-snug">
              "{isRegistering ? 'Mulailah perjalanan Anda membangun warisan digital.' : 'Setiap mahakarya dimulai dengan sebuah visi.'}"
            </h2>
            <p className="text-lg font-light">- Culturstone Admin</p>
          </div>
        </div>

        {/* Right Side: Login/Register Form */}
        <div className="w-full md:w-1/2 bg-white flex items-center justify-center p-8 sm:p-12 lg:p-16">
          <div className="w-full max-w-md bg-white p-8 rounded-2xl shadow-xl border border-gray-100 transform -translate-y-4 md:translate-y-0 transition-transform duration-500 ease-out">
            {/* Culturstone Logo */}
            <div className="flex items-center justify-center gap-2 mb-8">
                <div className="w-10 h-10 bg-teal-500 rounded-lg rotate-45 flex items-center justify-center"><div className="w-5 h-5 bg-white -rotate-45"></div></div>
                <div className="font-serif font-bold text-3xl tracking-widest text-slate-800">CULTUR<span className="text-teal-500">STONE</span></div>
            </div>

            <h2 className="text-3xl font-bold mb-8 text-center text-slate-800">
              {isRegistering ? 'Daftar Akun Admin' : 'Masuk ke Admin Panel'}
            </h2>

            <form onSubmit={isRegistering ? handleRegister : handleLogin} className="space-y-6">
              <div>
                <label htmlFor="username" className="sr-only">Username</label>
                <input 
                  name="username" 
                  placeholder="Username" 
                  className="w-full p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 shadow-sm transition-all duration-200" 
                  required 
                />
              </div>
              <div>
                <label htmlFor="password" className="sr-only">Password</label>
                <input 
                  type="password" 
                  name="password" 
                  placeholder="Password" 
                  className="w-full p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 shadow-sm transition-all duration-200" 
                  required 
                />
              </div>
              
              <button 
                type="submit" 
                className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors duration-200"
              >
                {isRegistering ? 'Daftar Sekarang' : 'Masuk'}
              </button>
            </form>

            <p className="mt-8 text-center text-gray-600">
              {isRegistering ? 'Sudah punya akun?' : 'Belum punya akun?'}
              <button
                type="button"
                onClick={() => setIsRegistering(!isRegistering)}
                className="ml-2 text-teal-600 hover:text-teal-700 font-semibold"
              >
                {isRegistering ? 'Login' : 'Daftar'}
              </button>
            </p>
          </div>
        </div>
    </div>
  );
}