# Swatantra

**Swatantra** adalah platform blockchain hibrida (UTXO + Ed25519 + Keccak256 + PoW-EMA) dengan service layer AI (Decentralized Inference Oracle & Federated Learning) dan starter kits untuk digital governance dan supply chain.

## Visi

Menjadi infrastruktur socio-technical untuk pemerintahan digital, pendidikan, dan ekonomi inklusif yang ditenagai oleh AI dan terdesentralisasi.

## Fitur Utama

Platform Swatantra dibangun di atas beberapa lapisan modular:

1.  **Core Blockchain**: Fondasi yang aman dan stabil menggunakan model UTXO, konsensus Proof-of-Work dengan penyesuaian kesulitan dinamis (EMA), dan VM untuk smart contract berbasis Lua.
2.  **AI Service Layer**: Fitur unggulan yang menyediakan layanan AI terdesentralisasi:
    -   **Decentralized Inference Oracle**: Sebuah "pasar" on-chain di mana pengguna dapat meminta eksekusi model AI oleh jaringan node yang terdesentralisasi.
    -   **Federated Learning**: Sebuah protokol untuk melatih model AI secara kolaboratif tanpa memusatkan data pelatihan yang sensitif.
3.  **Application Layer**: "Starter Kits" siap pakai untuk mempercepat pengembangan aplikasi di atas Swatantra, seperti untuk kasus penggunaan tata kelola digital dan transparansi rantai pasok.

## Token (SWT)

Ekosistem Swatantra ditenagai oleh token utilitas `SWT`, yang digunakan untuk:
-   Membayar layanan inferensi di AI Oracle.
-   Memberi hadiah kepada partisipan Federated Learning.
-   Melakukan staking untuk menjadi node berkualitas (inference & aggregator).
-   Membayar biaya layanan di marketplace starter kits.

## Memulai Cepat (Quick Start)

Panduan ini akan membantu Anda menjalankan node Swatantra lokal.

### 1. Prasyarat
- Go (versi 1.18+)

### 2. Kompilasi

```bash
go build
```

Ini akan menghasilkan executable `swatantra-node`.

### 3. Jalankan Jaringan Lokal

```bash
# 1. Buat wallet pertama Anda
./swatantra-node createwallet
# Salin alamat yang muncul (ini adalah ALAMAT_ANDA)

# 2. Buat blockchain baru
./swatantra-node createblockchain -address [ALAMAT_ANDA]

# 3. Periksa saldo Anda
./swatantra-node getbalance -address [ALAMAT_ANDA]
```

## Dokumentasi

Untuk panduan yang lebih mendalam, referensi API, dan tutorial, silakan kunjungi direktori dokumentasi kami:

- **[Dokumentasi Lengkap](./docs/index.md)**