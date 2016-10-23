package gowa


const (
	PERM_READ_ONLY = iota
	PERM_RW
)

type User struct{

	Email		string
	Passwd		string
	Permission	uint8
}

func (u User) IsValid() bool{
	if u.Email != "" && u.Passwd != ""{
		return true;
	}
	return false;
}

func (u User) Create() error{
	db, _:= GM.GetSession();

	err := db.Create(&u).Error;
	if err != nil {
		return err;
	};
	return nil;
}

func (u User) Delete() error{
	db, _:= GM.GetSession();

	err := db.Delete(&u).Error;
	if err != nil {
		return err;
	};
	return nil;
}