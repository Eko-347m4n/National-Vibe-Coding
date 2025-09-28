-- fl_market.lua
-- Smart contract untuk mengoordinasikan tugas Federated Learning.

function init()
    log("Kontrak Federated Learning Market diinisialisasi.")
    set_value("task_count", "0")
end

-- Buat tugas FL baru
function create_task(model_name, min_participants)
    log("Menerima permintaan untuk membuat tugas FL baru untuk model: " .. model_name)

    local task_id = tonumber(get_value("task_count") or "0") + 1

    set_value("task:" .. task_id .. ":model", model_name)
    set_value("task:" .. task_id .. ":min_participants", min_participants)
    set_value("task:" .. task_id .. ":status", "OPEN")
    set_value("task:" .. task_id .. ":participant_count", "0")
    set_value("task:" .. task_id .. ":submission_count", "0")

    set_value("task_count", tostring(task_id))

    log("Tugas FL baru dibuat dengan ID: " .. task_id)
    return tostring(task_id)
end

-- Bergabung dengan tugas FL yang ada
function join_task(task_id, participant_address, stake_amount)
    -- Verifikasi Stake (P4-T1)
    local required_stake = 10 -- Contoh: butuh minimal 10 koin stake
    if tonumber(stake_amount) < required_stake then
        log("Error: Peserta " .. participant_address .. " tidak memiliki cukup stake. Dibutuhkan: " .. required_stake .. ", dimiliki: " .. stake_amount)
        return
    end

    local status = get_value("task:" .. task_id .. ":status")
    if status ~= "OPEN" then
        log("Error: Tugas " .. task_id .. " tidak lagi terbuka untuk pendaftaran.")
        return
    end

    log("Node " .. participant_address .. " (stake: "..stake_amount..") bergabung dengan tugas " .. task_id)

    local p_count = tonumber(get_value("task:" .. task_id .. ":participant_count") or "0") + 1
    set_value("task:" .. task_id .. ":participant:" .. p_count, participant_address)
    set_value("task:" .. task_id .. ":participant_count", tostring(p_count))

    local min_p = tonumber(get_value("task:" .. task_id .. ":min_participants"))
    if p_count >= min_p then
        set_value("task:" .. task_id .. ":status", "TRAINING")
        log("Jumlah peserta minimum tercapai. Status tugas " .. task_id .. " diubah menjadi TRAINING.")
    end
end

-- Kirimkan hash dari model yang sudah dilatih
function submit_hash(task_id, model_hash, node_address)
    local status = get_value("task:" .. task_id .. ":status")
    if status ~= "TRAINING" then
        log("Error: Tugas " .. task_id .. " tidak dalam fase TRAINING.")
        return
    end

    log("Menerima hash model untuk tugas " .. task_id .. " dari node " .. node_address)

    local s_count = tonumber(get_value("task:" .. task_id .. ":submission_count") or "0") + 1
    set_value("task:" .. task_id .. ":submission:" .. s_count .. ":node", node_address)
    set_value("task:" .. task_id .. ":submission:" .. s_count .. ":hash", model_hash)
    set_value("task:" .. task_id .. ":submission_count", tostring(s_count))

    local p_count = tonumber(get_value("task:" .. task_id .. ":participant_count"))
    if s_count >= p_count then
        set_value("task:" .. task_id .. ":status", "AGGREGATING")
        log("Semua peserta telah mengirimkan hash. Status tugas " .. task_id .. " diubah menjadi AGGREGATING.")
    end
end
