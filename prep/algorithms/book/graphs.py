# Vertex (Node)
# Edge (connection between nodes)
# Weight optional cost to traverse an Edge


# You can represent a graph as a 2D matrix
# x = None
# [      v1
#    [1, 2, 2]
#v2  [x, 4, x],  v1 connects v2 with a weight of 4
#    [4, 2, 2]   never really seen this before
# ]  # this is also not efficent
# Rare are graphs this interconnected in practice
# So the matrix isn't usually used

# That doesn't mean problems won't use matrix to represent graphs though
# For example a 2D matrix can represent a "Map" of objects and you may
# want to consider the matrix as representing a type of graph.

# Adjacency List
# what you typically think of as a graph
# a Graph objects with GraphNodes that have a list of connected nodes


class Vertex():
    def __init__(self, key):
        self.key = key
        self.neighbors = {}

    def add_neighbor(self, neighbor, weight=0):
        self.neighbors[neighbor] = weight

    def __str__(self):
        return '{} neighbors: {}'.format(
                self.key,
                [x.key for x in self.neighbors]
        )

    def get_connections(self):
        return self.neighbors.keys()


    def get_weight(self, neighbor):
        return self.neighbors[neighbor]


class Graph():
    def __init__(self):
        self.verticies = {}

    def add_vertex(self, vertex):
        self.verticies[vertex.key] = vertex

    def get_vertex(self, key):
        self.verticies
        try:
            return self.verticies[key]
        except KeyError:
            return None

    def __contains__(self, key):
        return key in self.verticies

    def add_edge(self, from_key, to_key, weight=0):
        if from_key not in self.verticies:
            self.add_vertex(Vertex(from_key))
        if to_key not in self.verticies:
            self.add_vertex(Vertex(to_key))
        self.verticies[from_key].add_neighbor(self.verticices[to_key], weight)

    def get_vertices(self):
        return self.verticies.keys()

    def __iter__(self):
        return iter(self.verticies.values())






