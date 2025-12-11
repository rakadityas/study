package graph

import "testing"

func numIslands(grid [][]byte) int {
    if len(grid)==0 { return 0 }
    m, n := len(grid), len(grid[0])
    dirs := [][2]int{{1,0},{-1,0},{0,1},{0,-1}}
    var dfs func(i, j int)
    dfs = func(i, j int) {
        if i<0||i>=m||j<0||j>=n||grid[i][j]!='1' { return }
        grid[i][j] = '0'
        for _, d := range dirs { dfs(i+d[0], j+d[1]) }
    }
    count := 0
    for i:=0;i<m;i++ { for j:=0;j<n;j++ {
        if grid[i][j]=='1' { count++; dfs(i,j) }
    }}
    return count
}

func TestNumIslands(t *testing.T) {
    g := [][]byte{
        []byte("11110"),
        []byte("11010"),
        []byte("11000"),
        []byte("00000"),
    }
    if numIslands(g) != 1 { t.Fatalf("expected 1") }
}

