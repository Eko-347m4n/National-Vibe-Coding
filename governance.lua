-- governance.lua
-- Starter kit contract for digital governance: proposals, voting, and budget tracking.

function init()
    log("Digital Governance Kit Contract Initialized.")
    set_value("proposal_count", "0")
    set_value("total_budget", "10000") -- Anggaran awal (contoh)
    set_value("spent_budget", "0")
end

-- -- PROPOSAL & VOTING FUNCTIONS -- --

-- Create a new proposal for voting.
-- @param title The title of the proposal.
-- @param description A short description of what is being proposed.
-- @return proposal_id The ID of the newly created proposal.
function create_proposal(title, description)
    log("Creating new proposal: " .. title)
    local count = tonumber(get_value("proposal_count"))
    local proposal_id = count + 1

    set_value("proposal:" .. proposal_id .. ":title", title)
    set_value("proposal:" .. proposal_id .. ":desc", description)
    set_value("proposal:" .. proposal_id .. ":votes_for", "0")
    set_value("proposal:" .. proposal_id .. ":votes_against", "0")
    set_value("proposal:" .. proposal_id .. ":status", "VOTING")

    set_value("proposal_count", tostring(proposal_id))
    return tostring(proposal_id)
end

-- Cast a vote on a proposal.
-- @param proposal_id The ID of the proposal to vote on.
-- @param vote The vote, either "FOR" or "AGAINST".
-- @param voter_address The address of the voter (for uniqueness, not yet enforced).
function vote(proposal_id, vote, voter_address)
    local status = get_value("proposal:" .. proposal_id .. ":status")
    if status ~= "VOTING" then
        log("Error: Proposal " .. proposal_id .. " is not open for voting.")
        return
    end

    if vote == "FOR" then
        local current_votes = tonumber(get_value("proposal:" .. proposal_id .. ":votes_for"))
        set_value("proposal:" .. proposal_id .. ":votes_for", tostring(current_votes + 1))
        log("Vote FOR cast on proposal " .. proposal_id .. " by " .. voter_address)
    elseif vote == "AGAINST" then
        local current_votes = tonumber(get_value("proposal:" .. proposal_id .. ":votes_against"))
        set_value("proposal:" .. proposal_id .. ":votes_against", tostring(current_votes + 1))
        log("Vote AGAINST cast on proposal " .. proposal_id .. " by " .. voter_address)
    else
        log("Error: Invalid vote. Must be 'FOR' or 'AGAINST'.")
    end
end

-- Get the details and status of a proposal.
-- @param proposal_id The ID of the proposal.
-- @return A JSON string with the proposal details.
function get_proposal(proposal_id)
    local title = get_value("proposal:" .. proposal_id .. ":title")
    if title == nil then
        return json.encode({ error = "Proposal not found." })
    end

    local desc = get_value("proposal:" .. proposal_id .. ":desc")
    local votes_for = get_value("proposal:" .. proposal_id .. ":votes_for")
    local votes_against = get_value("proposal:" .. proposal_id .. ":votes_against")
    local status = get_value("proposal:" .. proposal_id .. ":status")

    local proposal_data = {
        id = proposal_id,
        title = title,
        description = desc,
        votes_for = tonumber(votes_for),
        votes_against = tonumber(votes_against),
        status = status
    }
    
    return json.encode(proposal_data)
end

-- -- BUDGET & AI ORACLE FUNCTIONS -- --

-- Get the current budget status.
-- @return A JSON string with budget details.
function get_budget()
    local total = tonumber(get_value("total_budget"))
    local spent = tonumber(get_value("spent_budget"))
    local remaining = total - spent

    local budget_data = {
        total = total,
        spent = spent,
        remaining = remaining
    }
    return json.encode(budget_data)
end

-- This is a placeholder showing how the kit would integrate with an AI Oracle.
-- In a real scenario, this would be called by an off-chain daemon monitoring budget transactions.
-- @param transaction_data A description of a transaction to be analyzed.
-- @return A simulated response from an AI model.
function check_anomaly(transaction_data)
    log("Calling AI Oracle to check for anomaly in: " .. transaction_data)
    
    -- In a real implementation, this would involve creating a job on the inference_market contract.
    -- For this demo, we simulate the result.
    local is_anomalous = false
    if string.find(transaction_data, "urgent") then
        is_anomalous = true
    end

    if is_anomalous then
        log("AI Oracle response: ANOMALY DETECTED.")
        return "ANOMALY_DETECTED"
    else
        log("AI Oracle response: Looks normal.")
        return "NORMAL"
    end
end

-- Panggil init() saat deploy
init()
