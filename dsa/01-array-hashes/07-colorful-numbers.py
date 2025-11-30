# https://tutorialhorizon.com/algorithms/colorful-numbers/
# time complexity: O(n^2)
# space complexity: O(n^2)

def is_colorful(number):
    num_str = str(number)
    length = len(num_str)
    product_set = set()

    for i in range(length):
        product = 1
        for j in range(i, length):
            product *= int(num_str[j])
            if product in product_set:
                return False
            product_set.add(product)
    
    return True



