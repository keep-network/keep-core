from mesa import Agent, Model
from mesa.time import RandomActivation
import numpy as np

class Node(Agent):
    """ Node: One hardware device used to stake tokens on the network. 
    Each node will create virtual stakers proportional to
    the number of tokens owned by the node """
    def __init__(self, unique_id, node_id, model, tickets, failure_percent, death_percent):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "node"
        self.node_id = node_id
        self.num_tickets = int(tickets)
        self.ticket_list = []
        self.connection_status = "not connected" #change later to event - currently used for node failure process
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"
        self.connection_delay = np.random.randint(0,1) #uniform randomly assigned connection delay step value
        self.mainloop_fork_delay = np.random.randint(0,2) #uniform randomly assigned connection delay step value
        self.timer = self.model.timer
        self.node_failure_percent = failure_percent
        self.node_death_percent = death_percent
        self.failure = False
        self.death = False
    
    def step(self):
        #connect to chain
        if self.connection_delay>0:
            self.connection_delay -=1
        else:
            self.connection_status = "connected"
        #once connected fork the main loop
            if self.mainloop_fork_delay>0:
                self.mainloop_fork_delay -=1
            else: 
                self.mainloop_status = "forked"
            
        #simulate node failure
        self.failure = np.random.randint(0,100) < self.node_failure_percent
        self.death = np.random.randint (0,1000) < self.node_death_percent

    def advance(self):
        if self.failure or self.death:
            self.node_disconnect()

        #print(str("node " + str(self.id) + "status " + self.mainloop_status ))
        #print("Mainloop fork delay = " + str(self.mainloop_fork_delay))
        #print("Mainloop_status = " + self.mainloop_status)

    def generate_tickets(self):
        #generates tickets using the uniform distribution
        self.stake_status = "staked"
        self.ticket_list = np.random.random_sample(self.num_tickets)
        
    def node_disconnect(self):
        # disconnect the node from the network; causes it to go through the entire reconnection sequence in the next step
        self.connection_status = "not connected"
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"
        if self.death == False: # does not reset the failure trigger if the death trigger is true
            self.failure = False


class Group(Agent):
    """ A Group """
    def __init__(self, unique_id, group_id, model, members, expiry):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.group_id = group_id
        self.type = "group"
        self.members = members
        self.last_signature = "none"
        self.status = "Active"
        self.expiry = expiry # of steps before expiration
        self.timer = self.model.timer
        self.model.newest_id +=1
        self.model.newest_group_id +=1
        self.ownership_distr = self.calculate_ownership_distr()

    def step(self):
        """ At each step check if the group as expired """
        self.expiry -=1
        #print('group ID '+ str(self.id) + ' expiry ' + str(self.expiry))
        if self.expiry == 0: 
            self.status = "Expired"
        

    def advance(self):
        pass

    def calculate_ownership_distr(self):
        temp_distr = np.zeros(self.model.num_nodes)
        for node in self.members:    
            temp_distr[node.node_id] +=1 # increments by 1 for each node index everytime it exists in the member list, at each step
        return temp_distr

class Signature(Agent):
    def __init__(self, unique_id, signature_id, model, group_object):
        super().__init__(unique_id, model)
        self.group = group_object
        self.id = unique_id
        self.type = "signature"
        self.delay = np.random.poisson(self.model.signature_delay) #delay between when it is triggered and when it hits the chain
        self.start_signature_process = False
        self.end_signature_process = False
        self.ownership_distr = []
        self.signature_success = False
        self.model.newest_id +=1
        self.model.newest_signature_id +=1

    def step(self):
        if not self.end_signature_process: 
            if not self.start_signature_process:
                print("Starting signature process for signature ID = "+ str(self.id))
                self.start_signature_process =True
            elif self.delay <=0:
                temp_distr = np.zeros(self.model.num_nodes)
                print("     Checking for active nodes in randomly selected group")
                active_count = np.zeros(self.model.num_nodes)
                for node in self.group.members:
                    active_count[node.node_id] = (node.mainloop_status=="forked") #adds 1 to the index matching the node id if the node is active
                    temp_distr[node.node_id] += (node.mainloop_status=="forked") #counts the node in the group distr only if it's active

                print(sum(active_count))
                self.ownership_distr = temp_distr


                if sum(active_count)>= self.model.signature_threshold: 
                    print("         signature successful")
                    self.model.unsuccessful_signature_events.append(0)
                    self.signature_success = True
                else:
                    print("         signature unsuccessful")
                    self.model.unsuccessful_signature_events.append(1)
                self.end_signature_process = True
            else:
                print("Signature ID " + str(self.id) + " Delay = "+str(self.delay))
                self.delay -=1

    def advance(self):
        pass


        

