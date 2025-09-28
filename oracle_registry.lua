-- oracle_registry.lua
-- Smart contract untuk mengelola pendaftaran node inference.

-- Inisialisasi state saat kontrak pertama kali di-deploy
function init()
    log("Kontrak Oracle Registry diinisialisasi.")
    set_value("node_count", "0")
end

-- Daftarkan node baru
-- @param node_address Alamat dari node yang ingin mendaftar
function register(node_address)
    if node_address == nil or node_address == "" then
        log("Error: Alamat node tidak boleh kosong.")
        return
    end

    log("Mendaftarkan node baru: " .. node_address)

    local count_str = get_value("node_count")
    local count = tonumber(count_str)
    local new_count = count + 1

    -- Simpan alamat node dengan key berindeks
    set_value("node_" .. tostring(new_count), node_address)
    -- Perbarui jumlah node
    set_value("node_count", tostring(new_count))

    log("Pendaftaran berhasil. Total node sekarang: " .. tostring(new_count))
end

-- Ambil daftar semua node yang terdaftar
-- @return String berisi daftar alamat node, dipisahkan koma
function get_nodes()
    log("Mengambil daftar node...")
    local count_str = get_value("node_count")
    if count_str == nil then
        return ""
    end

    local count = tonumber(count_str)
    if count == 0 then
        return ""
    end

    local nodes = {}
    for i = 1, count do
        local node_addr = get_value("node_" .. tostring(i))
        table.insert(nodes, node_addr)
    end

    -- Gabungkan alamat menjadi satu string
    return table.concat(nodes, ",")
end