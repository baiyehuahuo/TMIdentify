import matplotlib.pyplot as plt



# 无法显示中文
plt.rcParams["font.sans-serif"]=["SimHei"]
plt.rcParams["font.family"]="sans-serif"
#无法显示负号
plt.rcParams['axes.unicode_minus'] =False

config = {
    "font.family": 'serif',
    "font.serif": ['simsun'],
    "font.size": 12,
    "mathtext.fontset": 'stix',
}
plt.rcParams.update(config)

# plt.subplot(1, 2, 1)
x = [5,290,431,1075,1139,1408,1426,1551,1554,1557,1560,1575,1674,1937,1989,1997,2753,2924,2928,3067,3096,3303,3559,3661,3701,3740,3786,3789,3932,3949,3950,3998,4002,4006,4007,4020,4104,4225,4248,4250,4406,4473,4633,4677,4683,4706,4739,4792,4793,4857,4980,5252,5428,5677,6277]
y = [1,1,1,2,4,1,1,10,38,41,51,1,1,1,4,1,5,5,6,1,3,1,1,11,2,1,38,16,2,6,9,5,15,1,127,31,1,2,3,1,1,2,1,2,17,2,1,4,1,3,1,1,6,1,10]
total = plt.plot(x, y, label="total")
# y = [2,1,1,1,2,1,15,10,1,4,1,6,5,1,0,1,0,2,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,5,0,1,0,8,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,16,0,0,0,0,0,0,0,0,0,0,0,0,6,1]
# unrecognized = plt.plot(x, y, label="unrecognized")
# plt.grid(color="k", linestyle=":")
# plt.legend(loc='upper left')
size = 16
plt.xlabel("会话流前5个数据包负载总长度", fontsize=16)
plt.ylabel("频率", fontsize=16)

# plt.subplot(1, 2, 2)
# x = [1, 2, 3, 4]
# y = [1, 4, 9, 16]
# plt.plot(x, y)
plt.tight_layout()
plt.savefig("DPI_Distribution_5.pdf", bbox_inches = 'tight', dpi=80)
plt.show()