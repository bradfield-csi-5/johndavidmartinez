cases = [
    ["(()()()())", True],
    ["(((())))", True],
    ["(()((())()))", True],
    ["((((((())", False],
    ["()))", False],
    ["(()()(()", False],
]

# This is O(n) time as the entire param list must be traversed
# This is O(n) space as the worst case is all open parens
def is_balanced(parens):
    s = [] #stack
    for p in parens:
        if p == "(":
            # Open paren add to stack
            s.append(p)
        if p == ")":
            # Close paren try to close from stack
            if len(s) == 0:
                return False # Nothing in Stack
            peek = s[-1]
            if peek == "(":
                # Found what we were looking for
                s.pop()
    return len(s) == 0 # If all parens match stack should be empty


if __name__ == "__main__":
    for c in cases:
        assert is_balanced(c[0]) == c[1]



# You can imagine how we'd expand this to support general paraens
PAIRING = {
        "{": "}",
        "[": "]",
        "(": ")"
}
