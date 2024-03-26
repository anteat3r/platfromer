import os
import time
import keyboard
import random
from termcolor import colored

player = (0, 0)
cam = (0, 0)

# blocks = [(10, 5), (20, 0), (30, 2), (35, 4)]
blocks = [(random.randint(10, 100), (random.randint(0, 3))) for x in range(20)]
spikes = [(random.randint(10, 100), (random.randint(0, 3))) for x in range(30)]


mapa = []

player_y_speed = 0

def is_colliding():
    right = False
    left = False
    up = False
    down = False

    for block in blocks:
        if player[0]+2 == block[0] and round(player[1]) in [block[1]+y for y in range(-1, 2)]:
            right = True
    
    for block in blocks:
        if player[0]-2 == block[0] and round(player[1]) in [block[1]+y for y in range(-1, 2)]:
            left = True
    
    for block in blocks:
        if round(player[1])-2 == block[1] and player[0] in [block[0]+y for y in range(-1, 2)]:
            down = True
    
    for block in blocks:
        if round(player[1])+2 == block[1] and player[0] in [block[0]+y for y in range(-1, 2)]:
            up = True
    
    return up, right, down, left




def generate_map():
    for block in blocks:
        for x in range(-1, 2):
            for y in range(-1, 2):
                mapa.append((block[0]+x, block[1]+y))


    for x in range(-100, 101):
        for y in range(-5, 0):
            mapa.append((x, y))


def print_sceen():
    os.system("cls")
    screen = ''
    for y in range(round(cam[1])+10, round(cam[1])-10, -1):
        row = ""
        for x in range(round(cam[0])-20, round(cam[0])+20):
            if (x, y) == (round(player[0]), round(player[1])):
                row += colored("pp", "green")
            elif (x, y) in spikes:
                row += colored("AA", "red")
            elif (x, y) in mapa:
                row += "XX"
            else:
                row += "  "
        screen += row + "\n"
    print(screen)



generate_map()
while True:
    cam = (cam[0] + (player[0] - cam[0])*0.2, cam[1] + (player[1] - cam[1])*0.2)
    u, r, d, l = is_colliding()
    if (round(player[0]), round(player[1])) in spikes:
        break
    if keyboard.is_pressed("d") and not r:
        player = (player[0]+1, player[1])
    if keyboard.is_pressed("a") and not l:
        player = (player[0]-1, player[1])
    if keyboard.is_pressed("w") and (player[1] == 0 or d) and not u:
        player = (player[0], player[1]+1)
        player_y_speed = 0.8
        d = False
    
    player_y_speed -= 0.1
    
    if player[1] <= 0 or d:
        if player[1] <= 0:
            player = (player[0], 0)
        player_y_speed = 0
    player = (player[0], player[1]+player_y_speed) 
    print_sceen()
    time.sleep(0.03)
            
print("Game over")
