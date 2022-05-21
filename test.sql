SELECT * FROM  "wikipedia" LIMIT 10;
SELECT FLOOR(__time to HOUR) AS HourTime, SUM(deleted) AS LinesDeleted FROM wikipedia GROUP BY 1 LIMIT 10;
SELECT __time, added, channel, user FROM  "wikipedia" LIMIT 10;
