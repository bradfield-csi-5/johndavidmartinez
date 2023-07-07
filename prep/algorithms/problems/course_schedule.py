from collections import defaultdict

# This is an implementation of Khan's algorithm
# That omits the topiological sort in acyclic cases
# https://en.wikipedia.org/wiki/Topological_sorting

# O(j) where j is head nodes
# Is there a faster way to do this?
# Keeping track of the indegree as you go?
# https://leetcode.com/problems/course-schedule-ii/solutions/982542/python3-kahns-algorithm-beats-90-and-dfs-topological-sort/
def has_incoming_edge(graph, key):
    for k in graph.keys():
        if key in graph[k]:
            return True
    return False

def solution(num_courses, prerequisites):
    # First construct a graph
    graph = defaultdict(set)
    # Nodes that have no incoming edge
    head_nodes = set()
    for prereq in prerequisites:
        dependency, dependent = prereq
        graph[dependent].add(dependency)
        # Dependent may be a head node
        head_nodes.add(dependent)
        try:
            # Dependency is definitely not a head node
            # as it has an incoming edge by definition
            head_nodes.remove(dependency)
        except KeyError:
            pass

    # In some edge cases there is still a false head node in our list
    head_nodes = set([n for n in head_nodes if not has_incoming_edge(graph, n)])

    # No head nodes, there's a cycle at the top-level
    if not head_nodes:
        return False

    while head_nodes:
        head = head_nodes.pop()
        while graph[head]:
            dependency = graph[head].pop()
            # Check if dependency is now head node
            if not has_incoming_edge(graph, dependency):
                head_nodes.add(dependency)

    return not bool([x for x in graph.values() if len(x) > 0])



if __name__ == "__main__":
    # Test Suite
    tests = [
            [[[1, 0], [0, 1]], 2, False],
            [[[2, 1], [1, 0]], 3, True],
            [[[2, 1], [1, 0], [2, 2]], 3, False]
    ]
    for t in tests:
        prereq, num_courses, expected = t
        assert solution(num_courses, prereq) == expected

