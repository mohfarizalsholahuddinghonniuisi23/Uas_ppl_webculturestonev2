import React from 'react';
import { X, Phone, ArrowRight } from 'lucide-react';

export const ProductModal = ({ product, onClose, t, language }) => (
    <div className="fixed inset-0 bg-black/70 z-[60] flex items-center justify-center p-4 backdrop-blur-md animate-fade-in">
        <div className="bg-white/90 backdrop-blur-lg max-w-4xl w-full rounded-3xl overflow-hidden grid md:grid-cols-2 relative shadow-2xl border border-white/30">
            <button onClick={onClose} className="absolute top-5 right-5 bg-white/30 hover:bg-white/50 text-slate-800 rounded-full p-2 transition z-10 shadow-md"><X/></button>
            <div className="h-full bg-slate-100 flex items-center justify-center p-4">
                <img src={product.ImageURL} className="w-full h-full object-contain rounded-xl shadow-lg" alt={product.Name} />
            </div>
            <div className="p-8 md:p-12 flex flex-col justify-center">
                <h2 className="font-serif text-4xl md:text-5xl font-bold mb-4 text-slate-800 leading-tight">{product.Name}</h2>
                <div className="grid grid-cols-2 gap-4 text-base bg-slate-50 p-5 rounded-xl mb-8 border border-slate-100 text-slate-700 shadow-sm">
                    <div><span className="block text-slate-500 text-sm">{t.modal.size}</span><b className="font-semibold">{product.Size}</b></div>
                    <div><span className="block text-slate-500 text-sm">{t.modal.finish}</span><b className="font-semibold">{product.Finishing}</b></div>
                    <div className="col-span-2"><span className="block text-slate-500 text-sm">{t.modal.quality}</span><b className="font-semibold">{product.Quality}</b></div>
                </div>
                <p className="text-slate-600 mb-10 leading-relaxed text-lg">{language === 'id' ? product.Description : (product.DescriptionEn || product.Description)}</p>
                <a href={`https://wa.me/6285320022112?text=Order: ${product.Name}`} target="_blank" className="w-full bg-emerald-500 text-white text-center py-4 rounded-xl font-bold flex justify-center gap-3 items-center hover:bg-emerald-600 transition shadow-lg hover:shadow-xl group">
                    <Phone size={20} className="group-hover:scale-105 transition-transform"/> {t.modal.btn_wa}
                </a>
            </div>
        </div>
    </div>
);

export const CategoryModal = ({ category, onClose, language, navigate }) => (
    <div className="fixed inset-0 bg-black/70 z-[60] flex items-center justify-center p-4 backdrop-blur-md animate-fade-in">
        <div className="bg-white/90 backdrop-blur-lg max-w-3xl w-full rounded-3xl overflow-hidden relative shadow-2xl flex flex-col md:flex-row border border-white/30">
            <button onClick={onClose} className="absolute top-5 right-5 bg-white/30 hover:bg-white/50 text-slate-800 rounded-full p-2 z-20 transition shadow-md"><X size={20}/></button>
            <div className="md:w-1/2 h-64 md:h-auto flex items-center justify-center p-4 bg-slate-100">
                <img src={category.ImageURL} className="w-full h-full object-contain rounded-xl shadow-lg" alt={language === 'id' ? category.Name : (category.NameEn || category.Name)} />
            </div>
            <div className="md:w-1/2 p-8 md:p-12 flex flex-col justify-center">
                <div className="text-sm font-bold text-teal-500 uppercase tracking-widest mb-2">Kategori Batu</div>
                <h2 className="font-serif text-3xl md:text-4xl font-bold mb-4 text-slate-800 leading-tight">{language === 'id' ? category.Name : (category.NameEn || category.Name)}</h2>
                <div className="w-16 h-1 bg-teal-500 mb-6 rounded-full"></div>
                <p className="text-slate-600 leading-relaxed text-lg mb-8">{language === 'id' ? category.Description : (category.DescriptionEn || category.Description)}</p>
                <button onClick={() => {onClose(); navigate('products');}} className="self-start inline-flex items-center gap-2 px-6 py-3 bg-teal-500 text-white font-bold rounded-xl shadow-md hover:bg-teal-600 transition-colors group">
                    Lihat Koleksi <ArrowRight size={18} className="group-hover:translate-x-1 transition-transform"/>
                </button>
            </div>
        </div>
    </div>
);

export const PortfolioModal = ({ portfolio, onClose, language }) => (
    <div className="fixed inset-0 bg-black/70 z-[60] flex items-center justify-center p-4 backdrop-blur-md animate-fade-in">
        <div className="bg-white/90 backdrop-blur-lg max-w-5xl w-full rounded-3xl overflow-hidden relative shadow-2xl flex flex-col md:flex-row max-h-[90vh] border border-white/30">
            <button onClick={onClose} className="absolute top-5 right-5 bg-white/30 hover:bg-white/50 text-slate-800 rounded-full p-2 z-20 transition shadow-md"><X size={20}/></button>
            <div className="w-full md:w-2/3 h-64 md:h-auto flex items-center justify-center p-4 bg-slate-100">
                <img src={portfolio.ImageURL} className="w-full h-full object-contain rounded-xl shadow-lg" alt={portfolio.Title}/>
            </div>
            <div className="w-full md:w-1/3 p-8 md:p-10 flex flex-col justify-center bg-white overflow-y-auto">
                <div className="text-sm font-bold text-teal-500 uppercase tracking-widest mb-2">Proyek Portofolio</div>
                <h2 className="font-serif text-2xl md:text-3xl font-bold mb-4 text-slate-800 leading-tight">{portfolio.Title}</h2>
                <div className="w-16 h-1 bg-teal-500 mb-6 rounded-full"></div>
                <p className="text-slate-600 leading-relaxed text-base">{language === 'id' ? portfolio.Description : (portfolio.DescriptionEn || portfolio.Description)}</p>
            </div>
        </div>
    </div>
);