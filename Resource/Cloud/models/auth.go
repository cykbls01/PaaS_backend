package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	ID 				int 		`gorm:"primary_key" json:"id"`
	Username 		string 		`json:"username"`
	Password		string 		`json:"password"`
	Authority 		int 		`json:"authority"`
	Name			string 		`json:"name"`
	Email			string 		`json:"email"`
	Profile			string		`json:"profile"`
}

func IsUserExist(username string)bool{
	var user Auth
	db.Table("auth").Select("*").Where("username = ?",username).First(&user)
	return user.ID !=0
}

func CheckAuth(username, password string) int {
	var auth Auth
	db.Select("*").Where(Auth{Username : username}).First(&auth)
	err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password))
	if auth.ID > 0 {
		if err != nil{
			fmt.Println(err)
			return 0
		}
		return auth.Authority
	}
	return 0
}

func GetUserInfo(username string) Auth{
	var auth Auth
	db.Select("*").Where("username = ?",username).First(&auth)
	return auth
}

func GetUserInfoById(userid int)Auth{
	var auth Auth
	db.Select("*").Where("id = ?",userid).First(&auth)
	return auth
}

func CheckPassword(username string,password string)bool{
	var auth Auth
	db.Select("*").Where(Auth{Username:username}).First(&auth)
	err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(password))
	if err != nil{
		return false
	}
	return true
}

func UpdatePassword(username string,newpasswd string){
	var auth Auth
	db.Select("*").Where(Auth{Username:username}).First(&auth)
	hash,_ := bcrypt.GenerateFromPassword([]byte(newpasswd),bcrypt.DefaultCost)
	auth.Password = string(hash)
	db.Save(&auth)
}

func UpdateInfo(username string,email string,profile string){
	var auth Auth
	db.Select("*").Where(Auth{Username:username}).First(&auth)
	auth.Profile = profile
	auth.Email = email
	db.Save(&auth)
}

func UpdateAuthority(username string,authority int){
	var auth Auth
	db.Select("*").Where(Auth{Username:username}).First(&auth)
	auth.Authority = authority
	db.Save(&auth)
}

func CreateUser(username string,password string,name string,authoriry int){
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var auth Auth
	auth.Username = username
	auth.Password = string(hash)
	auth.Name = name
	auth.Authority = authoriry
	db.Save(&auth)
}

func GetAllUsers()[]Auth{
	var users []Auth
	db.Table("auth").Select("username").Find(&users)
	return users
}

func DeleteUser(username string){
	db.Table("auth").Where("username = ?",username).Delete(Auth{})
}