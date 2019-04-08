from mesa import Model
from mesa.time import SimultaneousActivation
import agent
import numpy as np

class Beacon_Model(Model):
    """The model"""
    def __init__(self, nodes, token_distribution):
        self.num_nodes = nodes
        self.schedule = SimultaneousActivation(self)
        self.relay_request = []

        #create relay trigger
        relay_trigger = agent.Relay_Trigger(1,self)
        self.schedule.add(relay_trigger)

        #create ticket distribution
        tickets = np.ones(nodes)*10 #this can be any distribution we decide

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, tickets[i], relay_trigger)
            self.schedule.add(node)
        
    def step(self):
        '''Advance the model by one step'''
        #generate relay requests
        self.schedule.step()
