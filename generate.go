package sudokgo

import (
	"fmt"
	"math/rand"
	"time"
)

func (s *Sudoku) generateSolution() bool {
	s.Reset()

	rand.Seed(time.Now().UnixNano())
	for y := 0; y < rowSize; y++ {
		for x := 0; x < rowSize; x++ {
			if s.possGrid[x][y].count == 0 {
				return false
			}

			z := rand.Intn(rowSize)
			for !s.possGrid[x][y].possibilities[z] {
				z = rand.Intn(rowSize)
			}

			s.Grid[x][y] = z + 1
			s.possGrid[x][y].possibilities[z] = false
			s.possGrid[x][y].count--
			s.updatePossibilities(x, y)
		}
	}

	s.Print()
	return true
}

func (s *Sudoku) Generate(target int) bool {
	var puzzle [rowSize][rowSize]int
	retScore := 0
	tries := 0

	broken := false
	passes := false
	multiples := 0
	for !passes {
		z := 1
		for !s.generateSolution() {
			z++
		}
		fmt.Println(z, "attempts to get a solution")

		/* generate was OK, now eliminate */
		retScore = 0
		tries = 0
		soln := s.Grid
		broken = false
		for !passes && tries < 25 {
			/* reset for another pass */
			s.Grid = soln
			s.Print()
			s.loadPossGrid()
			tries++

			/* precull some cells */
			for y := 0; y < rowSize; y++ {
				for x := 0; x < rowSize; x++ {
					if rand.Intn(3) == 0 {
						s.Grid[x][y] = -1
						for z := 0; z < rowSize; z++ {
							s.possGrid[x][y].possibilities[z] = true
						}
						s.possGrid[x][y].count = rowSize
					}
				}
			}

			for retScore < target {
				x := rand.Intn(rowSize)
				y := rand.Intn(rowSize)

				if s.Grid[x][y] == -1 {
					continue
				}

				s.Grid[x][y] = -1
				for z := 0; z < rowSize; z++ {
					s.possGrid[x][y].possibilities[z] = true
				}
				s.possGrid[x][y].count = rowSize
				puzzle = s.Grid

				score, err := s.Solve()
				if err != nil {
					fmt.Println("UNSOLVABLE - rolling back")
					broken = true
					break
				}
				if score >= target {
					fmt.Println("Complexity met - rolling back and returning")
					passes = true
					break
				}
				/* roll back to before Solve so we can aggregate eliminations */
				s.Grid = puzzle
				retScore = score
			}

			if !passes && !broken {
				fmt.Println("Generated puzzle not complex enough", Difficulty(retScore), "trying another")
			}
		}
		if !passes {
			fmt.Println("Tried too many times, generating a new soluition")
			multiples++
		}
	}
	s.Grid = puzzle

	fmt.Println("Generated grid with difficulty", Difficulty(retScore), "after",
		tries+(multiples*25), "attempts.")
	return true
}
