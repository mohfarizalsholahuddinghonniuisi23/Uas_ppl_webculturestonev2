import React, { useState } from 'react';
import { Search } from 'lucide-react';

export default function Products({ t, products, setSelectedProduct }) {
  const [searchQuery, setSearchQuery] = useState("");

  const filteredProducts = products.filter(p => 
      p.Name.toLowerCase().includes(searchQuery.toLowerCase()) || 
      p.Category.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="pt-32 min-h-screen bg-stone-50 px-6 animate-fade-in">
        <div className="max-w-7xl mx-auto">
            <div className="text-center mb-12">
                <h2 className="font-serif text-4xl font-bold mb-4 text-stone-800">{t.products.title}</h2>
                <div className="max-w-md mx-auto mt-6 relative">
                    <input type="text" placeholder={t.products.search_ph} className="w-full pl-10 pr-4 py-3 rounded-full border border-stone-200 outline-none focus:ring-2 focus:ring-[#4EC5C1] text-stone-800" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)}/>
                    <Search className="absolute left-3 top-3.5 text-stone-400" size={18}/>
                </div>
            </div>
            <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
                {filteredProducts.length > 0 ? filteredProducts.map(p => (
                    <div key={p.ID} className="bg-white rounded-xl shadow-sm hover:shadow-xl transition overflow-hidden group cursor-pointer" onClick={() => setSelectedProduct(p)}>
                        <div className="h-64 overflow-hidden"><img src={p.ImageURL} className="w-full h-full object-cover group-hover:scale-110 transition"/></div>
                        <div className="p-5">
                            <div className="text-xs font-bold text-[#4EC5C1] mb-1 uppercase">{p.Category}</div>
                            <h3 className="font-serif font-bold text-xl mb-2 text-stone-800">{p.Name}</h3>
                        </div>
                    </div>
                )) : <div className="col-span-4 text-center text-stone-500">Produk tidak ditemukan.</div>}
            </div>
        </div>
    </div>
  );
}