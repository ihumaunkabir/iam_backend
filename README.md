### Indentity Access Management
Simple implementation in Golang and MongoDB.

#### Clone
```
git clone https://github.com/ihumaunkabir/iam_backend.git
```

#### Installation
```go
cd iam_backend
go mod tidy
```

#### Run
```go
go run main.go
```
```Listening and serving HTTP on :8080```

#### API Endpoints
Registration
```go
POST   /api/v1/register  
```
```json
{
	"username": "usernamerequired",
	"email": "requiredemail",
	"password": "requiredmin6length"
}
```
Login
```go
POST   /api/v1/login  
```
```json
{
	"username": "usernamerequired",
	"password": "requiredmin6length"
}
```

#### Data Model
```go
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email        string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Roles        []string           `bson:"roles" json:"roles"`
	Active       bool               `bson:"active" json:"active"`
	LastLogin    *time.Time         `bson:"last_login" json:"last_login"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
```

#### License
MIT

#### Acknowledgement
This repo was assisted with Claude 3.7 Sonnet. 
