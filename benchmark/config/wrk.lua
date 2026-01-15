-- wrk Lua Script for Advanced Load Testing
-- Run with: wrk -t12 -c400 -d30s --latency -s wrk.lua http://localhost:3000/

-- Global variables
local counter = 0
local threads = {}

-- Setup function (called once per thread)
function setup(thread)
   thread:set("id", counter)
   table.insert(threads, thread)
   counter = counter + 1
end

-- Initialize request
function init(args)
   local r = {}
   r.headers = {}
   r.headers["Content-Type"] = "application/json"
   r.headers["User-Agent"] = "wrk-benchmark"
   return r
end

-- Generate request
function request()
   -- Randomly select endpoint
   local endpoints = {
      { method = "GET", path = "/" },
      { method = "GET", path = "/json" },
      { method = "POST", path = "/json", body = '{"id":1,"name":"Test","email":"test@example.com"}' },
      { method = "GET", path = "/user/123" },
      { method = "GET", path = "/query?name=test" }
   }
   
   local endpoint = endpoints[math.random(#endpoints)]
   
   if endpoint.method == "POST" then
      return wrk.format(endpoint.method, endpoint.path, nil, endpoint.body)
   else
      return wrk.format(endpoint.method, endpoint.path)
   end
end

-- Response callback
function response(status, headers, body)
   if status ~= 200 then
      print("Error: HTTP " .. status)
   end
end

-- Done callback (called once after all requests)
function done(summary, latency, requests)
   io.write("========================================\n")
   io.write("wrk Load Test Results\n")
   io.write("========================================\n\n")
   
   io.write(string.format("Duration:        %d ms\n", summary.duration / 1000))
   io.write(string.format("Requests:        %d\n", summary.requests))
   io.write(string.format("Bytes:           %d\n", summary.bytes))
   io.write(string.format("Errors:          %d (connect: %d, read: %d, write: %d, timeout: %d)\n",
      summary.errors.connect + summary.errors.read + summary.errors.write + summary.errors.timeout,
      summary.errors.connect, summary.errors.read, summary.errors.write, summary.errors.timeout))
   
   io.write("\nLatency Distribution:\n")
   io.write(string.format("  Min:           %.2f ms\n", latency.min / 1000))
   io.write(string.format("  Mean:          %.2f ms\n", latency.mean / 1000))
   io.write(string.format("  Max:           %.2f ms\n", latency.max / 1000))
   io.write(string.format("  StdDev:        %.2f ms\n", latency.stdev / 1000))
   
   io.write("\nPercentiles:\n")
   for _, p in pairs({ 50, 75, 90, 95, 99, 99.9, 99.99 }) do
      local n = latency:percentile(p)
      io.write(string.format("  p%.2f:         %.2f ms\n", p, n / 1000))
   end
   
   io.write("\nThroughput:\n")
   io.write(string.format("  Requests/sec:  %.2f\n", summary.requests / (summary.duration / 1000000)))
   io.write(string.format("  Transfer/sec:  %.2f KB\n", (summary.bytes / 1024) / (summary.duration / 1000000)))
   
   io.write("\n========================================\n")
end
