import axios from 'axios';
import { Edit, Trash2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useToast } from '../../components/UI/ToastNotification'; // Import useToast
import { API_URL } from '../../config';

const TestimonialManagement = ({ editingItem, setEditingItem }) => {
    const [testimonials, setTestimonials] = useState([]);
    const [portfolios, setPortfolios] = useState([]);
    const [formData, setFormData] = useState({
        client_name: '',
        testimonial_text: '',
        image: null,
        video: null,
        portfolio_id: '',
    });
    const [loading, setLoading] = useState(true);
    const showToast = useToast(); // Initialize useToast

    // Effect to update form data when editingItem prop changes
    useEffect(() => {
        if (editingItem) {
            setFormData({
                client_name: editingItem.ClientName || '',
                testimonial_text: editingItem.TestimonialText || '',
                image: null, // File inputs must be reset
                video: null,
                portfolio_id: editingItem.PortfolioID || '',
            });
        } else {
            setFormData({ client_name: '', testimonial_text: '', image: null, video: null, portfolio_id: '' });
        }
    }, [editingItem]);

    useEffect(() => {
        fetchData();
    }, [editingItem]); // Re-fetch data after save/cancel

    const fetchData = async () => {
        try {
            setLoading(true);
            const token = localStorage.getItem('token');
            const headers = { 'Authorization': `Bearer ${token}` };

            const [testimonialsRes, portfoliosRes] = await Promise.all([
                axios.get(`${API_URL}/testimoni`, { headers }),
                axios.get(`${API_URL}/portfolios`, { headers }) // Portfolios might be public, but for admin context, ensure auth.
            ]);
            setTestimonials(testimonialsRes.data);
            setPortfolios(portfoliosRes.data);
        } catch (err) {
            showToast('Gagal mengambil data: ' + err.message, 'error');
            console.error('Error fetching data:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const handleFileChange = (e) => {
        const { name, files } = e.target;
        setFormData((prev) => ({ ...prev, [name]: files[0] }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const data = new FormData();
        data.append('client_name', formData.client_name);
        data.append('testimonial_text', formData.testimonial_text);
        data.append('portfolio_id', formData.portfolio_id || '0');
        if (formData.image) {
            data.append('image', formData.image);
        }
        if (formData.video) {
            data.append('video', formData.video);
        }

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                showToast('Token tidak ditemukan. Silakan login kembali.', 'error');
                return;
            }

            const url = editingItem 
                ? `${API_URL}/testimoni/${editingItem.ID}`
                : `${API_URL}/testimoni`;
            
            const method = editingItem ? 'PUT' : 'POST';

            // Gunakan fetch API untuk FormData (lebih reliable)
            const res = await fetch(url, {
                method: method,
                headers: {
                    'Authorization': `Bearer ${token}`
                    // Jangan set Content-Type, browser akan set otomatis dengan boundary
                },
                body: data
            });

            if (!res.ok) {
                const errorData = await res.json().catch(() => ({ error: res.statusText }));
                throw new Error(errorData.error || `HTTP ${res.status}`);
            }

            const result = await res.json();
            showToast(editingItem ? 'Testimoni berhasil diperbarui!' : 'Testimoni berhasil ditambahkan!', 'success');
            setEditingItem(null);
            setFormData({ client_name: '', testimonial_text: '', image: null, video: null, portfolio_id: '' });
            fetchData();
        } catch (err) {
            showToast('Gagal menyimpan testimoni: ' + err.message, 'error');
            console.error('Error saving testimoni:', err);
        }
    };

    const handleDelete = async (id) => {
        if (window.confirm('Apakah Anda yakin ingin menghapus testimoni ini?')) {
            try {
                const token = localStorage.getItem('token');
                await axios.delete(`${API_URL}/testimoni/${id}`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                showToast('Testimoni berhasil dihapus!', 'success');
                fetchData();
            } catch (err) {
                showToast('Gagal menghapus testimoni: ' + (err.response?.data?.error || err.message), 'error');
                console.error('Error deleting testimoni:', err.response?.data || err);
            }
        }
    };

    if (loading) return <div className="text-slate-600">Memuat Testimoni...</div>;

    return (
        <div className="grid md:grid-cols-3 gap-6">
            {/* Form Column */}
            <div className="bg-white p-6 rounded-2xl shadow-lg border border-gray-100 h-fit sticky top-6">
                <h3 className="font-bold mb-4 text-2xl text-slate-800">{editingItem ? 'Edit Testimoni' : 'Tambah Testimoni'}</h3>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="client_name" className="block text-sm font-medium text-gray-700 mb-1">Nama Klien / Proyek</label>
                        <input
                            type="text"
                            id="client_name"
                            name="client_name"
                            value={formData.client_name}
                            onChange={handleChange}
                            className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm"
                            required
                        />
                    </div>
                    <div>
                        <label htmlFor="testimonial_text" className="block text-sm font-medium text-gray-700 mb-1">Teks Testimoni (Opsional)</label>
                        <textarea
                            id="testimonial_text"
                            name="testimonial_text"
                            value={formData.testimonial_text}
                            onChange={handleChange}
                            rows="4"
                            className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm"
                        ></textarea>
                    </div>
                    <div>
                        <label htmlFor="portfolio_id" className="block text-sm font-medium text-gray-700 mb-1">Proyek Terkait (Opsional)</label>
                        <select
                            id="portfolio_id"
                            name="portfolio_id"
                            value={formData.portfolio_id}
                            onChange={handleChange}
                            className="w-full border p-3 rounded-lg focus:ring-teal-500 focus:border-teal-500 shadow-sm"
                        >
                            <option value="">Pilih Proyek</option>
                            {portfolios.map((portfolio) => (
                                <option key={portfolio.ID} value={portfolio.ID}>
                                    {portfolio.Title}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label htmlFor="image" className="block text-sm font-medium text-gray-700 mb-1">Gambar</label>
                        <input
                            type="file"
                            id="image"
                            name="image"
                            accept="image/*"
                            onChange={handleFileChange}
                            className="w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100"
                        />
                    </div>
                    <div>
                        <label htmlFor="video" className="block text-sm font-medium text-gray-700 mb-1">Video (Opsional)</label>
                        <input
                            type="file"
                            id="video"
                            name="video"
                            accept="video/*"
                            onChange={handleFileChange}
                            className="w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-teal-50 file:text-teal-700 hover:file:bg-teal-100"
                        />
                    </div>
                    <button
                        type="submit"
                        className="w-full bg-teal-500 text-white py-3 rounded-lg font-bold shadow-md hover:bg-teal-600 transition-colors"
                    >
                        {editingItem ? 'Update Testimoni' : 'Tambah Testimoni'}
                    </button>
                    {editingItem && (
                        <button
                            type="button"
                            onClick={() => setEditingItem(null)}
                            className="w-full bg-gray-200 text-slate-700 py-3 rounded-lg font-bold shadow-sm hover:bg-gray-300 transition-colors mt-2"
                        >
                            Batal Edit
                        </button>
                    )}
                </form>
            </div>

            {/* List Column */}
            <div className="md:col-span-2 space-y-4">
                <h2 className="text-2xl font-bold mb-4 text-slate-800">Daftar Testimoni</h2>
                <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tipe</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Media</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {testimonials.map((item, index) => (
                                <tr key={item.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.ClientName}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                        {item.TestimonialText ? 
                                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">Klien</span> : 
                                          <span className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">Proyek</span>
                                        }
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                                        {item.ImagePath && <img src={item.ImagePath} alt={item.ClientName} className="w-16 h-16 object-cover rounded-md inline-block" />}
                                        {item.VideoPath && <a href={item.VideoPath} target="_blank" rel="noopener noreferrer" className="text-blue-500 ml-2">Lihat Video</a>}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                        <div className="flex items-center justify-end space-x-2">
                                            <button
                                                onClick={() => setEditingItem(item)}
                                                className="text-blue-500 hover:text-blue-700 mr-3"
                                            >
                                                <Edit size={18}/>
                                            </button>
                                            <button
                                                onClick={() => handleDelete(item.ID)}
                                                className="text-red-500 hover:text-red-700"
                                            >
                                                <Trash2 size={18}/>
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default TestimonialManagement;