# HireLoop

The HireLoop Backend is a RESTful API for a modern Job Board platform. It allows employers to post jobs, applicants to apply with resumes, and provides search & filter functionality for job seekers.

## 🚀 Features

#### 📝 Employers can create, update, and delete job posts.

#### 👨‍💼 Applicants can browse jobs and apply with resumes.

#### 🔍 Full-text search and filtering (by location, type, salary).

#### 🔒 Role-based access control (Employer, Applicant).

#### 📂 File upload support for resumes.

## 🛠 Tech Stack

- Backend: Go (Golang) + Gorilla Mux
- Database: PostgreSQL
- Auth: JWT-based authentication
- File Storage: Local (for resumes)


## ⚙️ Installation

#### ✅ Prerequisites

- Go 1.21+ installed
- PostgreSQL installed & running


📥 Clone Repository

```bash
git clone https://github.com/your-username/hireloop-backend.git
cd hireloop-backend
```

⚙️ Configuration

Create a .env file in the root directory of the project and then copy all the configuration variables from .env.example into your new .env file and update the values as needed.


📦 Install Dependencies

go mod tidy

🏃 Run Server
```bash
go run main.go
```
Server will be running at http://localhost:<port>
🗄 Database

Run migrations to set up tables:

psql -U your_user -d hireloop -f migrations.sql

📬 API Docs

API documentation is available in the docs/ folder or via Postman collection.
📜 License

MIT License.