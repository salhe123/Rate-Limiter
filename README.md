# Task 3: Rate Limiter (Go with Gin Framework)

This project implements the "Rate Limiter" task as part of the Eagle Point AI technical assessment. It is a Go application built with the Gin Framewrok.

## Technology Stack

- **Go:** You Can Use Version 1.22 later
- **Gin Web Framework:** For building the REST API.

## Rate Limiting Logic

The rate limiter is configured to allow:
- **Limit:** 5 requests each
- **Time Window:** per 60 seconds restart the request limit
- **Per User:** Tracked by a `X-User-ID` header.

If the limit is exceeded within the time window, subsequent requests from the same user will receive a `Status Code 429 Too Many Requests` HTTP status.

## API Endpoints

The application exposes the following endpoints:

### Rate-Limited Endpoints
These endpoints are subject to the 5 requests per 60 seconds per user limit.

- **URL:** `/api/resource`
- **HTTP Method:** `GET`
- **Requires Header:** `X-User-ID` (e.g., `X-User-ID: user123`)
- **Description:** A sample resource that demonstrates rate-limited access.

- **URL:** `/api/data`
- **HTTP Method:** `POST`
- **Requires Header:** `X-User-ID` (e.g., `X-User-ID: user123`)
- **Description:** A sample data processing endpoint that demonstrates rate-limited access for POST requests.

### Public Endpoint
This endpoint is not subject to any rate limiting.

- **URL:** `/public`
- **HTTP Method:** `GET`
- **Description:** A publicly accessible resource.

## How to Run the Application

### Prerequisites

- Go (v1.22 or later)

### Steps

1.  **Navigate to the project directory:**
    ```bash
    cd rate-limiter-go
    ```

2.  **Run the Go application:**
    ```bash
    go run main.go
    ```
    The application will start on `http://localhost:8081`. You will see a log message in your terminal indicating that the service is running.

## How to Test the Rate Limiter

Once the application is running, you can use `curl` from your terminal to test the different endpoints.

### Testing Rate-Limited Endpoints (`/api/resource`, `/api/data`)

To test the rate limiting, you must include the `X-User-ID` header. Replace `user123` with any desired user identifier.

1.  **Successful Requests (First 5 within 60 seconds):**
    Repeat the following `curl` command 5 times quickly (within 60 seconds):
    
    **For GET `/api/resource`:**
    ```bash
    curl -H "X-User-ID: user123" http://localhost:8081/api/resource
    ```
    Expected response: `{"message":"Access granted to resource!"}` (HTTP 200 OK)

    **For POST `/api/data`:**
    ```bash
    curl -X POST -H "X-User-ID: user123" -H "Content-Type: application/json" -d "{}" http://localhost:8081/api/data
    ```
    Expected response: `{"message":"Data processed!"}` (HTTP 200 OK)

2.  **Blocked Request (6th request within 60 seconds):**
    After making 5 successful requests, send one more request immediately:

    **For GET `/api/resource`:**
    ```bash
    curl -H "X-User-ID: user123" http://localhost:8081/api/resource
    ```
    Expected response: `{"error":"Too many requests. Please try again later."}` (HTTP 429 Too Many Requests)

    **For POST `/api/data`:**
    ```bash
    curl -X POST -H "X-User-ID: user123" -H "Content-Type: application/json" -d "{}" http://localhost:8081/api/data
    ```
    Expected response: `{"error":"Too many requests. Please try again later."}` (HTTP 429 Too Many Requests)

3.  **Rate Limit Reset:**
    Wait for 60 seconds after the last blocked request, then try another request. It should succeed again as the window resets.

### Testing Public Endpoint (`/public`)

This endpoint is not rate-limited. You can call it any number of times.

```bash
curl http://localhost:8081/public
```
Expected response: `{"message":"This is a public resource, not rate-limited."}` (HTTP 200 OK)
