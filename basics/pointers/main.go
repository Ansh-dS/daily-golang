package main
import "fmt"

func main(){
	var num int= 5
	fmt.Println(num)

	var pointer *int= &num
	fmt.Println(pointer)
	fmt.Println(&pointer)

}
