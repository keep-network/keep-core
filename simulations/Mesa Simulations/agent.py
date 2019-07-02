from mesa import Agent, Model
from mesa.time import RandomActivation
import numpy as np

class Node(Agent):
    """ Node: One hardware device used to stake tokens on the network. 
    Each node will create virtual stakers proportional to
    the number of tokens owned by the node """
    def __init__(self, unique_id, node_id, model, tickets, 
    failure_percent, death_percent, node_connection_delay,
    node_mainloop_connection_delay, misbehaving_nodes):
        super().__init__(unique_id, model)
        self.id = unique_id
        self.type = "node"
        self.node_id = node_id
        self.num_tickets = int(tickets)
        self.ticket_list = []
        self.connection_status = "not connected" #change later to event - currently used for node failure process
        self.mainloop_status = "not forked"
        self.stake_status = "not staked"
        self.connection_delay = np.random.randint(0,node_connection_delay) #uniform randomly assigned connection delay step value
        self.mainloop_fork_delay = np.random.randint(0,node_mainloop_connection_delay) #uniform randomly assigned connection delay step value
        self.timer = self.model.timer
        self.node_failure_percent = failure_percent
        self.node_death_percent = death_percent
        self.failure = False
        self.death = False
        self.dkg_misbehaving = False
        self.malicious = np.random.randint (0,100) < misbehaving_nodes  # sets the node as malicious baed on the misbehaving % in model parameters
    
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
        self.death = np.random.randint (0,100) < self.node_death_percent

        #disconnect the node if failure occurs
        if self.failure or self.death:
            self.node_disconnect()

        #refresh active nodes list
        self.model.refresh_connected_nodes_list()

    def advance(self):
        pass


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
        self.status = "dkg" # status types: dkg, compromised, active, expired
        self.expiry = expiry # of steps before expiration
        self.timer = self.model.timer
        self.model.newest_id +=1
        self.model.newest_group_id +=1
        self.ownership_distr = 0
        self.malicious_percent = 0
        self.offline_percent = 0
        self.compromised_percent = 0
        self.process_complete = False
        self.dkg_block_delay = self.model.dkg_block_delay

        self.calculate_ownership_distr()


    def step(self):
        """ block delay to simulate the multiple DKG steps"""
        if self.status == "dkg":
            if self.dkg_block_delay>=0:
                self.dkg_block_delay -=1 # counts down the block delay
                # based on DKG process we check for missing/malicious nodes at stage 11
                if self.dkg_block_delay >=11:
                    self.offline_percent = self.calculate_offline()/sum(self.ownership_distr) # calculates % nodes offline during dkg
                    self.compromised_percent = self.malicious_percent + self.offline_percent
            else:
                self.status = "active"
        elif self.status == "active":
            """ At each step check if the group has expired """
            self.expiry -=1
            if self.expiry <= 0: 
                self.status = "expired"


        # refresh the active group list
        self.model.refresh_active_group_list()
        
    def advance(self):
        pass

    def calculate_ownership_distr(self):
        temp_distr = np.zeros(self.model.num_nodes)
        temp_malicious_count = 0
        for node in self.members:    
            temp_distr[node.node_id] +=1 # increments by 1 for each node index everytime it exists in the member list, at each step
            if node.malicious:
                temp_malicious_count +=1
        self.malicious_percent = temp_malicious_count/sum(temp_distr)
        self.ownership_distr = temp_distr


    def calculate_offline(self):
        offline_count = 0
        for node in self.members:
            if node.mainloop_status == "not forked": 
                offline_count +=1
            
        return offline_count

    

class Signature(Agent):
    def __init__(self, unique_id, signature_id, model, group_object):
        super().__init__(unique_id, model)
        self.group = group_object
        self.id = unique_id
        self.signature_id = signature_id
        self.type = "signature"
        self.status = "started"
        self.delay = np.random.poisson(self.model.signature_delay) #delay between when it is triggered and when it hits the chain
        self.ownership_distr = []
        self.model.newest_id +=1 # increments the model agent ID by 1 after a new signature is created
        self.model.newest_signature_id +=1 #increments the signature ID by one after a new signature is created
        self.signature_process_complete = False
        self.block_delay_complete = False
        self.dominator_id = -1
        self.dominator_percent = 0
        self.offline_percent = 0
        self.signature_failure = False

    def step(self):
        #signature
        if not self.block_delay_complete:
            if self.delay >0:
                self.delay -=1
            else :
                self.block_delay_complete = True
        elif not self.signature_process_complete :
            self.signature_process()
            self.signature_process_complete = True
            self.status = "complete"

    def advance(self):
        pass

    def signature_process(self):
        # Calculates ownership data just before the signature is complete
        temp_signature_distr = np.zeros(self.model.num_nodes)
        ("Signature ID "+str(self.signature_id)+" performing signature")

        self.model.refresh_connected_nodes_list()
        for i,node_tickets in enumerate(self.group.ownership_distr): # checks if the node has a non-zero ownership, i is the node id
            if node_tickets > 0:
                for node in self.model.active_nodes:
                    if node.node_id == i : 
                        temp_signature_distr[i] = node_tickets
        self.ownership_distr = temp_signature_distr
        failed_list = np.array(self.group.ownership_distr)-np.array(self.ownership_distr)
        self.offline_percent = sum(failed_list)/sum(self.ownership_distr)
        self.dominator_percent = (sum(failed_list) + max(self.ownership_distr))/sum(self.group.ownership_distr) # adds the failed node virtual stakers and max node virtual stakers
        self.signature_failure = (sum(failed_list)+max(self.ownership_distr)) > (1-self.model.max_malicious_threshold)

        if self.dominator_percent*100 > self.model.max_malicious_threshold:
            self.dominator_id = np.argmax(self.ownership_distr)

        # calculate misbehaving nodes
    





        

