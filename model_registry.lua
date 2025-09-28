-- model_registry.lua
-- Smart contract untuk me-registrasi hash dan lokasi model AI.

function init()
    log("Kontrak Model Registry diinisialisasi.")
    -- Tidak perlu inisialisasi state khusus saat ini
end

-- Publikasikan model baru atau versi baru dari model yang sudah ada.
-- @param name Nama unik untuk model (misal: 'risk-assessment-v1')
-- @param model_hash Hash dari file model (misal: Keccak256)
-- @param location Lokasi penyimpanan file model (misal: IPFS CID atau URL)
function publish_model(name, model_hash, location)
    if name == nil or name == "" then
        log("Error: Nama model tidak boleh kosong.")
        return
    end
    if model_hash == nil or model_hash == "" then
        log("Error: Hash model tidak boleh kosong.")
        return
    end

    log("Mempublikasikan model: " .. name)

    -- Gunakan nama model sebagai kunci untuk menyimpan detailnya
    local hash_key = "model:" .. name .. ":hash"
    local location_key = "model:" .. name .. ":location"

    set_value(hash_key, model_hash)
    set_value(location_key, location)

    log("Model " .. name .. " berhasil dipublikasikan.")
end

-- Ambil detail dari sebuah model.
-- @param name Nama model yang akan dicari.
-- @return String berisi hash dan lokasi, dipisahkan koma.
function get_model(name)
    if name == nil or name == "" then
        return ""
    end

    log("Mengambil detail untuk model: " .. name)

    local hash_key = "model:" .. name .. ":hash"
    local location_key = "model:" .. name .. ":location"

    local model_hash = get_value(hash_key)
    local location = get_value(location_key)

    if model_hash == nil then
        return "Model tidak ditemukan."
    end

    return model_hash .. "," .. location
end

-- Panggil init() saat deploy
init()
