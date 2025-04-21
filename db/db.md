# Database Schema Documentation

This document outlines the database schema for the GPU Scheduler system. It includes the table definitions, relationships, and a visual representation of the schema.

---

## Entity-Relationship Diagram (ERD)
```mermaid
erDiagram
    USERS {
        INT id
        VARCHAR email
        VARCHAR name
        VARCHAR password
        DATETIME signup_date
        BOOLEAN is_admin
        BOOLEAN is_whitelisted
    }
    REQUESTS {
        INT id
        INT user_id
        INT requested_time
        ENUM gpu_size
        INT num_gpus
        ENUM priority
        VARCHAR server_name
        ENUM status
        DATETIME start_time
        DATETIME end_time
        DATETIME created_at
    }
    GPUS {
        CHAR id
        VARCHAR server_name
        INT gpu_number
        VARCHAR manufacturer
        VARCHAR model_name
        INT vram_size_mb
    }
    GPU_USAGE {
        INT id
        INT user_id
        INT request_id
        DATETIME start_time
        DATETIME end_time
        INT actual_usage_time
    }
    WHITELIST {
        INT id
        VARCHAR email
    }

    USERS ||--o{ REQUESTS : "has many"
    USERS ||--o{ GPU_USAGE : "has many"
    REQUESTS ||--o{ GPU_USAGE : "has many"
    GPUS ||--o{ REQUESTS : "can be assigned to"
```

---

## Tables

### 1. **Users Table**
Stores information about users who sign up for the system.

| Column Name      | Data Type        | Description                                      |
|------------------|------------------|--------------------------------------------------|
| `id`             | INT             | Primary Key, unique identifier for each user.    |
| `email`          | VARCHAR(255)    | User's email address, must be unique.            |
| `name`           | VARCHAR(255)    | Full name of the user.                           |
| `password`       | VARCHAR(255)    | Hashed password for authentication.              |
| `signup_date`    | DATETIME        | Timestamp of when the user signed up.            |
| `is_admin`       | BOOLEAN         | Indicates if the user is an admin.               |
| `is_whitelisted` | BOOLEAN         | Indicates if the user is allowed to sign up.     |

---

### 2. **Requests Table**
Tracks GPU access requests made by users.

| Column Name      | Data Type        | Description                                      |
|------------------|------------------|--------------------------------------------------|
| `id`             | INT             | Primary Key, unique identifier for each request. |
| `user_id`        | INT             | Foreign Key, references `users.id`.              |
| `requested_time` | INT             | Requested GPU time in hours.                     |
| `gpu_size`       | ENUM            | GPU size: 'small', 'medium', 'large'.            |
| `num_gpus`       | INT             | Number of GPUs requested.                        |
| `priority`       | ENUM            | Priority: 'low', 'medium', 'high', 'emergency'.  |
| `server_name`    | VARCHAR(255)    | Server name, or NULL for "any server".           |
| `status`         | ENUM            | Request status: 'pending', 'approved', etc.      |
| `start_time`     | DATETIME        | Scheduled start time for GPU usage.              |
| `end_time`       | DATETIME        | Scheduled end time for GPU usage.                |
| `created_at`     | DATETIME        | Timestamp of when the request was created.       |

---

### 3. **GPUs Table**
Tracks available GPUs and their specifications.

| Column Name      | Data Type        | Description                                      |
|------------------|------------------|--------------------------------------------------|
| `id`             | CHAR(32)        | Primary Key, MD5 hash of `server_name` + `gpu_number`. |
| `server_name`    | VARCHAR(255)    | Name of the server hosting the GPU.             |
| `gpu_number`     | INT             | GPU index on the server (e.g., 0, 1, 2).         |
| `manufacturer`   | VARCHAR(255)    | Manufacturer of the GPU (e.g., NVIDIA, AMD).     |
| `model_name`     | VARCHAR(255)    | Model name of the GPU (e.g., RTX 3090).          |
| `vram_size_mb`   | INT             | VRAM size in megabytes (e.g., 8192 for 8GB).     |

---

### 4. **GPU Usage Table**
Tracks actual GPU usage for auditing and analytics.

| Column Name      | Data Type        | Description                                      |
|------------------|------------------|--------------------------------------------------|
| `id`             | INT             | Primary Key, unique identifier for each record.  |
| `user_id`        | INT             | Foreign Key, references `users.id`.              |
| `request_id`     | INT             | Foreign Key, references `requests.id`.           |
| `start_time`     | DATETIME        | Actual start time of GPU usage.                  |
| `end_time`       | DATETIME        | Actual end time of GPU usage.                    |
| `actual_usage_time` | INT          | Actual usage time in minutes.                    |

---

### 5. **Whitelist Table**
Stores emails of users allowed to sign up.

| Column Name      | Data Type        | Description                                      |
|------------------|------------------|--------------------------------------------------|
| `id`             | INT             | Primary Key, unique identifier for each record.  |
| `email`          | VARCHAR(255)    | Email address of the whitelisted user.           |

---

## Relationships

1. **Users → Requests**:
   - A user can make multiple requests.
   - `requests.user_id` references `users.id`.

2. **Users → GPU Usage**:
   - A user can have multiple GPU usage records.
   - `gpu_usage.user_id` references `users.id`.

3. **Requests → GPU Usage**:
   - A request can result in multiple GPU usage records.
   - `gpu_usage.request_id` references `requests.id`.

4. **GPUs**:
   - GPUs are uniquely identified by a combination of `server_name` and `gpu_number`.

---

## Example Usage and Instructions

### 1. **Create the Database**
Run the SQL script provided in the `create_db.sql` file to create the database tables:
```bash
mysql -u <username> -p < create_db.sql
```

### 2. **Insert Data**
Use the CSV import utility to populate the `users` and `gpus` tables.

#### Insert Mode
To append data to a table:
```bash
go run csv_import.go -file=gpus.csv -table=gpus -mode=insert
```

#### Overwrite Mode
To overwrite existing data in a table:
```bash
go run csv_import.go -file=users.csv -table=users -mode=overwrite
```

### 3. **Query the Database**
Use SQL queries to interact with the database. For example:
- Find available GPUs:
  ```sql
  SELECT * FROM gpus WHERE id = 'available';
  ```
- View pending requests:
  ```sql
  SELECT * FROM requests WHERE status = 'pending';
  ```

### 4. **Update GPU Status**
Update the status of GPUs as they are assigned or go into maintenance:
```sql
UPDATE gpus SET status = 'in_use' WHERE id = MD5(CONCAT('server1', '0'));
```

---

Let me know if you need further updates!
