import sys
import re

text = open(sys.argv[1]).read()
text = re.sub(r"｜(.+?)《.+?》", r"\1", text)
text = re.sub(r"《.+?》", r"", text)
text = re.sub(r"(.)ゝ", r"\1\1", text)
text = re.sub(r"(..)／＼", r"\1", text)

print(text)
