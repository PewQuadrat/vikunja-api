package models

// Team holds a team object
type Team struct {
	ID          int64  `xorm:"int(11) autoincr not null unique pk" json:"id" param:"team"`
	Name        string `xorm:"varchar(250) not null" json:"name"`
	Description string `xorm:"varchar(250)" json:"description"`
	CreatedByID int64  `xorm:"int(11) not null INDEX" json:"-"`

	CreatedBy User        `xorm:"-" json:"created_by"`
	Members   []*TeamUser `xorm:"-" json:"members"`

	Created int64 `xorm:"created" json:"created"`
	Updated int64 `xorm:"updated" json:"updated"`

	CRUDable `xorm:"-" json:"-"`
	Rights   `xorm:"-" json:"-"`
}

// TableName makes beautiful table names
func (Team) TableName() string {
	return "teams"
}

// AfterLoad gets the created by user object
func (t *Team) AfterLoad() {
	// Get the owner
	t.CreatedBy, _ = GetUserByID(t.CreatedByID)

	// Get all members
	x.Select("*").
		Table("users").
		Join("INNER", "team_members", "team_members.user_id = users.id").
		Where("team_id = ?", t.ID).
		Find(&t.Members)
}

// TeamMember defines the relationship between a user and a team
type TeamMember struct {
	ID     int64 `xorm:"int(11) autoincr not null unique pk" json:"id"`
	TeamID int64 `xorm:"int(11) not null INDEX" json:"team_id" param:"team"`
	UserID int64 `xorm:"int(11) not null INDEX" json:"user_id" param:"user"`
	Admin  bool  `xorm:"tinyint(1) INDEX" json:"admin"`

	Created int64 `xorm:"created" json:"created"`
	Updated int64 `xorm:"updated" json:"updated"`

	CRUDable `xorm:"-" json:"-"`
	Rights   `xorm:"-" json:"-"`
}

// TableName makes beautiful table names
func (TeamMember) TableName() string {
	return "team_members"
}

// TeamUser is the team member type
type TeamUser struct {
	User  `xorm:"extends"`
	Admin bool `json:"admin"`
}

// GetTeamByID gets a team by its ID
func GetTeamByID(id int64) (team Team, err error) {
	if id < 1 {
		return team, ErrTeamDoesNotExist{id}
	}

	exists, err := x.Where("id = ?", id).Get(&team)
	if err != nil {
		return
	}
	if !exists {
		return team, ErrTeamDoesNotExist{id}
	}

	return
}

// ReadOne implements the CRUD method to get one team
func (t *Team) ReadOne() (err error) {
	*t, err = GetTeamByID(t.ID)
	return
}

// ReadAll gets all teams the user is part of
func (t *Team) ReadAll(user *User) (teams interface{}, err error) {
	all := []*Team{}
	err = x.Select("teams.*").
		Table("teams").
		Join("INNER", "team_members", "team_members.team_id = teams.id").
		Where("team_members.user_id = ?", user.ID).
		Find(&all)

	return all, err
}