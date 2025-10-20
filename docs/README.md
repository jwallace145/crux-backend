# Crux Backend API Documentation

This directory contains the OpenAPI specification and interactive API documentation for the Crux Backend.

## Interactive API Documentation

The API documentation is powered by [Scalar](https://github.com/scalar/scalar), providing a beautiful and interactive interface to explore and test the API endpoints.

### Accessing the Documentation

When the server is running, you can access the interactive API documentation at:

```
http://localhost:3000/docs
```

Or if deployed:

```
https://your-domain.com/docs
```

### Features

The interactive documentation provides:

- **Browse All Endpoints**: View all available API endpoints organized by tags (Health, Authentication, Users, Climbs)
- **Try It Out**: Test API endpoints directly from the browser with a built-in API client
- **Request/Response Examples**: See example requests and responses for each endpoint
- **Schema Documentation**: Detailed documentation of all request and response schemas
- **Authentication Testing**: Test authenticated endpoints using cookie-based JWT authentication
- **Search**: Quick search to find specific endpoints or schemas
- **Multiple Environments**: Switch between local and deployed servers for testing

### Switching Between Environments

The documentation supports testing against multiple server environments:

**Available Servers:**
- **Local Development** (`http://localhost:3000`) - Test local changes during development
- **Development Server** (`http://dev-api.cruxproject.io`) - Test deployed changes in the dev environment

**How to Switch:**
1. Open the API documentation at `http://localhost:3000/docs`
2. Look for the **server dropdown** at the top of the documentation interface
3. Select the environment you want to test against:
   - Choose "Local development server" to test your local changes
   - Choose "Development server" to test deployed changes at dev-api.cruxproject.io
4. All subsequent API requests in the "Try It Out" feature will be sent to the selected server

This allows you to:
- **Test locally** before deploying changes
- **Verify deployed changes** without switching documentation sources
- **Compare behavior** between local and deployed environments

## OpenAPI Specification

The OpenAPI specification file is available at:

```
http://localhost:3000/docs/openapi.yaml
```

You can use this specification with various tools:

- Import into Postman or Insomnia
- Generate client SDKs using OpenAPI Generator
- Use with other API documentation tools

## Local Development

### Starting the Server

```bash
# Start the server
go run main.go

# The server will log the documentation URL on startup:
# "Starting Crux API server" address="0.0.0.0:3000" docs="http://localhost:3000/docs"
```

### Updating the Documentation

To update the API documentation:

1. Edit the OpenAPI specification file: `docs/openapi.yaml`
2. Restart the server
3. Refresh the documentation page in your browser

The OpenAPI specification follows the [OpenAPI 3.0.3 specification](https://swagger.io/specification/).

## API Overview

### Available Endpoints

#### Health
- `GET /health` - Health check endpoint

#### Authentication
- `POST /login` - Authenticate with username/email and password
- `POST /logout` - Log out and revoke session
- `POST /refresh` - Refresh access token

#### Users
- `POST /users` - Create a new user account
- `GET /users` - Get authenticated user profile

#### Climbs
- `POST /climbs` - Log a new climbing activity
- `GET /climbs` - Retrieve climbs with optional date filtering

### Authentication

The API uses JWT tokens stored in HTTP-only cookies:

- **Access Token**: Short-lived token for API requests (automatically included in cookies)
- **Refresh Token**: Long-lived token to obtain new access tokens (automatically included in cookies)

When testing authenticated endpoints in the documentation:

1. First call `POST /login` to authenticate
2. The response will set cookies automatically
3. Subsequent requests will include the authentication cookies
4. Use `POST /refresh` to get a new access token when needed
5. Call `POST /logout` to end the session

### Response Format

All API endpoints return a standardized response structure:

```json
{
  "service_name": "crux-backend",
  "version": "1.0.0",
  "environment": "development",
  "api_name": "endpoint_name",
  "request_id": "unique-request-id",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "success",
  "message": "Human-readable message",
  "data": {
    // Endpoint-specific data
  }
}
```

For errors:

```json
{
  "service_name": "crux-backend",
  "version": "1.0.0",
  "environment": "development",
  "api_name": "endpoint_name",
  "request_id": "unique-request-id",
  "timestamp": "2024-01-01T12:00:00Z",
  "status": "error",
  "message": "Human-readable error message",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": {}
  }
}
```

## Customization

### Scalar Theme

The documentation UI uses the "purple" theme with a modern layout. You can customize this by editing the Scalar configuration in `internal/services/docs/docs_handler.go`:

```javascript
data-configuration='{"theme":"purple","layout":"modern","showSidebar":true}'
```

Available themes: `default`, `alternate`, `moon`, `purple`, `solarized`, `bluePlanet`, `saturn`, `kepler`, `mars`, `deepSpace`

Available layouts: `modern`, `classic`

## Production Deployment

In production, ensure the OpenAPI specification is included in your deployment:

1. The `docs/openapi.yaml` file should be included in your Docker image or deployment package
2. The documentation endpoints (`/docs` and `/docs/openapi.yaml`) are accessible
3. Update the server URL in the OpenAPI spec to point to your production domain

## Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [Scalar Documentation](https://github.com/scalar/scalar)
- [Crux Backend Repository](https://github.com/jwallace145/crux-backend)
