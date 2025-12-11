package graph

import "testing"

func maxAreaOfIsland(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    var dfs func(i, j int) int
    dfs = func(i, j int) int {
        if i<0||i>=m||j<0||j>=n||grid[i][j]==0 { return 0 }
        grid[i][j] = 0
        return 1 + dfs(i+1,j) + dfs(i-1,j) + dfs(i,j+1) + dfs(i,j-1)
    }
    best := 0
    for i:=0;i<m;i++ { for j:=0;j<n;j++ {
        if grid[i][j]==1 { area := dfs(i,j); if area>best { best = area } }
    }}
    return best
}

func TestMaxAreaOfIsland(t *testing.T) {
    g := [][]int{{0,0,1,0,0,0,0,1,0,0,0,0,0},{0,0,0,0,0,0,0,1,1,1,0,0,0},{0,1,1,0,1,0,0,0,0,0,0,0,0},{0,1,0,0,1,1,0,0,1,0,1,0,0},{0,1,0,0,1,1,0,0,1,1,1,0,0},{0,0,0,0,0,0,0,0,0,0,1,0,0},{0,0,0,0,0,0,0,1,1,1,0,0,0},{0,0,0,0,0,0,0,1,1,0,0,0,0}}
    if maxAreaOfIsland(g) != 6 { t.Fatalf("expected 6") }
}

