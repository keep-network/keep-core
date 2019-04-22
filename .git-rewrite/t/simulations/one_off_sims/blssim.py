import numpy as np  


def bls(seed):
    count = 0
    while seed <= 100:
        grps = np.random.exponential(1,100) # Change this to apprpriate "collusion distr"
        if grps[seed] >  0.95 and grps[seed] < 1:
            seed = 1005
        else:
            seed = np.random.randint(0,100)
            count += 1

    return count
        


for i in range(1,100):
    a = []
    seed = np.random.randint(0,100)
    a.append(bls(seed))
print("mean", np.mean(a))


 