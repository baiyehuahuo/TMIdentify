import numpy as np
import matplotlib.pyplot as plt

config = {
    "font.family": 'serif',
    "font.serif": ['simsun', 'Times New Roman'],
    "font.size": 12,
    "mathtext.fontset": 'stix',
}
plt.rcParams.update(config)


N = 2
ind = (0.2, 0.4)

# plt.xticks(np.arange(0, 10, 1))
plt.xlabel('分类方法')
plt.xlim(0.1, 0.5)
plt.xticks(ind, ('DPI', 'DT'))
plt.ylabel('flow占比')
plt.yticks(np.arange(0, 1.1, 0.1))

tmTotal = 1681
dpiTotal = 791
mlTotal = 1512
dpiClassify = 695
mlClassify = 1273
# both = 715 34 749

Correct = (dpiClassify/tmTotal, mlClassify/tmTotal)
Identify = (dpiTotal/tmTotal, mlTotal/tmTotal)
print(dpiTotal/tmTotal, dpiClassify/dpiTotal)
print(mlTotal/tmTotal, mlClassify/mlTotal)
print()
print(Correct)
print(((dpiTotal-dpiClassify)/tmTotal, (mlTotal-mlClassify)/tmTotal))
print((tmTotal-dpiTotal)/tmTotal, (tmTotal-mlTotal)/tmTotal)
Total = (1, 1)

width = 0.1
p1 = plt.bar(ind, Total, width, color=(0.8, 0.8, 1))#'#A9A9A9') 
p2 = plt.bar(ind, Identify, width, color=(0.6, 0.6, 1))#'#808080')  
p3 = plt.bar(ind, Correct, width, color=(0.4, 0.4, 1))#'#696969')
plt.legend((p1[0], p2[0], p3[0]), ('Total', 'Classify', 'Correct'), bbox_to_anchor=(0.142,1), ncol=3, handlelength=1.8)
plt.tight_layout()
plt.savefig("stack.pdf", bbox_inches = 'tight', dpi=800)
plt.show()