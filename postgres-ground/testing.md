postgres=#  SELECT * FROM pg_stat_wal;

 wal_records | wal_fpi | wal_bytes | wal_buffers_full | wal_write | wal_sync | wal_write_time | wal_sync_time |          stats_reset
-------------+---------+-----------+------------------+-----------+----------+----------------+---------------+-------------------------------
        2697 |    1895 |   8866201 |               40 |        98 |       56 |              0 |             0 | 2024-11-23 10:56:20.383887+00
(1 row)

`pg_stat_wal` provides the wal stats of the postgres server.

Data directory is `/var/lib/postgresql/data`




pg_tool to visualize wal flies -  `pg_waldump`

Decoding WAL entry in postgres : 




