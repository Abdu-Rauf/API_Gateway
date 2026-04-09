local tokens = {}

-- load the jwt in lua tables before processing all the requests
function init(args) 
    for line in io.lines("auth/jwt.txt") do
        table.insert(tokens,line)
    end
end


function request() 
    -- select a random token from 1-len(tokens)
    local token = tokens[math.random(1,#tokens)]

    return wrk.format("GET", "/", { ["Authorization"] = "Bearer " .. token })
end