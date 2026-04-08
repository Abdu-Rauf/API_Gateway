local tokens = {}

-- load the jwt in lua tables before processing all the requests
function init(args) do
    for line in io.lines("benchmark/benchmark.lua") do
        tokens.insert(tokens,line)
    end
end


function request() do
    -- select a random token from 1-len(tokens)
    local token = tokens[math.random(1,#tokens)]

    return wrk.format("GET", "/", { ["Authorization"] = "Bearer " .. token })
end