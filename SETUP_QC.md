# 🚀 Panduan Setup Lengkap - Culturestone Web (untuk QC / Tester)

Dokumen ini adalah panduan **langkah demi langkah** untuk menjalankan aplikasi Culturestone secara lokal dari awal, termasuk cara login ke panel admin.

---

## 📋 Prasyarat (Wajib Terinstall)

Pastikan semua tools berikut sudah terinstall di komputer Anda:

| Tool | Versi Minimal | Link Download |
|------|--------------|---------------|
| **Go** | 1.21+ | https://go.dev/dl/ |
| **Node.js** | 18+ | https://nodejs.org/ |
| **MySQL** | 8.0+ | https://dev.mysql.com/downloads/ |
| **Git** | Terbaru | https://git-scm.com/ |

Cek instalasi dengan membuka terminal dan ketik:
```bash
go version
node --version
npm --version
mysql --version
```

---

## 📥 Langkah 1: Clone Repository

```bash
git clone https://github.com/<username>/<nama-repo>.git
cd culturestone-web
```

---

## 🗄️ Langkah 2: Setup Database MySQL

### 2.1. Buka MySQL (bisa via terminal atau MySQL Workbench)

Jika via terminal:
```bash
mysql -u root -p
```
*(tekan Enter jika password kosong)*

### 2.2. Buat Database

```sql
CREATE DATABASE IF NOT EXISTS culturestone CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

---

## ⚙️ Langkah 3: Setup Konfigurasi Backend (.env)

Masuk ke folder backend:
```bash
cd backend
```

Buat file `.env` dengan isi berikut (copy-paste):

```env
# Database MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=culturestone

# JWT Secret
API_SECRET=CULTURESTONE_JWT_SECRET_2024_SUPER_AMAN

# Port Server
PORT=8080

# Base URL
BASE_URL=http://localhost:8080
```

> ⚠️ PENTING: Sesuaikan DB_PASS dengan password MySQL Anda. Jika MySQL tidak ada password, biarkan kosong.

---

## ▶️ Langkah 4: Jalankan Backend

Dari folder `backend/`, jalankan:

```bash
go run main.go
```

Jika berhasil, terminal akan menampilkan:
```
Berhasil konek ke MySQL: root@localhost:3306/culturestone
AutoMigrate selesai — semua tabel siap.
Server berjalan di http://localhost:8080
```

> ⚠️ Jika ada error koneksi database, periksa kembali konfigurasi .env (host, port, user, password).

---

## 🔑 Langkah 5: Buat Akun Admin (WAJIB - Hanya Sekali)

Setelah backend berjalan, buka terminal **baru** dan jalankan perintah berikut untuk membuat akun admin.

**Menggunakan PowerShell (Windows):**

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/api/register" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"username": "admin", "password": "Admin123"}'
```

**Atau menggunakan Postman / Insomnia:**
- Method: `POST`
- URL: `http://localhost:8080/api/register`
- Body (raw JSON):
```json
{
  "username": "admin",
  "password": "Admin123"
}
```

Jika berhasil, respons:
```json
{
  "message": "Admin berhasil didaftarkan",
  "username": "admin"
}
```

Jika muncul `"Registrasi admin tidak diizinkan. Admin sudah terdaftar."` berarti akun admin sudah ada, langsung ke Langkah 6.

> Catatan: Password harus min 8 karakter, ada huruf besar, huruf kecil, dan angka. Contoh: `Admin123`

---

## 💻 Langkah 6: Jalankan Frontend

Buka terminal **baru**, masuk ke folder frontend:

```bash
cd frontend
npm install
npm run dev
```

Frontend akan berjalan di: **http://localhost:5173**

---

## 🔐 Langkah 7: Login ke Panel Admin

1. Buka browser → akses: **http://localhost:5173/login**
2. Masukkan credentials:
   - **Username:** `admin`
   - **Password:** `Admin123`
3. Klik tombol **Login**
4. Anda akan diarahkan ke **Dashboard Admin**

---

## 📁 Dua Terminal yang Harus Aktif

```
Terminal 1 (Backend):
  cd backend
  go run main.go
  → Berjalan di: http://localhost:8080

Terminal 2 (Frontend):
  cd frontend
  npm run dev
  → Berjalan di: http://localhost:5173
```

---

## 🧪 Credentials untuk Testing

| Akun | Username | Password |
|------|----------|----------|
| Admin | `admin` | `Admin123` |

---

## ❓ Troubleshooting

| Error | Solusi |
|-------|--------|
| `Gagal koneksi database MySQL` | Pastikan MySQL berjalan & `.env` sudah dikonfigurasi |
| `port already in use` | Port 8080/5173 sudah dipakai, matikan proses lain |
| `go: command not found` | Install Go dari https://go.dev/dl/ |
| `npm: command not found` | Install Node.js dari https://nodejs.org/ |
| Halaman login blank / error CORS | Pastikan backend sudah berjalan di port 8080 |
| `Invalid credentials` saat login | Ulangi Langkah 5 untuk membuat akun admin |

---

*Jika ada pertanyaan, hubungi developer via grup atau email.*
