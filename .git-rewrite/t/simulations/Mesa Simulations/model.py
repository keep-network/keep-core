from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, ticket_distribution, active_group_threshold, group_size, signature_threshold, group_expiry):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = False
        self.active_groups = []
        self.active_group_threshold = active_group_threshold
        self.signature_threshold = signature_threshold
        self.group_size = group_size
        self.ticket_distribution = ticket_distribution
        self.newest_group_id = 0
        self.group_expiry = group_expiry
        self.timer = 0

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, self.ticket_distribution[i])
            self.schedule.add(node)

        #bootstrap active groups
        for i in range(active_group_threshold):
            self.active_groups.append(self.group_registration())
            print(self.active_groups[i].id)

        #add groups to the scheduler
        for group in self.active_groups:
            self.schedule.add(group)


    def step(self):
        '''Advance the model by one step'''
        #check how many active groups are available
        self.refresh_active_list()
        print(len(self.active_groups))
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
        max_tickets = int(max(self.ticket_distribution))
        for node in node_list:
            adjusted_ticket_list = []
            node.generate_tickets()
            adjusted_ticket_list = np.concatenate([node.ticket_list[self.newest_group_id],np.ones(int(max_tickets)-len(node.ticket_list[self.newest_group_id]))])  #adds 2's the ends of the list so that the 2D array can have equal length rows
            ticket_list.append(adjusted_ticket_list)

        #iteratively add group members by lowest value
        while len(group_members) <= self.group_size:

            min_index = np.where(ticket_list == np.min(ticket_list)) # find the index of the minimum value in the array
            for i,index in enumerate(min_index[0]): #if there are repeated values, iterate through and add the indexes to the group
                group_members.append(index)
                ticket_list[index][min_index[1][i]] = 2 # Set the value of the ticket to a high value so it doesn't get counted again
        
        #create a group agent which can track expiry, sign, etc
        group_object = agent.Group(self.newest_group_id, self, group_members, self.group_expiry, self.signature_threshold)

        return group_object

    def refresh_active_list(self):
        temp_list = []

        for group in self.active_groups:
            if group.status == "Active":
                temp_list.append(group)
        
        self.active_groups = temp_list



        

        







