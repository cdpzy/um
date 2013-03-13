package pg

import (
	"github.com/golibs/um"
	"math/rand"
	"testing"
)

// TestCreateValidUser makes sure a the user manager can handle a valid user creation
func TestCreateValidUser(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)

	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	status := rand.Int31()
	user, err := manager.CreateUser("testcreatevaliduser", "test@example.com", status)
	if err != nil {
		panic(err)
	}
	if user == nil {
		t.Error("CreateUser did not return a valid User structure")
		t.FailNow()
	}
	if user.Id() == 0 {
		t.Error("CreateUser did not give the User structure a valid ID number")
		t.Fail()
	}
	if user.UserName() != "testcreatevaliduser" {
		t.Errorf("UserName is '%s' instead of 'testcreatevaliduser' as expected", user.UserName())
		t.Fail()
	}
	if user.EmailAddr() != "test@example.com" {
		t.Errorf("EmailAddr is '%s' instead of 'test@example.com'", user.EmailAddr())
		t.Fail()
	}
	if user.Status() != status {
		t.Errorf("Status is '%d' instead of '%d' as expected", user.Status(), status)
		t.Fail()
	}

	props := map[string]interface{}{
		"id":         user.Id(),
		"user_name":  user.UserName(),
		"email_addr": user.EmailAddr(),
		"status":     user.Status(),
	}
	if !assertRecord(session, "um_users", props) {
		t.Error("Cannot find the coresponding user record in the database")
		t.FailNow()
	}
}

// TestUserNameExists makes sure UserNameExists returns true iff there's a user record in the database
func TestUserNameExists(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	exists, err := manager.UserNameExists("fixtureuser1")
	if err != nil {
		panic(err)
	}
	if exists != true {
		t.Error("fixtureuser1 exists in the database")
	}
	exists, err = manager.UserNameExists("FIXTUREUSER2")
	if err != nil {
		panic(err)
	}
	if exists != true {
		t.Error("FIXTUREUSER2 (case insensitive) exists in the database")
	}
	exists, err = manager.UserNameExists("fixtureuser3doesnotexist")
	if err != nil {
		panic(err)
	}
	if exists != false {
		t.Error("fixtureuser3doesnotexist does not exist in the database")
	}
}

// TestFindByIdSuccess makes sure FindById returns the correct existing user
func TestFindByIdSuccess(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1)
	if err != nil {
		panic(err)
	}
	if user == nil {
		t.Error("FindById did not return a valid user for ID #1")
		t.FailNow()
	}
	if id := user.Id(); id != 1 {
		t.Errorf("User ID returns %d, but 1 is expected", id)
		t.Fail()
	}
	if userName := user.UserName(); userName != "fixtureuser1" {
		t.Errorf("User name returns '%s', but 'fixtureuser1' is expected", userName)
		t.Fail()
	}
}

// TestFindByIdFail makes sure FindById returns nil and an error when the user does not exist
func TestFindByIdFail(t *testing.T) {
	session := testSetup()
	defer testTearDown(session)
	manager, err := um.Open("postgres", c_testDns)
	if err != nil {
		panic(err)
	}
	defer manager.Close()

	user, err := manager.FindById(1000)
	if !(err != nil && user == nil) {
		t.Error("FindById need to return an error and a nil-um.User")
		t.FailNow()
	}
}