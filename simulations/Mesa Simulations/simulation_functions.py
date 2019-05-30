import random
import datetime
import numpy as np
import copy
import pandas as pd

def min_index(ticket_array, group_size):
    array = copy.deepcopy(ticket_array)
    group = sorted(array)[0:group_size] # generates a sorted list of min values of length = group size
    #print(group)
    indexes = [] #initializes the array of indexes for min ticket values
    
    for ticket in group: #iterates through each ticket value in the sorted list
        
        ticket_index = np.where(array==(ticket)) #finds the index with the ticket value
        #print("ticket_index = " + str(ticket_index))
        indexes.append(ticket_index[0][0]) # adds the index value to the array of indexes
        #print("indexes = " + str(indexes))
        array[ticket_index[0][0]] = 1 #sets the vaue at that index to the max value of 1 to address the problem of repeatd values
        
    #print(ticket_array_temp)
    return sorted(indexes)

def preprocess_tickets(runs, total_tickets):
# Pre-processing ticket arrays
# runs = number of simulation runs
# total_tickets = total # of tickets (virtual stakers)
    tickets=[]
    for i in range(0, runs):
        tickets.append(np.random.random_sample(int(total_tickets)))
    return tickets

def preprocess_groups(tickets, runs, group_size):
# Pre-processing groups
    group_members = []
    for i in range(0, runs):
        group_members.append(min_index(tickets[i], group_size)) # finds the index of group members with min ticket values
    return group_members

def create_cdf(nodes,ticket_distr):
# Create CDF's - used to determine max ownership ticket index
    cdf = np.zeros(nodes)
    for node,ticketmax in enumerate(ticket_distr):
        
        cdf[node]=sum(ticket_distr[0:node+1])
    return cdf

def group_distr(runs, nodes, group_members, cdf):
# function to calculate group ownership distribution
    total_group_distr = np.zeros(nodes)
    max_owned = np.zeros(runs)
    group_distr_matrix = np.zeros((runs,nodes))
    for i in range(runs):
        group_distr = np.zeros(nodes)
        group_distr[1] = sum(group_members[i]<cdf[0])
        for j in range(1,nodes):
            group_distr[j] = sum(group_members[i]<cdf[j])-sum(group_members[i]<cdf[j-1])
        max_owned[i] = max(group_distr)/sum(group_distr)
        total_group_distr +=group_distr
        group_distr_matrix[i] = group_distr #saves the group ticket distribution for each run
        print(group_distr_matrix[i])
    return total_group_distr, max_owned, group_distr_matrix

def node_failures(nodes, runs, node_failure_percent):
# pre-processes failed nodes
    failed_nodes = np.random.rand(runs, nodes) < node_failure_percent
    return failed_nodes