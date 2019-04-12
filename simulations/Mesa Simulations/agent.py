from mesa import Agent, Model
from mesa.time import RandomActivation
import numpy as np

class Node(Agent):
    """ Node: One hardware device used to stake tokens on the network. 
    Each node will create virtual stakers proportional to
    the number of tokens owned by the node 
    
    Attributes:
    unique_id: unique int
    token_amount: int value of tokens staked by node
    node_status: status of node can be - not connected, connected
    """
    def __init__(self, unique_id, model, tickets, failure_percent, death_percent):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "node"
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
    def __init__(self, unique_id, model, members, expiry):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "group"
        self.members = members
        self.last_signature = "none"
        self.status = "Active"
        self.expiry = expiry # of steps before expiration
        self.timer = self.model.timer
        self.model.newest_id +=1

    def step(self):
        """ At each step check if the group as expired """
        self.expiry -=1
        #print('group ID '+ str(self.id) + ' expiry ' + str(self.expiry))
        if self.expiry == 0: 
            self.status = "Expired"

    def advance(self):
        pass

class Signature(Agent):
    def __init__(self, unique_id, model, group_object):
        super().__init__(unique_id, model)
        self.group = group_object
        self.delay = np.random.poisson(6) #delay between when it is triggered and when it hits the chain
        self.start_signature_process = False

    def step(self):
        if not self.start_signature_process:
            print("Starting signature process:")
            self.start_signature_process =True
        elif self.delay <=0:
            print("     Checking for active nodes in randomly selected group")
            active_count = []
            for node in self.group.members:
                active_count.append(node.mainloop_status=="forked") 
            
            print(sum(active_count))

            if sum(active_count)>= self.model.signature_threshold: 
                print("         signature successful")
                self.model.unsuccessful_signature_events.append(0)
            else:
                print("         signature unsuccessful")
                self.model.unsuccessful_signature_events.append(1)
        else:
            self.delay -=1


        

