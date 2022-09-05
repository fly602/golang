package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Name     string
	Age      uint16
	BirthDay uint32
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)
}

func main() {

	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:53306)/mytest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln("db connect err,", err)
	}
	log.Println("db connect success...")
	defer db.Close()
	db.LogMode(true)

	u := User{
		Name:     "fly",
		Age:      18,
		BirthDay: uint32(time.Now().Unix()),
	}
	log.Printf("Get u=%+v\n", u)
	if !db.HasTable(&u) {
		db.CreateTable(&u)
		log.Println("NewRecord", db.NewRecord(u))
		db.Create(&u)
	}
	log.Printf("Get u=%+v\n", u)
	log.Println("NewRecord", db.NewRecord(u))
	var fly = new(User)
	db.Where("Name = ?", "fly").First(fly)
	log.Println("Get first", fly)
	var count int
	// if db.Model(&User{}).Where("Name = ?", "cc").Count(&count); count == 0 {
	if db.Model(&User{}).Where(&User{Name: "cc"}).Count(&count); count == 0 {
		log.Println("Get count=", count)
		cc := User{
			Name:     "fly",
			Age:      18,
			BirthDay: uint32(time.Now().Unix()),
		}
		scope := db.NewScope(&cc)
		scope.SetColumn("Name", "cc")
		db.Create(&cc)
	}
	users := new([]User)
	db.Find(users)
	log.Println("Lookup users", users)
	users = new([]User)
	db.Where("Name in (?)", []string{"cc", "fly"}).Find(users)
	log.Println("Lookup users", users)
	users = new([]User)
	db.Where("Name like ?", "f%").Find(users)
	log.Println("Lookup table like 'f' =", users)
	users = new([]User)
	db.Where("ID between ? and ?", 2, 3).Find(users)
	log.Println("Lookup ID between 2 and 3", users)
	db.Where(map[string]interface{}{"Name": "fly", "Age": 18}).Find(users)
	log.Println("Lookup Name=fly, Age=18", users)
	db.Where([]int{1, 2, 3, 4}).Find(users)
	log.Println("Lookup ID in 1,2,3,4=", users)
	db.Where([]string{"1123", "2"}).Find(users)
	user := new(User)
	db.First(user, "ID = ?", 3)
	db.Find(users, "ID in (?)", []int{2, 3, 1})
	db.Set("gorm:query_option", "FOR UPDATE").Find(user, 2)
	user = new(User)
	db.FirstOrInit(user, &User{
		Model: gorm.Model{
			ID: 10,
		},
		Name: "fsz"})
	log.Println(user)
	db.Attrs(User{Age: 2}).FirstOrInit(user)
	log.Println(user)
	user.ID = 1
	user.Name = "fly"
	db.Assign(User{Age: 30}).FirstOrInit(user)
	log.Println(user)
	user = new(User)
	user.Name = "fsz"
	user.Age = 2
	user.ID = 3
	db.Attrs(User{BirthDay: uint32(time.Now().Unix())}).FirstOrCreate(user)
	users = new([]User)
	db.Where("age > (?)", db.Table("users").Where("id > ?", 1).Select("AVG(age)").QueryExpr()).Where("id > ?", 1).Find(users)
	users = new([]User)
	db.Order("age asc").Order("name desc").Limit(2).Offset(1).Find(users)
	db.Table("users").Count(&count)
	log.Println(users, count)
	rows, err := db.Table("users").Select("date(created_at) as date,sum(age) as ages").Group("created_at").Rows()
	if err != nil {
		log.Fatalln("Group test err,", err)
	}
	defer rows.Close()
	for rows.Next() {

	}
	names := new([]string)
	db.Find(&users).Pluck("name", names)
	log.Println("Get names=", names)

	result := new([]struct {
		Name string
		age  uint16
	})
	db.Table("users").Select("name,age").Scan(result)
	log.Println(result)
	user = new(User)
	user.Name = "fly"
	db.First(user)
	// 'delete' and 'update' test
	// user.Age = 30
	// db.Save(user).Update("age", 29)
	// log.Println(db.Model(&User{Model: gorm.Model{ID: 2}}).Update("age", 28).RowsAffected)
	// user = new(User)
	// user.Name = "fsz"
	// db.Where("Name = ?", "fsz").Delete(user)
	// db.Unscoped().Delete(&User{})
	// type Profile struct {
	// 	gorm.Model
	// 	High      uint16
	// 	User      User `gorm:"foreignkey:UserRefer"`
	// 	UserRefer string
	// }
	// db.AutoMigrate(&Profile{})
	// CreateUser(db)
	// db.DropTable(&User{})
	// db.DB().Ping()
}

func CreateUser(db *gorm.DB) (err error) {
	// 注意在事务中要使用 tx 作为数据库句柄
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}
	if err = tx.Create(&User{
		Name:     "wyh",
		Age:      50,
		BirthDay: uint32(time.Now().Unix()),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Create(&User{
		Name:     "fxm",
		Age:      55,
		BirthDay: uint32(time.Now().Unix()),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if tx.Commit().Error != nil {
		tx.Rollback()
		return err
	}
	return nil
}
