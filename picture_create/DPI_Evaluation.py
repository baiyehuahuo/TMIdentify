import matplotlib.pyplot as plt
import numpy as np
from pyparsing import line


plt.legend(loc='lower right')
#无法显示负号
plt.rcParams['axes.unicode_minus'] = False

config = {
    "font.family": 'serif',
    "font.serif": ['simsun', 'Times New Roman'],
    "font.size": 12,
    "mathtext.fontset": 'stix',
}
plt.rcParams.update(config)
start = 1
end = 5
mark = ""
x = np.arange(start, end+1, 1)
def plotScores(scores, color, label, linestyle):
    
    Accuracy = scores[:, 0]
    Precision = scores[:, 1]
    Recall = scores[:, 2]
    F1 = scores[:, 3]

    plt.subplot(2,2,1)
    setSubPlot(Accuracy, "准确率", label, "(a)", linestyle)

    plt.subplot(2,2,2)
    plt.ylim(0, 1)
    setSubPlot(Precision, "精确率", label, "(b)", linestyle)

    plt.subplot(2,2,3)
    plt.ylim(0, 1)
    setSubPlot(Recall, "召回率", label, "(c)", linestyle)

    plt.subplot(2,2,4)
    plt.ylim(0, 1)
    setSubPlot(F1, "F1点数", label, "(d)", linestyle)

def setSubPlot(scores, ylabel, label, title, linestyle):
    plt.xlabel("数据包个数")
    plt.ylabel(ylabel)
    plt.grid(axis='y', linestyle='dashdot')
    plt.xticks(np.arange(start, end+1, 1))
    plt.title(title, y=-0.5, fontsize=12)
    plt.plot(x, scores, marker = mark, color = color, label = label, linestyle = linestyle)

values = { 
    '#1f77b4': {
        'scores': np.array([
			[0.997209,1.000000,0.675542,0.806356],
			[0.997806,1.000000,0.758427,0.862620],
			[0.997726,1.000000,0.747782,0.855693],
			[0.999011,1.000000,0.858804,0.924039],
			[0.998001,1.000000,0.827723,0.905742],
        ]),
        'baseThreshold': 0.1,
        'linestyle': '--',
    },
    '#ff7f0e': {
        'scores': np.array([
			[0.997098,1.000000,0.662539,0.797020],
			[0.996950,1.000000,0.664170,0.798200],
			[0.998926,1.000000,0.880862,0.936658],
			[0.998697,1.000000,0.813953,0.897436],
			[0.997817,1.000000,0.811881,0.896175],
        ]),
        'baseThreshold': 0.2,
        'linestyle': '-', 
    },
    '#2ca02c': {
        'scores': np.array([
			[0.995068,1.000000,0.426625,0.598090],
			[0.994791,1.000000,0.426342,0.597812],
			[0.993293,0.743961,0.390368,0.512053],
			[0.998394,0.890572,0.878738,0.884615],
			[0.990098,0.565836,0.629703,0.596064],
        ]),
        'baseThreshold': 0.3,
        'linestyle': '--', 
    }, 
    '#d62728': {
        'scores': np.array([
			[0.995068,1.000000,0.426625,0.598090],
			[0.994791,1.000000,0.426342,0.597812],
			[0.993190,0.611561,0.670469,0.639661],
			[0.991737,0.427419,0.528239,0.472511],
			[0.990006,0.561837,0.629703,0.593838],
        ]),
        'baseThreshold': 0.4,
        'linestyle': '-', 
    },
    '#9467bd': {
        'scores': np.array([
			[0.991399,0.000000,0.000000,0.000000],
			[0.995493,0.849351,0.612360,0.711643],
			[0.995098,0.763158,0.661597,0.708758],
			[0.986291,0.262376,0.528239,0.350606],
			[0.983022,0.373650,0.685149,0.483578],
        ]),
        'baseThreshold': 0.5,
        'linestyle': '--', 
    }
}

for color, val in values.items():
    plotScores(val['scores'], color, val['baseThreshold'], val['linestyle'])


plt.subplots_adjust(wspace = 0.35, hspace = 0.5)
# plt.tight_layout()
plt.legend(bbox_to_anchor=(1.05,2.85),borderaxespad=0.4,ncol=5,)  #绘制表示框，右下角绘制
plt.savefig("DPI_Evaluation_0.5.pdf", bbox_inches = 'tight', dpi=80)
plt.show()