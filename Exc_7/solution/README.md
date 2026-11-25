
# Your decision: which architecture and why
I will use Architecture B (Lake of truth) since it ensures that all records are consistant and only appear once.

Architecture A doesn't fulfill the reqirements for duplicated since it cannot save the file before writing (maybe something like a WAL would be better suited for this task, but even then)
Architecture C works, but only on the outside, since duplicates still get saved, which isn't the end of the world, but unneccasary.

Architecture B will stage all changes before actually commiting them. e.g.: uncommited directories only contain unfinished/interrupted writes.


# Explain your result
Architecture A saves a file after writing, but if the write failes, then the new records are in the system, but the writer doesn't know that.

e.g.: There is no protection for in-progres writes, which is apparent based on the output of the run (4. run stops WHILE writing) duplicates 32182 records.


Architecture C will do the same as Architecture A (since they are basically the same thing), but the SQL querry used to actually read the data will handle the duplicates.

This makes it look like it worked, but internally, all records are still saved, they just don't show up on the output.

This is inefficient with space and adds overhead to the querries.


Architecture B will write new records, but not commit (e.g.: publish) them until it's done writing.

This Architecture workes and has very small overhead since nothing about the querry / writing process actually changes.

However, if the lake is gigantic, then the initial querry to determine the state will become relatively expensive.

-> for gigantic lakes, a simple WAL would probably be a better choice.

-> The WAL would contain the item that failed, which could then be written, along with everything else in the WAL during recovery.

-> This would also make concurrency easier (e.g.: because of the state file) and faster, since every reader wouldn't need to scan the entire lake to determine its offset.


# Runs

## run_once_state_file

reset: lake, staging and state deleted

| situation | in lake | duplicates | missing | not loaded | run | check |
| --- | --- | --- | --- | --- | --- | --- |
| 1 first run | 2,000,000 | 0 | 0 | 0 | 7.3s | PASS |
| 2 new records arrive, run again | 2,100,000 | 0 | 0 | 0 | 2.2s | PASS |
| 3 nothing new, run again | 2,100,000 | 0 | 0 | 0 | 0.1s | PASS |
| 4 run stops WHILE writing: right after  | 2,132,182 | 0 | 0 | 67,818 |  | PASS |
|   ... then the next run | 2,232,182 | 32,182 | 0 | 0 | 2.2s | *FAIL* |
| 5 run stops AFTER writing: right after | 2,332,182 | 32,182 | 0 | 0 |   | *FAIL* |
|   ... then the next run | 2,432,182 | 132,182 | 0 | 0 | 2.2s | *FAIL* |

Exactly once after every run: *NO* — see the rows marked *FAIL*
~~~yaml
[('1 first run',
  {'topic': 2000000,
   'lake': 2000000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 4,
   'mb': 111.330808,
   'ok': True,
   'problems': []},
  7.281781911849976),
 ('2 new records arrive, run again',
  {'topic': 2100000,
   'lake': 2100000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 7,
   'mb': 112.434078,
   'ok': True,
   'problems': []},
  2.177886486053467),
 ('3 nothing new, run again',
  {'topic': 2100000,
   'lake': 2100000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 7,
   'mb': 112.434078,
   'ok': True,
   'problems': []},
  0.05332660675048828),
 ('4 run stops WHILE writing: right after',
  {'topic': 2200000,
   'lake': 2132182,
   'dup': 0,
   'miss': 0,
   'lag': 67818,
   'files': 8,
   'mb': 112.790925,
   'ok': True,
   'problems': []},
  None),
 ('  ... then the next run',
  {'topic': 2200000,
   'lake': 2232182,
   'dup': 32182,
   'miss': 0,
   'lag': 0,
   'files': 11,
   'mb': 113.893776,
   'ok': False,
   'problems': ['32,182 duplicates']},
  2.1733381748199463),
 ('5 run stops AFTER writing: right after',
  {'topic': 2300000,
   'lake': 2332182,
   'dup': 32182,
   'miss': 0,
   'lag': 0,
   'files': 14,
   'mb': 114.979269,
   'ok': False,
   'problems': ['32,182 duplicates']},
  None),
 ('  ... then the next run',
  {'topic': 2300000,
   'lake': 2432182,
   'dup': 132182,
   'miss': 0,
   'lag': 0,
   'files': 17,
   'mb': 116.064762,
   'ok': False,
   'problems': ['132,182 duplicates']},
  2.152933359146118)]
~~~
## run_once_lake_truth

reset: lake, staging and state deleted

| situation | in lake | duplicates | missing | not loaded | run | check |
| --- | --- | --- | --- | --- | --- | --- |
| 1 first run | 2,300,000 | 0 | 0 | 0 | 10.8s | PASS |
| 2 new records arrive, run again | 2,400,000 | 0 | 0 | 0 | 3.0s | PASS |
| 3 nothing new, run again | 2,400,000 | 0 | 0 | 0 | 0.1s | PASS |
| 4 run stops WHILE writing: right after | 2,400,000 | 0 | 0 | 100,000 |  | PASS |
|   ... then the next run | 2,500,000 | 0 | 0 | 0 | 2.1s | PASS |
| 5 run stops AFTER writing: right after | 2,500,000 | 0 | 0 | 100,000 |  | PASS |
|   ... then the next run | 2,600,000 | 0 | 0 | 0 | 2.7s | PASS |

Exactly once after every run: *YES*
~~~yaml
[('1 first run',
  {'topic': 2300000,
   'lake': 2300000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 4,
   'mb': 116.481432,
   'ok': True,
   'problems': []},
  10.804900169372559),
 ('2 new records arrive, run again',
  {'topic': 2400000,
   'lake': 2400000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 8,
   'mb': 117.627664,
   'ok': True,
   'problems': []},
  2.9500675201416016),
 ('3 nothing new, run again',
  {'topic': 2400000,
   'lake': 2400000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 8,
   'mb': 117.627664,
   'ok': True,
   'problems': []},
  0.06372261047363281),
 ('4 run stops WHILE writing: right after',
  {'topic': 2500000,
   'lake': 2400000,
   'dup': 0,
   'miss': 0,
   'lag': 100000,
   'files': 8,
   'mb': 117.627664,
   'ok': True,
   'problems': []},
  None),
 ('  ... then the next run',
  {'topic': 2500000,
   'lake': 2500000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 11,
   'mb': 118.707137,
   'ok': True,
   'problems': []},
  2.090510129928589),
 ('5 run stops AFTER writing: right after',
  {'topic': 2600000,
   'lake': 2500000,
   'dup': 0,
   'miss': 0,
   'lag': 100000,
   'files': 11,
   'mb': 118.707137,
   'ok': True,
   'problems': []},
  None),
 ('  ... then the next run',
  {'topic': 2600000,
   'lake': 2600000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 15,
   'mb': 119.863576,
   'ok': True,
   'problems': []},
  2.666128158569336)]
~~~

## run_once_rem_dup_read

reset: lake, staging and state deleted

| situation | in lake | duplicates | missing | not loaded | run | check |
| --- | --- | --- | --- | --- | --- | --- |
| 1 first run | 2,600,000 | 0 | 0 | 0 | 12.1s | PASS |
| 2 new records arrive, run again | 2,700,000 | 0 | 0 | 0 | 2.9s | PASS |
| 3 nothing new, run again | 2,700,000 | 0 | 0 | 0 | 0.1s | PASS |
| 4 run stops WHILE writing: right after | 2,732,110 | 0 | 0 | 67,890 |  | PASS |
|   ... then the next run | 2,800,000 | 0 | 0 | 0 | 2.8s | PASS |
| 5 run stops AFTER writing: right after | 2,900,000 | 0 | 0 | 0 |  | PASS |
|   ... then the next run | 2,900,000 | 0 | 0 | 0 | 1.5s | PASS |

Exactly once after every run: YES
~~~yaml
[('1 first run',
  {'topic': 2600000,
   'lake': 2600000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 4,
   'mb': 121.627614,
   'ok': True,
   'problems': []},
  12.120233297348022),
 ('2 new records arrive, run again',
  {'topic': 2700000,
   'lake': 2700000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 8,
   'mb': 122.76641,
   'ok': True,
   'problems': []},
  2.944634437561035),
 ('3 nothing new, run again',
  {'topic': 2700000,
   'lake': 2700000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 8,
   'mb': 122.76641,
   'ok': True,
   'problems': []},
  0.05396294593811035),
 ('4 run stops WHILE writing: right after',
  {'topic': 2800000,
   'lake': 2732110,
   'dup': 0,
   'miss': 0,
   'lag': 67890,
   'files': 9,
   'mb': 123.122744,
   'ok': True,
   'problems': []},
  None),
 ('  ... then the next run',
  {'topic': 2800000,
   'lake': 2800000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 13,
   'mb': 124.228835,
   'ok': True,
   'problems': []},
  2.7637085914611816),
 ('5 run stops AFTER writing: right after',
  {'topic': 2900000,
   'lake': 2900000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 15,
   'mb': 125.275244,
   'ok': True,
   'problems': []},
  None),
 ('  ... then the next run',
  {'topic': 2900000,
   'lake': 2900000,
   'dup': 0,
   'miss': 0,
   'lag': 0,
   'files': 17,
   'mb': 126.321653,
   'ok': True,
   'problems': []},
  1.5325474739074707)]
~~~