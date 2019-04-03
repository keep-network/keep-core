from mesa import Agent, Model
from mesa.time import RandomActivation

class Node(Agent):
       """ Node """
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)

class Virtial_Stakers(Agent):
        """ Virtual Stakers aka Tickets """
       def __init__(self, unique_id, model):
        super().__init__(unique_id, model)

class Group(Agent):
       """ A Group """
    def __init__(self, unique_id, model):
        super().__init__(unique_id, model)

class Beacon_Model(Model):
    """The model"""
    def __init__(self, N):