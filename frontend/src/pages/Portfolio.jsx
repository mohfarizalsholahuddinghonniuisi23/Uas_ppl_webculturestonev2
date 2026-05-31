import { Search } from 'lucide-react';
import { useState } from 'react';

export default function Portfolio({ t, portfolios, language, setSelectedPortfolio }) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredPortfolios = portfolios.filter(pf => 
      (pf.Title || '').toLowerCase().includes(searchQuery.toLowerCase()) ||
      (pf.Description || '').toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="pt-32 min-h-screen bg-white px-6 animate-fade-in">
        <div className="max-w-7xl mx-auto">
            <div className="text-center mb-12">
                <span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t.portfolio.subtitle}</span>
                <h2 className="font-serif text-4xl font-bold mt-2 text-stone-800">{t.portfolio.title || 'Portofolio'}</h2>
                <div className="max-w-md mx-auto mt-6 relative">
                    <input type="text" placeholder="Cari proyek..." className="w-full pl-10 pr-4 py-3 rounded-full border border-stone-200 outline-none focus:ring-2 focus:ring-[#4EC5C1] text-stone-800" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)}/>
                    <Search className="absolute left-3 top-3.5 text-stone-400" size={18}/>
                </div>
            </div>
            <div className="grid md:grid-cols-3 gap-6">
                {filteredPortfolios.length > 0 ? filteredPortfolios.map(pf => (
                    <div key={pf.ID} className="relative group overflow-hidden rounded-xl h-80 shadow-lg cursor-pointer" onClick={() => setSelectedPortfolio(pf)}>
                        <img src={pf.ImageURL} className="w-full h-full object-cover transition duration-700 group-hover:scale-110"/>
                        <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent flex items-end p-6 opacity-90 transition">
                            <div>
                                <h4 className="font-bold text-xl text-white">{pf.Title}</h4>
                                <p className="text-sm text-stone-300 mt-1">{language==='id'?pf.Description:(pf.DescriptionEn||pf.Description)}</p>
                            </div>
                        </div>
                    </div>
                )) : <div className="col-span-3 text-center text-stone-500">Proyek tidak ditemukan.</div>}
            </div>
        </div>
    </div>
  );
}