import matplotlib.pyplot as plt

colors = ['#4473C5', '#325596', '#4369B2', '#8DA3D4']
config = {
    "font.family": 'serif',
    "font.serif": ['simsun', 'Times New Roman'],
    "font.size": 18,
    "mathtext.fontset": 'stix',
}

plt.rcParams.update(config)

labels = '只有DT适用', '两者均适用', '只有DPI适用' , '两者均不适用'
sizes = [1512-749, 749, 791-749, 1681-1554]
otherSizes = [15102-12859,12859,51034-12859,291721-53277]
print(sizes)
print(otherSizes)
# Pie chart, where the slices will be ordered and plotted counter-clockwise:

explode = (0, 0.1, 0, 0)  # only "explode" the 2nd slice (i.e. 'Hogs')

fig1, ax1 = plt.subplots()
# ax1.pie(otherSizes, explode=explode, labels=labels, autopct='%1.1f%%',startangle=90, colors=colors)
# ax1.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.

labels = '只有DT识别', '两者均识别', '只有CART识别' , '两者均不识别'
sizes = [68,589,58,34]
ax1.pie(sizes, explode=explode, labels=labels, autopct='%1.1f%%',startangle=90, colors=colors)
ax1.axis('equal')  # Equal aspect ratio ensures that pie is drawn as a circle.
plt.show()

# plt.savefig("pie_other.pdf", bbox_inches = 'tight', dpi=800)
# plt.show()