# 微信跳一跳辅助工具

### 依赖

- adb
- 安卓手机

### 使用方法

1. 在bin目录中直接下载pc的对应版本，mac和linux请给予可执行权限
2. 手机开启开发者模式，通过usb连接到pc上
3. 手机打开游戏并开始游戏
4. pc在命令行直接运行可执行程序 `./jump-xxx-xxx -s 1.32`
5. s参数为跳跃速度系数，需要根据你的手机分辨率做调整，我的手机是1080*2220 用1.32这个默认数值刚好合适

### 参考思路

python版 https://github.com/wangshub/wechat_jump_game

kotlin版 https://github.com/uglyer/wechat_jump_ai_kotlin

golang版 https://github.com/faceair/youjumpijump
