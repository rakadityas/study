-- https://leetcode.com/problems/product-sales-analysis-i/description

select p.product_name, s.year, s.price
from Sales s
left join Product p
on s.product_id = p.product_id