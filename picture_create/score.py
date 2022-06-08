import matplotlib.pyplot as plt
TP = 1453
FP = 192
TN = 52703
FN = 101
Accuray = (TP+TN)/(TP+FP+TN+FN)
Precision = TP/(TP+FP)
Recall = TP/(TP+FN)
F1 = (2*Precision*Recall)/(Precision + Recall)
print(Accuray, Precision, Recall, F1)

scores = [Accuray, Precision, Recall, F1]
labels = ["Accuray", "Precision", "Recall", "F1"]
plt.bar(range(len(scores)), scores, tick_label = labels)
plt.show()