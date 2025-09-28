-- oracle_registry.lua
-- Smart contract untuk mengelola pendaftaran node inference melalui staking.

function init()
    log("Kontrak Oracle Registry diinisialisasi.")
    set_value("node_count", "0")
end

-- Fungsi ini dipanggil oleh sistem blockchain saat transaksi stake terdeteksi
-- @param node_address Alamat node yang melakukan stake
-- @param amount Jumlah yang di-stake
function register_stake(node_address, amount)
    log("Mencatat stake untuk node: " .. node_address .. " sejumlah: " .. amount)

    local stake_key = "stake:" .. node_address
    local current_stake = tonumber(get_value(stake_key) or "0")
    local new_stake = current_stake + tonumber(amount)

    -- Jika ini adalah pendaftaran pertama, tambah jumlah node
    if current_stake == 0 and new_stake > 0 then
        local count = tonumber(get_value("node_count") or "0") + 1
        set_value("node_count", tostring(count))
        set_value("node_" .. count, node_address)
    end

    -- Perbarui jumlah stake
    set_value(stake_key, tostring(new_stake))
    log("Stake untuk " .. node_address .. " berhasil diperbarui menjadi: " .. new_stake)
end

-- Ambil jumlah stake dari sebuah node
-- @param node_address Alamat node
-- @return Jumlah stake sebagai string
function get_stake(node_address)
    local stake_key = "stake:" .. node_address
    return get_value(stake_key) or "0"
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

    return table.concat(nodes, ",")
end