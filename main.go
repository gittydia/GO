package main

import "practice/chatting"

func main() {
	//for loop calling the functions
    var u1 chatting.SigninU
    var u2 chatting.LoginU
    
    chatting.Signin(&u1)
    chatting.Login(&u2, &u1)
}