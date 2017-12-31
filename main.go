package main

import (
	"fmt"
	"os/exec"
	"log"
	"strconv"
	"time"
	"os"
	"image/png"
	"image"
	"math"
	"flag"
	"image/color"
)

// 速度系数s
var speed = flag.Float64("s", 1.325, "速度系数，不同手机分辨率需要多次尝试")
// 执行次数
var count = 0

// 截取手机屏幕画面
func screenshot() (image.Image, error) {
	err := exec.Command("adb", "shell", "screencap", "-p", "/sdcard/wxjump.png").Run()
	if err != nil {
		return nil, fmt.Errorf("screenshot fail: %v",err.Error())
	}
	err = exec.Command("adb", "pull", "/sdcard/wxjump.png", ".").Run()
	if err != nil {
		return nil, fmt.Errorf("pull screenshot fail: %v",err.Error())
	}
	f, err := os.Open("wxjump.png")
	if err != nil {
		return nil, fmt.Errorf("open wxjump.png fail: %v",err.Error())
	}
	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode wxjump.png fail: %v",err.Error())
	}
	return img, nil
}

// 颜色判断是否相等
func isColorRight(a, b [3]int, difference float64) bool {
	for i:=0; i<3; i++ {
		if math.Abs(float64(a[i]-b[i])) > difference {
			return false
		}
	}
	return true
}

// 获得坐标点的rgb色值
func getRGB(img image.Image, x int, y int) [3]int {
	c := img.At(x, y)
	var r, g, b int
	switch c.(type) {
	case color.RGBA:
		r = int(c.(color.RGBA).R)
		g = int(c.(color.RGBA).G)
		b = int(c.(color.RGBA).B)
	case color.RGBA64:
		r = int(c.(color.RGBA64).R)
		g = int(c.(color.RGBA64).G)
		b = int(c.(color.RGBA64).B)
	case color.NRGBA:
		r = int(c.(color.NRGBA).R)
		g = int(c.(color.NRGBA).G)
		b = int(c.(color.NRGBA).B)
	case color.NRGBA64:
		r = int(c.(color.NRGBA64).R)
		g = int(c.(color.NRGBA64).G)
		b = int(c.(color.NRGBA64).B)
	}
	res := [3]int{int(r),int(g),int(b)}
	return res
}

// 获得小人的坐标
func findMe(img image.Image) ([2]int, error) {
	// 获得图片大小
	size := img.Bounds().Size()
	// 上方空白区域高度 TODO
	top := size.Y / 4 - 100
	// 保存结果
	res := [2]int{}
	// 小人的默认RGB颜色
	meColor := [3]int{54, 52, 92}
	// 找寻坐标
	for y := top; y < size.Y; y++ {
		line := 0 // 颜色匹配宽度
		for x := 0; x < size.X; x++ {
			curColor := getRGB(img, x, y)
			if isColorRight(meColor, curColor, 10) {
				line++
			} else {
				if x-line > 10 && line > 30 {
					res[0] = x - line/2
					res[1] = y
					return res, nil
				}
				line = 0
			}
		}
	}
	return res, fmt.Errorf("not found me point")
}

// 获得目标坐标
func findTarget(img image.Image, mePoint [2]int) ([2]int, error) {
	// 获得图片大小
	size := img.Bounds().Size()
	// 上方分数区域高度
	top := size.Y / 4 - 100
	// 保存结果
	res := [2]int{}
	// 找寻坐标
	for y := top; y < size.Y; y++ {
		// 颜色匹配宽度
		line := 0
		// 获得当前行的背景颜色
		bgColor := getRGB(img, size.X-10, y)
		for x := 0; x < size.X; x++ {
			curColor := getRGB(img, x, y)
			if !isColorRight(bgColor, curColor, 10) {
				line++
			} else {
				if x-line > 10 && // 不是在屏幕最左边
					line > 36 &&  // 匹配平台宽度大于35
						((x-line/2) < (mePoint[0]-20) || (x-line/2) > (mePoint[0]+20)) { // 和小人的x点不重合
					res[0] = x - line/2
					res[1] = y
					return res, nil
				}
				line = 0
			}
		}
	}
	return res, fmt.Errorf("not found target point")
}

// 获得按压时间
func getPressTime(me, target [2]int) int {
	// 勾股定理计算两点之间距离 * 速度参数
	l := math.Abs(math.Sqrt(math.Pow(float64(target[0]-me[0]), 2)+math.Pow(float64(target[1]-me[1]), 2)))
	time := l * (*speed)
	log.Printf("第%v次点击 距离:%v 速度:%v 起点:%v 目标点:%v 模拟按压时间:%vms", count, int(l), (*speed), me, target, int(time))
	return int(time)
}

// 模拟按压操作
func press(time int) error {
	_, err := exec.Command("adb", "shell", "input", "swipe", "310", "400", "310", "400", strconv.Itoa(time)).Output()
	if err != nil {
		return fmt.Errorf("press fail:%v", err.Error())
	}
	return nil
}

func main() {
	flag.Parse()
	fmt.Println("微信跳一跳辅助程序 start")
	for {
		count++
		// 截图
		img, err := screenshot()
		if err != nil {
			log.Fatal(err)
		}
		// 获得小人坐标
		mePoint, err := findMe(img)
		if err != nil {
			log.Fatal(err)
		}
		// 获得目标平台坐标
		targetPoint,err := findTarget(img, mePoint)
		if err != nil {
			log.Fatal(err)
		}
		// 计算按压时间
		pressTime := getPressTime(mePoint, targetPoint)
		// 模拟按压操作
		err = press(pressTime)
		if err != nil {
			log.Fatal(err)
		}
		// 时间间隔
		time.Sleep(1500 * time.Millisecond)
	}
}
