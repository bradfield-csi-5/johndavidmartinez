

def convert_to_binary(num):
    bin_s = []
    while num:
        right_most_bit = num % 2
        bin_s.append(str(right_most_bit))
        num //= 2
    # we can now "pop" them off to find the num
    # or just reverse O(n)
    # Now you could pop it off and create a new list
    # To illustrate why a stack solves this problem
    return ''.join(reversed(bin_s))


print(convert_to_binary(10))
print(convert_to_binary(4))

# You can generalize this across bases
# For example a hex converter could be made out of it

DIGITS = '0123456789abcdef'
