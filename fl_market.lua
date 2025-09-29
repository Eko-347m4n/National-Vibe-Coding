-- fl_market.lua
-- Smart contract untuk mengoordinasikan tugas Federated Learning (dengan reward).

function init()
    log("Kontrak Federated Learning Market diinisialisasi.")
    set_value("task_count", "0")
end

-- Buat tugas FL baru dan kunci dana hadiah di dalamnya.
function create_task(model_name, min_participants, reward_fund)
    log("Menerima permintaan untuk membuat tugas FL baru untuk model: " .. model_name)

    local task_id = tonumber(get_value("task_count") or "0") + 1

    set_value("task:" .. task_id .. ":model", model_name)
    set_value("task:" .. task_id .. ":min_participants", min_participants)
    set_value("task:" .. task_id .. ":status", "OPEN")
    set_value("task:" .. task_id .. ":participant_count", "0")
    set_value("task:" .. task_id .. ":submission_count", "0")
    -- Kunci dana hadiah di dalam state kontrak
    set_value("task:" .. task_id .. ":reward_fund", reward_fund)

    set_value("task_count", tostring(task_id))

    log("Tugas FL baru dibuat dengan ID: " .. task_id .. " dengan dana hadiah: " .. reward_fund)
    return tostring(task_id)
end

-- Bergabung dengan tugas FL yang ada
function join_task(task_id, participant_address, stake_amount)
    local required_stake = 10
    if tonumber(stake_amount) < required_stake then
        log("Error: Stake tidak cukup.")
        return
    end

    local status = get_value("task:" .. task_id .. ":status")
    if status ~= "OPEN" then
        log("Error: Tugas tidak terbuka untuk pendaftaran.")
        return
    end

    log("Node " .. participant_address .. " bergabung dengan tugas " .. task_id)

    local p_count = tonumber(get_value("task:" .. task_id .. ":participant_count") or "0") + 1
    set_value("task:" .. task_id .. ":participant:" .. p_count, participant_address)
    set_value("task:" .. task_id .. ":participant_count", tostring(p_count))

    local min_p = tonumber(get_value("task:" .. task_id .. ":min_participants"))
    if p_count >= min_p then
        set_value("task:" .. task_id .. ":status", "TRAINING")
        log("Jumlah peserta minimum tercapai. Status tugas diubah menjadi TRAINING.")
    end
end

-- Kirimkan hash dari model yang sudah dilatih
function submit_hash(task_id, model_hash, node_address)
    local status = get_value("task:" .. task_id .. ":status")
    if status ~= "TRAINING" then
        log("Error: Tugas tidak dalam fase TRAINING.")
        return
    end

    log("Menerima hash terenkripsi untuk tugas " .. task_id .. " dari " .. node_address)

    local s_count = tonumber(get_value("task:" .. task_id .. ":submission_count") or "0") + 1
    set_value("task:" .. task_id .. ":submission:" .. s_count .. ":node", node_address)
    set_value("task:" .. task_id .. ":submission:" .. s_count .. ":hash", model_hash)
    set_value("task:" .. task_id .. ":submission_count", tostring(s_count))

    local p_count = tonumber(get_value("task:" .. task_id .. ":participant_count"))
    if s_count >= p_count then
        set_value("task:" .. task_id .. ":status", "AGGREGATING")
        log("Semua peserta telah mengirimkan hash. Status diubah menjadi AGGREGATING.")
    end
end

-- Finalisasi tugas, catat hasil agregasi, dan distribusikan hadiah.
function finalize_task(task_id, aggregated_hash)
    local status = get_value("task:" .. task_id .. ":status")
    if status ~= "AGGREGATING" then
        log("Error: Tugas belum siap untuk diagregasi.")
        return
    end

    log("Finalisasi tugas " .. task_id .. ". Hash hasil agregasi: " .. aggregated_hash)
    set_value("task:" .. task_id .. ":aggregated_hash", aggregated_hash)
    set_value("task:" .. task_id .. ":status", "COMPLETED")

    -- Logika Distribusi Hadiah
    local total_reward = tonumber(get_value("task:" .. task_id .. ":reward_fund"))
    local submission_count = tonumber(get_value("task:" .. task_id .. ":submission_count"))
    if submission_count == 0 then return end

    local reward_per_node = total_reward / submission_count
    log("Total hadiah: " .. total_reward .. ". Dibagikan kepada " .. submission_count .. " peserta.")
    log("Hadiah per peserta: " .. reward_per_node)

    for i = 1, submission_count do
        local node_address = get_value("task:" .. task_id .. ":submission:" .. i .. ":node")
        -- Di dunia nyata, ini akan membuat transaksi keluar.
        -- Di sini, kita hanya mencatatnya sebagai event.
        log("EVENT:DISTRIBUTE_REWARD: to=" .. node_address .. ", amount=" .. reward_per_node)
    end

    log("Tugas " .. task_id .. " selesai dan hadiah telah didistribusikan (simulasi).")
end
