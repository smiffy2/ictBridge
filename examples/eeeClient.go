package main

import (
	"fmt"

	"github.com/smiffy2/ictBridge"
	"time"
)

func main() {


	env := "TonyTest"
	effect := "Going good"

	client := ictBridge.CreateIctBridgeClient("35.204.80.128","7331")

	client.AddEffectListener(env)
	client.SubmitEffectMessage(env,effect)
	time.Sleep(5 * time.Second)
	var result string 
	var err error
	for result == "" {
	        time.Sleep(1 * time.Second)
		result,err = client.PollEffect(env)
		fmt.Println("Looping ...")
		if(err != nil) {
			fmt.Println(err)
			result = "Error" 
		}
	}
	fmt.Println("Result = " + result)
}

