-- https://leetcode.com/problems/replace-employee-id-with-the-unique-identifier/description

select emuni.unique_id, em.name
from Employees em 
left join EmployeeUNI emuni
on em.id = emuni.id