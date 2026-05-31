import React, { useState, useEffect } from 'react';
import { Trash2, Edit, Reply, ShoppingCart, DollarSign, Package, Users, MessageSquare } from 'lucide-react';
import Sidebar from './Sidebar';
import VisitorDashboard from './VisitorDashboard';
import TestimonialManagement from './TestimonialManagement';
import { API_URL } from '../../config';

const dummyStats = {
  totalSales: 'Rp 12.5 Jt',
  newOrders: 15,
  totalProducts: 245,
  totalCustomers: 120,
};

const dummyRecentOrders = [
  { id: 1, customer: 'John Doe', item: 'Marmer Carrara (5m²)', status: 'Pending', date: '2023-10-26' },
  { id: 2, customer: 'Jane Smith', item: 'Granit Hitam (3m²)', status: 'Completed', date: '2023-10-25' },
  { id: 3, customer: 'Peter Jones', item: 'Kerajinan Meja Batu', status: 'Cancelled', date: '2023-10-24' },
  { id: 4, customer: 'Alice Brown', item: 'Patung Buddha Marmer', status: 'Pending', date: '2023-10-23' },
  { id: 5, customer: 'Bob White', item: 'Lantai Teraso (10m²)', status: 'Completed', date: '2023-10-22' },
];

export default function AdminPanel({ onLogout }) {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [stats, setStats] = useState(null);
  const [products, setProducts] = useState([]);
  const [portfolios, setPortfolios] = useState([]);
  const [categories, setCategories] = useState([]);
  const [messages, setMessages] = useState([]);
  const [editingItem, setEditingItem] = useState(null);

  useEffect(() => {
      if (activeTab === 'dashboard') fetchStats();
      if (activeTab === 'products') { fetchProducts(); fetchCategories(); }
      if (activeTab === 'portfolios') fetchPortfolios();
      if (activeTab === 'categories') fetchCategories();
      if (activeTab === 'messages') fetchMessages();
      if (activeTab === 'testimonials') {}
  }, [activeTab]);

  const fetchStats = async () => {
    try {
      const token = localStorage.getItem('token');
      const headers = { 'Authorization': `Bearer ${token}` };
      const res = await fetch(`${API_URL}/admin/stats`, { headers });
      const data = await res.json();
      if (data.error) {
        setStats(null);
      } else {
        setStats(data);
      }
    } catch (e) { console.error(e); }
  };

  const fetchProducts = async () => { try { const res = await fetch(`${API_URL}/products`); setProducts(await res.json()); } catch (e) { console.error(e); } };
  const fetchPortfolios = async () => { try { const res = await fetch(`${API_URL}/portfolios`); setPortfolios(await res.json()); } catch (e) { console.error(e); } };
  const fetchCategories = async () => { try { const res = await fetch(`${API_URL}/categories`); setCategories(await res.json()); } catch (e) { console.error(e); } };
  
  const fetchMessages = async () => {
    try {
      const token = localStorage.getItem('token');
      const headers = { 'Authorization': `Bearer ${token}` };
      const res = await fetch(`${API_URL}/admin/messages`, { headers });
      const data = await res.json();
      if (data.error) {
        setMessages([]);
      } else {
        setMessages(data);
      }
    } catch (e) { console.error(e); }
  };

  const handleDelete = async (endpoint, id, refreshFunc) => { 
    if(confirm("Hapus data ini?")) { 
      try {
        const token = localStorage.getItem('token');
        await fetch(`${API_URL}/admin/${endpoint}/${id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } }); 
        refreshFunc(); 
      } catch (e) { console.error(e); }
    }
  };
  
  const handleProductSubmit = async (e) => {
    e.preventDefault(); 
    const formData = new FormData(e.target);
    const token = localStorage.getItem('token');
    const url = editingItem ? `${API_URL}/admin/products/${editingItem.ID}` : `${API_URL}/admin/products`;
    const method = editingItem ? 'PUT' : 'POST';
    try {
      await fetch(url, { method, body: formData, headers: { 'Authorization': `Bearer ${token}` } }); 
      setEditingItem(null); 
      e.target.reset(); 
      fetchProducts();
    } catch (e) { console.error(e); }
  };

  const handleCategorySubmit = async (e) => {
    e.preventDefault(); 
    const formData = new FormData(e.target);
    const token = localStorage.getItem('token');
    const url = editingItem ? `${API_URL}/admin/categories/${editingItem.ID}` : `${API_URL}/admin/categories`;
    const method = editingItem ? 'PUT' : 'POST';
    try {
      await fetch(url, { method, body: formData, headers: { 'Authorization': `Bearer ${token}` } }); 
      setEditingItem(null); 
      e.target.reset(); 
      fetchCategories();
    } catch (e) { console.error(e); }
  };

  const handlePortfolioSubmit = async (e) => {
    e.preventDefault(); 
    const formData = new FormData(e.target);
    const token = localStorage.getItem('token');
    const url = editingItem ? `${API_URL}/admin/portfolios/${editingItem.ID}` : `${API_URL}/admin/portfolios`;
    const method = editingItem ? 'PUT' : 'POST';
    try {
      await fetch(url, { method, body: formData, headers: { 'Authorization': `Bearer ${token}` } }); 
      setEditingItem(null); 
      e.target.reset(); 
      fetchPortfolios();
    } catch (e) { console.error(e); }
  };

  const handleReplyMessage = (msg) => { window.location.href = `mailto:${msg.Email}?subject=Balasan Culturstone&body=Halo ${msg.Name}, menanggapi pesan Anda...`; };

  return (
    <div className="flex min-h-screen bg-slate-50">
        <Sidebar activeTab={activeTab} setActiveTab={(t) => {setActiveTab(t); setEditingItem(null);}} onLogout={onLogout} />
        
        <div className="flex-1 flex flex-col">
            <header className="sticky top-0 z-40 flex items-center justify-between p-6 bg-white/70 backdrop-blur-lg shadow-sm border-b border-gray-100">
                <h1 className="text-3xl font-bold text-slate-800">{activeTab === 'dashboard' ? "Dashboard Admin" : activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}</h1>
                <div className="flex items-center space-x-4">
                    <span className="text-slate-600">Admin User</span>
                    <button onClick={onLogout} className="px-4 py-2 bg-red-500 text-white rounded-lg shadow-md hover:bg-red-600 transition-colors">Logout</button>
                </div>
            </header>

            <main className="flex-1 p-8 overflow-y-auto">
                {activeTab === 'dashboard' && (
                    <>
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
                            <div className="bg-white p-6 rounded-2xl shadow-lg border border-gray-100 flex items-center space-x-4">
                                <div className="p-3 bg-blue-100 rounded-full">
                                    <ShoppingCart className="text-blue-600" size={24} />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Produk</p>
                                    <h2 className="text-2xl font-bold text-slate-800">{stats?.products || 'N/A'}</h2>
                                    <p className="text-xs text-emerald-600">+5% bulan lalu</p>
                                </div>
                            </div>
                            <div className="bg-white p-6 rounded-2xl shadow-lg border border-gray-100 flex items-center space-x-4">
                                <div className="p-3 bg-yellow-100 rounded-full">
                                    <Package className="text-yellow-600" size={24} />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Portofolio</p>
                                    <h2 className="text-2xl font-bold text-slate-800">{stats?.portfolios || 'N/A'}</h2>
                                    <p className="text-xs text-emerald-600">+8% bulan lalu</p>
                                </div>
                            </div>
                            <div className="bg-white p-6 rounded-2xl shadow-lg border border-gray-100 flex items-center space-x-4">
                                <div className="p-3 bg-emerald-100 rounded-full">
                                    <MessageSquare className="text-emerald-600" size={24} />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Pesan Masuk</p>
                                    <h2 className="text-2xl font-bold text-slate-800">{stats?.messages || 'N/A'}</h2>
                                    <p className="text-xs text-red-600">-3% minggu ini</p>
                                </div>
                            </div>
                            <div className="bg-white p-6 rounded-2xl shadow-lg border border-gray-100 flex items-center space-x-4">
                                <div className="p-3 bg-purple-100 rounded-full">
                                    <Users className="text-purple-600" size={24} />
                                </div>
                                <div>
                                    <p className="text-sm text-gray-500">Testimoni</p>
                                    <h2 className="text-2xl font-bold text-slate-800">{stats?.testimoni || 'N/A'}</h2>
                                    <p className="text-xs text-emerald-600">+12% bulan lalu</p>
                                </div>
                            </div>
                        </div>

                        <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                            <h2 className="text-2xl font-bold text-slate-800 mb-6">Pesan Terbaru</h2>
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                                        <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                                        <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Pesan</th>
                                        <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
                                        <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                    {Array.isArray(messages) && messages.length > 0 ? (
                                        messages.slice(0, 5).map((msg, index) => (
                                            <tr key={msg.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{msg.Name}</td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{msg.Email}</td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 truncate max-w-xs">{msg.Message}</td>
                                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{new Date(msg.CreatedAt).toLocaleDateString()}</td>
                                                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                    <button onClick={() => handleReplyMessage(msg)} className="text-blue-500 hover:text-blue-700 mr-3"><Reply size={18}/></button>
                                                    <button onClick={() => handleDelete('messages', msg.ID, fetchMessages)} className="text-red-500 hover:text-red-700"><Trash2 size={18}/></button>
                                                </td>
                                            </tr>
                                        ))
                                    ) : (
                                        <tr>
                                            <td colSpan="5" className="px-6 py-4 text-center text-gray-500">Tidak ada pesan terbaru.</td>
                                        </tr>
                                    )}
                                </tbody>
                            </table>
                        </div>
                    </>
                )}

                {activeTab === 'categories' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <h2 className="text-2xl font-bold mb-6 text-slate-800">Manajemen Kategori</h2>
                        <div className="grid md:grid-cols-3 gap-6">
                            <div className="bg-white p-6 rounded-xl shadow-sm h-fit sticky top-6">
                                <h3 className="font-bold mb-4 text-slate-800">{editingItem ? "Edit" : "Tambah"} Kategori</h3>
                                <form onSubmit={handleCategorySubmit} className="space-y-3">
                                    <input name="name" placeholder="Nama" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Name} required/>
                                    <textarea name="description" placeholder="Deskripsi" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Description}/>
                                    <input type="file" name="image" className="text-sm w-full file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100"/>
                                    <button className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors">Simpan</button>
                                </form>
                            </div>
                            <div className="md:col-span-2 space-y-4">
                                <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
                                    <table className="min-w-full divide-y divide-gray-200">
                                        <thead className="bg-gray-50">
                                            <tr>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Gambar</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Deskripsi</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody className="bg-white divide-y divide-gray-200">
                                            {categories.map((c, index) => (
                                                <tr key={c.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                                    <td className="px-6 py-4 whitespace-nowrap"><img src={c.ImageURL} className="w-12 h-12 object-cover rounded-md"/></td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{c.Name}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 truncate max-w-xs">{c.Description}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                        <button onClick={()=>setEditingItem(c)} className="text-blue-500 hover:text-blue-700 mr-3"><Edit size={18}/></button>
                                                        <button onClick={()=>handleDelete('categories', c.ID, fetchCategories)} className="text-red-500 hover:text-red-700"><Trash2 size={18}/></button>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {activeTab === 'products' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <h2 className="text-2xl font-bold mb-6 text-slate-800">Manajemen Produk</h2>
                        <div className="grid md:grid-cols-3 gap-6">
                            <div className="bg-white p-6 rounded-xl shadow-sm h-fit sticky top-6">
                                <h3 className="font-bold mb-4 text-slate-800">{editingItem ? "Edit" : "Tambah"} Produk</h3>
                                <form onSubmit={handleProductSubmit} className="space-y-3">
                                    <input name="name" placeholder="Nama Produk" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Name} required/>
                                    <select name="category" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Category} required>
                                        <option value="">Pilih Kategori</option>
                                        {categories.map(c => <option key={c.ID} value={c.Name}>{c.Name}</option>)}
                                    </select>
                                    <input name="size" placeholder="Ukuran" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Size}/>
                                    <input name="quality" placeholder="Kualitas" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Quality}/>
                                    <input name="finishing" placeholder="Finishing" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Finishing}/>
                                    <textarea name="description" placeholder="Deskripsi" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Description}/>
                                    <input type="file" name="image" className="text-sm w-full file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100"/>
                                    <button className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors">Simpan</button>
                                </form>
                            </div>
                            <div className="md:col-span-2 space-y-4">
                                <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
                                    <table className="min-w-full divide-y divide-gray-200">
                                        <thead className="bg-gray-50">
                                            <tr>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Gambar</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kategori</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Deskripsi</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Kualitas</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ukuran</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Finishing</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody className="bg-white divide-y divide-gray-200">
                                            {products.map((p, index) => (
                                                <tr key={p.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                                    <td className="px-6 py-4 whitespace-nowrap"><img src={p.ImageURL} className="w-12 h-12 object-cover rounded-md"/></td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{p.Name}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{p.Category}</td>
                                                    <td className="px-6 py-4 text-sm text-gray-600 truncate max-w-xs">{p.Description}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{p.Quality}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{p.Size}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{p.Finishing}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                        <div className="flex items-center justify-end space-x-2">
                                                            <button onClick={()=>setEditingItem(p)} className="text-blue-500 hover:text-blue-700 mr-3"><Edit size={18}/></button>
                                                            <button onClick={()=>handleDelete('products', p.ID, fetchProducts)} className="text-red-500 hover:text-red-700"><Trash2 size={18}/></button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
                
                {activeTab === 'portfolios' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <h2 className="text-2xl font-bold mb-6 text-slate-800">Manajemen Portofolio</h2>
                        <div className="grid md:grid-cols-3 gap-6">
                            <div className="bg-white p-6 rounded-xl shadow-sm h-fit sticky top-6">
                                <h3 className="font-bold mb-4 text-slate-800">{editingItem ? "Edit" : "Tambah"} Portofolio</h3>
                                <form onSubmit={handlePortfolioSubmit} className="space-y-3">
                                    <input name="title" placeholder="Judul" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Title} required/>
                                    <textarea name="description" placeholder="Deskripsi" className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm" defaultValue={editingItem?.Description}/>
                                    <input type="file" name="image" className="text-sm w-full file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100"/>
                                    <button className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors">Simpan</button>
                                </form>
                            </div>
                            <div className="md:col-span-2 space-y-4">
                                <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
                                    <table className="min-w-full divide-y divide-gray-200">
                                        <thead className="bg-gray-50">
                                            <tr>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Gambar</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Judul</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Deskripsi</th>
                                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                            </tr>
                                        </thead>
                                        <tbody className="bg-white divide-y divide-gray-200">
                                            {portfolios.map((p, index) => (
                                                <tr key={p.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                                    <td className="px-6 py-4 whitespace-nowrap"><img src={p.ImageURL} className="w-12 h-12 object-cover rounded-md"/></td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{p.Title}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 truncate max-w-xs">{p.Description}</td>
                                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                        <div className="flex items-center justify-end space-x-2">
                                                            <button onClick={()=>setEditingItem(p)} className="text-blue-500 hover:text-blue-700 mr-3"><Edit size={18}/></button>
                                                            <button onClick={()=>handleDelete('portfolios', p.ID, fetchPortfolios)} className="text-red-500 hover:text-red-700"><Trash2 size={18}/></button>
                                                        </div>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            
                {activeTab === 'testimonials' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <TestimonialManagement editingItem={editingItem} setEditingItem={setEditingItem} />
                    </div>
                )}
                
                {activeTab === 'messages' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <h2 className="text-2xl font-bold text-slate-800 mb-6">Pesan Masuk</h2>
                        <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead className="bg-gray-50">
                                    <tr>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Pesan</th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                    {messages.map((m, index) => (
                                        <tr key={m.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{m.Name}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{m.Email}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 truncate max-w-xs">{m.Message}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{new Date(m.CreatedAt).toLocaleDateString()}</td>
                                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                <button onClick={() => handleReplyMessage(m)} className="text-blue-500 hover:text-blue-700 mr-3"><Reply size={18}/></button>
                                                <button onClick={() => handleDelete('messages', m.ID, fetchMessages)} className="text-red-500 hover:text-red-700"><Trash2 size={18}/></button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                )}

                {activeTab === 'visitors' && (
                    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
                        <VisitorDashboard />
                    </div>
                )}
            </main>
        </div>
    </div>
  );
}