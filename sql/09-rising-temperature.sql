-- https://leetcode.com/problems/rising-temperature/description

select today.id
from Weather today
join Weather yesterday
on datediff(today.recordDate, yesterday.recordDate) = 1
where today.temperature > yesterday.temperature