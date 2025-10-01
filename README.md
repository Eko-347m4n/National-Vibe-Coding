# Proyek Blockchain Swatantra

Selamat datang di Swatantra, sebuah platform blockchain yang dirancang untuk mendukung aplikasi terdesentralisasi dengan fokus pada tata kelola digital dan transparansi rantai pasok.

## Ringkasan Proyek

Proyek ini terdiri dari beberapa komponen utama:

1. **Swatantra Node**: Inti dari blockchain yang ditulis dalam Go. Komponen ini menangani semua logika fundamental seperti pembuatan blok, konsensus, manajemen transaksi, dan jaringan peer-to-peer.
2. **Smart Contracts**: Ditulis dalam Lua, kontrak-kontrak ini mendefinisikan logika bisnis untuk aplikasi terdesentralisasi (dApps). Saat ini, kami menyediakan dua contoh utama:
   * `governance.lua`: Untuk manajemen proposal, pemungutan suara, dan anggaran.
   * `supplychain.lua`: Untuk melacak aset dan mencatat peristiwa dalam rantai pasok.
3. **Starter Kits**: Dasbor web interaktif yang mempermudah interaksi dengan smart contract tanpa perlu menggunakan baris perintah secara ekstensif. Terdapat dua kit:
   * **Digital Governance Kit**: Antarmuka untuk `governance.lua`.
   * **Supply Chain Kit**: Antarmuka untuk `supplychain.lua`.

## Cara Memulai

### Prasyarat

- Go (versi 1.18 atau lebih baru)
- Buka `index.html` di browser Anda.

### Menjalankan Node

1. **Kompilasi Proyek**:

   ```bash
   go build -o swatantra-node main.go
   ```
2. **Buat Wallet Baru**:

   ```bash
   ./swatantra-node createwallet
   ```
3. **Mulai Node Blockchain**:
   Mulai node pertama Anda yang akan bertindak sebagai penambang awal.

   ```bash
   ./swatantra-node startnode -miner [ALAMAT_WALLET_ANDA]
   ```

### Menggunakan Starter Kits

1. Buka file `index.html` di root proyek pada browser Anda.
2. Pilih salah satu kit (Governance atau Supply Chain).
3. Ikuti instruksi di dalam dasbor:
   * **Deploy Kontrak**: Gunakan perintah `deploycontract` yang disediakan untuk men-deploy file `.lua` yang relevan.
   * **Interaksi**: Salin alamat kontrak yang dihasilkan ke dalam input yang tersedia dan gunakan tombol-tombol interaktif untuk memanggil fungsi-fungsi pada smart contract.

## Struktur Direktori

```
/
├── index.html             # Halaman utama
├── style.css              # Styling untuk halaman utama
├── README.md              # Anda sedang membaca ini
├── main.go                # Titik masuk utama aplikasi Go
├── go.mod, go.sum         # Dependensi Go
├── *.lua                  # Contoh smart contract
├── src/
│   ├── cli/                 # Logika untuk antarmuka baris perintah (CLI)
│   ├── core/                # Logika inti blockchain (blok, transaksi, PoW)
│   ├── crypto/              # Fungsi-fungsi kriptografi
│   ├── p2p/                 # Logika jaringan peer-to-peer
│   ├── vm/                  # Mesin virtual untuk eksekusi smart contract Lua
│   └── wallet/              # Manajemen wallet
├── governance-kit/        # Dasbor web untuk tata kelola
│   ├── index.html
│   ├── style.css
│   └── app.js
├── supplychain-kit/       # Dasbor web untuk rantai pasok
│   ├── index.html
│   ├── style.css
│   └── app.js
└── docs/                  # Dokumentasi proyek
```

## Kontribusi

Kami menyambut kontribusi dari siapa saja. Silakan buat *pull request* atau buka *issue* jika Anda menemukan masalah atau memiliki ide untuk perbaikan.
