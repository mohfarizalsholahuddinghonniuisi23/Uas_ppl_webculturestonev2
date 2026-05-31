import React from 'react';
import { Instagram, Facebook } from 'lucide-react';

export default function Footer() {
  return (
    <footer className="bg-stone-950 text-stone-400 py-8 text-center text-sm border-t border-stone-800 mt-auto">
        <div className="max-w-7xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center">
            <p>© 2026 Culturstone Indonesia. All Rights Reserved.</p>
            <div className="flex gap-4 mt-4 md:mt-0">
                <a href="#" className="hover:text-[#4EC5C1] transition"><Instagram size={18}/></a>
                <a href="#" className="hover:text-[#4EC5C1] transition"><Facebook size={18}/></a>
            </div>
        </div>
    </footer>
  );
}