-- supplychain.lua
-- Starter kit contract for supply chain asset tracking.

function init()
    log("Supply Chain Kit Contract Initialized.")
    set_value("asset_count", "0")
end

-- Create a new digital asset on the blockchain.
-- @param asset_id A unique identifier for the physical asset (e.g., batch number, serial number).
-- @param asset_type The type of asset (e.g., 'Kopi Arabika Gayo').
-- @param creator The address of the entity creating the asset.
-- @return The internal ID of the asset.
function create_asset(asset_id, asset_type, creator)
    log("Creating new asset: " .. asset_type .. " with ID: " .. asset_id)
    local count = tonumber(get_value("asset_count"))
    local internal_id = count + 1

    set_value("asset:" .. internal_id .. ":external_id", asset_id)
    set_value("asset:" .. internal_id .. ":type", asset_type)
    set_value("asset:" .. internal_id .. ":creator", creator)
    set_value("asset:" .. internal_id .. ":event_count", "0")

    set_value("asset_count", tostring(internal_id))
    
    -- Log the first event: CREATED
    log_event(tostring(internal_id), "CREATED", creator, "Asset registered on blockchain.")

    return tostring(internal_id)
end

-- Log a new event in the asset's lifecycle.
-- @param internal_id The internal ID of the asset.
-- @param event_name The name of the event (e.g., 'PACKAGED', 'SHIPPED').
-- @param actor The address of the entity performing the action.
-- @param details Additional details about the event.
function log_event(internal_id, event_name, actor, details)
    log("Logging event '" .. event_name .. "' for asset " .. internal_id)
    local event_count_str = get_value("asset:" .. internal_id .. ":event_count")
    local event_count = tonumber(event_count_str) + 1

    local event_key_prefix = "asset:" .. internal_id .. ":event:" .. event_count
    set_value(event_key_prefix .. ":name", event_name)
    set_value(event_key_prefix .. ":actor", actor)
    set_value(event_key_prefix .. ":details", details)
    -- Simple timestamp placeholder
    set_value(event_key_prefix .. ":timestamp", "now") 

    set_value("asset:" .. internal_id .. ":event_count", tostring(event_count))
end

-- Get the full history of an asset.
-- @param internal_id The internal ID of the asset.
-- @return A JSON string with the asset's details and its full event history.
function get_asset_history(internal_id)
    log("Fetching history for asset ID: " .. internal_id)
    local asset_type = get_value("asset:" .. internal_id .. ":type")
    if asset_type == nil then
        return "{ \"error\": \"Asset not found.\" }"
    end

    local external_id = get_value("asset:" .. internal_id .. ":external_id")
    local creator = get_value("asset:" .. internal_id .. ":creator")
    local event_count = tonumber(get_value("asset:" .. internal_id .. ":event_count"))

    local json = string.format("{\"internal_id\":\"%s\",\"external_id\":\"%s\",\"type\":\"%s\",\"creator\":\"%s\",\"history\":[",\n        internal_id, external_id, asset_type, creator)

    if event_count > 0 then
        for i = 1, event_count do
            local prefix = "asset:" .. internal_id .. ":event:" .. i
            local name = get_value(prefix .. ":name")
            local actor = get_value(prefix .. ":actor")
            local details = get_value(prefix .. ":details")
            
            json = json .. string.format("{\"event\":\"%s\",\"actor\":\"%s\",\"details\":\"%s\"}", name, actor, details)
            if i < event_count then
                json = json .. ","
            end
        end
    end

    json = json .. "]}"
    return json
end

-- Panggil init() saat deploy
init()
