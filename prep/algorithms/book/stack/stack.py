# Stack's "Abstract Data Type" interface
# Stack() creates new empty stack
# pop() removes from the top of the stack
# peek() returns the top item from the stack
# is_empty() returns a boolean representing wether the stack is empty
# size returns the number of items on the stack as an interger

# We can implement a stack in Python using a list
# Remember pop() on a list is O(1) and append is amoritized O(1)
# So we get the performance characteristics one would want out of a Stack

class Stack():
    def __init__(self):
        self.__s = []

    def push(self, item):
        self.__s.append(item) # O(1)

    def pop(self):
        if not self.is_empty():
            return self.__s.pop()

    def peek(self):
        if not self.is_empty(): # O(1)
            return self.__s[-1]

    def is_empty(self):
        return self.size() == 0 # O(1)

    def size(self):
        return len(self.__s)


#Demo
if __name__ == "__main__":
    s = Stack()
    s.push(1)
    assert s.peek() == 1
    item = s.pop()
    assert item == 1
    next_item = s.pop()
    assert next_item == None
    assert s.size() == 0
    assert s.is_empty()
