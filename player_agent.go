package main

type PlayerAgent interface {
	MakeMove(*GameState) Offset
}
