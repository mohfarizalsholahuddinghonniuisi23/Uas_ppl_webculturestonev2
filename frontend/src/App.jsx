import { useEffect, useState } from 'react';
import Footer from './components/layout/Footer';
import Navbar from './components/layout/Navbar';
import { CategoryModal, PortfolioModal, ProductModal } from './components/UI/Modals';
import { API_URL } from './config';
import AdminPanel from './pages/admin/AdminPanel';
import CompanyProfilePage from './pages/CompanyProfilePage';
import Home from './pages/Home';
import TestimonialsPage from './pages/Testimonials';
import { translations } from './translations';
import ContactPage from '/src/pages/Contact';
import Login from '/src/pages/Login';
import PortfolioPage from '/src/pages/Portfolio';
import ProductsPage from '/src/pages/Products';

export default function App() {
  const [currentPage, setCurrentPage] = useState('home'); 
  const [isAdminLoggedIn, setIsAdminLoggedIn] = useState(false);
  const [language, setLanguage] = useState('id'); 
  const t = translations[language] || {}; 

  // Data States
  const [products, setProducts] = useState([]);
  const [portfolios, setPortfolios] = useState([]);
  const [testimonials, setTestimonials] = useState([]);
  const [categories, setCategories] = useState([]); 

  // Modal States
  const [selectedProduct, setSelectedProduct] = useState(null); 
  const [selectedCategory, setSelectedCategory] = useState(null); 
  const [selectedPortfolio, setSelectedPortfolio] = useState(null);

  useEffect(() => {
    fetchPublicData();
    const token = localStorage.getItem('token');
    if (token) setIsAdminLoggedIn(true);
    
    const checkHash = () => {
      const path = window.location.pathname.toLowerCase();
      const hash = window.location.hash;

      // Check pathname first (e.g., /admin, /login)
      if (path === '/admin' || path === '/admin/') {
        const t = localStorage.getItem('token');
        setCurrentPage(t ? 'admin' : 'login');
      } else if (path === '/login' || path === '/login/') {
        setCurrentPage('login');
      }

      // Also check hash for backward compatibility (e.g., #admin, #login)
      if (hash === '#login') setCurrentPage('login');
      if (hash === '#admin') {
         const t = localStorage.getItem('token');
         setCurrentPage(t ? 'admin' : 'login');
      }
      if(hash) window.history.pushState("", document.title, window.location.pathname);
    };
    checkHash();
    window.addEventListener('hashchange', checkHash);
    return () => window.removeEventListener('hashchange', checkHash);
  }, []);

  useEffect(() => {
    const publicPages = ['home', 'products', 'portfolio', 'testimonials', 'companyprofilepage'];
    if (publicPages.includes(currentPage)) {
      fetchPublicData();
    }
  }, [currentPage]);

  const fetchPublicData = async () => {
    try {
      const [prodRes, portRes, testRes, catRes] = await Promise.all([
        fetch(`${API_URL}/products`),
        fetch(`${API_URL}/portfolios`),
        fetch(`${API_URL}/testimoni`),
        fetch(`${API_URL}/categories`),
      ]);

      const productsData = await prodRes.json();
      console.log("Products:", productsData);

      const portfolioData = await portRes.json();
      console.log("Portfolios:", portfolioData);

      const testimonialsData = await testRes.json();
      console.log("Testimonials:", testimonialsData);

      const categoriesData = await catRes.json();
      console.log("Categories:", categoriesData);

      setProducts(productsData);
      setPortfolios(portfolioData);
      setTestimonials(testimonialsData || []); // Ensure empty array if null
      setCategories(categoriesData);

    } catch (e) {
      console.error("Failed to fetch public data:", e);
    }
  };

  const navigate = (targetPage, targetId = null) => {
    setCurrentPage(targetPage);
    if (targetId) setTimeout(() => document.getElementById(targetId)?.scrollIntoView({ behavior: 'smooth' }), 100);
    else window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleLogout = () => { localStorage.removeItem('token'); setIsAdminLoggedIn(false); setCurrentPage('home'); };

  if (currentPage === 'admin' && isAdminLoggedIn) return <AdminPanel onLogout={handleLogout} />;
  if (currentPage === 'login') return <Login onLoginSuccess={() => {setIsAdminLoggedIn(true); setCurrentPage('admin');}} />;

  return (
    <div className="font-sans selection:bg-[#4EC5C1] selection:text-white text-stone-800 bg-white min-h-screen flex flex-col">
        <a href="https://wa.me/6285257121887?text=Halo,%20saya%20ingin%20bertanya%20tentang%20produk%20Culturstone." target="_blank" className="fixed bottom-6 right-6 z-50 bg-[#25D366] text-white p-4 rounded-full shadow-2xl flex items-center gap-3 hover:scale-110 transition group"><img src="https://upload.wikimedia.org/wikipedia/commons/6/6b/WhatsApp.svg" className="w-8 h-8"/></a>
        
        <Navbar navigate={navigate} currentPage={currentPage} language={language} setLanguage={setLanguage} t={t} />
        
        <main className="flex-grow">
          {currentPage === 'home' && (
            <Home 
              navigate={navigate} t={t} language={language}
              products={products} portfolios={portfolios} categories={categories} testimonials={testimonials}
              setSelectedProduct={setSelectedProduct} setSelectedCategory={setSelectedCategory} setSelectedPortfolio={setSelectedPortfolio}
            />
          )}
          
          {currentPage === 'products' && (
             <ProductsPage t={t} products={products} setSelectedProduct={setSelectedProduct} />
          )}

          {currentPage === 'portfolio' && (
             <PortfolioPage t={t} portfolios={portfolios} language={language} setSelectedPortfolio={setSelectedPortfolio} />
          )}

          {currentPage === 'testimonials' && (
             <TestimonialsPage t={t} testimonials={testimonials} />
          )}
          
          {currentPage === 'contact' && <ContactPage t={t} />}
          {currentPage === 'companyprofilepage' && <CompanyProfilePage t={t} />}
        </main>

        <Footer />
        
        {selectedProduct && <ProductModal product={selectedProduct} onClose={() => setSelectedProduct(null)} t={t} language={language} />}
        {selectedCategory && <CategoryModal category={selectedCategory} onClose={() => setSelectedCategory(null)} language={language} navigate={navigate}/>}
        {selectedPortfolio && <PortfolioModal portfolio={selectedPortfolio} onClose={() => setSelectedPortfolio(null)} language={language} />}
    </div>
  );
}