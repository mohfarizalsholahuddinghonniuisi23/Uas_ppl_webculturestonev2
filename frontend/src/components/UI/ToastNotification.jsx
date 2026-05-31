import React, { useState, useEffect, createContext, useContext, useCallback } from 'react';
import { CheckCircle, XCircle, Info, X } from 'lucide-react';

const ToastContext = createContext();

export const ToastProvider = ({ children }) => {
  const [toasts, setToasts] = useState([]);

  const showToast = useCallback((message, type = 'info', duration = 3000) => {
    const id = Date.now() + Math.random();
    setToasts((prevToasts) => [...prevToasts, { id, message, type }]);

    setTimeout(() => {
      setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
    }, duration);
  }, []);

  const removeToast = useCallback((id) => {
    setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={showToast}>
      {children}
      <div className="fixed top-4 right-4 z-[9999] space-y-3">
        {toasts.map((toast) => (
          <Toast key={toast.id} {...toast} onRemove={() => removeToast(toast.id)} />
        ))}
      </div>
    </ToastContext.Provider>
  );
};

export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

const Toast = ({ id, message, type, onRemove }) => {
  const icon = {
    success: <CheckCircle className="text-emerald-500" />,
    error: <XCircle className="text-red-500" />,
    info: <Info className="text-blue-500" />,
  }[type];

  const bgColor = {
    success: 'bg-emerald-50',
    error: 'bg-red-50',
    info: 'bg-blue-50',
  }[type];

  const borderColor = {
    success: 'border-emerald-300',
    error: 'border-red-300',
    info: 'border-blue-300',
  }[type];

  return (
    <div
      className={`relative flex items-center gap-3 p-4 pr-10 rounded-lg shadow-lg border ${bgColor} ${borderColor} transform transition-all duration-300 ease-out animate-slide-in-right`}
      role="alert"
    >
      {icon}
      <p className="text-sm font-medium text-slate-800 flex-1">{message}</p>
      <button onClick={onRemove} className="absolute top-2 right-2 text-slate-500 hover:text-slate-700">
        <X size={16} />
      </button>
    </div>
  );
};

// Add keyframes for slide-in-right animation to index.css if not already present
// @keyframes slide-in-right {
//   from { transform: translateX(100%); opacity: 0; }
//   to { transform: translateX(0); opacity: 1; }
// }
// .animate-slide-in-right {
//   animation: slide-in-right 0.3s ease-out forwards;
// }
