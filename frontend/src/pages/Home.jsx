import React from 'react';
import { ArrowRight, PlayCircle, MapPin, Briefcase, Star } from 'lucide-react';

export default function Home({ navigate, t, categories, portfolios, products, testimonials, setSelectedCategory, setSelectedPortfolio, setSelectedProduct, language }) {
  return (
    <div className="animate-fade-in">
        <header className="relative h-screen flex items-center justify-center overflow-hidden">
            <div className="absolute inset-0 z-0"><img src="https://images.unsplash.com/photo-1618221195710-dd6b41faaea6?q=80&w=2000" className="w-full h-full object-cover scale-105" /><div className="absolute inset-0 bg-black/40"></div></div>
            <div className="relative z-10 text-center px-6 max-w-4xl mt-10">
                <div className="inline-block border border-white/30 bg-white/20 backdrop-blur px-4 py-1 rounded-full text-xs font-bold text-white uppercase mb-6">{t.hero.badge}</div>
                <h1 className="font-serif text-5xl md:text-7xl font-bold text-white mb-6 drop-shadow-lg">{t.hero.title_1} <span className="text-[#4EC5C1]">{t.hero.title_2}</span> <br/> {t.hero.title_3}</h1>
                <div className="flex flex-col md:flex-row gap-4 justify-center mt-10">
                    <button onClick={() => navigate('products')} className="bg-[#4EC5C1] text-white px-8 py-4 rounded-full font-bold shadow-lg hover:bg-[#3ba8a5] inline-flex items-center gap-3 transition transform hover:-translate-y-1">{t.hero.btn_catalog} <ArrowRight/></button>
                    <button onClick={() => navigate('portfolio')} className="bg-white text-stone-900 px-8 py-4 rounded-full font-bold hover:bg-stone-100 inline-flex items-center gap-3 transition transform hover:-translate-y-1">{t.hero.btn_portfolio} <PlayCircle/></button>
                </div>
            </div>
        </header>
        <section id="section-about" className="py-20 md:py-24 bg-stone-50 px-6">
            <div className="max-w-7xl mx-auto">
                <div className="grid md:grid-cols-2 gap-16 items-center">
                    <div className="relative group">
                        <div className="absolute -top-4 -left-4 w-24 h-24 bg-[#4EC5C1] rounded-full opacity-20 group-hover:scale-110 transition duration-500"></div>
                        <img src="https://images.unsplash.com/photo-1600585154340-be6161a56a0c?w=800" className="relative z-10 rounded-3xl shadow-2xl w-full h-[500px] object-cover transition transform hover:scale-105 duration-700"/>
                    </div>
                    <div>
                        <div className="flex items-center gap-3 mb-4">
                            <div className="w-10 h-1 bg-[#4EC5C1]"></div>
                            <span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t.about.title}</span>
                        </div>
                        <h2 className="font-serif text-4xl md:text-5xl font-bold mb-6 text-stone-800 leading-tight">{t.about.subtitle}</h2>
                        <p className="text-stone-500 mb-8 text-lg leading-relaxed border-l-4 border-stone-200 pl-6">{t.about.desc}</p>
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                            <div className="bg-white p-6 rounded-xl shadow-sm border border-stone-100 hover:shadow-md transition">
                                <div className="p-3 bg-stone-50 rounded-full inline-block text-[#4EC5C1] mb-3"><MapPin size={24}/></div>
                                <h4 className="font-bold text-lg text-stone-800 mb-1">{t.about.feature_1}</h4>
                                <p className="text-sm text-stone-500">{t.about.feature_1_desc}</p>
                            </div>
                            <div className="bg-white p-6 rounded-xl shadow-sm border border-stone-100 hover:shadow-md transition">
                                <div className="p-3 bg-stone-50 rounded-full inline-block text-[#4EC5C1] mb-3"><Briefcase size={24}/></div>
                                <h4 className="font-bold text-lg text-stone-800 mb-1">{t.about.feature_2}</h4>
                                <p className="text-sm text-stone-500">{t.about.feature_2_desc}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
        <section className="py-20 bg-white px-6">
             <div className="max-w-6xl mx-auto text-center">
                <h2 className="font-serif text-4xl font-bold mb-4 text-stone-800">{t.category.title}</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-8 mt-12">
                    {categories.map(c => (
                        <div key={c.ID} className="group cursor-pointer" onClick={() => setSelectedCategory(c)}>
                            <div className="overflow-hidden rounded-xl shadow-md mb-6 h-48 bg-stone-100 relative">
                                <img src={c.ImageURL} className="w-full h-full object-cover transform group-hover:scale-110 transition duration-500"/>
                            </div>
                            <h3 className="text-xl font-bold text-stone-800 mb-2 group-hover:text-[#4EC5C1] transition">{c.Name}</h3>
                        </div>
                    ))}
                </div>
             </div>
        </section>
        <section className="py-20 bg-stone-50 px-6">
            <div className="max-w-7xl mx-auto">
                <div className="flex justify-between items-end mb-12">
                    <div><span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t.products.subtitle}</span><h2 className="font-serif text-4xl font-bold mt-2 text-stone-800">{t.products.title}</h2></div>
                    <button onClick={() => navigate('products')} className="text-stone-600 font-bold hover:text-[#4EC5C1] flex items-center gap-2">{t.products.btn_all} <ArrowRight size={18}/></button>
                </div>
                <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
                    {products.slice(0, 4).map(p => (
                        <div key={p.ID} className="bg-white rounded-xl shadow-sm hover:shadow-xl transition overflow-hidden group cursor-pointer" onClick={() => setSelectedProduct(p)}>
                            <div className="h-64 overflow-hidden"><img src={p.ImageURL} className="w-full h-full object-cover group-hover:scale-110 transition"/></div>
                            <div className="p-5"><div className="text-xs font-bold text-[#4EC5C1] mb-1 uppercase">{p.Category}</div><h3 className="font-serif font-bold text-xl mb-2 text-stone-800">{p.Name}</h3></div>
                        </div>
                    ))}
                </div>
            </div>
        </section>
        <section className="py-20 bg-white px-6">
             <div className="max-w-7xl mx-auto">
                <div className="text-center mb-12">
                    <span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t.testimonial.subtitle}</span>
                    <h2 className="font-serif text-4xl font-bold mt-2 text-stone-800">{t.testimonial.title}</h2>
                </div>
                <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
                    {testimonials.slice(0, 3).map((testimoni, i) => (
                        <div key={i} className="bg-white p-8 rounded-xl shadow-lg border border-stone-100 text-center transform hover:-translate-y-2 transition-transform duration-300">
                            {testimoni.VideoPath ? (
                                <video src={testimoni.VideoPath} controls className="w-full h-48 object-cover rounded-lg mx-auto mb-4"></video>
                            ) : testimoni.ImagePath ? (
                                <img src={testimoni.ImagePath} alt={testimoni.ClientName} className="w-full h-48 object-cover rounded-lg mx-auto mb-4" />
                            ) : (
                                <Star className="w-24 h-24 text-gray-400 mx-auto mb-4" />
                            )}

                            {testimoni.TestimonialText && testimoni.TestimonialText.trim() !== '' && (
                                <p className="italic text-stone-600 text-lg mt-4 mb-4">"{testimoni.TestimonialText}"</p>
                            )}
                            
                            <div className="mt-6">
                                <h4 className="font-bold text-xl text-stone-800">{testimoni.ClientName}</h4>
                                {testimoni.Portfolio && testimoni.Portfolio.Title && (
                                    <p className="text-sm text-stone-500 mt-1">Project: <span className="font-medium">{testimoni.Portfolio.Title}</span></p>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
             </div>
        </section>
        <section className="py-20 bg-stone-50 px-6">
             <div className="max-w-7xl mx-auto text-center">
                 <h2 className="font-serif text-3xl font-bold mb-6 text-stone-800">Punya Pertanyaan?</h2>
                 <button onClick={() => navigate('contact')} className="bg-[#4EC5C1] text-white px-8 py-4 rounded-full font-bold shadow-lg hover:bg-[#3ba8a5] inline-flex items-center gap-3 transition transform hover:-translate-y-1">
                     Hubungi Kami Sekarang <ArrowRight/>
                 </button>
             </div>
        </section>
    </div>
  );
}