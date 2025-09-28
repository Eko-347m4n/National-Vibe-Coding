-- counter.lua
-- Kontrak sederhana yang mengelola sebuah penghitung.

-- Fungsi ini hanya dijalankan sekali saat kontrak pertama kali diterbitkan.
function init()
    log("Kontrak Counter sedang diinisialisasi...")
    -- Atur nilai awal untuk 'counter'
    set_value("counter", "0")
    log("State 'counter' diinisialisasi menjadi 0.")
end

-- Fungsi untuk menambah nilai counter
function increment()
    log("Fungsi increment() dipanggil.")
    
    -- Ambil nilai saat ini
    local current_val_str = get_value("counter")
    
    -- Ubah ke angka, tambahkan 1
    local current_val = tonumber(current_val_str)
    local new_val = current_val + 1
    
    -- Simpan nilai baru
    set_value("counter", tostring(new_val))
    
    log("Nilai counter berhasil ditambah menjadi: " .. tostring(new_val))
end
