import "github.com/kuma807/knowledge_work_day2/displayGoroutine"
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go displayGoroutine.Watch(ctx, "testGoroutine")
	//監視したいゴールーチンをここに書く
	displayGoroutine.Show("testGoroutine")
}
