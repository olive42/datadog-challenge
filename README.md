Coding Challenge for Datadog

# Requirements

- Go (compiled and tested with Go 1.5.2)
  - uses github.com/hpcloud/tail apart from std library
- Ruby (tested with ruby 2.1.5 on Debian Jessie)

## Building

go build console.go

# Testing / Running

Short of using an Apache-style load generator, I used a w3c-style log
generator to generate an access_log file. In one terminal, run `ruby
log_generator.rb | tee access_log`. For better results, leave it
running for a few minutes before running the console program, or lower
the threshold.

In another terminal, run ./console -fname access_log -avgThreshold 50

After a few minutes, alerts should start coming in.

console_test.go contains some unittests, especially the one testing
the alert going on and off.

# Improvements

I used my best judgement on log contents and URLs to use; obviously a
real-world case will likely require more specific URLs; in particular,
I am surprised HTTP/1.1 access_logs do not contain the requested
servername in the URL ("GET http://www.domain.com/section/blah"), but
I have not looked too much into it.

There probably needs to be more timestamps in the console output. I
also did not get too creative as far as summary statistics go. I
picked 2-second samples arbitrarily.

In the first 2 minutes, the average hits displayed is not accurate as
the ring needs to accumulate data.

Code-wise, this should be in a package console and a very simple
package main should be in another file main.go; this would allow
better unittesting. (right now, change `package main` to `package
console` at the top of console.go and run go test).

The largest improvement would to use a time-series database or library
to store hits.

# Example output

```
[...]
Total hits: 8268, Avg hits over last 2m: 23.85
Top 1 hit: /tags, 2447 hits
Top 2 hit: /docs, 2418 hits
Top 3 hit: /articles, 2355 hits
Alerts:
High traffic generated an alert - hits = 119.48, triggered at Thu Feb 25 09:55:46 CET 2016

Total hits: 8385, Avg hits over last 2m: 23.87
Top 1 hit: /tags, 2490 hits
Top 2 hit: /docs, 2449 hits
Top 3 hit: /articles, 2384 hits
Alerts:
High traffic generated an alert - hits = 119.48, triggered at Thu Feb 25 09:55:46 CET 2016

Total hits: 8502, Avg hits over last 2m: 23.88
Top 1 hit: /tags, 2524 hits
Top 2 hit: /docs, 2485 hits
Top 3 hit: /articles, 2419 hits
Alerts:
High traffic generated an alert - hits = 119.48, triggered at Thu Feb 25 09:55:46 CET 2016

Total hits: 8622, Avg hits over last 2m: 23.83
Top 1 hit: /tags, 2566 hits
Top 2 hit: /docs, 2516 hits
Top 3 hit: /articles, 2456 hits
Alerts:
High traffic generated an alert - hits = 119.48, triggered at Thu Feb 25 09:55:46 CET 2016
High traffic alert recovered, triggered at Thu Feb 25 09:57:46 CET 2016
[...]
```