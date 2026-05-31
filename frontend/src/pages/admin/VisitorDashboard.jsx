import React, { useState, useEffect } from 'react';
import { API_URL } from '../../config';
import { useToast } from '../../components/UI/ToastNotification'; // Import useToast

export default function VisitorDashboard() {
  const [visitors, setVisitors] = useState([]);
  const [loading, setLoading] = useState(true);
  const showToast = useToast(); // Initialize useToast

  useEffect(() => {
    const fetchVisitors = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          showToast("Unauthorized: No token found.", 'error');
          setLoading(false);
          return;
        }

        const res = await fetch(`${API_URL}/admin/visitors`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        });

        if (!res.ok) {
          const errorData = await res.json();
          throw new Error(errorData.error || 'Failed to fetch visitor data');
        }

        const data = await res.json();
        setVisitors(data);
      } catch (err) {
        showToast("Error: " + err.message, 'error');
      } finally {
        setLoading(false);
      }
    };

    fetchVisitors();
  }, []);

  if (loading) {
    return <div className="p-4 text-center text-slate-600">Memuat data pengunjung...</div>;
  }

  return (
    <div className="bg-white p-8 rounded-2xl shadow-lg border border-gray-100">
      <h2 className="text-2xl font-bold mb-4 text-slate-800">Dashboard Pengunjung</h2>
      {visitors.length === 0 ? (
        <p className="text-slate-600">Belum ada data pengunjung.</p>
      ) : (
        <div className="overflow-x-auto bg-white rounded-2xl shadow-lg border border-gray-100">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Jumlah Kunjungan</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {visitors.map((visitor, index) => (
                <tr key={visitor.ID} className={index % 2 === 0 ? 'bg-white' : 'bg-gray-50 hover:bg-gray-100'}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{visitor.Date}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">{visitor.Count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
