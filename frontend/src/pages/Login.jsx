import React from 'react';
import { API_URL } from '../config';
import { useToast } from '../components/UI/ToastNotification';

export default function Login({ onLoginSuccess }) {
  const showToast = useToast();

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

  return (
    <div className="flex min-h-screen animate-fade-in">
        {/* Left Side: Visual & Quote */}
        <div 
          className="relative hidden md:flex md:w-1/2 bg-cover bg-center items-center justify-center p-8"
          style={{ backgroundImage: "url('/placeholder-marble-bg.jpg')" }}
        >
          <div className="absolute inset-0 bg-black bg-opacity-60"></div>
          <div className="relative z-10 text-white text-center max-w-lg">
            <h2 className="text-4xl font-bold italic mb-4 leading-snug">
              "Setiap mahakarya dimulai dengan sebuah visi."
            </h2>
            <p className="text-lg font-light">- Culturstone Admin</p>
          </div>
        </div>

        {/* Right Side: Login Form */}
        {/* BUG-01 FIX (Frontend): Form registrasi dihapus sepenuhnya dari halaman login. */}
        {/* Admin hanya bisa login. Pendaftaran admin baru dilakukan secara manual oleh developer. */}
        <div className="w-full md:w-1/2 bg-white flex items-center justify-center p-8 sm:p-12 lg:p-16">
          <div className="w-full max-w-md bg-white p-8 rounded-2xl shadow-xl border border-gray-100 transform -translate-y-4 md:translate-y-0 transition-transform duration-500 ease-out">
            {/* Culturstone Logo */}
            <div className="flex items-center justify-center gap-2 mb-8">
                <div className="w-10 h-10 bg-teal-500 rounded-lg rotate-45 flex items-center justify-center"><div className="w-5 h-5 bg-white -rotate-45"></div></div>
                <div className="font-serif font-bold text-3xl tracking-widest text-slate-800">CULTUR<span className="text-teal-500">STONE</span></div>
            </div>

            <h2 className="text-3xl font-bold mb-2 text-center text-slate-800">Masuk ke Admin Panel</h2>
            <p className="text-center text-gray-400 text-sm mb-8">Masukkan kredensial admin Anda untuk melanjutkan</p>

            <form onSubmit={handleLogin} className="space-y-6">
              <div>
                <label htmlFor="login-username" className="block text-sm font-semibold text-gray-700 mb-1">Username</label>
                <input
                  id="login-username"
                  name="username"
                  placeholder="Masukkan username"
                  className="w-full p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 shadow-sm transition-all duration-200"
                  required
                />
              </div>
              <div>
                <label htmlFor="login-password" className="block text-sm font-semibold text-gray-700 mb-1">Password</label>
                <input
                  id="login-password"
                  type="password"
                  name="password"
                  placeholder="Masukkan password"
                  className="w-full p-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-teal-500 focus:border-teal-500 shadow-sm transition-all duration-200"
                  required
                />
              </div>
              
              <button
                type="submit"
                className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors duration-200"
              >
                Masuk
              </button>
            </form>

            {/* BUG-01 FIX: Tombol 'Daftar' dihapus. Tidak ada akses publik ke form registrasi. */}
            <p className="mt-6 text-center text-xs text-gray-400">
              Hanya administrator yang berwenang dapat mengakses panel ini.
            </p>
          </div>
        </div>
    </div>
  );
}