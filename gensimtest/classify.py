print("initializing...")

from gensim.models import Word2Vec, KeyedVectors
from gensim.test.utils import datapath
from janome.tokenizer import Tokenizer
import numpy as np
import json
import sys

# prepare tokenizer
tokenizer = Tokenizer()

# load model
model = Word2Vec.load("word2vec.gensim.model")

# read vectors
lyric_vec = json.load(open("lyric_vec.json"))

# read target
target = open(sys.argv[1]).read()

words = []
for token in tokenizer.tokenize(target):
    type = token.part_of_speech.split(",")[0]
    if type not in ["助詞", "記号", "助動詞"]:
        words.append(token.base_form)

word_vecs = []
for w in words:
    try:
        word_vecs.append(model.wv[w])
    except KeyError:
        print("[-] Skipping word: {}".format(w))

# summing and normalize
x = np.sum(word_vecs, axis=0)
vec = x / np.linalg.norm(x)

print(vec)


def calc_similarity(v1, v2):
    return np.dot(v1, v2) / (np.linalg.norm(v1) * np.linalg.norm(v2))


similarities = []
# calc cos similarity
for lyric in lyric_vec:
    sim = calc_similarity(vec, lyric["vector"])
    similarities.append((sim, lyric["title"]))

similarities = sorted(similarities, reverse=True)
print(similarities[:10])
