import simpy
import simulation_components as sc

def event_test(env,event):
    while True:
        print("yeileding")
        yield env.process(event)



