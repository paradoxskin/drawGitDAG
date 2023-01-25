package main

import (
	"os"
	"fmt"
	"time"
	"math"
	"os/exec"
	"strings"
	"math/rand"
	"github.com/fogleman/gg"
)

type edge struct {
	to, next int
}
type Graph struct {
	edges []edge
	idx int
	head []int
	vis []int
	maxLevel int
	nodeInLevel [][]int
}
type node struct {
	x,y float64
	sha_1 string
	col [3]float64
	cntNodes []int
}
// some OI function
func (gp *Graph) init(n int) {
	gp.head = make([]int, n+1, n+1)
	gp.vis = make([]int, n+1, n+1)
	gp.nodeInLevel = append(gp.nodeInLevel, []int{})
	gp.edges = append(gp.edges, edge{})
}
func (gp *Graph) addedge(from, to int) {
	gp.idx++
	gp.edges = append(gp.edges, edge{to, gp.head[from]})
	gp.head[from] = gp.idx
}
func (gp *Graph) dfs(now, level int) {
	gp.vis[now] = 1
	for ; level > gp.maxLevel; gp.maxLevel++ {
		gp.nodeInLevel = append(gp.nodeInLevel, []int{})
	}
	gp.nodeInLevel[level] = append(gp.nodeInLevel[level], now)
	for i := gp.head[now]; i != 0; i = gp.edges[i].next {
		to := gp.edges[i].to
		if gp.vis[to] == 1 {
			continue
		}
		gp.dfs(to, level + 1)
	}
}
// some OS command i/o
func whereAmI() string {
	out, _ := exec.Command("pwd").Output()
	return strings.TrimSpace(string(out))
}
func getAllGitObjs() []string {
	var gitObjs []string
	out, _ := exec.Command("bash", "-c", "ls .git/objects").Output()
	dirs := strings.Split(string(out), "\n")
	for _, x := range(dirs) {
		out, _ = exec.Command("bash", "-c", fmt.Sprintf("ls .git/objects/%s", x)).Output()
		for _, y := range(strings.Split(string(out), "\n")) {
			if y == "" {
				continue
			}else if len(x+y) == 40 {
				// files in .git
				gitObjs = append(gitObjs, x+y)
			}else if len(y) - 4 > 0 && y[len(y) - 4:] == ".idx" {
				// commits in pack
				out, _ = exec.Command("bash", "-c", fmt.Sprintf("git verify-pack -v .git/objects/pack/%s|grep commit|cut -f 1 -d \" \"", y)).Output()
				for _, obj := range(strings.Split(string(out)[:len(string(out))], "\n")) {
					if obj == "" {
						continue
					}
					gitObjs = append(gitObjs, obj)
				}
			}
		}
	}
	return gitObjs
}
// judge if it is commit
func isCommit(sh string) int {
	cmd := fmt.Sprintf("git cat-file -t %s", sh)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	if string(out) == "commit\n" {
		return 1
	}
	return 0
}
func findParentsSha(sh string) []string {
	cmd := fmt.Sprintf("git cat-file -p %s|grep ^parent|cut -f 2 -d \" \"", sh)
	out, _ := exec.Command("bash", "-c", cmd).Output()
	return strings.Split(string(out), "\n")
}
// draw the picture
func drawArrow(dc *gg.Context, x1, y1, x2, y2 float64) {
	dc.SetRGB(0, 0, 0.8)
	dc.DrawLine(x1, y1, x2, y2)
	dc.SetLineWidth(2)
	dc.Stroke()
	cos30 := math.Cos(math.Pi / 6)
	sin30 := math.Sin(math.Pi / 6)
	x0 := x2 - x1
	y0 := y2 - y1
	linelen := math.Sqrt(x0 * x0 + y0 * y0)
	x0 /= linelen / 7
	y0 /= linelen / 7
	dc.DrawLine(x1, y1, x1 + x0 * cos30 - y0 * sin30, y1 + x0 * sin30 + y0 * cos30)
	dc.DrawLine(x1, y1, x1 + x0 * cos30 + y0 * sin30, y1 - x0 * sin30 + y0 * cos30)
	dc.SetLineWidth(2)
	dc.Stroke()
}
func drawNode(dc *gg.Context, aNode *node){
	(*aNode).col = [3]float64{rand.Float64(), rand.Float64(), rand.Float64()} 
    dc.SetRGB(aNode.col[0], aNode.col[1], aNode.col[2])
	dc.DrawCircle(aNode.x, aNode.y, 7)
    dc.Fill()
}
func drawText(dc *gg.Context, aNode *node){
    dc.SetRGB(aNode.col[0], aNode.col[1], aNode.col[2])
	dc.DrawString(aNode.sha_1, aNode.x, aNode.y - 10)
	dc.Fill()
}
func drawConnect(dc *gg.Context, startNode, endNode *node){
	dx := endNode.x - startNode.x
	dy := endNode.y - startNode.y
	dLen := math.Sqrt(dx * dx + dy * dy)
	drawArrow(dc, startNode.x + dx / dLen * 10, startNode.y + dy / dLen * 10, endNode.x - dx / dLen * 10, endNode.y - dy / dLen * 10)
}

func main() {

	// step.1 input the abs path of git resp dir
	//var dirct string
	//fmt.Scanf("%s", &dirct)

	// step.2 change dir to input
	now := whereAmI()
	//os.Chdir(dirct)
	os.Chdir("/home/paradoxd/.vim")
	//os.Chdir("/tmp/tmpp")

	// step.3 get all gitobjs in .git
	allGitObjs := getAllGitObjs()

	// step.4 & 5 find all cmtObjs build a DAG for cmtObjs
	// may be add the branch info?
	sha2ID := make(map[string]int)
	ID2sha := make(map[int]string)
	id := 0

	for _, x := range(allGitObjs) {
		if isCommit(x) == 1 {
			id++
			sha2ID[x] = id
			ID2sha[id] = x
		}
	}
	var start int
	var gp Graph
	nodes := make([]node, 1, 1)
	gp.init(id)
	for i := 1; i <= id; i++ {
		sha := ID2sha[i]
		parentsSha := findParentsSha(sha)
		nodes = append(nodes, node{sha_1: sha[:6]})
		for _, parentSha := range(parentsSha) {
			parentId := sha2ID[parentSha]
			gp.addedge(parentId, i)
			if parentSha == "" {
				continue
			}
			if parentSha == "" {
				continue
			}
			nodes[i].cntNodes = append(nodes[i].cntNodes, parentId)
		}
		if len(parentsSha) == 0 {
			start = i
		}
	}
	gp.dfs(start, 0)
	fmt.Println(gp.nodeInLevel)
	for i := 1; i <= id; i++ {
		fmt.Println(ID2sha[i][:6], " -> ", findParentsSha(ID2sha[i]))
	}
	// step.6 分配位置
	for i := 1; i <= gp.maxLevel; i++ {
		num := len(gp.nodeInLevel[i])
		piece := math.Pi / float64(2 * (num + 1))
		for id, node := range(gp.nodeInLevel[i]) {
			nodes[node].x = 30 * float64(i) * math.Cos(piece * float64(id + 1))
			nodes[node].y = 30 * float64(i) * math.Sin(piece * float64(id + 1))
		}
	}

	// step.7 draw the DAG
	os.Chdir(now)
	rand.Seed(time.Now().Unix())
	dc := gg.NewContext(800, 800)
    dc.SetRGB(0.9, 0.9, 0.9)
	dc.LoadFontFace("/usr/share/fonts/SourceCodePro/Sauce Code Pro Black Nerd Font Complete Mono.ttf", 11)
	dc.Clear()
	for i := 1; i < len(nodes); i++ {
		drawNode(dc, &nodes[i])
		for _, j := range(nodes[i].cntNodes) {
			drawConnect(dc, &nodes[j], &nodes[i])
		}
	}
	for i := 1; i < len(nodes); i++ {
		drawText(dc, &nodes[i])
	}
	filename := fmt.Sprintf("%d.png", time.Now().Unix())
	if err := dc.SavePNG(filename);err != nil {
		fmt.Print(err)
	}
}
