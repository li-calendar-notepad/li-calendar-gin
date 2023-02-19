package models

// 文件
type File struct {
	BaseModel
	UserId   uint   `json:"userId"`                       // 用户ID
	FileName string `json:"fileName" gorm:"varchar(255)"` // 文件名
	Src      string `json:"src" gorm:"varchar(2000)"`     // 储存路径
	User     User
}

// 添加一个文件记录
func (m *File) AddFile(userId uint, fileName, src string) (File, error) {
	file := File{
		UserId:   userId,
		FileName: fileName,
		Src:      src,
	}
	err := Db.Create(&file).Error

	return file, err
}
