-- utils.lua
-- Pustaka utilitas untuk smart contract Swatantra

json = {}

-- Meng-encode tabel Lua sederhana (satu tingkat, tanpa nested table/array) menjadi string JSON.
-- @param tbl Tabel yang akan di-encode.
-- @return String JSON.
function json.encode(tbl)
    if tbl == nil then return "null" end

    local parts = {}
    for k, v in pairs(tbl) do
        local key_str = string.format("\"%s\"", tostring(k))
        local val_str

        if type(v) == "string" then
            val_str = string.format("\"%s\"", v)
        elseif type(v) == "boolean" then
            val_str = tostring(v)
        else -- Anggap sebagai angka
            val_str = tostring(v)
        end
        table.insert(parts, key_str .. ":" .. val_str)
    end

    return "{" .. table.concat(parts, ",") .. "}"
end

