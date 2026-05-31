import React from 'react';
import { Award, TrendingUp, Users, Phone } from 'lucide-react'; 

export default function CompanyProfilePage({ t }) {
  return (
    <div className="bg-white text-slate-800 min-h-screen">
      {/* Navbar handled in App.jsx, assume it's sticky with blur */}

      <main className="container mx-auto px-4 py-16 sm:px-6 lg:px-8">
        {/* Hero Section */}
        <section className="relative flex flex-col md:flex-row items-center justify-between py-20 lg:py-32 mb-24">
          <div className="md:w-1/2 text-center md:text-left z-10">
            <h1 className="text-5xl lg:text-7xl font-extrabold leading-tight tracking-tighter mb-6 bg-gradient-to-r from-teal-600 to-emerald-500 text-transparent bg-clip-text">
              {t.companyProfileHeroMainTitle}
            </h1>
            <p className="text-xl lg:text-2xl text-slate-600 max-w-lg mx-auto md:mx-0">
              {t.companyProfileHeroSubtitle}
            </p>
          </div>
          <div className="md:w-1/2 relative flex justify-center mt-12 md:mt-0 z-0">
            <div className="relative w-full max-w-md lg:max-w-xl">
              {/* Main Aesthetic Photo */}
              <img
                src="/placeholder-hero-company.jpg" 
                alt="Culturstone Craftsmanship"
                className="w-full h-auto rounded-3xl shadow-xl transform rotate-3 hover:rotate-0 transition-all duration-300 ease-in-out"
              />
              {/* Overlapping Quality Badge */}
              <div className="absolute -bottom-8 -left-8 md:-bottom-12 md:-left-12 bg-white p-6 rounded-2xl shadow-2xl backdrop-blur-sm bg-opacity-80 border border-gray-200 transform rotate-[-5deg] hover:rotate-0 transition-all duration-300 ease-in-out z-20">
                <div className="flex items-center space-x-3">
                  <Award size={36} className="text-teal-500" />
                  <div>
                    <p className="text-xl font-bold text-teal-600">{t.companyProfileQualityBadgeTitle}</p>
                    <p className="text-sm text-slate-500">{t.companyProfileQualityBadgeSubtitle}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          {/* Background Gradient Blobs */}
          <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-emerald-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-2000"></div>
          <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-teal-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-blob animation-delay-4000"></div>
        </section>

        
        <section className="mb-24">
            <h2 className="text-4xl font-bold text-center text-slate-800 mb-12">{t.companyProfileWhoWeAreTitle}</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-12 items-center bg-white p-8 rounded-3xl shadow-2xl border border-gray-100">
                <div className="prose lg:prose-xl text-slate-700 leading-relaxed order-2 md:order-1">
                    <p className="text-lg font-medium">{t.companyProfilePageParagraph1}</p>
                    <p>{t.companyProfilePageParagraph2}</p>
                </div>
                <div className="order-1 md:order-2">
                    <img
                        src="/cslogonw.png" 
                        alt="Who We Are"
                        className="w-full h-auto rounded-2xl shadow-lg transform -rotate-2 hover:rotate-0 transition-transform duration-300 ease-in-out"
                    />
                </div>
            </div>
            <div className="max-w-5xl mx-auto prose lg:prose-xl text-slate-700 leading-relaxed mt-8 bg-white p-8 rounded-3xl shadow-2xl border border-gray-100">
                <p>{t.companyProfilePageParagraph3}</p>
                <p>{t.companyProfilePageParagraph4}</p>
            </div>
        </section>

        {/* Feature Grid */}
        <section className="mb-24">
          <h2 className="text-4xl font-bold text-center text-slate-800 mb-12">{t.companyProfileFeaturesTitle}</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* Feature Card 1 */}
            <div className="bg-white p-8 rounded-2xl shadow-xl border border-gray-100 transform hover:-translate-y-2 transition-transform duration-200 ease-in-out">
              <Award size={48} className="text-teal-500 mb-6 mx-auto" />
              <h3 className="text-2xl font-semibold text-slate-800 mb-4 text-center">{t.companyProfileFeature1Title}</h3>
              <p className="text-slate-600 leading-relaxed text-center">{t.companyProfileFeature1Text}</p>
            </div>
            {/* Feature Card 2 */}
            <div className="bg-white p-8 rounded-2xl shadow-xl border border-gray-100 transform hover:-translate-y-2 transition-transform duration-200 ease-in-out">
              <TrendingUp size={48} className="text-teal-500 mb-6 mx-auto" />
              <h3 className="text-2xl font-semibold text-slate-800 mb-4 text-center">{t.companyProfileFeature2Title}</h3>
              <p className="text-slate-600 leading-relaxed text-center">{t.companyProfileFeature2Text}</p>
            </div>
            {/* Feature Card 3 */}
            <div className="bg-white p-8 rounded-2xl shadow-xl border border-gray-100 transform hover:-translate-y-2 transition-transform duration-200 ease-in-out">
              <Users size={48} className="text-teal-500 mb-6 mx-auto" />
              <h3 className="text-2xl font-semibold text-slate-800 mb-4 text-center">{t.companyProfileFeature3Title}</h3>
              <p className="text-slate-600 leading-relaxed text-center">{t.companyProfileFeature3Text}</p>
            </div>
          </div>
        </section>

        
        
      </main>
    </div>
  );
}