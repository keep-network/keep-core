#%% Change working directory from the workspace root to the ipynb file location. Turn this addition off with the DataScience.changeDirOnImportExport setting
import os
try:
	os.chdir(os.path.join(os.getcwd(), 'keep-core/simulations'))
	print(os.getcwd())
except:
	pass

#%%
import numpy as np
import matplotlib.pyplot as plt

total_tickets =0.5*100*10000
Totals = np.zeros(100)
GroupMax = np.zeros(100)
MaxOwnership = np.zeros(100)
MaxOwnershipCount = 0
GroupSize = []
for i in range(0,1000):
    tickets = np.random.random_sample(int(total_tickets))
    winning_tickets = []
    i = 0
    i_val =[]
    j_val =[]
    for j in range(0,100,1):
        winning_tickets.append(sum(tickets[i:i+100*j]<0.00025))# this assigns ownership to nodes by index
        i_val.append(i)
        j_val.append(j*100)
        i = i+100*j
        j +=1
    GroupSize.append(sum(winning_tickets))
    GroupMax[np.argmax(winning_tickets)] +=1
    MaxOwnership[np.argmax(winning_tickets)]=winning_tickets[np.argmax(winning_tickets)]/sum(winning_tickets)
    MaxOwnershipCount += np.round(winning_tickets[np.argmax(winning_tickets)]/sum(winning_tickets))
    Totals = Totals + winning_tickets
    


#%%
ind = np.arange(len(winning_tickets))
hist = plt.bar(ind,Totals)


#%%
hist2 = plt.bar(ind,GroupMax)


#%%
print("Maximum times 51% ownership of a group in 1000 runs = "+str(MaxOwnershipCount))
print("Group Size Range, Max= "+str(GroupSize[np.argmax(GroupSize)])+ "Min= "+ str(GroupSize[np.argmin(GroupSize)]))
print("Maximum count 51% ownership of a group = "+str(MaxOwnership))
print("Maximum ownership  = "+str(MaxOwnership[np.argmax(MaxOwnership)]))


