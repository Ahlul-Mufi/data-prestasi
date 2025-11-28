🚀 Sistem Pelaporan Prestasi Mahasiswa – Backend API
UAS Pemrograman Backend Lanjut – Backend REST API
👤 Identitas Mahasiswa
| Keterangan | Data       |
| ---------- | ---------- |
| **Nama**   | AHLUL MUFI |
| **NIM**    | 434231078  |
| **Kelas**  | TI-C2      |

📌 Deskripsi Project
Aplikasi Backend REST API untuk mengelola pelaporan prestasi mahasiswa, dilengkapi dengan sistem autentikasi, verifikasi dosen wali, dan integrasi database ganda (PostgreSQL + MongoDB).
Fitur Utama:
🔐 Role Based Access Control (RBAC)
🔑 Autentikasi JWT
🗂️ Pelaporan prestasi dinamis (MongoDB)
👨‍🏫 Verifikasi prestasi oleh dosen wali
👥 Manajemen pengguna (admin, mahasiswa, dosen)
📊 Dashboard statistik dasar
📎 Upload lampiran prestasi
📘 Dokumentasi standar SRS (Software Requirement Specification)

🧱 Arsitektur Sistem
🗄️ Database
| Jenis                        | Kegunaan                                     |
| ---------------------------- | -------------------------------------------- |
| **PostgreSQL (Relasional)**  | Users, roles, permissions, metadata prestasi |
| **MongoDB (Non-relasional)** | Detail dinamis prestasi & lampiran           |

🔐 Role & Akses (RBAC)
| Role           | Hak Akses                                        |
| -------------- | ------------------------------------------------ |
| **Admin**      | Mengelola data pengguna, memantau semua prestasi |
| **Mahasiswa**  | Membuat & submit prestasi, upload lampiran       |
| **Dosen Wali** | Melihat prestasi bimbingan, verifikasi/menolak   |

🔄 Alur Sistem
Mahasiswa membuat laporan prestasi
Mahasiswa mengirim (submit) prestasi
Dosen wali melihat daftar prestasi mahasiswa bimbingan
Dosen memverifikasi / menolak prestasi
Admin dapat melihat seluruh histori prestasi

🛠️ Teknologi yang Digunakan
| Teknologi       | Deskripsi                                             |
| --------------- | ----------------------------------------------------- |
| **Golang**      | Bahasa utama backend                                  |
| **Fiber / Gin** | Framework HTTP (pilih salah satu sesuai implementasi) |
| **PostgreSQL**  | Basis data relasional                                 |
| **MongoDB**     | Penyimpanan prestasi dinamis                          |
| **JWT**         | Autentikasi                                           |


