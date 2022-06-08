# from turtle import color
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
from matplotlib.ticker import MultipleLocator


plt.rcParams['axes.unicode_minus'] = False  # 用来正常显示负号
matplotlib.use('TkAgg')

config = {
    "font.family": 'serif',
    "font.serif": ['simsun'],
    "font.size": 24,
    "mathtext.fontset": 'stix',
}
plt.rcParams.update(config)
def plotCM(classes, matrix, coordinate):
    """classes: a list of class names"""
 
    # Normalize by row
    matrix = matrix.astype(np.float)
    linesum = matrix.sum(1)
    linesum = np.dot(linesum.reshape(-1, 1), np.ones((1, matrix.shape[1])))
    matrix /= linesum
 
    fig = plt.figure()
    ax = fig.add_subplot(coordinate)
    cax = ax.matshow(matrix)
    fig.colorbar(cax)
 
    ax.xaxis.set_major_locator(MultipleLocator(1))
    ax.yaxis.set_major_locator(MultipleLocator(1))
 
    # for i in range(matrix.shape[0]):
    ax.text(0, 0, str('%.2f' % (matrix[0, 0] * 100)), va='center', ha='center')
    ax.text(0, 1, str('%.2f' % (matrix[0, 1] * 100)), va='center', ha='center', color="white")
    ax.text(1, 0, str('%.2f' % (matrix[1, 0] * 100)), va='center', ha='center', color="white")
    ax.text(1, 1, str('%.2f' % (matrix[1, 1] * 100)), va='center', ha='center')

    ax.set_xticklabels(['']+classes)
    ax.set_yticklabels(['']+classes)
 
if __name__ == '__main__':
    classes = ["TM", "Not_TM"]
    savename = ["confusion_matrix.svg", 
                "confusion_matrix_dpi.svg",
                "confusion_matrix_dt.svg"]

    matrixs = [np.array([[715, 34],[91, 12768]]),
               np.array([[657, 92],[0, 12859]]),
               np.array([[647, 102], [91, 12768]])]
    for i in range(len(savename)):
        plotCM(classes, matrixs[i], 111)
        plt.tight_layout()
        plt.savefig(savename[i], bbox_inches = 'tight', dpi=800)
        plt.show()
    


    