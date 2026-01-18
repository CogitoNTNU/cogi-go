package db

import (
	"time"

	"github.com/CogitoNTNU/cogi-go/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID              uuid.UUID     `db:"user_id"`
	FirstName       string         `db:"first_name"`
	LastName        string         `db:"last_name"`
	Description     string        `db:"description"`
	Nickname        string         `db:"nickname"`
	Email           string         `db:"email"`
	Phone           string         `db:"phone"`
	Gender          model.GenderType     `db:"gender"` 
	GithubURL       *string        `db:"github_url"`
	LinkedinURL     *string        `db:"linkedin_url"`
	KaggleURL       *string        `db:"kaggle_url"`
	HuggingfaceURL  *string        `db:"huggingface_url"`
	Password       string         `db:"password"`
	Avatar          *string        `db:"avatar"`
	ImagePermission *time.Time     `db:"image_permission"` 
	FoodPreference  pq.StringArray `db:"food_preference"` 
}

func (u *User) ToModel() *model.User {
	return &model.User{
		Id:             u.ID,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Description:    u.Description,
		Nickname:       u.Nickname,
		Email:          u.Email,
		Phone:          u.Phone,
		Gender:        u.Gender,
		GithubUrl:     u.GithubURL,
		LinkedinUrl:   u.LinkedinURL,
		KaggleUrl:     u.KaggleURL,
		HuggingfaceUrl: u.HuggingfaceURL,
		Avatar:        u.Avatar,
		ImagePermission: u.ImagePermission,
		FoodPreference: u.FoodPreference,
	}
}

type Project struct {
	ID          uuid.UUID `db:"project_id"`
	Title 	 string    `db:"title"`
	GithubURL  *string   `db:"github_url"`
	Logo 	*string   `db:"logo"`
	Playable bool     `db:"playable"`
	Released bool     `db:"released"`
	ActiveProject bool     `db:"active_project"`
	ProjectURL *string   `db:"project_url"`
}

func (p *Project) ToModel() *model.Project {
	return &model.Project{
		Id:          p.ID,
		Title:       p.Title,
		GithubUrl:   p.GithubURL,
		Logo:        p.Logo,
		Playable:    p.Playable,
		Released:    p.Released,
		ActiveProject: p.ActiveProject,
		ProjectUrl:  p.ProjectURL,
	}
}

// Usikker på om []byte er den beste måten å lagre en JSONB på.
type Article struct {
	ID        uuid.UUID `db:"article_id"`
	Content   []byte    `db:"content"`
	AuthorID  uuid.UUID `db:"author_id"`
	ProjectID uuid.UUID  `db:"project_id"`
	CreatedAt *time.Time  `db:"created_at"`
	UpdatedAt *time.Time  `db:"updated_at"`
}

func (a *Article) ToModel() *model.Article {
	return &model.Article{
		Id:        a.ID,
		Content: a.Content,
		AuthorId:  a.AuthorID,
		ProjectId: a.ProjectID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type Application struct {
	ID        uuid.UUID `db:"application_id"`
	UserID    uuid.UUID `db:"user_id"`
	ProjectID1 uuid.UUID `db:"project_id_1"`
	ProjectID2 uuid.UUID `db:"project_id_2"`
	ProjectID3 uuid.UUID `db:"project_id_3"`
	ApplicationText string    `db:"application_text"`
	AppliedAt      *time.Time `db:"applied_at"`
	ModifiedAt      *time.Time `db:"modified_at"`
	IsMember       bool      `db:"is_member"`
}

func (a *Application) ToModel() *model.Application {
	return &model.Application{
		Id:              a.ID,
		UserId:          a.UserID,
		ProjectId1:     a.ProjectID1,
		ProjectId2:     a.ProjectID2,
		ProjectId3:     a.ProjectID3,
		ApplicationText: a.ApplicationText,
		AppliedAt:      a.AppliedAt,
		ModifiedAt:     a.ModifiedAt,
		IsMember:       a.IsMember,
	}
}

type Sponsor struct {
	ID        uuid.UUID `db:"sponsor_id"`
	Name      string    `db:"name"`
	Logo      *string   `db:"logo"`
	Website   *string   `db:"website"`
	Description string   `db:"description"`
	SponsorLevel model.SponsorLevel `db:"level"`
}

func (s *Sponsor) ToModel() *model.Sponsor {
	return &model.Sponsor{
		Id:          s.ID,
		Name:        s.Name,
		Logo:        s.Logo,
		Website:     s.Website,
		Description: s.Description,
		Level:       s.SponsorLevel,
	}
}

type ProjectSponsor struct {
	ProjectID uuid.UUID `db:"project_id"`
	SponsorID uuid.UUID `db:"sponsor_id"`
	StartDate *time.Time `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
}

func (ps *ProjectSponsor) ToModel() *model.ProjectSponsor {
	return &model.ProjectSponsor{
		ProjectID: ps.ProjectID,
		SponsorID: ps.SponsorID,
		StartDate: ps.StartDate,
		EndDate:   ps.EndDate,
	}
}

type UserActivity struct {
	ID       uuid.UUID    `db:"activity_id"`
	UserID   uuid.UUID    `db:"user_id"`
	Activity []*time.Time `db:"activity_time"`
}

func (ua *UserActivity) ToModel() *model.UserActivity {
	return &model.UserActivity{
		Id:       ua.ID,
		UserId:   ua.UserID,
		Activity: ua.Activity,
	}
}

type Team struct {
    ID       uuid.UUID `db:"team_id"`
    Semester model.Semester `db:"semester"`
    Year     int       `db:"year"`
}

func (t *Team) ToModel() *model.Team {
    return &model.Team{
        Id:       t.ID,
        Semester: t.Semester,
        Year:     t.Year,
    }
}

type ProjectTeam struct {
    Team                
    ProjectID uuid.UUID `db:"project_id"`
}

func (pt *ProjectTeam) ToModel() *model.ProjectTeam {
    return &model.ProjectTeam{
        Team: model.Team{
            Id:       pt.ID,
            Semester: pt.Semester,
            Year:     pt.Year,
        },
        ProjectId: pt.ProjectID,
    }
}

type AdministrasjonTeam struct {
    Team                
    GroupType model.GroupType `db:"group_type"`
}

func (at *AdministrasjonTeam) ToModel() *model.AdministrasjonTeam {
    return &model.AdministrasjonTeam{
        Team: model.Team{
            Id:       at.ID,
            Semester: at.Semester,
            Year:     at.Year,
        },
        GroupType: at.GroupType,
    }
}

type Role struct {
	ID   uuid.UUID `db:"role_id"`
	Name string    `db:"role_name"`
}

func (r *Role) ToModel() *model.Role {
	return &model.Role{
		Id:   r.ID,
		Name: r.Name,
	}
}

type TeamMembers struct {
	TeamID uuid.UUID `db:"team_id"`
	UserID uuid.UUID `db:"user_id"`
	RoleID uuid.UUID `db:"role_id"`
}

func (tm *TeamMembers) ToModel() *model.TeamMembers {
	return &model.TeamMembers{
		TeamID: tm.TeamID,
		UserID: tm.UserID,
		RoleID: tm.RoleID,
	}
}

type Achivement struct {
	ID          uuid.UUID `db:"achivement_id"`
	Title       string    `db:"title"`
	IconURL	*string   `db:"icon_url"`
}

func (a *Achivement) ToModel() *model.Achivement {
	return &model.Achivement{
		Id:    a.ID,
		Title: a.Title,
		IconUrl: a.IconURL,
	}
}

type UserAchivement struct {
	ID 		uuid.UUID `db:"user_achivement_id"`
	UserID 	uuid.UUID `db:"user_id"`
	AchivementID uuid.UUID `db:"achivement_id"`
	DateEarned *time.Time `db:"date_earned"`
}

func (ua *UserAchivement) ToModel() *model.UserAchivement {
	return &model.UserAchivement{
		Id:          ua.ID,
		UserId:      ua.UserID,
		AchivementId: ua.AchivementID,
		DateEarned: ua.DateEarned,
	}
}

type TempApplication struct {
	Id              uuid.UUID
	Email           string
	PhoneNumber     string
	Projects        []string
	ApplicationText string
	CreatedAt      *time.Time
}

func (ta *TempApplication) ToModel() *model.TempApplication {
	return &model.TempApplication{
		Id:              ta.Id,
		Email:           ta.Email,
		PhoneNumber:     ta.PhoneNumber,
		Projects:        ta.Projects,
		ApplicationText: ta.ApplicationText,
		CreatedAt:      ta.CreatedAt,
	}
}