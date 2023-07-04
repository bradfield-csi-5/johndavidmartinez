# List
list indexing O(1)

appending to lists is "amortized" O(1)
Beacuse list expansion which is O(n) at expansion time only happens when the list is exhausted
Since the list size is tripled when exhausted we say this is O(1) on an amortized basis

amortize means "gradually write off the initial cost of an asset over a period of time"

list popping the end is O(1)
because you're only shifting the index once
however popping at a certin position is O(N)
because you're shifting the index n times where n is the elements to the right of the index you popped

Iteration is O(n). Intuitive as you're traversing n items

Slicing
Slicing is O(k) since you must iterate over the slice of size k and make a new array

Multiplying
Multiplying a list is O(nk) n being the size of the list and k being your multiple

reversing
O(n) since you have to reposition n elements

Sorting
Sorting is O(n log n) because the most efficent sorting algorithms known like
merge sort and quick sort have this time complexity.

# Dictionaries
contains O(1)
get O(1)
iteration O(N)

dictionaries are hashmaps that have constant lookup and assignment.
That's because all that's required is hashing a key to access them a single constant cost operation.
Hashmaps can have collisions that degrade performance and their underlying arrays do need to resize
but since these are amortized over time we consider that O(1) like arrays

Python ref: https://wiki.python.org/moin/TimeComplexity
