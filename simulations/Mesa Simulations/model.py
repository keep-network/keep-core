from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, token_distribution, active_group_threshold, group_size):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = False
        self.active_groups = []
        self.active_group_threshold = active_group_threshold
        self.group_size = group_size
        self.timer = 0

        #create ticket distribution
        tickets = np.ones(nodes)*10 #this can be any distribution we decide

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, tickets[i])
            self.schedule.add(node)

        #bootstrap active groups
        for i in range(active_group_threshold):
            self.group_registration()

    def step(self):
        '''Advance the model by one step'''
        #generate relay requests
        self.relay_request = bool(np.random.randint(0,1))

        if self.relay_request:
            print(self.active_groups[np.random.randint(len(self.active_groups))]) # print a random group from the active list- change this to signing mechanism later
            self.group_registration()
        else:
            print("No relay request")
        self.timer += 1

        #advance the agents
        self.schedule.step()

    def group_registration(self):
        node_list = []
        ticket_list = []
        group_members = []

        # find all the agents that are nodes and add them to a list
        for i in self.schedule.agents:
            if i.type == "node":
                node_list.append(i)

        # make each node generate tickets and save them to a list
        for node in node_list:
            node.generate_tickets()
            ticket_list.append(node.ticket_list[self.timer])

        #iteratively add group members by lowest value
        while len(group_members) < self.group_size:
            min_index = np.where(ticket_list == np.min(ticket_list)) # find the index of the minimum value in the array
            for i,index in enumerate(min_index[0]): #if there are repeated values, iterate through and add the indexes to the group
                group_members.append(index)
                ticket_list[index][min_index[1][i]] = 2 # Set the value of the ticket to a high value so it doesn't get counted again

        return group_members

        

        







