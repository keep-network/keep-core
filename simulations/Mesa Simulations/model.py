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

        #create nodes
        for i in range(nodes):
            node = agent.Node(i, self, tickets,  )
            self.schedule.add(node)
        
    def step(self):
        '''Advance the model by one step'''
        #generate relay requests
        self.schedule.step()
