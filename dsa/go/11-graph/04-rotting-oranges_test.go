package graph

import "testing"

func orangesRotting(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    q := [][2]int{}
    fresh := 0
    for i:=0;i<m;i++ { for j:=0;j<n;j++ {
        if grid[i][j]==2 { q = append(q, [2]int{i,j}) }
        if grid[i][j]==1 { fresh++ }
    }}
    minutes := 0
    dirs := [][2]int{{1,0},{-1,0},{0,1},{0,-1}}
    for len(q)>0 && fresh>0 {
        size := len(q)
        for s:=0;s<size;s++ {
            cell := q[0]; q = q[1:]
            for _,d := range dirs {
                ni,nj := cell[0]+d[0], cell[1]+d[1]
                if ni>=0&&ni<m&&nj>=0&&nj<n&&grid[ni][nj]==1 {
                    grid[ni][nj]=2; fresh--; q = append(q, [2]int{ni,nj})
                }
            }
        }
        minutes++
    }
    if fresh>0 { return -1 }
    return minutes
}

func TestOrangesRotting(t *testing.T) {
    g := [][]int{{2,1,1},{1,1,0},{0,1,1}}
    if orangesRotting(g) != 4 { t.Fatalf("expected 4") }
}

