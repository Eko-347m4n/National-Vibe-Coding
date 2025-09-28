-- inference_market.lua
-- Smart contract yang berfungsi sebagai papan pekerjaan untuk tugas inferensi AI.

function init()
    log("Kontrak Inference Market diinisialisasi.")
    set_value("job_count", "0")
end

-- Buat permintaan inferensi baru
-- @param model_name Nama model yang terdaftar di model_registry
-- @param input_data Data input untuk model
-- @param reward_amount Jumlah hadiah (placeholder)
-- @return job_id ID dari pekerjaan yang baru dibuat
function request_inference(model_name, input_data, reward_amount)
    log("Menerima permintaan inferensi untuk model: " .. model_name)

    local count_str = get_value("job_count")
    local job_id = tonumber(count_str) + 1

    -- Simpan detail pekerjaan
    set_value("job:" .. job_id .. ":model", model_name)
    set_value("job:" .. job_id .. ":input", input_data)
    set_value("job:" .. job_id .. ":reward", reward_amount)
    set_value("job:" .. job_id .. ":status", "OPEN")
    set_value("job:" .. job_id .. ":response_count", "0")

    -- Perbarui job count
    set_value("job_count", tostring(job_id))

    log("Pekerjaan inferensi baru dibuat dengan ID: " .. job_id)
    return tostring(job_id)
end

-- Kirimkan hasil inferensi untuk sebuah pekerjaan
-- @param job_id ID pekerjaan
-- @param result Hasil inferensi
-- @param node_address Alamat node yang mengirimkan hasil
function submit_response(job_id, result, node_address)
    log("Menerima hasil untuk pekerjaan ID: " .. job_id .. " dari node: " .. node_address)

    local status = get_value("job:" .. job_id .. ":status")
    if status ~= "OPEN" then
        log("Error: Pekerjaan " .. job_id .. " tidak lagi terbuka.")
        return
    end

    local resp_count_str = get_value("job:" .. job_id .. ":response_count")
    local resp_count = tonumber(resp_count_str) + 1

    -- Simpan detail respons
    set_value("job:" .. job_id .. ":response:" .. resp_count .. ":node", node_address)
    set_value("job:" .. job_id .. ":response:" .. resp_count .. ":result", result)
    set_value("job:" .. job_id .. ":response_count", tostring(resp_count))

    log("Hasil untuk pekerjaan " .. job_id .. " berhasil disimpan.")
end

-- Ambil detail lengkap dari sebuah pekerjaan
-- @param job_id ID pekerjaan
-- @return String JSON yang berisi detail pekerjaan dan semua hasilnya
function get_job(job_id)
    log("Mengambil detail untuk pekerjaan ID: " .. job_id)

    local model = get_value("job:" .. job_id .. ":model")
    if model == nil then
        return "{ \"error\": \"Pekerjaan tidak ditemukan.\" }"
    end

    local input = get_value("job:" .. job_id .. ":input")
    local reward = get_value("job:" .. job_id .. ":reward")
    local status = get_value("job:" .. job_id .. ":status")
    local resp_count = tonumber(get_value("job:" .. job_id .. ":response_count"))

    -- Format ke JSON secara manual (Lua tidak punya library JSON built-in)
    local json = "{ "
    json = json .. "\"job_id\": \"" .. job_id .. "\", "
    json = json .. "\"model\": \"" .. model .. "\", "
    json = json .. "\"input\": \"" .. input .. "\", "
    json = json .. "\"reward\": \"" .. reward .. "\", "
    json = json .. "\"status\": \"" .. status .. "\", "
    json = json .. "\"responses\": [ "

    if resp_count > 0 then
        for i = 1, resp_count do
            local node = get_value("job:" .. job_id .. ":response:" .. i .. ":node")
            local result = get_value("job:" .. job_id .. ":response:" .. i .. ":result")
            json = json .. "{ \"node\": \"" .. node .. "\", \"result\": \"" .. result .. "\" }"
            if i < resp_count then
                json = json .. ", "
            end
        end
    end

    json = json .. " ] }"

    return json
end


-- Panggil init() saat deploy
init()
