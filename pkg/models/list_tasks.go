//  Vikunja is a todo-list application to facilitate your life.
//  Copyright 2018 Vikunja and contributors. All rights reserved.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <https://www.gnu.org/licenses/>.

package models

import "code.vikunja.io/web"

// ListTask represents an task in a todolist
type ListTask struct {
	ID            int64   `xorm:"int(11) autoincr not null unique pk" json:"id" param:"listtask"`
	Text          string  `xorm:"varchar(250)" json:"text" valid:"runelength(3|250)"`
	Description   string  `xorm:"varchar(250)" json:"description" valid:"runelength(0|250)"`
	Done          bool    `xorm:"INDEX" json:"done"`
	DueDateUnix   int64   `xorm:"int(11) INDEX" json:"dueDate"`
	RemindersUnix []int64 `xorm:"JSON TEXT" json:"reminderDates"`
	CreatedByID   int64   `xorm:"int(11)" json:"-"` // ID of the user who put that task on the list
	ListID        int64   `xorm:"int(11) INDEX" json:"listID" param:"list"`
	RepeatAfter   int64   `xorm:"int(11) INDEX" json:"repeatAfter"`

	Created int64 `xorm:"created" json:"created"`
	Updated int64 `xorm:"updated" json:"updated"`

	CreatedBy User `xorm:"-" json:"createdBy" valid:"-"`

	web.CRUDable `xorm:"-" json:"-"`
	web.Rights   `xorm:"-" json:"-"`
}

// TableName returns the table name for listtasks
func (ListTask) TableName() string {
	return "tasks"
}

// GetTasksByListID gets all todotasks for a list
func GetTasksByListID(listID int64) (tasks []*ListTask, err error) {
	err = x.Where("list_id = ?", listID).Find(&tasks)
	if err != nil {
		return
	}

	// No need to iterate over users if the list doesn't has tasks
	if len(tasks) == 0 {
		return
	}

	// Get all users and put them into the array
	var userIDs []int64
	for _, i := range tasks {
		found := false
		for _, u := range userIDs {
			if i.CreatedByID == u {
				found = true
				break
			}
		}

		if !found {
			userIDs = append(userIDs, i.CreatedByID)
		}
	}

	var users []User
	err = x.In("id", userIDs).Find(&users)
	if err != nil {
		return
	}

	for in, task := range tasks {
		for _, u := range users {
			if task.CreatedByID == u.ID {
				tasks[in].CreatedBy = u
				break
			}
		}

		// obsfucate the user password
		tasks[in].CreatedBy.Password = ""
	}

	return
}

// GetListTaskByID returns all tasks a list has
func GetListTaskByID(listTaskID int64) (listTask ListTask, err error) {
	if listTaskID < 1 {
		return ListTask{}, ErrListTaskDoesNotExist{listTaskID}
	}

	exists, err := x.ID(listTaskID).Get(&listTask)
	if err != nil {
		return ListTask{}, err
	}

	if !exists {
		return ListTask{}, ErrListTaskDoesNotExist{listTaskID}
	}

	u, err := GetUserByID(listTask.CreatedByID)
	if err != nil {
		return
	}
	listTask.CreatedBy = u

	return
}
