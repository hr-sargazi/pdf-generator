# PDF Generator Service

## Overview
This project is a Go-based PDF generation service that converts HTML templates into PDF documents. It uses Go's native packages (`net/http`) to create a REST API and `chromedp` to render HTML and generate PDFs using Chromium. The service accepts a `multipart/form-data` request with an HTML template file and JSON data, renders the HTML with the provided data, and generates a PDF. The project follows a clean architecture pattern, separating concerns into layers (handlers, services, models, and infrastructure), and is containerized using Docker.

The service includes comprehensive unit and integration tests to ensure reliability and correctness. Tests are executed as part of the Docker build process to catch issues early. The service is deployed as `pdf-generator:latest`.

## Project Structure
```
pdf-generator/
├── internal/
│   ├── handlers/          # HTTP handlers (presentation layer)
│   │   ├── pdf_handler.go
│   │   └── pdf_handler_test.go
│   ├── infrastructure/    # External dependencies (infrastructure layer)
│   │   ├── chromedp.go
│   │   └── chromedp_client_test.go
│   ├── models/            # Data models (domain layer)
│   │   └── pdf_request.go
│   ├── services/          # Business logic (application layer)
│   │   ├── pdf_service.go
│   │   └── pdf_service_test.go
├── templates/             # Sample HTML templates (for testing)
│   └── service_request.html
├── vendor/                # Vendored Go dependencies
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker build configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
├── main.go                # Main application code
├── main_test.go           # Tests for the main package
└── README.md              # Project documentation
```

## Prerequisites
- **Docker**: Required to build and run the containerized service.
- **Docker Compose**: Required to manage the service using `docker-compose.yml`.
- **curl** (optional): For testing the API endpoint.
- **Postman** (optional): For testing the API endpoint with a GUI.
- **Go** (optional): Required if running tests locally or modifying the code.

## Setup Instructions

### 1. Clone the Repository
Clone the repository to your local machine:
```bash
git clone <repository-url>
cd pdf-generator
```

### 2. Build the Docker Image (Optional)
If you need to rebuild the image locally, ensure you have the `Dockerfile` and run:
```bash
docker build -t pdf-generator:latest .
```

This command will:
- Build the Go application.
- Run all tests (unit and integration) during the build process, ensuring the application is thoroughly validated.
- Create the final container image with Chromium for runtime PDF generation.

Alternatively, you can pull the pre-built image from the registry:
```bash
docker pull pdf-generator:latest
```

### 3. Run the Service with Docker Compose
Use the provided `docker-compose.yml` to start the service:
```bash
docker-compose up -d
```

This will start the `pdf-generator` service on port `8080`.

### 4. Verify the Service is Running
Check the status of the container:
```bash
docker-compose ps
```

View the logs to ensure the service started successfully:
```bash
docker-compose logs
```

## Testing

### Overview
The project includes a comprehensive test suite to ensure reliability and correctness:
- **Unit Tests**: Cover individual components (`handlers`, `services`, `infrastructure`, `main`) in isolation, mocking dependencies like `chromedp` to avoid requiring Chromium.
- **Integration Tests**: Test the full PDF generation flow with a real Chromium instance, ensuring the application works end-to-end.
- **Docker Build Testing**: Tests are executed as part of the Docker build process, catching issues before the image is created.

### Running Tests Locally
To run tests locally, ensure you have Go installed and Chromium available (for integration tests):
1. **Install Chromium** (if running integration tests):
   On macOS:
   ```bash
   brew install chromium
   ```
   On Ubuntu:
   ```bash
   sudo apt-get install chromium-browser
   ```
   On Alpine (if running in a container):
   ```bash
   apk add chromium
   ```

2. **Set Environment Variables** (if running integration tests):
   ```bash
   export CHROME_PATH=/path/to/chromium-browser
   export RUN_INTEGRATION_TESTS=true
   ```
   Replace `/path/to/chromium-browser` with the actual path (e.g., `/usr/bin/chromium-browser` on Ubuntu).

3. **Run Tests**:
   ```bash
   go test ./...
   ```
   To run with verbose output:
   ```bash
   go test -v ./...
   ```
   To measure test coverage:
   ```bash
   go test -cover ./...
   ```

### Test Structure
- **main_test.go**: Tests the `main` package, verifying the HTTP server setup and handler registration by sending requests to the `/generate-pdf` endpoint.
- **handlers/pdf_handler_test.go**: Tests the `PDFHandler`, covering successful PDF generation, invalid methods, missing template files, invalid JSON data, and error handling.
- **services/pdf_service_test.go**: Tests the `PDFService`, ensuring HTML rendering and PDF generation logic works correctly, including edge cases like empty templates and invalid data.
- **infrastructure/chromedp_client_test.go**: Tests the `ChromedpClient`, including unit tests for Chrome path selection and an integration test for PDF generation with Chromium.

### Testing in Docker Build
The Dockerfile includes a testing step in the builder stage:
- **Chromium Installation**: Chromium is installed in the builder stage to support integration tests.
- **Test Execution**: The command `RUN RUN_INTEGRATION_TESTS=true go test ./...` runs all tests, including integration tests, during the build.
- **Failure Handling**: If any test fails, the build process will stop, ensuring only a working application is built into the image.

Example build command:
```bash
docker build -t pdf-generator:latest .
```

If tests fail, the build output will show the failure details, allowing you to debug and fix issues before deployment.

## Usage

### API Endpoint
The service exposes a single endpoint for PDF generation:

- **URL**: `POST /generate-pdf`
- **Content-Type**: `multipart/form-data`
- **Response**: PDF file (`application/pdf`)

### Request Body
The endpoint expects a `multipart/form-data` request with the following fields:
- `template_file`: An HTML template file (e.g., `service_request.html`) that defines the structure of the PDF.
- `data`: A JSON string containing the data to populate the template.

#### Example HTML Template (`service_request.html`)
The `templates/service_request.html` file in the repository can be used as a template. It expects data fields like `customer_name`, `customer_number`, etc. Here’s a simplified example:
```html
<!DOCTYPE html>
<html lang="fa">
<head>
    <meta charset="UTF-8">
    <title>Service Request</title>
    <style>
        body { font-family: Arial, sans-serif; direction: rtl; text-align: right; }
        .container { width: 100%; max-width: 800px; margin: 0 auto; padding: 20px; }
        h1 { color: #3C4750; }
        p { color: #576977; }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.customer_name}}</h1>
        <p>شماره مشتری: {{.customer_number}}</p>
        <p>{{.service_request_details}}</p>
    </div>
</body>
</html>
```

#### Example JSON Data
The `data` field should be a JSON string matching the fields in your template:
```json
{
    "customer_name": "علی سامانیان",
    "customer_number": "679994AF9",
    "customer_banker_name": "فرداد درگاهی",
    "customer_has_ubank_contract": "خیر",
    "service_request_type": "انتقال وجه",
    "service_request_title": "انتقال وجه به حساب سرمایه گذاری",
    "service_request_number": "139904110471",
    "service_request_date": "1399/04/10",
    "service_request_time": "11:50",
    "service_request_status": "تکمیل شده",
    "service_request_details": "انتقال وجه از حساب سرمایه گذاری به حساب جاری شما با شماره حساب 1042049506 و مبلغ 1,000,000,000 ریال در تاریخ 1399/04/10 با کد رهگیری YAFADMYF YAFADMYF"
}
```

### Testing with Postman
1. **Create a New Request in Postman**:
   - Open Postman and create a new request.
   - Set the method to `POST`.
   - Set the URL to `http://localhost:8080/generate-pdf`.

2. **Configure the Request Body**:
   - Go to the **Body** tab and select `form-data`.
   - Add two key-value pairs:
     - Key: `template_file`, Value: Select the `service_request.html` file from your local system.
     - Key: `data`, Value: Paste the JSON string (as text):
       ```json
       {"customer_name":"علی سامانیان","customer_number":"679994AF9","customer_banker_name":"فرداد درگاهی","customer_has_ubank_contract":"خیر","service_request_type":"انتقال وجه","service_request_title":"انتقال وجه به حساب سرمایه گذاری","service_request_number":"139904110471","service_request_date":"1399/04/10","service_request_time":"11:50","service_request_status":"تکمیل شده","service_request_details":"انتقال وجه از حساب سرمایه گذاری به حساب جاری شما با شماره حساب 1042049506 و مبلغ 1,000,000,000 ریال در تاریخ 1399/04/10 با کد رهگیری YAFADMYF YAFADMYF"}
       ```

3. **Send the Request**:
   - Click the **Send** button.
   - Postman will send the request, and you should receive a PDF file in the response.

4. **Save the Response**:
   - Click **Save Response** and choose **Save to a file** to save the PDF as `dynamic_document.pdf`.
   - Open the PDF to verify the content.

### Testing with `curl`
You can also test the endpoint using `curl`:
```bash
curl -X POST http://localhost:8080/generate-pdf \
    -F "template_file=@/path/to/service_request.html" \
    -F "data={\"customer_name\":\"علی سامانیان\",\"customer_number\":\"679994AF9\",\"customer_banker_name\":\"فرداد درگاهی\",\"customer_has_ubank_contract\":\"خیر\",\"service_request_type\":\"انتقال وجه\",\"service_request_title\":\"انتقال وجه به حساب سرمایه گذاری\",\"service_request_number\":\"139904110471\",\"service_request_date\":\"1399/04/10\",\"service_request_time\":\"11:50\",\"service_request_status\":\"تکمیل شده\",\"service_request_details\":\"انتقال وجه از حساب سرمایه گذاری به حساب جاری شما با شماره حساب 1042049506 و مبلغ 1,000,000,000 ریال در تاریخ 1399/04/10 با کد رهگیری YAFADMYF YAFADMYF\"}" \
    --output dynamic_document.pdf
```

## Debugging
- **View Logs**: Check for errors in the container logs:
  ```bash
  docker-compose logs
  ```
- **Access the Container**: Debug inside the container:
  ```bash
  docker exec -it pdf-generator sh
  ```
- **Test Failures**: If the Docker build fails due to test errors, check the build output for details. Common issues include:
  - Missing Chromium (for integration tests): Ensure Chromium is installed in the builder stage.
  - Test configuration: Verify `CHROME_PATH` and `RUN_INTEGRATION_TESTS` environment variables are set correctly.

## Notes
- The service requires Chromium to generate PDFs. The `CHROME_PATH` environment variable is set in both the Dockerfile and `docker-compose.yml` to point to `/usr/bin/chromium-browser`.
- The generated PDFs are in A4 format (8.27 x 11.69 inches) with the background included.
- The service uses Go’s native `net/http` package for HTTP handling, following a clean architecture pattern.
- Tests are executed during the Docker build to ensure the application is reliable before deployment.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details (if applicable).
