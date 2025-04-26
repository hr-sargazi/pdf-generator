# PDF Generator Service

## Overview
This project is a Go-based PDF generation service that converts HTML templates into PDF documents. It uses the `gin` framework to create a REST API, `chromedp` to render HTML and generate PDFs using Chromium, and Docker for containerization. The service accepts JSON input containing service request details and generates a PDF based on an HTML template.

The project is containerized and deployed as `reg.hamsaa.ir/pdf-generator:latest`.

## Project Structure
```
pdf-generator/
├── templates/
│   └── service_request.html  # HTML template for PDF rendering
├── vendor/                   # Vendored Go dependencies
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker build configuration
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── main.go                   # Main application code
└── README.md                 # Project documentation
```

## Prerequisites
- **Docker**: Required to build and run the containerized service.
- **Docker Compose**: Required to manage the service using `docker-compose.yml`.
- **curl** (optional): For testing the API endpoint.

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
docker build -t reg.hamsaa.ir/pdf-generator:latest .
```

Alternatively, you can pull the pre-built image from the registry:
```bash
docker pull reg.hamsaa.ir/pdf-generator:latest
```

If the registry requires authentication:
```bash
docker login reg.hamsaa.ir
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

## Usage

### API Endpoint
The service exposes a single endpoint for PDF generation:

- **URL**: `POST /generate-pdf`
- **Content-Type**: `application/json`
- **Response**: PDF file (`application/pdf`)

### Request Body
The endpoint expects a JSON payload with the following structure:
```json
{
    "customer_name": "string",
    "customer_number": "string",
    "customer_banker_name": "string",
    "customer_has_ubank_contract": "string",
    "service_request_type": "string",
    "service_request_title": "string",
    "service_request_number": "string",
    "service_request_date": "string",
    "service_request_time": "string",
    "service_request_status": "string",
    "service_request_details": "string"
}
```

### Example Request
Create a file named `request.json` with sample data:
```json
{
    "customer_name": "علی سلطانمهر",
    "customer_number": "679994AF9",
    "customer_banker_name": "فرزاد براکی",
    "customer_has_ubank_contract": "دارد",
    "service_request_type": "درخواست",
    "service_request_title": "انتقال وجه به سپرده ارزی شعبه سامان جهت تسویه مالیات",
    "service_request_number": "139904110471",
    "service_request_date": "1399/04/10",
    "service_request_time": "11:50",
    "service_request_status": "دروسافت",
    "service_request_details": "انتقال وجه از سپرده دلاری شعبه سامان به حساب مالیاتی جهت تسویه مالیات با کد اقتصادی YAFADMFY YAFADMFY به مبلغ 1,000,000,000 ريال در تاریخ 1399/04/10"
}
```

Send the request using `curl`:
```bash
curl -X POST http://localhost:8080/generate-pdf \
    -H "Content-Type: application/json" \
    -d @request.json \
    --output service_request.pdf
```

### Response
The response will be a PDF file named `service_request.pdf`. Open the file to verify the content.

## Testing in Docker
The service is designed to run in a Docker environment with Chromium for PDF rendering. The `docker-compose.yml` file simplifies management.

### Start the Service
```bash
docker-compose up -d
```

### Test the Endpoint
Use the `curl` command above to test the `/generate-pdf` endpoint.

### Stop the Service
```bash
docker-compose down
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
- **Font Issues**: If Persian characters are not rendered correctly, ensure the `Vazir` font is included in the Docker image and referenced in `service_request.html`.

## Notes
- The service requires Chromium to generate PDFs. The `CHROME_PATH` environment variable is set in `docker-compose.yml` to point to `/usr/bin/chromium-browser`.
- The `service_request.html` template uses the `Vazir` font for Persian text. Ensure the font is available in the container for accurate rendering.
- The generated PDFs are in A4 format (8.27 x 11.69 inches) with the background included.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details (if applicable).