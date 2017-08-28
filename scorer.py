import numpy as np
from pprint import pprint
from sklearn.preprocessing import normalize
import spacy

nlp = spacy.load('en')

def getRankedMatches(texts, queryString):
    """
    scores texts against query string and returns top 10 results
    all unicode
    """
    docs = np.array([normalize(nlp(doc).vector[:,np.newaxis]).ravel() for doc in texts])
    query = normalize(nlp(queryString).vector[:,np.newaxis]).ravel()
    scores = docs.dot(query.T)
    return scores
