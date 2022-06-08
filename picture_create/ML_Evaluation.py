import matplotlib.pyplot as plt
import numpy as np

# 无法显示中文
plt.rcParams["font.sans-serif"]=["SimHei"]
plt.rcParams["font.family"]="sans-serif"
#无法显示负号
plt.rcParams['axes.unicode_minus'] = False

config = {
    "font.family": 'serif',
    "font.serif": ['simsun', 'Times New Roman'],
    "font.size": 12,
    "mathtext.fontset": 'stix',
}
plt.rcParams.update(config)

plt.legend(loc='lower right')

# plt.subplot(1, 2, 1)
x = [5, 6, 7, 8, 9, 10, 11, 12 ,13 ,14 ,15]

scores = np.array([
    [0.979392,0.628666,0.881300,0.733849],
    [0.989897,0.864419,0.810902,0.836806],
    [0.974106,0.775266,0.797538,0.786244],
    [0.979549,0.946429,0.697852,0.803351],
    [0.948943,0.713612,0.861309,0.780535],
    [0.944597,0.670481,0.887207,0.763767],
    [0.966869,0.906303,0.871417,0.888518],
    [0.971845,0.910953,0.895013,0.902913],
    [0.964966,0.916304,0.877211,0.896332],
    [0.960354,0.873072,0.881437,0.877235],
    [0.958353,0.856965,0.916223,0.885604],
])
Accuracy = scores[:, 0]
Precision = scores[:, 1]
Recall = scores[:, 2]
F1 = scores[:, 3]
mark = ""
# F1 = [0.7225470128584505, 0.8559549537325235, 0.7783549794826123, 0.7412278956231481, 0.7763947888742115, 0.7663548405546547, 0.8864403198041921, 0.8988171343643255, 0.9040084428566313, 0.8725825035567292, 0.8751559550683534]
plt.subplot(2,2,1)
plt.ylim(0.6, 1)
plt.plot(x, Accuracy, marker = mark)
plt.xlabel("数据包个数")
plt.ylabel("准确率")
plt.xticks(np.arange(5, 15, 1))
# plt.title("数据包-准确率")

plt.subplot(2,2,2)
plt.ylim(0.6, 1)
plt.plot(x, Precision, marker = mark)
plt.xlabel("数据包个数")
plt.ylabel("精确率")
plt.xticks(np.arange(5, 15, 1))
# plt.title("数据包-精确率")

plt.subplot(2,2,3)
plt.ylim(0.6, 1)
plt.plot(x, Recall, marker = mark)
plt.xlabel("数据包个数")
plt.ylabel("召回率")
plt.xticks(np.arange(5, 15, 1))
# plt.title("数据包-召回率")

plt.subplot(2,2,4)
plt.ylim(0.6, 1)
plt.plot(x, F1, marker = mark)
plt.xlabel("数据包个数")
plt.ylabel("F1点数")
plt.xticks(np.arange(5, 15, 1))
# plt.title("数据包-F1点数")

plt.subplots_adjust(wspace = 0.5, hspace = 0.5)
plt.suptitle("机器学习评估")
plt.show()