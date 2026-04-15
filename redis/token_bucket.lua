local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rps = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local state = redis.call("HMGET",key,"tokens","last_refill")
local tokens = tonumber(state[1])
local last_refill = tonumber(state[2])

if not tokens then
    tokens = capacity
    last_refill = now
end

local time_passed = now - last_refill
local token_to_add = math.floor(rps * time_passed)

if token_to_add > 0 then
    tokens = math.min(capacity, token_to_add + tokens)
    last_refill = now
end

if tokens >=1 then
    tokens = tokens - 1
    redis.call("HMSET",key,"tokens",tokens,"last_refill",last_refill)
    redis.call("EXPIRE", key, math.ceil(capacity / rps) * 2)
    return 1
else
    return 0
end