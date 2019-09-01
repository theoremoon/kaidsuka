print("initializing...")

from gensim.models import Word2Vec, KeyedVectors
from gensim.test.utils import datapath
from janome.tokenizer import Tokenizer
import numpy as np
import json

# prepare tokenizer
tokenizer = Tokenizer()
# load model
model = Word2Vec.load("word2vec.gensim.model")
# read lyrics
lyrics = json.load(open("lyrics.json"))
# memo: set of types={'接頭詞', '動詞', '名詞', '助動詞', '助詞', '記号', '感動詞', '形容詞', '接続詞', 'フィラー', '連体詞', '副詞'}

print("DONE!")

lyric_vec = []

for z in lyrics:
    print(z["title"])
    # get all meaningful words
    lyric_words = []
    for line in z["lyrics"]:
        for token in tokenizer.tokenize(line):
            type = token.part_of_speech.split(",")[0]
            if type not in ["助詞", "記号", "助動詞"]:
                lyric_words.append(token.base_form)

    # get word vectors
    word_vecs = []
    for w in lyric_words:
        try:
            word_vecs.append(model.wv[w])
        except KeyError:
            print("[-] Skipping word: {}".format(w))

    # summing and normalize
    x = np.sum(word_vecs, axis=0)
    vec = x / np.linalg.norm(x)
    lyric_vec.append({"title": z["title"], "vector": vec, "phrases": z["phrases"]})


class MyEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, np.integer):
            return int(obj)
        elif isinstance(obj, np.floating):
            return float(obj)
        elif isinstance(obj, np.ndarray):
            return obj.tolist()
        else:
            return super(MyEncoder, self).default(obj)


# save vector
with open("lyric_vec.json", "w") as f:
    json.dump(lyric_vec, f, cls=MyEncoder, ensure_ascii=False, indent=2)
