// This is an example of the design pattern "Specification"
// see: https://en.wikipedia.org/wiki/Specification_pattern

package main

import (
	"fmt"
	"strings"
)

type UserType int

const (
	Personal UserType = iota
	Admin
	SuperAdmin
)

type User struct {
	Type   UserType
	Name   string
	Locked bool
}

func (u User) String() string {
	names := map[UserType]string{
		Personal:   "PERSONAL",
		Admin:      "ADMIN",
		SuperAdmin: "SUPER ADMIN",
	}
	return fmt.Sprintf("%s (Type:%v Locked:%t)", u.Name, names[u.Type], u.Locked)
}

type SpecificationUser interface {
	IsSatisfiedBy(u *User) bool

	// or use: IsSatisfiedBy(object interface{}) bool
	// to use a template in interface{} type, also redefine specifications And, Or, Not, etc.
}

//And
type AndSpecification struct {
	specs []SpecificationUser
}

func And(specs ...SpecificationUser) *AndSpecification {
	return &AndSpecification{
		specs: specs,
	}
}
func (s *AndSpecification) IsSatisfiedBy(u *User) bool {
	for _, s := range s.specs {
		if !s.IsSatisfiedBy(u) {
			return false
		}
	}
	return true
}

// Or
type OrSpecification struct {
	specs []SpecificationUser
}

func Or(specs ...SpecificationUser) *OrSpecification {
	return &OrSpecification{
		specs: specs,
	}
}
func (s *OrSpecification) IsSatisfiedBy(u *User) bool {
	for _, s := range s.specs {
		if s.IsSatisfiedBy(u) {
			return true
		}
	}
	return false
}

// Not
type NotSpecification struct {
	spec SpecificationUser
}

func Not(spec SpecificationUser) *NotSpecification {
	return &NotSpecification{
		spec: spec,
	}
}

func (s *NotSpecification) IsSatisfiedBy(u *User) bool {
	return !s.spec.IsSatisfiedBy(u)
}

//Specification type
type TypeSpecification struct {
	typ UserType
}

func (s *TypeSpecification) IsSatisfiedBy(u *User) bool {
	return s.typ == u.Type
}

//Specification name: too short
type NameLengthSpecification struct {
	l int
}

func NameShort(l int) *NameLengthSpecification {
	return &NameLengthSpecification{
		l: l,
	}
}
func (s *NameLengthSpecification) IsSatisfiedBy(u *User) bool {
	return len(u.Name) <= s.l
}

// SpecificationUserName
type NameSpecification struct {
	name string
}

func Name(name string) *NameSpecification {
	return &NameSpecification{
		name: strings.ToLower(name),
	}
}
func (s *NameSpecification) IsSatisfiedBy(u *User) bool {
	return strings.ToLower(u.Name) == s.name
}

//SpecificationLocked
type LockedSpecification struct{}

func (s *LockedSpecification) IsSatisfiedBy(u *User) bool {
	return u.Locked
}

// Predefined rules
var (
	IsPersonal   = &TypeSpecification{typ: Personal}
	IsAdmin      = &TypeSpecification{typ: Admin}
	IsSuperAdmin = &TypeSpecification{typ: SuperAdmin}

	AnyAdmin      = Or(IsAdmin, IsSuperAdmin)
	NotAdmin      = Not(AnyAdmin)
	NotSuperAdmin = Not(IsSuperAdmin)

	IsNameShort4 = NameShort(4)

	Locked    = &LockedSpecification{}
	NotLocked = Not(Locked)

	ValidNameNotAdmin = And(Not(AnyAdmin), NotLocked, Not(IsNameShort4))
)

func UserIsSatisfiedBy(u *User, spec SpecificationUser) bool {
	return spec.IsSatisfiedBy(u)
}

func checkAccess(spec SpecificationUser, name string, handler func()) func(*User) error {
	return func(user *User) error {
		if !spec.IsSatisfiedBy(user) {
			return fmt.Errorf("%s: access denied, user: %v", name, user)
		}
		fmt.Printf("%s: access granted, user: %v\n", name, user)
		handler()
		return nil
	}
}

func main() {
	user := &User{
		Type: Admin,
		Name: "Alex",
	}

	supAdmin := &User{
		Type: SuperAdmin,
		Name: "SuperAlex",
	}

	fmt.Printf("%s: Any Admin? %v\n", user, UserIsSatisfiedBy(user, AnyAdmin))
	fmt.Printf("%s: Any SuperAdmin? %v\n", user, UserIsSatisfiedBy(user, IsSuperAdmin))

	handlerSecret := checkAccess(IsSuperAdmin, "high level", func() {
		fmt.Println("execute handlerSecret")
	})

	if err := handlerSecret(user); err != nil {
		fmt.Println(err)
	}

	if err := handlerSecret(supAdmin); err != nil {
		fmt.Println(err)
	}

	BooFooLocked := &User{
		Type:   Personal,
		Name:   "BooFooLocked",
		Locked: true,
	}
	BooFoo := &User{
		Type:   Personal,
		Name:   "BooFoo",
		Locked: false,
	}
	handlerOnlyValidUser := checkAccess(ValidNameNotAdmin, "onlyValidUser", func() {
		fmt.Println("execute handlerOnlyValidUser")
	})
	if err := handlerOnlyValidUser(user); err != nil {
		fmt.Println(err)
	}
	if err := handlerOnlyValidUser(BooFooLocked); err != nil {
		fmt.Println(err)
	}
	if err := handlerOnlyValidUser(BooFoo); err != nil {
		fmt.Println(err)
	}

}
