package main 

import (
	"fmt"
	"time"
	"context"

)


func GetCombos(numCards int, target int, dbm *DatabaseManager) []string {
//	ctx := context.Background()
 //combosID := fmt.Sprintf("%s-%s", numCards, target)
 // var combos []string
	//err := dbm.pool.QueryRow(ctx, `SELECT cards FROM combos WHERE id = "$1"`, combosID).Scan(&combos)
	combos :=[]string{"1-2-3-4", "2-3-4-4"}
	return combos	
}

func GetGameIndex(table string, gameID string, dbm *DatabaseManager) int {
	ctx:= context.Background()
	var gameIndex int
	err:= dbm.pool.QueryRow(ctx, `SELECT game_index FROM $1 WHERE id="$2"`, table, gameID).Scan(&gameIndex)
	if err!= nil {
		fmt.Println("there was an error gettign the game index")
	}
	return gameIndex
}

func CalculateScore(time time.Time) int {
	return 1
}

func CalculateSkipScore() int {
	return -1
}
