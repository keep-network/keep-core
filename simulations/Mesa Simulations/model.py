from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, ticket_distribution, active_group_threshold, group_size, signature_threshold, group_expiry, group_formation_threshold):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = False
        self.active_groups = []
        self.active_nodes = []
        self.active_group_threshold = active_group_threshold # number of groups that will always be maintained in an active state
        self.signature_threshold = signature_threshold
        self.group_size = group_size
        self.ticket_distribution = ticket_distribution
        self.newest_group_id = 0
        self.group_expiry = group_expiry
        self.bootstrap_complete = False # indicates when the initial active group list bootstrap is complete
        self.group_formation_threshold = group_formation_threshold # min nodes required to form a group
        self.timer = 0

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, self.ticket_distribution[i])
            self.schedule.add(node)


    def step(self):
        '''Advance the model by one step'''
        print("newest group id" + str(self.newest_group_id))
        
        #check how many active nodes are available
        self.refresh_connected_nodes_list()
        print("Number of nodes in the forked state = " + str(len(self.active_nodes)))


        #bootstrap active groups as nodes become available. Can only happen once enough nodes are online
        if self.bootstrap_complete == False:
            if len(self.active_nodes)> self.group_formation_threshold:
                for i in range(self.active_group_threshold):
                    self.active_groups.append(self.group_registration())
                    #print(self.active_groups[i].id)
            

        #check how many active groups are available
        self.refresh_active_group_list()
        print('number of active groups = ' + str(len(self.active_groups)))
        
        #generate relay requests
        self.relay_request = np.random.choice([True,False])

        if self.relay_request:
            try:
                print('selecting group at random')
                print(self.active_groups[np.random.randint(len(self.active_groups))]) # print a random group from the active list- change this to signing mechanism later
            except:
                print('no active groups available')

            print('registering new group')
            self.group_registration()
        else:
            print("No relay request")
        self.timer += 1

        #advance the agents
        self.schedule.step()

    def group_registration(self):
        ticket_list = []
        group_members = []

        if len(self.active_nodes)<self.group_formation_threshold:
            print("Not enough nodes to register a group")

        else:
            # make each node generate tickets and save them to a list
            max_tickets = int(max(self.ticket_distribution))
            for node in self.active_nodes:
                adjusted_ticket_list = []
                node.generate_tickets()
                adjusted_ticket_list = np.concatenate([node.ticket_list,np.ones(int(max_tickets)-len(node.ticket_list))])  #adds 2's the ends of the list so that the 2D array can have equal length rows
                ticket_list.append(adjusted_ticket_list)

            #iteratively add group members by lowest value
            while len(group_members) <= self.group_size:

                min_index = np.where(ticket_list == np.min(ticket_list)) # find the index of the minimum value in the array
                for i,index in enumerate(min_index[0]): #if there are repeated values, iterate through and add the indexes to the group
                    group_members.append(index)
                    ticket_list[index][min_index[1][i]] = 2 # Set the value of the ticket to a high value so it doesn't get counted again
            
            #create a group agent which can track expiry, sign, etc
            group_object = agent.Group(self.newest_group_id, self, group_members, self.group_expiry, self.signature_threshold)

            #add group to schedule
            self.schedule.add(group_object)

            #add group to active group list
            self.active_groups.append(group_object)
            
            return group_object

    def refresh_active_group_list(self):
        temp_list = []

        for group in self.active_groups:
            if group.status == "Active":
                temp_list.append(group)
        
        self.active_groups = temp_list

    def refresh_connected_nodes_list(self):
        print("refreshing active nodes list")
        temp_active_node_list = []
        for agent in self.schedule.agents:
            print("agent id = " + str(agent.id))
            if agent.type == "node":
                if agent.mainloop_status == "forked": 
                    print("forked agent ID" + str(agent.id))
                    temp_active_node_list.append(agent) #adds the node to the active list only if it is in the forked state
        self.active_nodes = temp_active_node_list





        

        







