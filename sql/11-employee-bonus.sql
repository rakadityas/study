-- https://leetcode.com/problems/employee-bonus

select emp.name, bn.bonus
from Employee emp
left join Bonus bn
on emp.empId = bn.empId
where bn.bonus is Null or bn.bonus < 1000