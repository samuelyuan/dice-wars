# Dice Wars

A remake of the classic Dice Wars Flash game by GAMEDESIGN. This version is built with Go and Ebiten.

## How to play

1. Choose the number of players (2-8). You are always player 1; the rest are CPU opponents.
2. Click **Start!** to begin.
3. On your turn, click one of your territories with **2+ dice**, then click an adjacent enemy territory to attack.
4. Dice are rolled for both sides; higher total wins. On victory, all but one die move to the conquered territory; on defeat, all but one die are lost.
5. At the end of your turn, click **End Turn** to receive reinforcements equal to your largest connected territory group.
6. Conquer the map to win.

## Controls

- **Left click** - select territory / attack / UI buttons
- **End Turn** - finish your turn and receive reinforcements
- **Auto** - let the AI play your turn
- **Menu** - return to the main menu
- **Cheat** (hidden) - click the small area near the top-left of the game screen to toggle boosted dice rolls for human players

## Build & run

Requires Go 1.26+.

```bash
go mod tidy
go run .
```
