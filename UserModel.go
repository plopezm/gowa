package gowa


const (
	PERM_READ_ONLY = iota
	PERM_RW
)

type GowaUser struct{

	Email		string
	Passwd		string
	Permission	uint8
}

func (u GowaUser) IsValid() bool{
	if u.Email != "" && u.Passwd != ""{
		return true;
	}
	return false;
}

func (u GowaUser) Create() error{
	db, _:= GM.getSession();

	_, err := db.Insert(&u);
	if err != nil {
		return err;
	};
	return nil;
}

func (u GowaUser) Delete() error{
	db, _:= GM.getSession();

	_, err := db.Remove(&u);
	if err != nil {
		return err;
	};
	return nil;
}