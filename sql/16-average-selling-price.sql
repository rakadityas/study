-- https://leetcode.com/problems/average-selling-price/

select p.product_id, coalesce(round(sum(p.price * us.units)/sum(us.units), 2),0) as average_price
from Prices p
left join UnitsSold us
on p.product_id = us.product_id
and us.purchase_date between start_date and end_date
group by p.product_id