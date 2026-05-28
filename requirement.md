# Authentication Platform

# Overview

## Project Background

โปรเจกต์นี้เป็นระบบ Authentication Platform สำหรับจัดการการสมัครสมาชิกและเข้าสู่ระบบของผู้ใช้งาน โดยออกแบบให้รองรับ REST API และสามารถนำไปต่อยอดร่วมกับระบบอื่นในอนาคตได้

ระบบถูกออกแบบให้รองรับผู้ใช้งานจำนวนมาก มีโครงสร้างที่สามารถ maintain และ scale ได้ง่าย รวมถึงสามารถแยก deploy แบบ Microservice ได้ในอนาคต

ปัจจุบันระบบอยู่ในช่วงเริ่มต้นของการพัฒนา จึงมีเฉพาะฟีเจอร์ Register และ Login ก่อน เพื่อใช้เป็นพื้นฐานสำหรับระบบอื่นในอนาคต

นอกจากนี้ เนื่องจากระบบยังอยู่ในช่วงพัฒนาและทดสอบการใช้งาน ระบบจะเปิดให้บริการเฉพาะช่วงเวลา 06:00 น. ถึง 23:00 น. เท่านั้น และมีแผนเปิดให้บริการเต็มรูปแบบตลอด 24 ชั่วโมงในอนาคตเมื่อระบบมีความสมบูรณ์มากขึ้น

---

# Business Requirement

Product Owner ต้องการระบบที่สามารถ:

* ให้ผู้ใช้งานสมัครสมาชิก
* ให้ผู้ใช้งานเข้าสู่ระบบ
* รองรับ REST API สำหรับ Frontend และ Mobile Application
* รองรับ JWT Authentication
* รองรับผู้ใช้งานจำนวนมาก
* มีโครงสร้างระบบที่สามารถขยายเป็น Microservice ได้ในอนาคต
* รองรับการ maintain และ scale ระบบได้ง่าย
* จัดเก็บรหัสผ่านแบบปลอดภัยด้วย password_hash และ password_salt
* เปิดให้บริการระบบในช่วงเวลา 06:00 น. ถึง 23:00 น.
* รองรับการขยายเวลาให้บริการแบบเต็มรูปแบบในอนาคต

---

# Functional Scope

ระบบประกอบด้วยความสามารถหลักดังนี้

## 1. User Management

* สมัครสมาชิก
* เข้าสู่ระบบ

## 2. Authentication

* Generate JWT Token
* Validate User Credential
* Hash Password พร้อม Salt ก่อนจัดเก็บ

## 3. System Architecture

* รองรับ RESTful API
* รองรับการแยก Service ในอนาคต
* รองรับการ scale ระบบ

## 4. Service Availability

* ระบบเปิดให้บริการระหว่างเวลา 06:00 น. ถึง 23:00 น.
* หากมีการเรียกใช้งานนอกช่วงเวลาที่กำหนด ระบบต้องแจ้งสถานะว่าอยู่นอกเวลาทำการ

---

# Non-Functional Requirement

* ระบบต้องรองรับ RESTful API
* รองรับผู้ใช้งานจำนวนมาก
* สามารถ deploy แบบ Microservice ได้ในอนาคต
* รองรับการ maintain และ scale ระบบได้ง่าย
* มีมาตรฐานด้าน Security สำหรับข้อมูลผู้ใช้งาน
* รองรับ JWT Authentication
* รองรับการจัดการ Configuration ผ่าน ENV และ YAML
* รหัสผ่านต้องถูกจัดเก็บในรูปแบบ password_hash และ password_salt เท่านั้น
* ระบบต้องสามารถกำหนดช่วงเวลาเปิด-ปิดให้บริการได้
* ระบบต้องรองรับการเปลี่ยนแปลงเวลาให้บริการในอนาคตโดยไม่กระทบต่อระบบหลัก

---

# Technology Requirement

| Technology       | Requirement             |
| ---------------- | ----------------------- |
| Language         | Golang                  |
| Framework        | Echo v4                 |
| ORM              | GORM                    |
| Database         | MySQL                   |
| Cache            | Redis                   |
| Authentication   | JWT                     |
| Migration Tool   | golang-migrate          |
| Configuration    | YAML + ENV              |
| Containerization | Docker + Docker Compose |

---

# User Register API

## Endpoint

```http
POST /api/v1/register
```

## Request Body

```json
{
  "email": "customer@example.com",
  "password": "123456",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "0812345678"
}
```

## Success Response

```json
{
  "code": 1000,
  "message": "Register Success"
}
```

## Error Response

```json
{
  "code": 4001,
  "message": "Email already exists"
}
```

```json
{
  "code": 5001,
  "message": "Service is available between 06:00 and 23:00"
}
```

---

# User Login API

## Endpoint

```http
POST /api/v1/login
```

## Request Body

```json
{
  "email": "customer@example.com",
  "password": "123456"
}
```

## Success Response

```json
{
  "code": 1000,
  "message": "Success",
  "data": {
    "access_token": "jwt-token"
  }
}
```

## Error Response

```json
{
  "code": 4002,
  "message": "Invalid email or password"
}
```

```json
{
  "code": 5001,
  "message": "Service is available between 06:00 and 23:00"
}
```

---

# Validation Requirement

## User Validation

* email ต้องไม่ซ้ำ
* email ต้องอยู่ใน format ที่ถูกต้อง
* password อย่างน้อย 6 ตัวอักษร
* phone_number ต้องเป็นตัวเลข
* first_name ห้ามเป็นค่าว่าง
* last_name ห้ามเป็นค่าว่าง
* password ต้องถูก hash พร้อม salt ก่อนบันทึกลงฐานข้อมูล
* ระบบต้องอนุญาตให้ใช้งานเฉพาะช่วงเวลา 06:00 น. ถึง 23:00 น.

## Authentication Validation

* email และ password ต้องถูกต้อง
* ระบบต้อง generate JWT Token หลัง Login สำเร็จ
* ระบบต้องนำ password_salt มารวมกับ password ก่อนตรวจสอบ password_hash
* หากมีการเรียกใช้งานนอกเวลาที่กำหนด ระบบต้องตอบกลับด้วยข้อความแจ้งนอกเวลาทำการ

---

# Redis Requirement

หลัง Login สำเร็จ ให้เก็บ session ลง Redis

## Redis Key

```text
session:{user_id}
```

## Example Value

```json
{
  "user_id": 1,
  "email": "customer@example.com"
}
```

## TTL

```text
24 hours
```

---

# Database Design

## users table

| Column        | Type            | Constraint  |
| ------------- | --------------- | ----------- |
| id            | bigint unsigned | primary key |
| email         | varchar(255)    | unique      |
| password_hash | varchar(255)    | not null    |
| password_salt | varchar(255)    | not null    |
| first_name    | varchar(100)    | not null    |
| last_name     | varchar(100)    | not null    |
| phone_number  | varchar(20)     | not null    |
| created_at    | datetime        | not null    |
| updated_at    | datetime        | not null    |
| deleted_at    | datetime        | nullable    |

---

# Project Structure

```text
project-root/
├── app/
│   ├── cmd/
│   │   └── main.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── internal/
│       ├── handler/
│       │   ├── auth_handler.go
│       │   └── routes.go
│       │
│       ├── service/
│       │   └── auth_service.go
│       │
│       ├── repository/
│       │   └── user_repository.go
│       │
│       ├── model/
│       │   └── user.go
│       │
│       └── middleware/
│           └── auth_middleware.go
├── config/
│   ├── config.yaml
│   └── secret.env
|
├── migrations/
|
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

# Structure Description

| Folder/File             | Description                                                              |
| ----------------------- | ------------------------------------------------------------------------ |
| app/cmd/main.go         | จุดเริ่มต้นของ application สำหรับ bootstrap server และ dependency ต่าง ๆ |
| app/config/config.go    | จัดการโหลด configuration จาก ENV หรือไฟล์ config                         |
| app/internal/handler    | จัดการ HTTP Request/Response และเรียกใช้งาน service                      |
| app/internal/service    | จัดการ business logic ของระบบ                                            |
| app/internal/repository | จัดการการเข้าถึงฐานข้อมูล MySQL หรือ Redis                               |
| app/internal/model      | เก็บโครงสร้าง model/entity และ request/response model                    |
| app/internal/middleware | Middleware ของระบบ เช่น JWT Authentication หรือ Logging                  |
| config/config.yaml      | ไฟล์ configuration หลักของระบบ                                           |
| config/secret.env       | เก็บ environment variable และ secret ของระบบ                             |
| migrations              | เก็บ database migration script                                           |
| Dockerfile              | สำหรับ build docker image ของ application                                |
| docker-compose.yml      | สำหรับรัน service ต่าง ๆ เช่น API, MySQL และ Redis                       |
| Makefile                | รวม command สำหรับช่วย build/run/test project                            |
| go.mod                  | จัดการ Go module dependency                                              |
| go.sum                  | checksum ของ dependency                                                  |
| README.md               | เอกสารอธิบายโปรเจกต์ วิธีรัน และ requirement                             |


---

# Required Packages

```bash
go get github.com/labstack/echo/v4
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/redis/go-redis/v9
```

---

# Environment Example

## config.yaml
```yaml
log:
  level: info
  env: local

app:
  name: xxx
  project-id: xxx

server:
  address: ":4000"
  time-zone: "Asia/Bangkok"

database:
  primary:
    name: xxx
    ssl-mode: true
    max-idle-connections: 32
    max-open-connections: 64
    max-life-time: 10m

jwt: 
  expire_in: 24h

```

## secret.env

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your-secret-key
```

---

# Bonus Tasks (Optional)

หากทำเสร็จแล้ว สามารถลองเพิ่ม:

* Forgot Password
* Refresh Token
* Unit Test
* Graceful Shutdown

---

# Evaluation Criteria

| Criteria         | Description                        |
| ---------------- | ---------------------------------- |
| Clean Code       | โค้ดอ่านง่าย                       |
| Layer Separation | แยก Handler / Service / Repository |
| Validation       | ตรวจสอบข้อมูลถูกต้อง               |
| Security         | Hash Password + JWT                |
| Database Design  | ออกแบบเหมาะสม                      |
| Error Handling   | Handle Error ถูกต้อง               |

---

Good luck 🚀
