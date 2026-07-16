import React, { useState } from 'react';
import { MapPin, Phone, Mail, Clock, Send } from 'lucide-react';
import { API_URL } from '../config';
import { useToast } from '../components/UI/ToastNotification'; 

export default function Contact({ t }) {
  const [contactForm, setContactForm] = useState({name: '', phone: '', email: '', message: ''});
  const [isSending, setIsSending] = useState(false);
  const showToast = useToast(); // Initialize useToast

  const handleContactSubmit = async (e) => { 
      e.preventDefault(); setIsSending(true); 
      try { 
          const res = await fetch(`${API_URL}/contact`, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(contactForm) }); 
          if(res.ok) { 
              showToast("Pesan Terkirim! Admin akan membalas via Email.", 'success'); 
          } else {
              showToast("Gagal mengirim pesan.", 'error'); // Use toast for failure
          }
      } catch(e) {
          showToast("Error koneksi backend: " + e.message, 'error'); // Use toast for connection error
      } 
      setIsSending(false); 
  }

  return (
    <div className="pt-32 min-h-screen bg-white px-6 animate-fade-in">
        <div className="max-w-7xl mx-auto mb-20">
            <div className="text-center mb-16">
                <span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t?.nav?.contact || "Kontak"}</span>
                <h2 className="font-serif text-4xl md:text-5xl font-bold mt-2 text-stone-800">{t?.contact?.title || "Hubungi Kami"}</h2>
            </div>
            
            <div className="grid md:grid-cols-3 gap-12">
                {/* Informasi Kontak */}
                <div className="md:col-span-1 space-y-8">
                    <div className="bg-stone-50 p-8 rounded-2xl border border-stone-100">
                        <h3 className="font-bold text-xl text-stone-800 mb-6">Informasi Kontak</h3>
                        <div className="space-y-6">
                            <div className="flex gap-4">
                                <div className="w-10 h-10 rounded-full bg-[#4EC5C1]/10 flex items-center justify-center text-[#4EC5C1] flex-shrink-0"><MapPin size={20}/></div>
                                <div><p className="font-bold text-stone-800">Alamat</p><p className="text-sm text-stone-500">Blumbang, Campurdarat, Tulungagung Regency, East Java 66272 Agung</p></div>
                            </div>
                            <div className="flex gap-4">
                                <div className="w-10 h-10 rounded-full bg-[#4EC5C1]/10 flex items-center justify-center text-[#4EC5C1] flex-shrink-0"><Phone size={20}/></div>
                                <div><p className="font-bold text-stone-800">Telepon</p><p className="text-sm text-stone-500">+6285257121887</p></div>
                            </div>
                            <div className="flex gap-4">
                                <div className="w-10 h-10 rounded-full bg-[#4EC5C1]/10 flex items-center justify-center text-[#4EC5C1] flex-shrink-0"><Mail size={20}/></div>
                                <div><p className="font-bold text-stone-800">Email</p><p className="text-sm text-stone-500">culturestone@gmail.com</p></div>
                            </div>
                            <div className="flex gap-4">
                                <div className="w-10 h-10 rounded-full bg-[#4EC5C1]/10 flex items-center justify-center text-[#4EC5C1] flex-shrink-0"><Clock size={20}/></div>
                                <div><p className="font-bold text-stone-800">Jam Buka</p><p className="text-sm text-stone-500">08.00 - 17.00</p></div>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Form Pesan */}
                <div className="md:col-span-2 bg-white p-8 md:p-12 rounded-2xl shadow-xl border border-stone-100">
                    <form onSubmit={handleContactSubmit} className="space-y-6">
                        <div className="grid md:grid-cols-2 gap-6">
                            {/* BUG-07 FIX: Nama hanya boleh huruf dan spasi */}
                            <div>
                                <label className="block text-sm font-bold text-stone-700 mb-2">{t?.contact?.form_name || "Nama"}</label>
                                <input
                                    required
                                    maxLength={100}
                                    className="w-full p-4 rounded-xl bg-stone-50 border border-stone-200 focus:ring-2 focus:ring-[#4EC5C1] outline-none"
                                    value={contactForm.name}
                                    placeholder="Contoh: Budi Santoso"
                                    onChange={e => {
                                        // Hanya izinkan huruf dan spasi
                                        const val = e.target.value.replace(/[^a-zA-Z\s\u00C0-\u024F]/g, '');
                                        setContactForm({...contactForm, name: val});
                                    }}
                                />
                            </div>
                            {/* BUG-06 FIX: No HP hanya boleh angka, +, -, spasi */}
                            <div>
                                <label className="block text-sm font-bold text-stone-700 mb-2">{t?.contact?.form_phone || "No. HP"}</label>
                                <input
                                    required
                                    type="tel"
                                    maxLength={20}
                                    className="w-full p-4 rounded-xl bg-stone-50 border border-stone-200 focus:ring-2 focus:ring-[#4EC5C1] outline-none"
                                    value={contactForm.phone}
                                    placeholder="Contoh: 08123456789"
                                    onChange={e => {
                                        // Hanya izinkan angka, +, -, spasi
                                        const val = e.target.value.replace(/[^0-9+\-\s]/g, '');
                                        setContactForm({...contactForm, phone: val});
                                    }}
                                />
                            </div>
                        </div>
                        <div><label className="block text-sm font-bold text-stone-700 mb-2">Email</label><input type="email" required className="w-full p-4 rounded-xl bg-stone-50 border border-stone-200 focus:ring-2 focus:ring-[#4EC5C1] outline-none" value={contactForm.email} onChange={e => setContactForm({...contactForm, email: e.target.value})}/></div>
                        <div><label className="block text-sm font-bold text-stone-700 mb-2">{t?.contact?.form_msg || "Pesan"}</label><textarea required rows="5" className="w-full p-4 rounded-xl bg-stone-50 border border-stone-200 focus:ring-2 focus:ring-[#4EC5C1] outline-none" value={contactForm.message} onChange={e => setContactForm({...contactForm, message: e.target.value})}></textarea></div>
                        <button disabled={isSending} className="w-full bg-[#4EC5C1] hover:bg-[#3ba8a5] text-white font-bold py-4 rounded-xl transition shadow-lg flex justify-center items-center gap-2">{isSending ? 'Mengirim...' : <><Send size={18}/> {t?.contact?.btn_send || "Kirim"}</>}</button>
                    </form>
                </div>
            </div>

            <div className="mt-16 w-full h-96 rounded-3xl overflow-hidden shadow-lg border border-stone-200 grayscale hover:grayscale-0 transition duration-700 group">
                <iframe 
                    title="Lokasi Culturstone"
                    width="100%" 
                    height="100%" 
                    
                    src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d5426.23006178796!2d111.85300936872747!3d-8.166993754585587!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e78e10af8aaed2d%3A0x1ae5b1280d74e838!2sCulturestone!5e0!3m2!1sen!2sid!4v1765049440438!5m2!1sen!2sid" 
                    style={{border:0}} 
                    allowFullScreen="" 
                    loading="lazy"
                    referrerPolicy="no-referrer-when-downgrade"
                ></iframe>
                <div className="hidden group-hover:flex absolute top-4 right-4 bg-white/90 backdrop-blur px-4 py-2 rounded-lg shadow text-xs font-bold text-stone-600 pointer-events-none">
                    <MapPin size={14} className="mr-1 text-[#4EC5C1]"/> Lokasi Kami
                </div>
            </div>

        </div>
    </div>
  );
}