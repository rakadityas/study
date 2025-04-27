-- https://leetcode.com/problems/managers-with-at-least-5-direct-reports/description/

-- USING SUBQUERY
--
-- select e1.name
-- from Employee e1
-- where (
--     select count(e2.id)
--     from Employee e2
--     where e2.managerId = e1.id
-- ) >= 5

select e1.name
from Employee e1
join Employee e2 
on e2.managerId = e1.id
group by e1.id, e1.name
having count(e2.id) >= 5