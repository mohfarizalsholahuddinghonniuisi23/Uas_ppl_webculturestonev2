import { Star } from 'lucide-react';

const TestimonialsPage = ({ testimonials, t }) => {
  return (
    <div className="pt-32 pb-24 min-h-screen bg-stone-50 px-6 animate-fade-in">
        <div className="max-w-7xl mx-auto">
            <div className="text-center mb-16">
                <span className="text-[#4EC5C1] font-bold tracking-widest uppercase text-sm">{t.testimonial.subtitle}</span>
                <h2 className="font-serif text-4xl md:text-5xl font-bold mt-2 text-stone-800">{t.testimonial.title}</h2>
            </div>
            
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
                {testimonials && Array.isArray(testimonials) && testimonials.length > 0 ? testimonials.map(item => (
                    <div key={item.ID} className="bg-white p-8 rounded-xl shadow-lg border border-stone-100 text-center transform hover:-translate-y-2 transition-transform duration-300">
                        {/* Display media (video or image) first */}
                        {item.VideoPath ? (
                            <video src={item.VideoPath} controls className="w-full h-48 object-cover rounded-lg mx-auto mb-4"></video>
                        ) : item.ImagePath ? (
                            <img src={item.ImagePath} alt={item.ClientName} className="w-full h-48 object-cover rounded-lg mx-auto mb-4" />
                        ) : (
                            // Fallback if no image or video (e.g., for client testimonials without client image)
                            <Star className="w-24 h-24 text-gray-400 mx-auto mb-4" />
                        )}

                        {/* Display testimonial text if available */}
                        {item.TestimonialText && item.TestimonialText.trim() !== '' && (
                            <p className="italic text-stone-600 text-lg mt-4 mb-4">"{item.TestimonialText}"</p>
                        )}
                        
                        {/* Display Client Name (or Project Name) */}
                        <div className="mt-6">
                            <h4 className="font-bold text-xl text-stone-800">{item.ClientName}</h4>
                            {item.Portfolio && item.Portfolio.Title && (
                                <p className="text-sm text-stone-500 mt-1">Project: <span className="font-medium">{item.Portfolio.Title}</span></p>
                            )}
                        </div>
                    </div>
                )) : (
                  <div className="col-span-full text-center text-stone-500 py-20">
                    <h3 className="text-2xl font-bold">Belum Ada Testimoni</h3>
                    <p className="mt-2">Jadilah yang pertama memberikan ulasan untuk produk kami!</p>
                  </div>
                )}
            </div>
        </div>
    </div>
  );
};

export default TestimonialsPage;
