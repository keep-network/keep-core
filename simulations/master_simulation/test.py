import simpy
import simulation_components as sc

env = simpy.Environment()
a=sc.Group(env, 1, 10, [0, 0, 0, 2, 5, 1, 7, 8, 9, 10])
