# Memulai dengan Swatantra

Panduan ini akan membantu Anda menginstal, menjalankan, dan berinteraksi dengan node Swatantra untuk pertama kalinya.

## 1. Prasyarat

- Pastikan Anda memiliki **Go (versi 1.18+)** terinstal.

## 2. Kompilasi

Clone repositori dan kompilasi binary utama:

```bash
go build
```

Ini akan menghasilkan sebuah executable bernama `swatantra-node`.

## 3. Menjalankan Jaringan Lokal

Alur kerja paling dasar adalah membuat wallet, membuat blockchain, dan kemudian berinteraksi dengannya.

### Langkah 3.1: Buat Wallet

Setiap interaksi yang mengubah state blockchain memerlukan alamat. Buat wallet pertama Anda:

```bash
./swatantra-node createwallet
# Output: Alamat baru Anda: 1... (SALIN ALAMAT INI)
```

### Langkah 3.2: Buat Blockchain (Genesis Block)

Gunakan alamat yang baru saja Anda buat untuk menerima hadiah dari blok pertama (genesis).

```bash
./swatantra-node createblockchain -address [ALAMAT_ANDA]
```

Blockchain Anda sekarang aktif dan berjalan secara lokal!

## 4. Interaksi Dasar

- **Periksa Saldo:**
  ```bash
  ./swatantra-node getbalance -address [ALAMAT_ANDA]
  ```

- **Kirim Koin:**
  Buat wallet kedua, lalu kirim koin dari wallet pertama.
  ```bash
  ./swatantra-node createwallet # Dapatkan ALAMAT_KEDUA
  ./swatantra-node send -from [ALAMAT_ANDA] -to [ALAMAT_KEDUA] -amount 10
  ```

## 5. Langkah Selanjutnya

Sekarang setelah Anda memiliki jaringan lokal yang berjalan, Anda bisa mulai berinteraksi dengan fitur-fitur yang lebih canggih.

- Lihat **[Referensi CLI](./cli-reference.md)** untuk daftar lengkap perintah.
- Pelajari cara mendeploy dan memanggil smart contract di **[Referensi Smart Contract](./contracts.md)**.
