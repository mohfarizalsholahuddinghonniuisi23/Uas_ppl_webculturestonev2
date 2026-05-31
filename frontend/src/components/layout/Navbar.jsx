import React, { useState } from 'react';
import { Menu, X, Globe } from 'lucide-react';

export default function Navbar({ navigate, currentPage, language, setLanguage, t }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <nav className="fixed top-0 w-full z-40 bg-white/95 backdrop-blur-sm border-b border-stone-100 shadow-sm transition-all text-stone-800">
      <div className="max-w-7xl mx-auto px-6 py-4 flex justify-between items-center">
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => navigate('home')}>
            <img src="/cslogonw.png" alt="Culturstone logo" className="h-10" />
            <div className="font-serif font-bold text-2xl tracking-widest text-stone-800">CULTURE <span className="text-[#4EC5C1]">STONE</span></div>
        </div>
        <div className="hidden md:flex gap-8 font-medium text-sm tracking-wide text-stone-600">
           <button onClick={() => navigate('home')} className={`hover:text-[#4EC5C1] transition ${currentPage==='home'?'text-[#4EC5C1] font-bold':''}`}>{t.nav.home}</button>
           <button onClick={() => navigate('companyprofilepage')} className={`hover:text-[#4EC5C1] transition ${currentPage==='companyprofilepage'?'text-[#4EC5C1] font-bold':''}`}>{t.nav.companyProfilePage}</button>
           <button onClick={() => navigate('products')} className={`hover:text-[#4EC5C1] transition ${currentPage==='products'?'text-[#4EC5C1] font-bold':''}`}>{t.nav.products}</button>
           <button onClick={() => navigate('portfolio')} className={`hover:text-[#4EC5C1] transition ${currentPage==='portfolio'?'text-[#4EC5C1] font-bold':''}`}>{t.nav.portfolio}</button>
           <button onClick={() => navigate('testimonials')} className={`hover:text-[#4EC5C1] transition ${currentPage==='testimonials'?'text-[#4EC5C1] font-bold':''}`}>{t.nav.testimonials}</button>
           <button onClick={() => navigate('contact')} className="hover:text-[#4EC5C1] transition">{t.nav.contact}</button>
        </div>
        <div className="flex items-center gap-4">
             <button onClick={() => setLanguage(language === 'id' ? 'en' : 'id')} className="flex items-center gap-1 text-xs font-bold text-stone-600 border border-stone-300 px-2 py-1 rounded hover:bg-stone-100 transition">
                <Globe size={14}/> {language.toUpperCase()}
            </button>
            <button className="md:hidden text-stone-800 p-1" onClick={() => setIsMenuOpen(!isMenuOpen)}>{isMenuOpen ? <X size={24}/> : <Menu size={24}/>}</button>
        </div>
      </div>
      {isMenuOpen && (
        <div className="md:hidden bg-white border-t border-stone-100 absolute w-full left-0 top-full shadow-2xl p-6 flex flex-col space-y-4 text-stone-800 text-center">
            <button onClick={() => {navigate('home'); setIsMenuOpen(false)}}>{t.nav.home}</button> 
            <button onClick={() => {navigate('companyprofilepage'); setIsMenuOpen(false)}}>{t.nav.companyProfilePage}</button>
            <button onClick={() => {navigate('products'); setIsMenuOpen(false)}}>{t.nav.products}</button> di 
            <button onClick={() => {navigate('portfolio'); setIsMenuOpen(false)}}>{t.nav.portfolio}</button> 
            <button onClick={() => {navigate('testimonials'); setIsMenuOpen(false)}}>{t.nav.testimonials}</button> 
            <button onClick={() => {navigate('contact'); setIsMenuOpen(false)}}>{t.nav.contact}</button>
        </div>
      )}
    </nav>
  );
}