local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local interval = tonumber(ARGV[2])

local time = redis.call('TIME')
local now = tonumber(time[1])

local tokens = tonumber(redis.call('HGET', key, 'tokens'))
local lastRefill = tonumber(redis.call('HGET', key, 'last_refill'))

if not tokens or not lastRefill then
    tokens = capacity
    lastRefill = now
end

local elapsed = now - lastRefill
local tokensToAdd = math.floor(elapsed / interval)

if tokensToAdd > 0 then
    tokens = math.min(capacity, tokens + tokensToAdd)
    lastRefill = lastRefill + (tokensToAdd * interval)
end

local allowed = 0

if tokens > 0 then
    tokens = tokens - 1
    allowed = 1
end

redis.call('HSET', key, 'tokens', tokens, 'last_refill', lastRefill)

redis.call('EXPIRE', key, capacity * interval)

return allowed
