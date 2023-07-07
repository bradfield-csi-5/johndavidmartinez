#!/usr/bin/python3
# Word Ladder program
import sys
import copy

from collections import defaultdict
from collections import deque
from itertools import product

def load_words_into_dict(filepath):
    words = {}
    with open(filepath, "rb") as f:
        for line in f.readlines():
            words[line.strip().decode("ascii")] = 1
    return words

def build_graph(words):
    buckets = defaultdict(list)
    graph = defaultdict(set)

    # Put words in 3 letter buckets
    for word in words:
        for i in range(len(word)):
            bucket = f"{word[:i]}_{word[i + 1:]}"
            buckets[bucket].append(word)

    # Add verticies and edges for words in same bucket
    for bucket, mutual_neighbors in buckets.items():
        for word1, word2 in product(mutual_neighbors, repeat=2):
            if word1 != word2:
                graph[word1].add(word2)
                graph[word2].add(word1)

    return graph

def traverse(graph, starting_vertex):
    visisted = set()
    queue = deque([[starting_vertex]])
    while queue:
        path = queue.popleft()
        vertex = path[-1]
        yield vertex, path
        for neighbor in graph[vertex] - visisted:
            visisted.add(neighbor)
            queue.append(path + [neighbor])

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage")
        exit(1)

    wordlist_filepath = sys.argv[1]
    starting_word = sys.argv[2]
    ending_word = sys.argv[3]

    print(f"Loading {wordlist_filepath}")
    words = load_words_into_dict(wordlist_filepath)
    if not starting_word in words:
        print(f"Err {starting_word} not in wordlist")
        exit(1)
    if not ending_word in words:
        print(f"Err {ending_word} not in wordlist")
        exit(1)

    # Build Graph
    g = build_graph(words)
    for vertex, path in traverse(g, starting_word):
        if vertex == ending_word:
            print(' -> '.join(path))











