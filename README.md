# PDF Generator Service

## Overview
This project is a Go-based PDF generation service that converts HTML templates into PDF documents. It uses Go's native packages (`net/http`) to create a REST API and `chromedp` to render HTML and generate PDFs using Chromium. The service accepts a `multipart/form-data` request with an HTML template file and JSON data, renders the HTML with the provided data, and generates a PDF. The project follows a clean architecture pattern, separating concerns into layers (handlers, services, models, and infrastructure), and is containerized using Docker.

The service is deployed as `reg.hamsaa.ir/pdf-generator:latest`.

## Project Structure
```
pdf-generator/
├── internal/
│   ├── handlers/          # HTTP handlers (presentation layer)
│   │   └── pdf_handler.go
│   ├── infrastructure/    # External dependencies (infrastructure layer)
│   │   └── chromedp.go
│   ├── models/            # Data models (domain layer)
│   │   └── pdf_request.go
│   ├── services/          # Business logic (application layer)
│   │   └── pdf_service.go
├── templates/             # Sample HTML templates (for testing)
│   └── service_request.html
├── vendor/                # Vendored Go dependencies
├── docker-compose.yml     # Docker Compose configuration
├── Dockerfile             # Docker build configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
├── main.go                # Main application code
└── README.md              # Project documentation
```

## Prerequisites
- **Docker**: Required to build and run the containerized service.
- **Docker Compose**: Required to manage the service using `docker-compose.yml`.
- **curl** (optional): For testing the API endpoint.
- **Postman** (optional): For testing the API endpoint with a GUI.

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
## Notes
- The service requires Chromium to generate PDFs. The `CHROME_PATH` environment variable is set in `docker-compose.yml` to point to `/usr/bin/chromium-browser`.
- The generated PDFs are in A4 format (8.27 x 11.69 inches) with the background included.
- The service uses Go’s native `net/http` package for HTTP handling, following a clean architecture pattern.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details (if applicable).