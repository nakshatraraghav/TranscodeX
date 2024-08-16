# TranscodeX

TranscodeX is a Golang-based microservice application designed for multimedia transcoding and processing, specifically focused on video and image transcoding and processing. The API allows users to upload media files, specify transcoding tasks, and retrieve the processed output, all managed via AWS services like S3, ECS Fargate.

## About the Micro-Service

TranscodeX allows users to seamlessly upload media files, specify desired processing tasks, and retrieve processed outputs. The API offers a range of transcoding options tailored to meet different media requirements, transcoding videos to streamable bitrate adjustable content,adjusting resolutions, or transforming images with various filters. The API is designed with scalability and security at its core, leveraging powerful AWS services for seamless file storage, processing, and task orchestration.

## Key Features

1. **Image Processing**:

   - **Resize**: Dynamically resize images to fit specific dimensions, preserving aspect ratios as needed. Ideal for adjusting images to different resolutions.
   - **Force Resize**: Forcefully resize images to exact dimensions, which may alter the aspect ratio if necessary.
   - **Rotate**: Rotate images by specified degrees to correct orientation or achieve desired visual effects.
   - **Convert Format**: Convert images between various formats such as JPEG, PNG, GIF, and more, ensuring compatibility with different platforms.
   - **Watermark**: Add custom watermarks to images to protect intellectual property or provide branding.
   - **Generate Thumbnail**: Create small preview versions of images, useful for displaying in galleries or lists.

2. **Video Transcoding**:

   - **Transcode to Multiple Resolutions**: Convert videos into various resolutions simultaneously, catering to different devices and bandwidth requirements. Resolutions can include 480p, 720p, 1080p, and 4K.
   - **Transcode to Resolution**: Convert videos to a specified resolution, providing flexibility in video output for various applications and devices.

3. **AWS Integration**:

   - **S3 Storage**: Utilize AWS S3 for secure, scalable, and durable storage of media files. Provides presigned URLs for secure, direct uploads and downloads of media files, ensuring that data is securely handled and efficiently stored.
   - **ECS Fargate Processing**: Leverage AWS ECS Fargate for serverless container-based processing. Each media processing task runs in isolated containers, ensuring high availability, scalability, and efficient resource management.
   - **SQS Queuing**: Integrate with AWS SQS for asynchronous task management. Tasks can be queued and processed in a distributed manner, which helps manage high loads and ensures reliable processing without bottlenecks.

4. **Security**:

   - **API Key Management**: Secure API access through API keys. Users can generate, manage, and revoke API keys to control access and maintain security.
   - **Session Management**: Implement secure user sessions using JWT (JSON Web Tokens). This provides stateless authentication and ensures secure interactions with the API.
   - **Encryption**: Encrypt all media files using server-side encryption with AWS S3. Ensure that API traffic is secured using HTTPS, protecting data during transmission.

5. **Scalability and Performance**:

   - **Auto-Scaling**: AWS ECS Fargate automatically scales processing resources based on demand. This ensures that the system can handle varying loads by spinning up additional containers as needed.
   - **Concurrency**: Utilize Golang’s goroutines for managing concurrent processing of media files. This design enables efficient handling of multiple requests simultaneously without compromising performance.

6. **Asynchronous Processing**:
   - **Non-Blocking Operations**: Handle media processing tasks asynchronously. Users can submit tasks and receive results at a later time, which is particularly beneficial for large video files or complex processing operations.
   - **Task Status Monitoring**: Provide endpoints to check the status of media processing tasks. Users can monitor the progress of their requests and receive updates on completion.

## Architecture Overview

### 1. Frontend Integration

- **Frontend**: The API serves as a backend service for media processing. The frontend application (if any) interacts with the API to manage user accounts, upload media files, create processing jobs, and retrieve results.

### 2. Backend Components

- **Backend Technology**: Developed in Golang, the backend API handles all core functionalities, including user management, API key generation, file uploads, and media processing requests.

### 3. AWS Services

#### 3.1 S3 Buckets

- **Input Bucket**:
  - **Purpose**: Stores raw media files uploaded by users.
  - **Usage**: Users upload files to this bucket via presigned URLs. The backend API generates these URLs for secure direct uploads.
- **Output Bucket**:
  - **Purpose**: Stores processed media files after transcoding or image processing tasks are completed.
  - **Usage**: The backend API retrieves and returns processed files from this bucket.

#### 3.2 SQS Queue

- **Purpose**: Manages and queues media processing tasks.
- **Usage**: When a media file upload is completed, a message is sent to the SQS queue to initiate processing.

#### 3.3 ECS Fargate

- **Purpose**: Runs Docker containers to handle media processing tasks.
- **Usage**: Each container performs specific media processing tasks, such as video transcoding or image resizing. ECS Fargate ensures that these containers scale automatically based on the workload.

#### 3.4 Lambda Functions

- **Purpose**: Executes event-driven functions.
- **Usage**: Lambda functions handle various tasks:
  - **Trigger Function**: Invoked by SQS messages to start processing.
  - **Processing Function**: Invoked by the trigger function to perform specific processing tasks.

#### 3.5 RDS Database

- **Purpose**: Provides a managed relational database for storing user data, API keys, session information, and job metadata.
- **Usage**: The backend API interacts with RDS to manage data and track job statuses.

#### 3.6 EC2 Instance

- **Purpose**: Hosts the backend API server.
- **Usage**: The API server handles user requests, interacts with other AWS services, and manages media processing workflows.

## Modules in the Micro-Service

### Backend

The **Backend** module of the Transcodex project is a critical component, developed in Go, that serves as the central hub for handling API requests and managing user interactions. It is responsible for authenticating users, generating API keys, managing media file uploads, and initiating processing jobs. The backend interfaces with AWS services to store media files in S3 buckets and track processing jobs using metadata stored in an RDS database. It also interacts with SQS queues to enqueue processing tasks, ensuring that media files are processed efficiently and reliably. By using a robust Go-based architecture, the backend can handle high-throughput requests and provide a reliable API for frontend applications or other clients.

### Infra

The **Infra** module is designed to manage the infrastructure setup for the Transcodex project using AWS CDK. It defines and deploys the essential AWS resources needed for the application's operation, including EC2 instances for running the backend server, S3 buckets for storing input and output media files, an RDS database for managing user and job data, and SQS queues for orchestrating processing tasks. Additionally, it sets up ECS Fargate clusters for running containerized worker tasks and Lambda functions for event-driven processing. By using AWS CDK, the Infra module allows for a declarative and automated approach to infrastructure management, ensuring that resources are provisioned consistently and efficiently.

### Lambda

The **Lambda** module handles serverless computation and processing tasks within the Transcodex project. It is responsible for executing functions in response to events, such as messages in the SQS queue. The primary Lambda functions are invoked when a new processing job is queued, triggering the execution of media processing tasks. This module leverages AWS Lambda’s ability to scale automatically with demand, ensuring that media files are processed quickly and efficiently without manual intervention. By using Lambda functions, Transcodex achieves a flexible and cost-effective solution for handling varying workloads, minimizing infrastructure management overhead.

### Worker

The **Worker** module is integral to the media processing workflow within Transcodex. Deployed on AWS ECS Fargate, the worker containers are tasked with performing the actual media transformations, such as video transcoding, image resizing, or format conversion. Each container is configured to handle specific processing tasks based on the instructions provided by the backend via SQS messages. The worker module benefits from ECS Fargate’s ability to scale container instances based on workload, ensuring that media processing jobs are completed efficiently and promptly. This module encapsulates the core processing logic, allowing the Transcodex system to handle large volumes of media files and complex transformations in a scalable manner.

## Detailed Workflow

### 1. User Interaction

1. **Account Creation**:

   - Users create an account via the API, which stores user details in the RDS database.

2. **Login and API Key Creation**:

   - Users log in using credentials. After successful authentication, they create API keys for authorization.

3. **Upload Media File**:

   - The backend API generates a presigned S3 URL.
   - Users use this URL to upload media files directly to the S3 input bucket.

4. **Add Processing Job**:
   - After uploading, users submit a job request to process the media file.
   - The backend API writes a message to the SQS queue, detailing the processing job requirements.

### 2. Media Processing

1. **Message Enqueuing**:

   - The message, which includes details about the media file and processing requirements, is placed in the SQS queue.

2. **Trigger Lambda Function**:

   - An SQS trigger invokes a Lambda function whenever a new message is available in the queue.
   - The trigger Lambda function reads the message and determines the processing task required.

3. **Process Lambda Function**:

   - The trigger Lambda function invokes another Lambda function specifically designed to perform the media processing tasks (transcoding, resizing, etc.).
   - This Lambda function interacts with ECS Fargate tasks to execute the processing job.

4. **ECS Fargate Processing**:
   - ECS Fargate launches containers based on the task definition to process the media file.
   - Containers execute the required transformations and save the processed files to the S3 output bucket.

### 3. Post-Processing

1. **Job Status Update**:

   - Once processing is complete, the Lambda function updates the job status in the RDS database.
   - The status reflects whether the job is complete, failed, or in progress.

2. **Retrieve Results**:
   - Users query the backend API to check the status of their processing job.
   - If completed, users receive a link to download the processed media file from the S3 output bucket.

### 4. Scalability and Reliability

- **Scalability**:

  - **ECS Fargate**: Automatically scales containerized applications based on the number of processing tasks.
  - **SQS**: Manages and distributes tasks across multiple Lambda functions or ECS Fargate tasks.

- **Reliability**:
  - **S3**: Provides durable and highly available storage for media files.
  - **RDS**: Offers a managed database with automated backups and high availability.
  - **Lambda and ECS Fargate**: Ensures serverless execution and scaling, reducing infrastructure management overhead.

---

## Codebase Architecture

```bash
├──  backend
│  ├──  backend.postman_collection.json
│  ├──  cmd
│  │  └──  main.go
│  ├──  db
│  │  ├──  db.go
│  │  └──  migration
│  │     ├──  000001_init_schema.down.sql
│  │     ├──  000001_init_schema.up.sql
│  │     ├──  000002_add_apikey_to_uploads_and_processing_jobs.down.sql
│  │     └──  000002_add_apikey_to_uploads_and_processing_jobs.up.sql
│  ├──  go.mod
│  ├──  go.sum
│  ├──  internal
│  │  ├──  aws
│  │  │  ├──  s3
│  │  │  │  └──  s3.go
│  │  │  └──  sqs
│  │  │     └──  sqs.go
│  │  ├──  config
│  │  │  └──  env.go
│  │  ├──  controllers
│  │  │  ├──  apikeys.controller.go
│  │  │  ├──  media.controller.go
│  │  │  ├──  session.controller.go
│  │  │  └──  user.controller.go
│  │  ├──  middlewares
│  │  │  ├──  auth.go
│  │  │  ├──  ensure-apikey.go
│  │  │  ├──  validate-apikey.go
│  │  │  └──  validate.go
│  │  ├──  routes
│  │  │  ├──  apikey.router.go
│  │  │  ├──  media.router.go
│  │  │  ├──  session.router.go
│  │  │  └──  user.router.go
│  │  ├──  schema
│  │  │  ├──  apikeys.schema.go
│  │  │  ├──  media.schema.go
│  │  │  ├──  session.schema.go
│  │  │  └──  user.schema.go
│  │  ├──  server
│  │  │  ├──  middlewares.go
│  │  │  ├──  routes.go
│  │  │  └──  server.go
│  │  └──  services
│  │     ├──  apikeys.service.go
│  │     ├──  media.service.go
│  │     ├──  session.service.go
│  │     └──  user.service.go
│  ├──  lib
│  │  └──  validator.go
│  ├──  Makefile
│  ├──  types
│  │  └──  ctx.go
│  └──  util
│     ├──  art.go
│     ├──  cookie.go
│     ├──  json.go
│     ├──  jwt.go
│     ├──  operations.go
│     └──  response.go
├──  infra
│  ├──  bin
│  │  ├──  infra.ts
│  │  └──  local.ts
│  ├──  cdk.json
│  ├──  config
│  │  └──  zenv.ts
│  ├──  constructs
│  │  ├──  ec2.ts
│  │  ├──  ecs.ts
│  │  ├──  lambda.ts
│  │  ├──  rds.ts
│  │  ├──  s3.ts
│  │  └──  sqs.ts
│  ├──  jest.config.js
│  ├──  lib
│  │  ├──  infra-stack.ts
│  │  └──  local-stack.ts
│  ├──  package.json
│  ├──  pnpm-lock.yaml
│  ├──  README.md
│  ├──  test
│  │  └──  infra.test.ts
│  └──  tsconfig.json
├──  lambda
│  ├──  bin
│  │  ├──  bootstrap
│  │  ├──  main
│  │  └──  main.zip
│  ├──  cmd
│  │  └──  lambda
│  │     └──  main.go
│  ├──  config
│  │  └──  env.go
│  ├──  go.mod
│  ├──  go.sum
│  ├──  internal
│  │  ├──  ecs
│  │  │  └──  ecsservice.go
│  │  ├──  handler
│  │  │  ├──  handler.go
│  │  │  └──  handler.local.go
│  │  ├──  rds
│  │  │  └──  rds.go
│  │  └──  secretsmanager
│  │     └──  secretsmanager.go
│  └──  lib
│     └──  validate.go
├──  README.md
└──  worker
   ├──  cmd
   │  └──  main.go
   ├──  config
   │  └──  env.go
   ├──  db
   │  └──  db.go
   ├──  Dockerfile
   ├──  go.mod
   ├──  go.sum
   ├──  internal
   │  ├──  application
   │  │  └──  application.go
   │  ├──  processors
   │  │  ├──  image
   │  │  │  └──  image.go
   │  │  └──  video
   │  │     └──  video.go
   │  ├──  s3
   │  │  └──  s3.go
   │  └──  service
   │     └──  job.service.go
   ├──  lib
   │  └──  validate.go
   └──  types

```

## Installation

### Infra

To set up the infrastructure for the Transcodex project, start by creating a Virtual Private Cloud (VPC) in AWS with private subnets and NAT gateways. This VPC will provide a secure network environment for your resources. Next, configure a security group for the ECS tasks, ensuring that it has no inbound access and allows all outbound traffic to meet security and connectivity requirements.

You need to prepare an environment file with the following variables to configure your AWS resources:

```bash
ENVIRONMENT=""
AWS_ACCOUNT_ID=""
AWS_REGION=""
TRANSCODEX_WORKER_IMAGE_URI=""
DATABASE_INSTANCE_IDENTIFIER=""
RDS_DATABASE_USERNAME=""
RDS_DATABASE_PASSWORD=""
BUCKET_NAME=""
ECS_CLUSTER_NAME=""
VPC_ID=""
CONNECTION_STRING=""
SUBNET_IDS=""
SECURITY_GROUP_ID=""
```

Make sure to include the VPC ID and security group ID in the environment file to ensure proper configuration. If you are working in local mode, you can also specify the connection string for the database to facilitate local development and testing.

**Available Scripts**

| **Script Name** | **Command**                                           | **Description**                                 |
| --------------- | ----------------------------------------------------- | ----------------------------------------------- |
| `build`         | `tsc`                                                 | Compile TypeScript code                         |
| `watch`         | `tsc -w`                                              | Watch for changes and recompile TypeScript code |
| `test`          | `jest`                                                | Run tests using Jest                            |
| `cdk`           | `cdk`                                                 | AWS CDK command-line tool                       |
| `deploy:local`  | `cdk deploy --app 'ts-node bin/local.ts'`             | Deploy local environment                        |
| `destroy:local` | `cdk destroy --app 'ts-node bin/local.ts' LocalStack` | Destroy local environment                       |
| `deploy:prod`   | `cdk deploy --app 'ts-node bin/infra.ts'`             | Deploy production environment                   |
| `destroy:prod`  | `cdk destroy --app 'ts-node bin/infra.ts' InfraStack` | Destroy production environment                  |

### Deployment

You can deploy your application both locally and in a production environment using the provided scripts. To deploy locally, use the `deploy:local` script, and to deploy in production, use the `deploy:prod` script.

### Backend

1. **Create a `.env.local` File**

Add the following environment variables to a `.env.local` file in your backend directory:

```bash
PORT=""
CONNECTION_STRING=""
JWT_PRIVATE_KEY=""
ACCESS_TOKEN_TTL=""
REFRESH_TOKEN_TTL=""
BUCKET_NAME=""
AWS_REGION=""
SQS_QUEUE_URL=""
```

2. **Run migrations using the Makefile**

You can run migrations and use other commands using the Makefile in the backend

**Makefile Commands**

| Target                | Command                                                                                                                                                                                                        |
| --------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `start`               | `@go run cmd/main.go`                                                                                                                                                                                          |
| `build`               | `@go build -o bin/api cmd/main.go`                                                                                                                                                                             |
| `run_pg`              | `@docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres`                                                                                                                           |
| `stop_pg`             | `@docker container stop postgres`                                                                                                                                                                              |
| `remove_pg_container` | `@docker container rm postgres`                                                                                                                                                                                |
| `createdb`            | `@docker exec -it postgres createdb --username=postgres --owner=postgres transcodex`                                                                                                                           |
| `dropdb`              | `@docker exec -it postgres dropdb --username=postgres transcodex`                                                                                                                                              |
| `migrateup`           | `@docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path=/migrations -database "postgresql://postgres:password@127.0.0.1:5432/transcodex?sslmode=disable" -verbose up`        |
| `migratedown`         | `@docker run --rm -v $(PWD)/db/migration:/migrations --network host migrate/migrate -path=/migrations -database "postgresql://postgres:password@127.0.0.1:5432/transcodex?sslmode=disable" -verbose down -all` |
| `clean`               | `@make stop_pg`<br>`@make remove_pg_container`                                                                                                                                                                 |

3. **Postman Collection**

You can use the postman collection provided in the repository to run test the api

### Lambda

To set up the Lambda function for your application, follow these steps:

1. **Create a Bin Directory**

   In your Lambda directory, create a `bin` folder to store the compiled binary executable:

   ```bash
   mkdir -p lambda/bin
   ```

2. **Build the Binary Executable**

   Compile your Go application into a binary executable that is compatible with the Linux operating system. Run the following command:

   ```bash
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/main cmd/lambda/main.go
   ```

   This command builds the Go code located in cmd/lambda/main.go and outputs the binary as bin/main.

3 **Add a Bootstrap File**
Create a bootstrap file that will serve as the entry point for the Lambda function. Add the following content:

```bash
#!/bin/sh
set -euo pipefail

# Execute your Go binary
./main
```

Make sure the bootstrap file is executable:

```bash
chmod +x lambda/bin/bootstrap
```

4. **Create a Deployment Package**

Zip the bootstrap file and the compiled binary into a deployment package named main.zip:

```bash
zip lambda/bin/main.zip bootstrap main
```

5. **Environment Variables**

These environment variables will be set automatically by the infrastructure deployment:

```bash
REGION_STRING=""
BUCKET_NAME=""
ECS_CLUSTER_NAME=""
ECS_TASK_DEFINITION=""
RDS_DATABASE_USERNAME=""
RDS_DATABASE_PASSWORD=""
DATABASE_INSTANCE_IDENTIFIER=""
```

### Worker

The Worker module processes media files based on the transformations specified. The environment variables loaded by the Lambda function for ECS tasks include:

```bash
  MEDIA_TYPE=""
  BUCKET_NAME=""
  OBJECT_KEY=""
  TRANSFORMATIONS=""
  CONNECTION_STRING=""
  UPLOAD_ID=""
```

**Transformation Commands**
| **Command** | **Module** | **Description** |
|---------------------------|--------------------------------------|----------------------------------------------------------------|
| **Image Transformations** | | |
| `RESIZE` | `ip.Resize` | Resize the image to specified dimensions. |
| `FORCE-RESIZE` | `ip.ForceResize` | Force resize the image, ignoring aspect ratio. |
| `ROTATE` | `ip.Rotate` | Rotate the image by a specified angle. |
| `CONVERT-FORMAT` | `ip.ConvertFormat` | Convert the image to a different format (e.g., JPEG, PNG). |
| `WATERMARK` | `ip.Watermark` | Add a watermark to the image. |
| `GENERATE-THUMBNAIL` | `ip.GenerateThumbnail` | Generate a thumbnail image from the original image. |
| **Video Transformations** | | |
| `TRANSCODE` | `vp.TranscodeToMultipleResolutions` | Transcode the video to multiple resolutions. |
| `TRANSCODE-RESOLUTION` | `vp.TranscodeToResolution` | Transcode the video to a specific resolution. |

**Dockerfile**

```bash
# Use Alpine as the base image for building the Go application
FROM golang:1.22.5-alpine as builder

# Install necessary build dependencies
RUN apk add --no-cache \
    build-base \
    git \
    pkgconfig \
    vips-dev \
    ffmpeg

# Set the working directory
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /app/bin/worker ./cmd/main.go

# Use a lightweight Alpine image for the runtime environment
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    vips \
    ffmpeg

# Set the working directory
WORKDIR /app

# Copy the built Go application from the builder stage
COPY --from=builder /app/bin/worker /app/worker

# Command to run the application
CMD ["/app/worker"]

```

## API Endpoints

### User Management

- **Create User**

  - `POST /users/`
  - Create a new user.

- **Delete User**
  - `DELETE /users`
  - Delete the current user.

### Session Management

- **Create Session (Login)**

  - `POST /sessions`
  - Create a new session (login) for the user.

- **Get Information About Current Session**

  - `GET /sessions`
  - Get information about the current session.

- **Get All Active Sessions**

  - `GET /sessions/all`
  - Get information about all active sessions.

- **Logout from Current Session**

  - `DELETE /sessions/`
  - Logout from the current session.

- **Log Out from All Devices**
  - `DELETE /sessions/all`
  - Log out from all devices and terminate all sessions.

### API Key Management

- **Get Active API Key**

  - `GET /apikeys`
  - Get the active API key for the user.

- **Generate API Key**

  - `POST /apikeys`
  - Generate a new API key for the user.

- **Revoke API Key**
  - `DELETE /apikeys`
  - Revoke the active API key for the user.

### Media Management

- **Create Upload**

  - `POST /media/upload`
  - Initiate an upload process and receive a presigned URL for uploading media files.

- **Create Processing Job**

  - `POST /media/process`
  - Create a processing job for the uploaded media, specifying the type of processing required.

- **Get Processing Job Status**

  - `GET /media/status/{job_id}`
  - Retrieve the status of a specific processing job using its job ID.

- **Download Processed Media**
  - `GET /media/download/{job_id}`
  - Download the processed media once the job is complete using its job ID.

## Contributing

Contributions are welcome! Here are the steps to contribute to the Sentrimetric API project:

1. Fork the repository.
2. Create a new branch for your feature or bug fix: `git checkout -b feature/your-feature-name` or `git checkout -b bugfix/your-bug-fix-name`.
3. Commit your changes: `git commit -m 'Add some feature'` or `git commit -m 'Fix some bug'`.
4. Push to the branch: `git push origin feature/your-feature-name` or `git push origin bugfix/your-bug-fix-name`.
5. Submit a pull request to the `main` branch of the original repository.

Please make sure to update tests, if applicable, and adhere to the existing code style and guidelines.

## License

This project is licensed under the [MIT License](LICENSE).

Feel free to modify this README.md file to fit your project's specific details and requirements.
