# FIFO

# We could use a python list but then we won't have
# the performance characteristics O(1) that we want
# since a python list.pop(0) is O(n)

# You can see with a little work we can
# make this queue and deque
# We're already tracking the head and the tail
# and the operation to remove and add from each
# is just two additional methods away
class QueueNode():
    def __init__(self):
        self.next = None
        self.prev = None
        self.val = None

class Queue():
    def __init__(self):
        self.size = 0
        self.head = None
        self.tail = None

    def enqueue(self, item):
        self.size += 1
        n = QueueNode()
        if self.head is None:
            self.head = n
            self.head.val = item
            self.tail = self.head
        else:
            self.head.next = QueueNode()
            self.head.next.prev = self.head
            self.head = self.head.next
            self.head.val = item

    def dequeue(self):
        if self.tail is None:
            return
        self.size -= 1
        val = self.tail.val
        if self.head == self.tail:
            self.head = None
            self.tail = None
        else:
            self.tail = self.tail.next
        return val

    def is_empty(self):
        return self.size == 0

    def size(self):
        return self.size


if __name__ == "__main__":
    q = Queue()
    q.enqueue(1)
    q.enqueue(2)
    q.enqueue(3)
    assert q.dequeue() == 1
    assert q.dequeue() == 2
    assert q.dequeue() == 3

