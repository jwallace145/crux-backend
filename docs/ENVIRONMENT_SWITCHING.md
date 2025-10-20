# Environment Switching Guide

This guide explains how to use the environment switching feature in the Scalar API documentation to test against different server environments.

## Overview

The API documentation supports testing against multiple server environments without needing to change the documentation source or configuration. This is achieved through OpenAPI's multi-server support, which Scalar automatically renders as a dropdown selector.

## Available Environments

### 1. Local Development Server
- **URL:** `http://localhost:3000`
- **Use Case:** Testing local changes during development
- **Requirements:** Local server must be running (`go run main.go`)
- **Best For:**
  - Developing new features
  - Debugging issues locally
  - Testing before committing changes

### 2. Development Server
- **URL:** `http://dev-api.cruxproject.io`
- **Use Case:** Testing deployed changes in the development environment
- **Requirements:** Changes must be deployed to dev environment
- **Best For:**
  - Verifying deployed changes
  - Integration testing
  - Demonstrating features to team members

## How to Switch Between Environments

### Step-by-Step Instructions

1. **Open the Documentation**
   ```
   http://localhost:3000/docs
   ```
   Or visit the deployed documentation at your GitHub Pages URL.

2. **Locate the Server Dropdown**
   - At the top of the Scalar documentation interface, you'll see a dropdown labeled "Server" or showing the currently selected server URL
   - The dropdown is typically located near the top-left of the interface

3. **Select Your Target Environment**
   - Click the dropdown to see all available servers
   - Select the environment you want to test against:
     - "Local development server - Use this to test local changes"
     - "Development server - Use this to test deployed changes"

4. **Test Your Endpoints**
   - Navigate to any endpoint in the documentation
   - Click "Try It" or the execute button
   - The request will be sent to the currently selected server
   - Responses will come from that environment

## Use Cases

### Testing a New Feature Locally

1. Start your local server: `go run main.go`
2. Open the docs: `http://localhost:3000/docs`
3. Select "Local development server" from the dropdown
4. Test your new endpoint with sample data
5. Verify the response matches expectations

### Verifying a Deployed Feature

1. Deploy your changes to dev: `make deploy`
2. Open the docs (local or deployed version)
3. Select "Development server" from the dropdown
4. Test the endpoint against the deployed environment
5. Confirm the feature works in the deployed context

### Comparing Local vs Deployed Behavior

1. Have both environments ready (local running, dev deployed)
2. Open the documentation
3. Test an endpoint with "Local development server" selected
4. Switch to "Development server"
5. Test the same endpoint with identical parameters
6. Compare responses to identify differences

## Tips & Best Practices

### For Local Testing

- **Always check the server is running** before testing
  ```bash
  curl http://localhost:3000/health
  ```
- **Use consistent test data** to ensure reproducible results
- **Check the browser console** for any CORS or network errors

### For Deployed Testing

- **Verify deployment status** before testing
  ```bash
  make ecs-status
  ```
- **Check deployment logs** if endpoints fail
  ```bash
  make ecs-logs
  ```
- **Remember authentication sessions** are separate between environments

### For Authentication Testing

When testing authenticated endpoints:

1. **Login First** - Use `POST /login` to authenticate
2. **Note the Environment** - Cookies are domain-specific
3. **Re-login if Switching** - If you switch environments, you'll need to login again
4. **Check Cookie Storage** - Browser dev tools → Application → Cookies

## Troubleshooting

### "Failed to fetch" Error

**Cause:** The selected server is not running or not accessible

**Solution:**
- For local: Ensure `go run main.go` is running
- For dev: Check deployment status with `make ecs-status`
- Check browser console for specific error messages

### CORS Errors

**Cause:** Request blocked by CORS policy

**Solution:**
- For local: CORS should be configured to allow `localhost:3000`
- For dev: Check CORS middleware configuration
- Verify the server's CORS settings in `internal/utils/middleware.go`

### Authentication Not Working

**Cause:** Cookies not being set or sent

**Solution:**
- Ensure you're using the correct server (cookies are domain-specific)
- Check browser dev tools → Application → Cookies
- Login again after switching servers
- Verify cookies are HTTP-only and properly configured

### Different Responses Between Environments

**Cause:** Data or configuration differences between environments

**Investigation Steps:**
1. Check database state in both environments
2. Verify environment variables match expectations
3. Check application logs in both environments
4. Ensure both environments are running the same code version

## Technical Details

### OpenAPI Servers Configuration

The server switching is implemented in the OpenAPI specification (`docs/openapi.yaml`):

```yaml
servers:
  - url: http://localhost:3000
    description: Local development server - Use this to test local changes
  - url: http://dev-api.cruxproject.io
    description: Development server - Use this to test deployed changes
```

Scalar automatically reads this configuration and provides a UI element to switch between servers.

### How Requests are Routed

When you select a server and make a request:

1. Scalar reads the selected server URL
2. Prepends it to the endpoint path (e.g., `http://localhost:3000/users`)
3. Sends the request to that full URL
4. Returns the response to the documentation UI

### Cookie Handling

- Cookies are automatically included in requests to the same domain
- Switching servers may require re-authentication
- HTTP-only cookies cannot be accessed via JavaScript (security feature)

## Adding New Environments

To add a new environment (e.g., staging or production):

1. Edit `docs/openapi.yaml`
2. Add a new server entry:
   ```yaml
   servers:
     - url: http://localhost:3000
       description: Local development server
     - url: http://dev-api.cruxproject.io
       description: Development server
     - url: https://staging-api.cruxproject.io
       description: Staging environment
   ```
3. Restart the server or redeploy
4. The new environment will appear in the dropdown automatically

## Related Documentation

- [Main README](../README.md) - Project overview and setup
- [API Documentation README](./README.md) - Complete API documentation guide
- [Scalar Documentation](https://guides.scalar.com/scalar/introduction) - Official Scalar docs
- [OpenAPI Specification](https://swagger.io/specification/) - OpenAPI standard

## Questions or Issues?

If you encounter issues with environment switching:

1. Check this guide's troubleshooting section
2. Review the OpenAPI specification in `docs/openapi.yaml`
3. Check browser console for error messages
4. Open a GitHub issue with details about the problem
