package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
	"time"
)

type Users struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Email              string    `gorm:"column:email;type:varchar(255)" json:"email"`
	Username           string    `gorm:"column:username;type:varchar(100);NOT NULL" json:"username"`
	PasswordHash       string    `gorm:"column:password_hash;type:varchar(255)" json:"passwordHash"`
	ProfilePicture     string    `gorm:"column:profile_picture;type:varchar(255)" json:"profilePicture"`
	NativeLanguageID   int64     `gorm:"column:native_language_id;type:bigint(20)" json:"nativeLanguageID"`
	LearningLanguageID int64     `gorm:"column:learning_language_id;type:bigint(20)" json:"learningLanguageID"`
	LanguageLevel      string    `gorm:"column:language_level;type:varchar(50)" json:"languageLevel"`
	Age                int       `gorm:"column:age;type:int(11)" json:"age"`
	Gender             string    `gorm:"column:gender;type:enum('male','female','other')" json:"gender"`
	Interests          string    `gorm:"column:interests;type:text" json:"interests"`
	CountryID          int64     `gorm:"column:country_id;type:bigint(20)" json:"countryID"`
	RegistrationDate   time.Time `gorm:"column:registration_date;type:timestamp;default:CURRENT_TIMESTAMP" json:"registrationDate"`
	LastLogin          time.Time `gorm:"column:last_login;type:timestamp" json:"lastLogin"`
	Status             string    `gorm:"column:status;type:enum('active','inactive','banned');default:active" json:"status"`
	EmailVerified      int       `gorm:"column:email_verified;type:tinyint(1);default:0" json:"emailVerified"`
	VerificationToken  string    `gorm:"column:verification_token;type:varchar(255)" json:"verificationToken"`
	TokenExpiration    time.Time `gorm:"column:token_expiration;type:timestamp" json:"tokenExpiration"`
}
