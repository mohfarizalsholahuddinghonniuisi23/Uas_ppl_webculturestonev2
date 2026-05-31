# Buku Panduan Aplikasi Culturstone

Dokumen ini adalah panduan lengkap untuk aplikasi web Culturstone, mencakup informasi bagi pengguna akhir dan panduan teknis untuk developer.

## 1. Pendahuluan

Aplikasi Culturstone adalah sebuah platform web yang berfungsi sebagai profil perusahaan, portofolio, dan katalog produk untuk bisnis yang bergerak di bidang batu alam. Aplikasi ini memiliki dua bagian utama:
-   **Frontend (Situs Publik):** Dapat diakses oleh semua pengunjung untuk melihat informasi perusahaan, produk, portofolio, dan testimoni.
-   **Backend (Panel Admin):** Area terbatas untuk admin mengelola konten yang ditampilkan di situs publik.

### Arsitektur & Teknologi
-   **Backend:**
    -   Bahasa: Go
    -   Database: SQLite (`culturstone.db`)
    -   Fungsi: Menyediakan REST API untuk mengelola data (produk, portofolio, testimoni, dll).
-   **Frontend:**
    -   Framework: React.js (dengan Vite)
    -   Styling: Tailwind CSS
    -   Fungsi: Menampilkan data dari backend dan menyediakan antarmuka bagi pengunjung dan admin.

---

## 2. Panduan Pengguna

### 2.1. Situs Publik

Pengunjung dapat mengakses halaman-halaman berikut:
-   **Home:** Halaman utama yang menampilkan ringkasan tentang perusahaan.
-   **Products:** Galeri produk yang ditawarkan oleh Culturstone.
-   **Portfolio:** Kumpulan proyek atau karya yang pernah dikerjakan.
-   **Company Profile:** Informasi detail mengenai perusahaan.
-   **Testimonials:** Ulasan dan testimoni dari klien.
-   **Contact:** Informasi kontak dan cara menghubungi perusahaan.

### 2.2. Panel Admin

Admin dapat mengelola konten website setelah melakukan login.
1.  **Login:** Akses halaman `/login` untuk masuk menggunakan kredensial admin.
2.  **Admin Panel:** Setelah berhasil login, admin akan diarahkan ke dasbor utama.
3.  **Manajemen Konten:** Dari panel admin, admin dapat melakukan:
    -   **Manajemen Testimoni:** Menambah, mengubah, atau menghapus data testimoni.
    -   **Manajemen Visitor:** Melihat data pengunjung (berdasarkan `VisitorDashboard.jsx`).
    -   (Asumsi) Manajemen Produk, Portofolio, dan Kategori berdasarkan struktur backend.

---

## 3. Panduan Developer

Bagian ini berisi instruksi teknis untuk menjalankan, mengembangkan, dan memelihara aplikasi.

### 3.1. Prasyarat

Pastikan perangkat Anda sudah terinstal:
-   **Go (GoLang):** Versi 1.18 atau lebih baru.
-   **Node.js:** Versi 18.x atau lebih baru.
-   **NPM:** Biasanya terinstal bersama Node.js.

### 3.2. Instalasi & Setup

**1. Clone Repository (jika belum ada):**
```bash
git clone <URL_REPOSITORY>
cd culturstone-web
```

**2. Backend (Go):**
Backend tidak memerlukan instalasi dependensi pihak ketiga yang kompleks (berdasarkan `go.mod`). Cukup jalankan dari direktori utama:
```bash
# Berpindah ke direktori backend
cd backend

# Menjalankan server backend
go run main.go
```
Server backend akan berjalan di port yang ditentukan di `main.go` (misalnya, `localhost:8080`).

**3. Frontend (React):**
Buka terminal baru dan jalankan perintah berikut dari direktori utama proyek:
```bash
# Berpindah ke direktori frontend
cd frontend

# Instal semua dependensi yang dibutuhkan
npm install

# Jalankan aplikasi dalam mode development
npm run dev
```
Aplikasi frontend akan berjalan dan dapat diakses melalui browser di alamat yang ditampilkan di terminal (biasanya `http://localhost:5173`).

### 3.3. Struktur Proyek

-   `backend/`: Berisi semua kode sumber untuk server Go.
    -   `main.go`: Titik masuk utama aplikasi backend.
    -   `controllers/`: Logika untuk menangani request HTTP.
    -   `models/`: Definisi struktur data (struct) dan interaksi database.
    -   `config/`: Konfigurasi, seperti koneksi database.
    -   `uploads/`: Direktori tempat file yang diunggah oleh admin disimpan.
-   `frontend/`: Berisi semua kode sumber untuk aplikasi React.
    -   `src/`: Direktori utama kode frontend.
    -   `pages/`: Komponen yang mewakili setiap halaman (misal: Home, Products).
    -   `components/`: Komponen UI yang dapat digunakan kembali (Navbar, Footer, Modal).
    -   `config.js`: Konfigurasi sisi klien, seperti URL endpoint API.
-   `culturstone.db`: File database SQLite.
-   `uploads/`: (Direktori root) Tampaknya merupakan versi lama atau alternatif dari `backend/uploads/`. Perlu dikonfirmasi mana yang aktif digunakan.

### 3.4. Manajemen File Unggahan

File yang diunggah melalui panel admin disimpan di `backend/uploads/`. Terdapat konvensi penamaan file yang jelas:
-   `cat_*.{jpg|png}`: Gambar untuk kategori produk.
-   `port_*.{jpg|png}`: Gambar untuk item portofolio.
-   `prod_*.{jpg|png}`: Gambar untuk produk.
-   `testimoni_img_*.{jpg|png}`: Gambar untuk testimoni.
-   `testimoni_video_*.mp4`: Video untuk testimoni.

---
Selesai.
