import numpy as np
import pandas as pd
from sklearn import linear_model
from sklearn import svm

from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
from sklearn import preprocessing
from sklearn.metrics import confusion_matrix, accuracy_score
import seaborn as sns

from sklearn.preprocessing import LabelEncoder
from sklearn.externals import joblib

def warn(*args, **kwargs):
    pass
import warnings
warnings.warn = warn

logreg = joblib.load('./logreg.pb1')
scaler = joblib.load('./scaler.pb1')

def combatproba(b):
    """predict probability battle"""
    

    battle = np.array(b)
    battle = scaler.transform(battle.reshape(1, -1))
    y_battle_prob = logreg.predict_proba(battle.reshape(1, -1))

    return round(y_battle_prob[0][1]*100,2)


import csv
with open('combats.csv', 'r') as f:
    reader = csv.reader(f, delimiter=';')
    your_list = list(reader)

for row in your_list:
    chance = str(combatproba(row[:9]))
    print(row[9] + ' has ' + chance + '% chance of winning')

