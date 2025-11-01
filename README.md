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
```bash
go mod tidy
```
🏃 Run Server
```bash
go build -o out && ./out
```
Server will be running at `http://localhost:<port>`

🗄 Database

Run migrations to set up tables:
`goose postgres "postgres://<user>:<password>@<host>:<port>/<database>" up`

📬 API Docs

Full API documentation is available here: [Swagger](https://hireloop.onrender.com/swagger/index.html)  

📜 License

MIT License.