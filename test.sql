SELECT * FROM  "wikiticker-2015-09-12-sampled" LIMIT 10;
SELECT FLOOR(__time to HOUR) AS HourTime, SUM(deleted) AS LinesDeleted FROM wikiticker-2015-09-12-sampled GROUP BY 1 LIMIT 10;
SELECT __time, added, channel, user FROM  "wikiticker-2015-09-12-sampled" WHERE added > 0 LIMIT 10;
