package wq

import (
	"bufio"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCase(t *testing.T) {
	ticker := time.NewTicker(time.Second * 5)
	ctx, canale := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("done")
		case <-ticker.C:
			fmt.Println("ticker")
		}
	}()
	for {
		time.Sleep(time.Second * 3)
		canale()
	}
}

func TestSGF(t *testing.T) {
	newSgf := "(;SZ[19];B[pj];W[jj];B[ml];W[nh];B[pl];W[ln];B[po];W[qg];B[rl];W[mf];B[jd];W[ig];B[hk];W[gi];B[gl];W[me];B[mp];W[fg];B[fo];W[lg];B[fn];W[oq];B[op];W[ep])"
	nn, _, _ := Load(newSgf)
	nn, _ = nn.PlayColour(Point(1, 1), BLACK, false)
	nn.SetValue("C", "B")
	fmt.Println(nn.Save())
}

type analysisData struct {
	Data KataGoAnalysis `json:"data"`
}

type resData struct {
	Code int `json:"code"`
	Data struct {
		Id             string `json:"id"`
		IsDuringSearch bool   `json:"isDuringSearch"`
		MoveInfos      []struct {
			Lcb           float64  `json:"lcb"`
			Move          string   `json:"move"`
			Order         int      `json:"order"`
			Prior         float64  `json:"prior"`
			Pv            []string `json:"pv"`
			ScoreLead     float64  `json:"scoreLead"`
			ScoreMean     float64  `json:"scoreMean"`
			ScoreSelfplay float64  `json:"scoreSelfplay"`
			ScoreStdev    float64  `json:"scoreStdev"`
			Utility       float64  `json:"utility"`
			UtilityLcb    float64  `json:"utilityLcb"`
			Visits        int      `json:"visits"`
			Winrate       float64  `json:"winrate"`
		} `json:"moveInfos"`
		Ownership []float64 `json:"ownership"`
		RootInfo  struct {
			CurrentPlayer string  `json:"currentPlayer"`
			ScoreLead     float64 `json:"scoreLead"`
			ScoreSelfplay float64 `json:"scoreSelfplay"`
			ScoreStdev    float64 `json:"scoreStdev"`
			SymHash       string  `json:"symHash"`
			ThisHash      string  `json:"thisHash"`
			Utility       float64 `json:"utility"`
			Visits        int     `json:"visits"`
			Winrate       float64 `json:"winrate"`
		} `json:"rootInfo"`
		TurnNumber int `json:"turnNumber"`
	} `json:"data"`
	Message string `json:"message"`
}

func TestBench(t *testing.T) {
	count := 0
	root := "/Users/zhangchaolong/project/computer-go-dataset/AI/"
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			thisSGF := string(content)
			kmre := regexp.MustCompile(`KM\[7\.5\]`)
			if kmre.MatchString(thisSGF) {
				rere := regexp.MustCompile(`[BW]\+([0-9]+(\.[0-9]+)?)`)
				if rere.MatchString(thisSGF) {
					scorere := regexp.MustCompile(`([BW])\+([0-9]+(\.[0-9]+)?)`)
					matches := scorere.FindAllStringSubmatch(thisSGF, -1)
					var score float64
					var err error
					for _, match := range matches {
						color := match[1]
						scoreStr := match[2]
						score, err = strconv.ParseFloat(scoreStr, 64)
						if err != nil {
							fmt.Println("解析得分值出错")
							break
						}
						if color == "W" {
							score = -score
						}
					}
					if score != 0 {
						count++
						checkScore(thisSGF, score)
					}
				}
			}
		}

		return nil
	})
	tom := "/Users/zhangchaolong/project/computer-go-dataset/Tom/out/"
	_ = filepath.Walk(tom, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			thisSGF := string(content)
			kmre := regexp.MustCompile(`KM\[7\.5\]`)
			if kmre.MatchString(thisSGF) {
				rere := regexp.MustCompile(`[BW]\+([0-9]+(\.[0-9]+)?)`)
				if rere.MatchString(thisSGF) {
					scorere := regexp.MustCompile(`([BW])\+([0-9]+(\.[0-9]+)?)`)
					matches := scorere.FindAllStringSubmatch(thisSGF, -1)
					var score float64
					var err error
					for _, match := range matches {
						color := match[1]
						scoreStr := match[2]
						score, err = strconv.ParseFloat(scoreStr, 64)
						if err != nil {
							fmt.Println("解析得分值出错")
							break
						}
						if color == "W" {
							score = -score
						}
					}
					if score != 0 {
						count++
						checkScore(thisSGF, score)
					}
				}
			}
		}

		return nil
	})
	file1, err := os.Open("/Users/zhangchaolong/project/computer-go-dataset/Professional/pro2000+.txt")
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}
	defer file1.Close()

	// 创建Scanner对象
	scanner1 := bufio.NewScanner(file1)
	for scanner1.Scan() {
		thisSGF := scanner1.Text()
		kmre := regexp.MustCompile(`KM\[7\.5\]`)
		if kmre.MatchString(thisSGF) {
			rere := regexp.MustCompile(`[BW]\+([0-9]+(\.[0-9]+)?)`)
			if rere.MatchString(thisSGF) {
				scorere := regexp.MustCompile(`([BW])\+([0-9]+(\.[0-9]+)?)`)
				matches := scorere.FindAllStringSubmatch(thisSGF, -1)
				var score float64
				var err error
				for _, match := range matches {
					color := match[1]
					scoreStr := match[2]
					score, err = strconv.ParseFloat(scoreStr, 64)
					if err != nil {
						fmt.Println("解析得分值出错")
						break
					}
					if color == "W" {
						score = -score
					}
				}
				if score != 0 {
					count++
					checkScore(thisSGF, score)
				}
			}
		}
	}
	file2, err := os.Open("/Users/zhangchaolong/project/computer-go-dataset/Professional/pro1940-1999.txt")
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}
	defer file2.Close()

	// 创建Scanner对象
	scanner2 := bufio.NewScanner(file2)
	for scanner2.Scan() {
		thisSGF := scanner2.Text()
		kmre := regexp.MustCompile(`KM\[7\.5\]`)
		if kmre.MatchString(thisSGF) {
			rere := regexp.MustCompile(`[BW]\+([0-9]+(\.[0-9]+)?)`)
			if rere.MatchString(thisSGF) {
				scorere := regexp.MustCompile(`([BW])\+([0-9]+(\.[0-9]+)?)`)
				matches := scorere.FindAllStringSubmatch(thisSGF, -1)
				var score float64
				var err error
				for _, match := range matches {
					color := match[1]
					scoreStr := match[2]
					score, err = strconv.ParseFloat(scoreStr, 64)
					if err != nil {
						fmt.Println("解析得分值出错")
						break
					}
					if color == "W" {
						score = -score
					}
				}
				if score != 0 {
					count++
					checkScore(thisSGF, score)
				}
			}
		}
	}
	fmt.Println(count)
}

func checkScore(sgfStr string, score float64) {
	gameId := uint(time.Now().UnixNano())
	b, _, _ := Load(sgfStr)
	board := b.Board()
	x := board.GetKataGoAnalysisData(gameId)
	request := analysisData{Data: *x}

	var response resData
	resp, err := req.
		R().
		SetBody(request).
		SetSuccessResult(&response).
		Post("http://10.18.28.10:37666/api/katago-analysis/analysis")
	if err != nil {
		return
	}
	if resp.IsSuccessState() {
		if response.Code == 0 {
			var bScore float64
			var wScore float64
			var controversyCount int
			for i := range response.Data.Ownership {
				if response.Data.Ownership[i] >= 0.5 {
					bScore++
				} else if response.Data.Ownership[i] < 0.5 && response.Data.Ownership[i] >= 0.25 {
					bScore += 0.5
				} else if response.Data.Ownership[i] > -0.5 && response.Data.Ownership[i] <= -0.25 {
					wScore += 0.5
				} else if response.Data.Ownership[i] <= -0.5 {
					wScore++
				} else {
					controversyCount++
				}
			}
			controversyAvgCount := controversyCount / 2
			bScore += float64(controversyAvgCount)
			wScore += float64(controversyAvgCount)
			controversyCount = controversyCount % 2
			if response.Data.RootInfo.CurrentPlayer == "B" {
				bScore += float64(controversyCount)
			} else {
				wScore += float64(controversyCount)
			}
			nowScore := bScore - wScore - board.KM
			if nowScore-1 != score && nowScore != score {
				fmt.Printf("bScore %f wScore %f km %f endScore %f sgfScore %f step %d sgf %s\n",
					bScore, wScore, board.KM, nowScore, score, board.Step, strings.ReplaceAll(sgfStr, "\n", ""))
				return
			} else {
				fmt.Println("same")
				return
			}
		} else {
			fmt.Printf("接口错误")
			return
		}
	}
}
