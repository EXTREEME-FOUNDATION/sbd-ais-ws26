# your Decision paragraph (Part 1)
Since the payment path is not owned, I would obv. choose 4, since the aforementioned problems of 1 (polling inconsistency) would be relatively bad for a fraud check.
The main concern Imo, is that there could be liability problems if fraud happens that is not detected.
There is an option to add a binary "acknowleged" column to every transaction, but that would still leave the second problem of querry load.

4 (Log based) would be best in this scenario. There could, however also be problems, regarding application access and data safety (e.g.: are transactions encrypted? At rest or also in production? Security implications of giving a application read access to everything / certain views?)
Also: if Debezium doesn't read PostgreSQL's WAL, then it will fill up (as mentioned in the table).


If we owned the application then an outbox would probably be the best idea. It enures that all data is either there on both the Storage table and in Kafka or none of them.
However, This additional writing and reading might be a problem (e.g.: a bank with millions of transactions a second might have problems with speed / overhead / duplicated data). I don't know how well this would work with a sharded database, but that might be a solution.






# for each of the three runs: prediction, velocity lines of the scoreboard (R3, R4), and why, in 2-3 sentences of your own (? that is not a english sentence)

## RUN 1 (Default)

### Prediction:
The anomaly model will pass all tests, since it is only one instance anyway. (e.g.: doesn't matter whats configured, as long as messages arrive in a timely manner)
Velocity's model will not have any false alerts / duplicates, since the config doesn't assign any partitions -> only one detector can run (because there is only one partition)
Latency will probably fail, since the one partition cannot keep up (lag grows)
frauds caught will be 100% if waited for. (If the program is terminated early, then the consumers won't have time to work through the backlog.

### Output
Run 1795020234: started 17:43:54, 150,298 transactions in 150 s (1,001/s), 16 frauds so far (still running)

Consumer groups on the CDC topic(s)   [R2]
  group                    topic                              partitions  members       lag
  anomaly-agent            fraud.public.transactions                   1        1     4,944
  velocity-agent           fraud.public.transactions                   1        3    79,078

Agent instances (telemetry of this run)   [R3, R4]
  instance         pid model      transactions  cardholders  out-of-order  lag (s)      seen
  ano-1         470807 anomaly         148,997       77,591             0      0.8       2s ago
  vel-1         470808 velocity         65,726       48,222             0     74.9       5s ago

Alerts on 'fraud-alerts' (committed, this run)

  velocity model (card testing)
    R3  frauds caught          3/7  (43%)   FAIL  (>= 90 %)
    R3  false alerts           0   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  58.1 s / 76.4 s   FAIL  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)

  anomaly model (account takeover)
    R3  frauds caught          8/8  (100%)   PASS  (>= 90 %)
    R3  false alerts           0   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  0.7 s / 0.9 s   PASS  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)

### Comment:
Nailed it.


## RUN 2 (6 partitions)

### Prediction
(anomaly, same as before; nothing changed)

Since there are now 6 partitions the timing should be much better (e.g.: no lag)
However, there might be a problem with how the logs are distributed
e.g.: normally kafka balances the logs to each partition (thats kind of the entire point of them)
so depending on if kafka knows on what merit they should be moved to certain partitions, they will be spread out equally, which will make the detector not work.


### Output
Run 1795025000: started 19:03:20, 110,264 transactions in 110 s (1,001/s), 11 frauds so far (still running)

Consumer groups on the CDC topic(s)   [R2]
  group                    topic                              partitions  members       lag
  anomaly-agent            fraud.public.transactions                   6        1     5,104
  velocity-agent           fraud.public.transactions                   6        3     5,104

Agent instances (telemetry of this run)   [R3, R4]
  instance         pid model      transactions  cardholders  out-of-order  lag (s)      seen
  ano-1         816445 anomaly         108,063       65,925             0      0.9       3s ago
  vel-1         816446 velocity         35,277       29,717             0      0.9       3s ago
  vel-2         816447 velocity         35,504       29,876             1      0.9       3s ago
  vel-3         816448 velocity         35,584       29,973             0      0.9       3s ago

Alerts on 'fraud-alerts' (committed, this run)

  velocity model (card testing)
    R3  frauds caught          0/5  (0%)   FAIL  (>= 90 %)
    R3  false alerts           0   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  0.9 s / 0.9 s   PASS  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)

  anomaly model (account takeover)
    R3  frauds caught          5/5  (100%)   PASS  (>= 90 %)
    R3  false alerts           0   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  0.5 s / 0.8 s   PASS  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)


### Comment:
Nailed it again.


## RUN 3

### Prediction:
(anomaly, same as every time)

Since this run adresses the exact concern that I had previously, I expect everything to pass.


### Output
Run 1795025403: started 19:10:03, 135,285 transactions in 135 s (1,001/s), 14 frauds so far (still running)

Consumer groups on the CDC topic(s)   [R2]
  group                    topic                              partitions  members       lag
  anomaly-agent            fraud.public.transactions                   6        1     1,001
  velocity-agent           fraud.public.transactions                   6        3     1,016

Agent instances (telemetry of this run)   [R3, R4]
  instance         pid model      transactions  cardholders  out-of-order  lag (s)      seen
  ano-1         846015 anomaly         132,984       73,421             0      0.9       4s ago
  vel-1         846016 velocity         42,150       23,913             0      0.9       5s ago
  vel-2         846017 velocity         42,325       23,936             0      0.9       4s ago
  vel-3         846018 velocity         41,821       23,745             0      0.9       4s ago

Alerts on 'fraud-alerts' (committed, this run)

  velocity model (card testing)
    R3  frauds caught          6/6  (100%)   PASS  (>= 90 %)
    R3  false alerts           0   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  0.7 s / 1.0 s   PASS  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)

  anomaly model (account takeover)
    R3  frauds caught          7/7  (100%)   PASS  (>= 90 %)
    R3  false alerts           1   PASS  (at most 10 % of the flagged transactions)
    R4  alert latency p50/p95  0.5 s / 0.9 s   PASS  (p95 <= 5 s after the commit)
        duplicate alerts       0   (info: alerts case management got twice)

### Comment:
Unsuprisingly, it worked.



# one Closing answer: Debezium chose two defaults that you had to change. What would have happened in production if nobody had checked them, and how could the team have noticed?
That really depends on how it was set up!
If they have propper logging infrastructure and tested their product well, then they are going to catch this right after trying the rollout.
e.g.: They would get an alert that the server is experiencing high lag.

Alternatively, depending on the system they might see the processor spike on only one core.

Alternatively, they would probably try out their own system and test it. E.g.: generate fake transaction data / try it out fr.
they would then see a missing alert. (or wouldn't see)

After rolling out the partition fix, it would be harder to catch.
e.g.: testing above is the only real option if it doesn't occur to them that kafka spreads the logs accross the partitions evenly.

TLDR: Testing and propper load testing is key. (e.g.: set the system up on a really weak VM and flood it with messages. see how it behaves.)


# your final connector, as `Exc_6/solution/fraud-connector.json`

