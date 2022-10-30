Idea of incremental map updates: update information incrementally after each move:

Used in:
* Strike detection
* Field detection

IDEAS:
* Track `unoccupied cells` inside a separate arary inside `GameState`, not to recalculate it all the time in players
* Precalculate circle mask and use it instead of calculating it every time
