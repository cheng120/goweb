package demo

import (
	demoService "goweb/app/provider/demo"
)

type UserDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func UserModelsToUserDTOs(models []UserModel) []UserDTO {
	ret := []UserDTO{}
	for _, model := range models {
		t := UserDTO{
			ID:   model.UserId,
			Name: model.Name,
		}
		ret = append(ret, t)
	}
	return ret
}

func StudentsToUserDTOs(students []demoService.Student) []UserDTO {
	ret := []UserDTO{}
	for _, student := range students {
		t := UserDTO{
			ID:   student.ID,
			Name: student.Name,
		}
		ret = append(ret, t)
	}
	return ret
}