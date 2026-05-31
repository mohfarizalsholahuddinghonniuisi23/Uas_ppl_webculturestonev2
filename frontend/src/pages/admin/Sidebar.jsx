import React from 'react';
import { LayoutDashboard, Package, MessageSquare, Star, Briefcase, Layers, LogOut, Users, Settings } from 'lucide-react'; // Added Users and Settings for potential future use

export default function Sidebar({ activeTab, setActiveTab, onLogout }) {
  const commonButtonClasses = "w-full flex items-center gap-3 p-3 rounded-xl transition-all duration-200 text-slate-300 hover:bg-slate-700 hover:text-white";
  const activeButtonClasses = "bg-slate-700 text-white font-semibold shadow-md";

  return (
    <div className="w-64 bg-slate-900 text-white p-6 flex flex-col sticky top-0 h-screen shadow-2xl">
        <div className="flex items-center gap-3 mb-10">
            <div className="w-8 h-8 bg-[#4EC5C1] rounded-lg rotate-45 flex items-center justify-center"><div className="w-4 h-4 bg-white -rotate-45"></div></div>
            <h2 className="text-2xl font-bold text-white tracking-wide">Cultur<span className="text-[#4EC5C1]">stone</span></h2>
        </div>
        <nav className="space-y-2 flex-1">
            <button onClick={() => setActiveTab('dashboard')} className={`${commonButtonClasses} ${activeTab==='dashboard' ? activeButtonClasses : ''}`}>
                <LayoutDashboard size={20} className={activeTab==='dashboard' ? "text-[#4EC5C1]" : ""}/> Dashboard
            </button>
            <button onClick={() => setActiveTab('categories')} className={`${commonButtonClasses} ${activeTab==='categories' ? activeButtonClasses : ''}`}>
                <Layers size={20} className={activeTab==='categories' ? "text-[#4EC5C1]" : ""}/> Kategori
            </button>
            <button onClick={() => setActiveTab('products')} className={`${commonButtonClasses} ${activeTab==='products' ? activeButtonClasses : ''}`}>
                <Package size={20} className={activeTab==='products' ? "text-[#4EC5C1]" : ""}/> Produk
            </button>
            <button onClick={() => setActiveTab('portfolios')} className={`${commonButtonClasses} ${activeTab==='portfolios' ? activeButtonClasses : ''}`}>
                <Briefcase size={20} className={activeTab==='portfolios' ? "text-[#4EC5C1]" : ""}/> Portofolio
            </button>
            <button onClick={() => setActiveTab('testimonials')} className={`${commonButtonClasses} ${activeTab==='testimonials' ? activeButtonClasses : ''}`}>
                <Star size={20} className={activeTab==='testimonials' ? "text-[#4EC5C1]" : ""}/> Testimoni
            </button>

            <button onClick={() => setActiveTab('messages')} className={`${commonButtonClasses} ${activeTab==='messages' ? activeButtonClasses : ''}`}>
                <MessageSquare size={20} className={activeTab==='messages' ? "text-[#4EC5C1]" : ""}/> Pesan
            </button>
            <button onClick={() => setActiveTab('visitors')} className={`${commonButtonClasses} ${activeTab==='visitors' ? activeButtonClasses : ''}`}>
                <Users size={20} className={activeTab==='visitors' ? "text-[#4EC5C1]" : ""}/> Pengunjung
            </button>
            {/* Added a placeholder for settings if needed */}
            {/* <button onClick={() => setActiveTab('settings')} className={`${commonButtonClasses} ${activeTab==='settings' ? activeButtonClasses : ''}`}>
                <Settings size={20} className={activeTab==='settings' ? "text-[#4EC5C1]" : ""}/> Pengaturan
            </button> */}
        </nav>
        <button onClick={onLogout} className="mt-8 flex items-center justify-center gap-2 p-3 rounded-xl bg-slate-700 text-red-400 hover:bg-slate-600 hover:text-red-300 transition-colors duration-200 shadow-md">
            <LogOut size={20} /> Keluar
        </button>
    </div>
  );
}