package models


const DbVer = 1


func Migrate() {
	_, err := db.Exec(`CREATE TABLE "config" ("id" PRIMARY KEY AUTOINCREMENT, "key" VARCHAR(64) NULL, "val" VARCHAR(1024) NULL)`)
	if err != nil {
		panic(err)
	}
}