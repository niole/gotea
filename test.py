"""
test parse
"""

import numpy as np
from pprint import pprint
from sklearn.preprocessing import normalize
import spacy

texts = [
u"The inspiration for this tea is Wuyi Xiao Zhong, a traditional and full-bodied black tea from Wuyishan (sometimes smoked with special pine wood). You might expect a bracing and full bodied tea, but when you taste this Tieguanyin varietal processed as a curled, roasted black tea, you find something quite different. The Liu Family has created a tea that is smooth, sweet, and rich. Full of mild cocoa notes with a hint of honey and fruit, this easy-going tea is like a cross between honey-cured burdock and Laoshan Black.",
u"This tea is a testament to the craft of the He Family. Picked in the He Familys fields while still covered with two layers of greenhouse protection, this 2017 spring green tea is carefully hand pressed by the family one leaf and bud at a time. The flat pressing process takes a full day of intense labor, and brings out more sweetness and a longer aftertaste. The Dragonwell-style finishing is a tribute to the fact that Laoshan tea's ancestor is the Dragonwell varietal Longing Qunti.",
u"Master Zhang is a true innovator. He doesnt make tea to follow trends. He experiments and takes risks to make tea better for the generations to come. This Original Wulong Revival uses the older Ben Shan varietal leaf and undergoes three times more careful hand turning and fluffing than modern Anxi oolong. For finishing, it is loosely rolled in the oldest style of oolong making that is half strip style and half ball, with many of the leaves more strip-style than rolled. Master Zhang describes the shape as a dragonfly. This hand processing and shaping yields a different tea- a genre of its own outside of Wuyi style, Guangdong style or Anxi style. The light roast is rewarding and brings out a unique savory sweet complexity we dont see in other teas from Master Zhang.",
u"Bai Mudan, or White Peony is usually picked after the earliest silver needle pickings and includes a blend of downy silver buds and nearly equally downy young leaves. The suspended down makes the tea thick and rich. The Weng family grows their tea high above Fuding in Hulin using traditional sustainable farming techniques and careful air-drying to lock in the flavor of the fresh leaf.",
u"Master Zhou's Gu Hua harvest is a careful blend of maocha from trees aged between one hundred and three hundred years old, picked for a balanced and rich full body and aroma. Gu Hua is the very early autumn harvest prized for its rich flavor and intense aroma. These truly wild trees grow in one of the oldest and most remote tea forests in the world, on Mt. Ailao. Every leaf is hand picked and carefully sun-dried without applying heat or using machinery for the most natural and pure flavor."
]

nlp = spacy.load('en')
docs = np.array([normalize(nlp(doc).vector[:,np.newaxis]).ravel() for doc in texts])
query = normalize(nlp(u"black pepper smoked wood melon cocoa honey fruit").vector[:,np.newaxis]).ravel()
pprint(docs.dot(query.T))
