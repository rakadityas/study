-- https://leetcode.com/problems/confirmation-rate

WITH confirmed AS (
    SELECT user_id, COUNT(*) AS confirmed
    FROM Confirmations
    WHERE action = 'confirmed'
    GROUP BY user_id
),
timeout AS (
    SELECT user_id, COUNT(*) AS timeout
    FROM Confirmations
    WHERE action = 'timeout'
    GROUP BY user_id
)
SELECT s.user_id, 
       coalesce(ROUND(
           IFNULL(cf.confirmed, 0) / 
           (IFNULL(cf.confirmed, 0) + IFNULL(t.timeout, 0)),
           2
       ), 0) AS confirmation_rate
FROM Signups s
LEFT JOIN confirmed cf ON s.user_id = cf.user_id
LEFT JOIN timeout t ON s.user_id = t.user_id
ORDER BY confirmation_rate, s.user_id;